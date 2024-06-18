import asyncio

from lib.db.jetstream_publisher import JetStreamPublisher
from lib.db.milvus_db import Milvus_DB
from lib.db.db_document import DocumentCRUD, Document
from lib.gen_types.semantic_data_pb2 import SemanticData, FileType
from lib.semantic.semantic_base import BaseSemantic
from lib.spider.spider_bs4 import BS4Spider  # Ensure you import the BS4Spider class correctly
import time, uuid, logging, datetime

from readiness_probe import ReadinessProbe


class URLSemantic(BaseSemantic):
    async def analyze(self, data: SemanticData, full_process_start_time: float, ack_wait: int,
                      cockroach_url: str) -> int:
        try:
            self.logger.info("Analyzing URL")
            start_time = time.time()  # Record the start time
            self.logger.info(f"Starting BS4Spider URL: {data.url}")

            spider = BS4Spider(data.url)

            # TODO: collected_data shall contain only the list of URLs in case of data.url_recursive
            collected_data = spider.process_page(data.url, data.url_recursive)

            collected_items = 0
            if not collected_data:
                self.logger.warning(f"😱 BS4Spider was not able to retrieve any content for {data.url}, switching to "
                                    f"SeleniumSpider")
                self.logger.warning(
                    "😱 SeleniumSpider is disabled, shall be re-enabled and tested as it is not working 100%")
                # self.logger.info(f"Starting SeleniumSpider for: {data.url}")
                # spider = SeleniumSpider(data.url)
                # collected_data = spider.process_page(data.url)

            chunking_session = uuid.uuid4()
            document_crud = DocumentCRUD(cockroach_url)

            if collected_data:
                # we store only the entities on the database
                # and we send a message to semantic for each collected data item
                # this way we avoid long run and NATS timeouts
                # step 1. prepare the ground delete all the documents with parent_id = data.document_id
                # delete previous added child (chunks) documents
                document_crud.delete_by_parent_id(data.document_id)
                if data.url_recursive:
                    #step 2 we isert all the docs
                    # Create Document objects from ChunkedItem objects
                    documents_to_insert = []
                    for item in collected_data:
                        documents_to_insert.append(
                            Document(parent_id=data.document_id, connector_id=data.connector_id, source_id="",
                                     url=item.url,
                                     signature="",
                                     chunking_session=chunking_session, analyzed=False))
                    documents_to_send = document_crud.insert_documents_batch(documents_to_insert)

                    # step 3 prepare to send the docs
                    publisher = JetStreamPublisher(subject=self.semantic_stream_subject,
                                                   stream_name=self.semantic_stream_name)
                    await publisher.connect()

                    for doc in documents_to_send:
                        self.logger.info(f"doc id {doc.id}")
                        semantic_data = SemanticData(
                            url=doc.url,
                            document_id=doc.id,
                            url_recursive=False,
                            # TODO ADD PARENT ID TO PROTO ????
                            connector_id=data.connector_id,
                            file_type=FileType.URL,
                            collection_name=data.collection_name)
                        await publisher.publish(semantic_data)
                        self.logger.info("✉️ sending message to jetstream")
                    await publisher.close()

                    collected_items = 1
                    # TODO: Important. Now, returning collected_items = 1 the status of the connector
                    # will be set to COMPLETED_SUCCESSFULLY which is not true because
                    # all the messages for each url still need to be processed
                    # also each of the single url will set again the connector to successfully completed
                    # we shall set successfully completed only when all the parent ids of the chunking session
                    # are analyzed. Invent something fancy like select top1 from docs where analyzed = false
                else:
                    collected_items = self.store_collected_data(data=data, document_crud=document_crud,
                                                                collected_data=collected_data,
                                                                chunking_session=chunking_session,
                                                                ack_wait=ack_wait,
                                                                full_process_start_time=full_process_start_time,
                                                                split_data=True)
            else:
                self.store_collected_data_none(data=data, document_crud=document_crud,
                                               chunking_session=chunking_session)

            self.log_end(collected_items, start_time)
            return collected_items
        except Exception as e:
            self.logger.error(f"❌ Failed to process semantic data: {e}")
