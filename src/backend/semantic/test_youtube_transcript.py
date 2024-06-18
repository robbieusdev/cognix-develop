import asyncio
import datetime
import logging
import os
import time

from lib.gen_types.semantic_data_pb2 import SemanticData, FileType
from lib.semantic.semantic_factory import SemanticFactory


# get log level from env
log_level_str = os.getenv('LOG_LEVEL', 'ERROR').upper()
log_level = getattr(logging, log_level_str, logging.INFO)
# get log format from env
log_format = os.getenv('LOG_FORMAT', '%(asctime)s - %(name)s - %(levelname)s - %(message)s')
# Configure logging
logging.basicConfig(level=log_level, format=log_format)
logger = logging.getLogger(__name__)


async def main():
    semantic_data = SemanticData()
    semantic_data.document_id = 123
    semantic_data.url = ("https://www.youtube.com/watch?v=UbDyjIIGaxQ")
    semantic_data.file_type = FileType.YT

    # performing semantic analysis on the source
    semantic = SemanticFactory.create_semantic_analyzer(semantic_data.file_type)

    entities_analyzed = semantic.analyze(data=semantic_data, full_process_start_time=time.time(),
                                         ack_wait=30, cockroach_url="")


if __name__ == "__main__":
    asyncio.run(main())
