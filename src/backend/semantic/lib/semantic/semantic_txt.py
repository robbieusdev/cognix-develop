from lib.db.milvus_db import Milvus_DB
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
from pathlib import Path
import time
import pymupdf

class TXTSemantic(BaseSemantic):
    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        self.analyze_doc(data=data, full_process_start_time=full_process_start_time, ack_wait=ack_wait,
                         cockroach_url=cockroach_url)