import asyncio
from nats.aio.client import Client as NATS

from lib.db.jetstream_publisher import JetStreamPublisher
from lib.gen_types.semantic_data_pb2 import SemanticData, FileType
from nats.errors import TimeoutError, NoRespondersError
from nats.js.api import StreamConfig, RetentionPolicy
from nats.js.errors import BadRequestError
import logging
import os
from dotenv import load_dotenv

# Load environment variables from .env file
load_dotenv()

nats_url = os.getenv('NATS_CLIENT_URL', 'nats://nats:4222')
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


def main():
    # Instantiate the publisher
    logger.info(f"{semantic_stream_name} - {semantic_stream_subject}")
    publisher = JetStreamPublisher(subject=semantic_stream_subject, stream_name=semantic_stream_name)

    # Connect to NATS
    publisher.connect()
    semntic_data = SemanticData(url="minio:tenant-c20a9f75-a363-40ea-86ef-eabcedbac7df:0493eb3a-4475-462e-a791-e47834ea7ba8-small.pdf",
                                document_id=976345414660126765,
                                connector_id= 975493320735424513,
                                file_type= FileType.PDF,
                                collection_name= "user_625ece7e042d4f40bd2588b16bec7be6")

    # Create a fake ChunkingData message
    # semntic_data = SemanticData(
    #     url="https://help.collaboard.app/sticky-notes",
    #     # url = "https://developer.apple.com/documentation/visionos/improving-accessibility-support-in-your-app",
    #     # url = "https://help.collaboard.app/what-is-collaboard",
    #     # url = "https://learn.microsoft.com/en-us/aspnet/core/tutorials/razor-pages/?view=aspnetcore-8.0",
    #     # url = "https://learn.microsoft.com/en-us/aspnet/core/tutorials/razor-pages/sql?view=aspnetcore-8.0&tabs=visual-studio",
    #     site_map="",
    #     search_for_sitemap=True,
    #     document_id=974396356851630081,
    #     file_type=FileType.URL,
    #     collection_name="user_id_998",
    #     model_name="sentence-transformers/paraphrase-multilingual-mpnet-base-v2",
    #     model_dimension=768
    # )

    logger.info(f"message being sent \n {semntic_data}")
    logger.info(f"{semantic_stream_name} - {semantic_stream_subject}")
    # Publish the message
    publisher.publish(semntic_data)

    # # Create a fake ChunkingData message
    # chunking_data = SemanticData(
    #     url="https://help.collaboard.app/extract-pages-from-a-document",
    #     site_map="",
    #     search_for_sitemap=True,
    #     document_id=993456788,
    #     file_type=FileType.URL,
    #     collection_name="tennant_id_3",
    #     model_name="sentence-transformers/paraphrase-multilingual-mpnet-base-v2",
    #     model_dimension=768
    # )

    # logger.info(f"message being sent \n {chunking_data}")

    # # Publish the message
    # await publisher.publish(chunking_data)

    #     # Create a fake ChunkingData message
    # chunking_data = SemanticData(
    #     url="https://help.collaboard.app/upload-images",
    #     site_map="",
    #     search_for_sitemap=True,
    #     document_id=993456788,
    #     file_type=FileType.URL,
    #     collection_name="tennant_id_5",
    #     model_name="sentence-transformers/paraphrase-multilingual-mpnet-base-v2",
    #     model_dimension=768
    # )

    # logger.info(f"message being sent \n {chunking_data}")

    # # Publish the message
    # await publisher.publish(chunking_data)

    # Close the connection
    publisher.close()


if __name__ == "__main__":
    main()
