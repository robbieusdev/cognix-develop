import time
import uuid
from typing import List, Dict, Optional
from urllib.parse import urlparse, parse_qs

from youtube_transcript_api import YouTubeTranscriptApi

from lib.db.db_document import DocumentCRUD
from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.chunked_item import ChunkedItem


class YTSemantic(BaseSemantic):
    def get_video_id(self, youtube_url: str) -> Optional[str]:
        """
        Extracts the video ID from a YouTube URL.

        Args:
            youtube_url (str): The URL of the YouTube video.

        Returns:
            Optional[str]: The video ID if found, else None.
        """
        parsed_url = urlparse(youtube_url)
        if parsed_url.hostname == 'youtu.be':
            return parsed_url.path[1:]
        if parsed_url.hostname in ['www.youtube.com', 'youtube.com']:
            if parsed_url.path == '/watch':
                return parse_qs(parsed_url.query).get('v', [None])[0]
            if parsed_url.path[:7] == '/embed/':
                return parsed_url.path.split('/')[2]
            if parsed_url.path[:3] == '/v/':
                return parsed_url.path.split('/')[2]
        return None

    def get_youtube_transcript(self, video_url: str) -> str | None: #Optional[List[Dict[str, str]]]:
        """
        Fetches the transcript of a YouTube video.

        Args:
            video_url (str): The URL of the YouTube video.

        Returns:
            Optional[List[Dict[str, str]]]: The list of transcript entries if successful, else None.
        """
        video_id = self.get_video_id(video_url)
        if not video_id:
            print("Invalid YouTube URL")
            return None
        try:
            transcript_list = YouTubeTranscriptApi.get_transcript(video_id)
            # return transcript_list
            transcript = ""
            for segment in transcript_list:
                transcript += segment['text'] + "\n"
            return transcript.strip()
        except Exception as e:
            self.logger.error(f"❌ {e}")
            return None

    def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int, cockroach_url: str) -> int:
        try:
            start_time = time.time()  # Record the start time
            self.logger.info(f"extracting transcript from: {data.url}")

            content = self.get_youtube_transcript(data.url)

            collected_items = 0
            chunking_session = uuid.uuid4()
            document_crud = DocumentCRUD(cockroach_url)

            if content:
                # TODO: VERY IMPORTANT
                # We need content's summarization

                collected_data = [ChunkedItem(url=data.url, content=content)]

                collected_items = self.store_collected_data(data=data, document_crud=document_crud,
                                                            collected_data=collected_data,
                                                            chunking_session=chunking_session,
                                                            ack_wait=ack_wait,
                                                            full_process_start_time=full_process_start_time,
                                                            split_data=True)
                self.logger.debug(f"transcript \n {content}")
            else:
                self.store_collected_data_none(data=data, document_crud=document_crud, chunking_session=chunking_session)

            self.log_end(collected_items, start_time)
            return collected_items
            # (if transcript: 1 else: 0)
        except Exception as e:
            self.logger.error(f"❌ Failed to process semantic data: {e}")