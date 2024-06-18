import asyncio
import logging
import os
import uuid

from dotenv import load_dotenv
from sqlalchemy import func

from lib.db.db_document import DocumentCRUD
from lib.db.db_connector import ConnectorCRUD, Status

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

# loading from env env
nats_url = os.getenv('NATS_CLIENT_URL', 'nats://127.0.0.1:4222')
nats_connect_timeout = int(os.getenv('NATS_CLIENT_CONNECT_TIMEOUT', '30'))
nats_reconnect_time_wait = int(os.getenv('NATS_CLIENT_RECONNECT_TIME_WAIT', '30'))
nats_max_reconnect_attempts = int(os.getenv('NATS_CLIENT_MAX_RECONNECT_ATTEMPTS', '3'))
semantic_stream_name = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_NAME', 'semantic')
semantic_stream_subject = os.getenv('NATS_CLIENT_SEMANTIC_STREAM_SUBJECT', 'chunk_activity')
semantic_ack_wait = int(os.getenv('NATS_CLIENT_SEMANTIC_ACK_WAIT', '3600'))  # seconds
semantic_max_deliver = int(os.getenv('NATS_CLIENT_SEMANTIC_MAX_DELIVER', '3'))

cockroach_url = os.getenv('COCKROACH_CLIENT_DATABASE_URL',
                          'cockroachdb://root:123@cockroach:26257/defaultdb?sslmode=disable')


async def main():
    logger.info(cockroach_url)

    crud = DocumentCRUD(cockroach_url)
    # Insert a new document
    # new_doc_id = crud.insert_document(
    #     parent_id=None,
    #     connector_id=1,
    #     source_id='unique_source_id',
    #     url='http://example.com',
    #     signature='signature_example',
    #     chunking_session=uuid.uuid4(),
    #     analyzed=False,
    #     creation_date=func.now(),
    #     last_update=None
    # )
    # print(f"Inserted document ID: {new_doc_id}")

    # Select the document
    document = crud.select_document(974396356851630081)
    logger.info(f"Selected document: {document}")

    # # Update the document
    # crud.update_document(new_doc_id, url='http://newexample.com')
    #
    # # Delete the document
    # crud.delete_document(new_doc_id)
    # print(f"Deleted document ID: {new_doc_id}")

    # # Connector operations
    connector_crud = ConnectorCRUD(cockroach_url)
    # new_connector_id = connector_crud.insert_connector(
    #     credential_id=None,
    #     name='Connector Name',
    #     type='Connector Type',
    #     connector_specific_config={},
    #     refresh_freq=3600,
    #     user_id=uuid.uuid4(),
    #     tenant_id=None,
    #     disabled=False,
    #     last_successful_index_date=None,
    #     last_attempt_status=None,
    #     total_docs_indexed=0,
    #     creation_date=func.now(),
    #     last_update=None,
    #     deleted_date=None
    # )
    # print(f"Inserted connector ID: {new_connector_id}")

    connector = connector_crud.select_connector(document.connector_id)
    logger.info(f"Selected connector: {connector}")
    connector_crud.update_connector(connector.id, last_attempt_status=Status.SCAN_COMPLETED_SUCCESSFULLY)
    connector = connector_crud.select_connector(document.connector_id)
    logger.info(f"Selected connector: {connector}")

    # connector_crud.delete_connector(new_connector_id)
    # print(f"Deleted connector ID: {new_connector_id}")
    #


if __name__ == "__main__":
    asyncio.run(main())
