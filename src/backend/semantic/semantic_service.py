#region imports
import asyncio
import logging
import os
import threading
import time
import datetime
from dotenv import load_dotenv
from nats.aio.msg import Msg

from lib.db.db_connector import ConnectorCRUD, Status
from lib.db.db_document import DocumentCRUD
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_factory import SemanticFactory
from lib.db.jetstream_event_subscriber import JetStreamEventSubscriber
from readiness_probe import ReadinessProbe
#endregion

#region .env and logs
# Load environment variables from .env file
load_dotenv()

# get log level from env
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(levelname)s - %(name)s - %(funcName)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)

logger = logging.getLogger(__name__)
logger.info(f"Logging configured with level {log_level_str} and format {log_format}")


# loading from env env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'semantic_activity')
semantic_ack_wait = int(os.getenv('NATS_CLIENT_SEMANTIC_ACK_WAIT', '3600'))  # seconds
semantic_max_deliver = int(os.getenv('NATS_CLIENT_SEMANTIC_MAX_DELIVER', '3'))

cockroach_url = os.getenv('COCKROACH_CLIENT_DATABASE_URL',
                          'postgres://root:123@cockroach:26257/defaultdb?sslmode=disable')
#endregion

# Define the event handler function
async def semantic_event(msg: Msg):
    start_time = time.time()  # Record the start time
    connector_id = 0
    entities_analyzed = 0
    try:
        logger.info("🔥 starting semantic analysis..")
        # Deserialize the message
        semantic_data = SemanticData()
        semantic_data.ParseFromString(msg.data)
        logger.info(f"message: \n {semantic_data}")
        if semantic_data.model_name == "":
            logger.error(f"❌ no model nameeeeeeee")
            semantic_data.model_name = "paraphrase-multilingual-mpnet-base-v2"
            semantic_data.model_dimension = 768
            logger.warning(f"😱 Addning model name and dimension manually remove this code ASAP")

        # verify document id is valid otherwise we cannot process the message
        if semantic_data.document_id <= 0:
            logger.error(f"❌ failed to process semantic data error: document_id must value must be positive")
        else:
            # see if doc exists
            document_crud = DocumentCRUD(cockroach_url)
            document = document_crud.select_document(semantic_data.document_id)

            if document:
                # needed for th finally block
                connector_id = document.connector_id

                # update connector's status
                connector_crud = ConnectorCRUD(cockroach_url)
                connector = connector_crud.select_connector(document.connector_id)
                last_successful_index_date = connector.last_successful_analyzed
                connector_crud.update_connector(document.connector_id,
                                                status=Status.PROCESSING,
                                                last_update=datetime.datetime.utcnow())

                # performing semantic analysis on the source
                semantic = SemanticFactory.create_semantic_analyzer(semantic_data.file_type)
                entities_analyzed = await semantic.analyze(data=semantic_data, full_process_start_time=start_time,
                                                     ack_wait=semantic_ack_wait, cockroach_url=cockroach_url)

                # if entities_analyzed == 0 this means no data was stored in the vector db
                # we shall find a way to tell the user, most likely put the message in the dead letter

                if entities_analyzed is None:
                    logger.error(f"❌ entities_analyzed is none!!!!!")
                    entities_analyzed = 0
                # updating again the connector
                a =connector_crud.update_connector(connector_id,
                                                status=Status.COMPLETED_SUCCESSFULLY,
                                                last_successful_analyzed=datetime.datetime.utcnow(),
                                                last_update=datetime.datetime.utcnow(),
                                                total_docs_analyzed=entities_analyzed # TODO: we are storing the total entities in total docs. one doc will probably generate more than one chunk
                                                )
                logger.info(a)
            else:
                logger.error(
                    f"❌ failed to process semantic data error: document_id {semantic_data.document_id} not valid")
        # Acknowledge the message when done
        await msg.ack_sync()
        logger.info(f"👍 message acknowledged successfully, total entities stored {entities_analyzed}")
    except Exception as e:
        error_message = str(e) if e else "Unknown error occurred"
        logger.error(f"❌ failed to process semantic data error: {error_message}")
        if msg:  # Ensure msg is not None before awaiting
            await msg.nak()
        try:
            if connector_id != 0:
                connector_crud = ConnectorCRUD(cockroach_url)
                connector_crud.update_connector(connector_id,
                                                status=Status.COMPLETED_WITH_ERRORS,
                                                last_update=datetime.datetime.utcnow())
        except Exception as e:
            error_message = str(e) if e else "Unknown error occurred"
            logger.error(f"❌ failed to process semantic data error: {error_message}")
    finally:
        end_time = time.time()  # Record the end time
        elapsed_time = end_time - start_time
        logger.info(f"⏰⏰ total elapsed time: {elapsed_time:.2f} seconds")


# TODO: IMPORTANT WHEN IT DOES NOT CONNECT TO COCKROACH IS PROCESSING!!!!!
# Andri shall make a fix query on the db if status is processing and max ack wait is more than last update
# then it means the process hanged but if max retries (from nats) has not reached it's limit
# then nats will post again the message
# orchestrator shall update to UNABLE_TO_PROCESS, if nats will post again the message this service will
# care to set the correct status

async def main():
    # Start the readiness probe server in a separate thread
    readiness_probe = ReadinessProbe()
    readiness_probe_thread = threading.Thread(target=readiness_probe.start_server, daemon=True)
    readiness_probe_thread.start()

    # circuit breaker for chunking
    # if for reason nats won't be available
    # semantic will wait till nats will be up again
    while True:
        logger.info("🛠️ service starting..")
        try:
            # subscribing to jet stream
            subscriber = JetStreamEventSubscriber(
                nats_url=nats_url,
                stream_name=semantic_stream_name,
                subject=semantic_stream_subject,
                connect_timeout=nats_connect_timeout,
                reconnect_time_wait=nats_reconnect_time_wait,
                max_reconnect_attempts=nats_max_reconnect_attempts,
                ack_wait=semantic_ack_wait,
                max_deliver=semantic_max_deliver,
                proto_message_type=SemanticData
            )

            subscriber.set_event_handler(semantic_event)
            await subscriber.connect_and_subscribe()

            # todo add an event to JetStreamEventSubscriber to signal that connection has been established
            logger.info("🚀 service started successfully")

            while True:
                await asyncio.sleep(1)

        except KeyboardInterrupt:
            logger.info("🛑 Service is stopping due to keyboard interrupt")
        except Exception as e:
            logger.exception(f"💀 recovering from a fatal error: {e}. The process will restart in 5 seconds..")
            await asyncio.sleep(5)


if __name__ == "__main__":
    asyncio.run(main())
