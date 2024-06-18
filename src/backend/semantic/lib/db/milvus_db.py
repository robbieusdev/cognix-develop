import logging
import os
import time
from typing import List

import grpc
from dotenv import load_dotenv
from numpy import int64
from pymilvus import connections, utility, FieldSchema, CollectionSchema, DataType, Collection

from lib.gen_types.semantic_data_pb2 import SemanticData
from lib.gen_types.embed_service_pb2 import EmbedRequest
from lib.gen_types.embed_service_pb2_grpc import EmbedServiceStub

# Load environment variables from .env file
load_dotenv()

# Get nats url from env 
milvus_alias = os.getenv("MILVUS_ALIAS", 'default')
milvus_host = os.getenv("MILVUS_HOST", "127.0.0.1")
milvus_port = os.getenv("MILVUS_PORT", "19530")
milvus_index_type = os.getenv("MILVUS_INDEX_TYPE", "DISKANN")
milvus_metric_type = os.getenv("MILVUS_METRIC_TYPE", "COSINE")

milvus_user = "root"
milvus_pass = "sq5/6<$Y4aD`2;Gba'E#"


embedder_grpc_host = os.getenv("EMBEDDER_GRPC_HOST", "localhost")
embedder_grpc_port = os.getenv("EMBEDDER_GRPC_PORT", "50051")


def is_connected():
    # TODO: get params from env (alias)
    return utility.connections.has_connection("default")


