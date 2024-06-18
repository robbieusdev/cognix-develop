import asyncio
from nats.aio.client import Client as NATS
from lib.gen_types.semantic_data_pb2 import SemanticData, FileType
from nats.errors import TimeoutError, NoRespondersError
from nats.js.api import StreamConfig, RetentionPolicy
from nats.js.errors import BadRequestError
import logging
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'semantic_activity')
semantic_ack_wait = int(os.getenv('NATS_CLIENT_SEMANTIC_ACK_WAIT', '3600'))  # seconds
semantic_max_deliver = int(os.getenv('NATS_CLIENT_SEMANTIC_MAX_DELIVER', '3'))

# get log level from env
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)


class JetStreamPublisher:
    def __init__(self, subject, stream_name):
        self.logger = None
        self.subject = subject
        self.stream_name = stream_name
        self.nc = NATS()
        self.js = None
        logger.info(f"{semantic_stream_name} - {semantic_stream_subject}")

    async def connect(self):
        # Connect to NATS
        await self.nc.connect(servers=[nats_url],
                              connect_timeout=nats_connect_timeout,
                              reconnect_time_wait=nats_reconnect_time_wait,
                              max_reconnect_attempts=nats_max_reconnect_attempts)
        # Create JetStream context
        self.js = self.nc.jetstream()

        # Create the stream configuration
        stream_config = StreamConfig(
            name=self.stream_name,
            subjects=[self.subject],
            # A work-queue retention policy satisfies a very common use case of queuing up messages that are intended
            # to be processed once and only once. https://natsbyexample.com/examples/jetstream/workqueue-stream/go
            retention=RetentionPolicy.WORK_QUEUE
        )

        try:
            await self.js.add_stream(stream_config)
        except BadRequestError as e:
            if e.code == 400:
                self.logger.info(
                    "Jetstream stream was using a different configuration. Destroying and recreating with the right "
                    "configuration")
                try:
                    await self.js.delete_stream(stream_config.name)
                    await self.js.add_stream(stream_config)
                    self.logger.info("Jetstream stream re-created successfully")
                except Exception as e:
                    self.logger.exception(f"Exception while deleting and recreating Jetstream: {e}")

    async def publish(self, message):
        try:
            await self.js.publish(self.subject, message.SerializeToString())
            logger.info("Message published successfully!")
        except NoRespondersError:
            logger.error("❌ No responders available for request")
        except TimeoutError:
            logger.error("❌ Request to JetStream timed out")
        except Exception as e:
            logger.error(f"❌ Failed to publish message: {e}")

    async def close(self):
        await self.nc.close()