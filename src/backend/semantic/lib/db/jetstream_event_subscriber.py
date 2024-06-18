import logging

from google.protobuf import message as _message
from nats.aio.client import Client as NATS
from nats.aio.msg import Msg
from nats.js.api import ConsumerConfig, StreamConfig, AckPolicy, DeliverPolicy, RetentionPolicy
from nats.js.errors import BadRequestError

from readiness_probe import ReadinessProbe


class JetStreamEventSubscriber:
    def __init__(self, nats_url: str, stream_name: str, subject: str,
                 connect_timeout: int, reconnect_time_wait: int,
                 max_reconnect_attempts: int, ack_wait: int,
                 max_deliver: int, proto_message_type: _message.Message):
        self.nats_url = nats_url
        self.stream_name = stream_name
        self.subject = subject
        self.connect_timeout = connect_timeout
        self.reconnect_time_wait = reconnect_time_wait
        self.ack_wait = ack_wait
        self.max_reconnect_attempts = max_reconnect_attempts
        self.max_deliver = max_deliver
        self.proto_message_type = proto_message_type
        self.event_handler = None
        self.nc = NATS()
        self.js = None  # needs to be created in connect_and_subscribe
        self.logger = logging.getLogger(self.__class__.__name__)

    async def connect_and_subscribe(self):
        try:
            # Connect to NATS
            # TODO: implement callbacks from connect <--------------------  DO THIS :) 
            self.logger.info(f"🔌 connecting to nats endpoint {self.nats_url} ..")

            await self.nc.connect(servers=[self.nats_url],
                                  connect_timeout=self.connect_timeout,
                                  reconnect_time_wait=self.reconnect_time_wait,
                                  max_reconnect_attempts=self.max_reconnect_attempts)
                                  # ping_interval=460000)
            # https://github.com/nats-io/nats.py/discussions/299

            self.logger.info(f"successfully connected {self.nats_url}")

            # Create JetStream context
            self.js = self.nc.jetstream()

            # Create the stream configuration
            stream_config = StreamConfig(
                name=self.stream_name,
                subjects=[self.subject],
                # A work-queue retention policy satisfies a very common use case of queuing up messages that are
                # intended to be processed once and only once.
                # https://natsbyexample.com/examples/jetstream/workqueue-stream/go
                retention=RetentionPolicy.WORK_QUEUE
                # retention=RetentionPolicy.LIMITS
            )

            try:
                await self.js.add_stream(stream_config)
            except BadRequestError as e:
                if e.code == 400:
                    self.logger.warning("😱 jetstream stream was using a different configuration. Destroying and "
                                        "recreating with the right configuration")
                    try:
                        await self.js.delete_stream(stream_config.name)
                        await self.js.add_stream(stream_config)
                        self.logger.info("jetstream stream re-created successfully")
                    except Exception as e:
                        self.logger.exception(f"❌ Exception while deleting and recreating Jetstream: {e}")
        except Exception as e:
            self.logger.exception(f"❌ {e}")
            raise e

        # Define consumer configuration
        consumer_config = ConsumerConfig(
            # durable_name="durable_chunkdata", do not set herem, it is set in pull_subscribe
            ack_wait=self.ack_wait,  # 3600 seconds
            max_deliver=self.max_deliver,
            ack_policy=AckPolicy.EXPLICIT,
            # DeliverPolicy.ALL is mandatory when setting  retention=RetentionPolicy.WORK_QUEUE for StreamConfig
            deliver_policy=DeliverPolicy.ALL,
        )

        # Subscribe to the subject
        try:
            self.logger.info(f"subscribing to jetstream {self.stream_name} - {self.subject} ..")
            psub = await self.js.pull_subscribe(
                subject=self.subject,
                stream=stream_config.name,
                durable="worker",
                config=consumer_config,
            )
            self.logger.info(f"successfully subscribed to jetstream {self.stream_name} - {self.subject}")

            # psub.fetch()
            while True:
                try:
                    # await asyncio.sleep(2)
                    # notifying the readiness probe that the service is alive
                    ReadinessProbe().update_last_seen()
                    msgs = await psub.fetch(1, timeout=5)
                    self.logger.info(msgs)
                    for msg in msgs:
                        # ack will be done once the process is completed
                        # await msg.ack_sync()
                        await self.message_handler(msg)
                        self.logger.info(msg)
                except TimeoutError:
                    self.logger.info("waiting for incoming events..")
                    pass
        except Exception as e:
            self.logger.error(f"❌ can't connect or subscribe to {self.nats_url} {self.stream_name} {self.subject} {e}")
            raise e

    async def message_handler(self, msg: Msg):
        try:
            if self.event_handler:
                await self.event_handler(msg)
        except Exception as e:
            self.logger.exception(f"❌ failed to process message: {e}")

    def set_event_handler(self, event_handler):
        self.event_handler = event_handler

    async def close(self):
        await self.nc.close()

    async def flush(self):
        await self.nc.flush(2)

    async def disconnected_event(self):
        self.logger.warning('😱 Got disconnected!')

    async def reconnected_event(self, nc: NATS) -> None:
        self.logger.warning(f'😱 Got reconnected to {nc.connected_url.netloc}')

    async def error_event(self, e: Exception) -> None:
        self.logger.error(f"❌there was an error: {e}")

    async def closed_event(self, nc: NATS) -> None:
        self.logger.info("connection closed")
