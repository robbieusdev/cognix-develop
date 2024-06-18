import time
import uuid

import pymupdf4llm

from lib.db.db_document import DocumentCRUD
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.markdown_extractor import MarkdownSectionExtractor
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.chunked_item import ChunkedItem
from minio import Minio
from minio.error import S3Error


class PDFSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        return self.analyze_doc(data=data, full_process_start_time=full_process_start_time, ack_wait=ack_wait,
                         cockroach_url=cockroach_url)