class Milvus_DB:
    def __init__(self):
        self.logger = logging.getLogger(self.__class__.__name__)

    def delete_by_document_id_and_parent_id(self, document_id: int64, collection_name: str):
        start_time = time.time()  # Record the start time
        self.logger.info(f"deleting all entities related to document {document_id}")
        try:
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port,
                user=milvus_user,
                password=milvus_pass
            )

            if utility.has_collection(collection_name):
                collection = Collection(collection_name)  # Get an existing collection.
                # collection.schema  # Return the schema.CollectionSchema of the collection.
                # collection.description  # Return the description of the collection.
                # collection.name  # Return the name of the collection.
                # collection.is_empty  # Return the boolean value that indicates if the collection is empty.
                self.logger.debug(f"collection: {collection_name} has {collection.num_entities} entities")

                # do not delete the entire collection
                # utility.drop_collection(collection_name)

                # Create expressions to find matching entities
                expr = f"document_id == {document_id} or parent_id == {document_id}"

                # Retrieve the primary keys of matching entities
                results = collection.query(expr, output_fields=["id"])
                ids_to_delete = [res["id"] for res in results]

                if ids_to_delete:
                    # Delete entities by their primary keys
                    delete_expr = f"id in [{', '.join(map(str, ids_to_delete))}]"
                    collection.delete(delete_expr)
                    collection.flush()
                    self.logger.debug(f"deleted documents with document_id or parent_id: {document_id}")
                else:
                    self.logger.debug(f"No documents found with document_id or parent_id: {document_id}")
        except Exception as e:
            self.logger.error(f"❌ failed to delete documents with document_id and parent_id {document_id}: {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

    def query(self, query: str, data: SemanticData) -> Collection:
        start_time = time.time()  # Record the start time
        try:
            # This way of adding data looks like extremely inefficient
            # We need to find a way to use the same connection across
            # different method calls
            # also not sure if and why the collection needs to be created every time
            # a pattern used by Milvus?
            # needs investigation
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port,
                user=milvus_user,
                password=milvus_pass
            )

            collection = Collection(name=data.collection_name)
            collection.load()

            # this makes a gRPC call to the embedding service
            # it's an architectural decision, so we have only one container handling
            # the embedding models. The container will require a significant amount of ram, cpu
            # and eventually it will run on gpu
            embedding = self.embedd(query, data.model_name)

            # fields = [
            #     FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
            #     FieldSchema(name="document_id", dtype=DataType.INT64),
            #     # text content expected format {"content":""}
            #     FieldSchema(name="content", dtype=DataType.JSON),
            #     FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=data.model_dimension),
            # ]

            result = collection.search(
                data=[embedding],  # Embed search value
                anns_field="vector",  # Search across embeddings
                # "M": 8, "ef": "top_k"
                param={"metric_type": f"{milvus_metric_type}",
                       "params": {"ef": 64}},
                limit=10,  # Limit to top_k results per search
                output_fields=["content"]
            )

            if self.logger.level == logging.DEBUG:
                answer = ""
                self.logger.debug("enumerating vector database results")
                for i, hits in enumerate(result):
                    for hit in hits:
                        # print(f"Query sentence: {sentences[i]}")
                        self.logger.debug(f"Nearest Neighbor Number {i}: {hit.entity.get('sentence')} ---- {hit.distance}\n")
                        answer = answer + hit.entity.get('sentence')
                self.logger.debug("end enumeration")
            return collection
        except Exception as e:
            self.logger.error(f"❌ {e}")
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"⏰ total elapsed time: {elapsed_time:.2f} seconds")

    def store_chunk(self, content: str, data: SemanticData, document_id: int64, parent_id: int64) -> bool:
        start_time = time.time()  # Record the start time
        success = False
        try:
            # self.ensure_connection()
            # if self.is_connected() == False:
            #     raise Exception("Connot connect to Milvus")

            # This way of adding data looks like extremely inefficient
            # We need to find a way to use the same connection across
            # different method calls 
            # also not sure if and why the collection needs to be created every time
            # a pattern used by Milvus?
            # needs investigation
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port,
                user=milvus_user,
                password=milvus_pass
            )

            fields = [
                FieldSchema(name="id", dtype=DataType.INT64, is_primary=True, auto_id=True),
                FieldSchema(name="document_id", dtype=DataType.INT64),
                FieldSchema(name="parent_id", dtype=DataType.INT64),
                # text content expected format {"content":""}
                FieldSchema(name="content", dtype=DataType.JSON),
                FieldSchema(name="vector", dtype=DataType.FLOAT_VECTOR, dim=data.model_dimension),
            ]

            # creating collection schema and adding the fields defined above
            schema = CollectionSchema(fields=fields, enable_dynamic_field=True)

            # creating collection based on the above schema
            collection = Collection(name=data.collection_name, schema=schema)

            # create the collection if needed
            # if not utility.has_collection(data.collection_name):
            # creating index params
            index_params = {
                "index_type": f"{milvus_index_type}",
                "metric_type": f"{milvus_metric_type}",
            }

            # adding the index to the collection
            collection.create_index(field_name="vector", index_params=index_params)

            # telling milvus to load the collection
            collection.load()

            # checksum = self.generate_checksum(content)

            # this makes a gRPC call to the embedding service
            # it's an architectural decision, so we have only one container handling
            # the embedding models. The container will require a significant amount of ram, cpu
            # and eventually it will run on gpu
            embedding = self.embedd(content, data.model_name)

            collection.insert([
                {
                    "document_id": document_id,
                    "parent_id": parent_id,
                    "content": f'{{"content":"{content}"}}',
                    "vector": embedding
                }
            ])

            collection.flush()
            success = True
            self.logger.debug(f"element successfully inserted in collection {data.collection_name}")
        except Exception as e:
            self.logger.error(f"❌ {e}")
            success = False
        finally:
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"⏰ store into vector db, including embedding generation, total elapsed time: {elapsed_time:.2f} seconds")
            return success

    def embedd(self, content_to_embedd: str, model: str) -> List[float]:
        # TODO: get padams fom env
        start_time = time.time()  # Record the start time
        with grpc.insecure_channel(f"{embedder_grpc_host}:{embedder_grpc_port}") as channel:
            stub = EmbedServiceStub(channel)

            self.logger.debug("Calling gRPC Service GetEmbed - Unary")

            # embed_request = EmbedRequest(content=content_to_embedd,
            #   model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
            embed_request = EmbedRequest(content=content_to_embedd, model=model)
            embed_response = stub.GetEmbeding(embed_request)

            self.logger.debug("GetEmbedding gRPC call received correctly")
            end_time = time.time()  # Record the end time
            elapsed_time = end_time - start_time
            self.logger.info(f"⏰ total elapsed time to create embedding: {elapsed_time:.2f} seconds")

            return list(embed_response.vector)

    def _connect(self):
        try:
            # TODO: get params from env
            connections.connect(
                alias=milvus_alias,
                host=milvus_host,
                # host='milvus-standalone'
                port=milvus_port,
                user=milvus_user,
                password=milvus_pass
            )

            self.logger.info(utility.connections.has_connection("defaul"))
            self.logger.info("Connected to Milvus")
        except Exception as e:
            self.logger.error(f"❌ Failed to connect to Milvus {e}")
            self.connection = None

    def ensure_connection(self):
        if not is_connected():
            self.logger.info("Reconnecting to Milvus")
            self._connect()
