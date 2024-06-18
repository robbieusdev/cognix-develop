from lib.gen_types.embed_service_pb2_grpc import EmbedServiceStub
from lib.gen_types.embed_service_pb2 import EmbedRequest
import grpc

def run():
    #with grpc.insecure_channel('127.0.0.1:50051') as channel:
    with grpc.insecure_channel('localhost:50051') as channel:
        stub = EmbedServiceStub(channel)
        print("Calling gRPC Service GetEmbed - Unary")

        for i in range(10000):
            content_to_embedd = "docker-compose up --build –  Jinna Baalu Apr 22, 2019 at 14:5 Just to clarify upper comment: docker-compose up --build rebuild all containers. Use docker-compose up --build <service_name> as stated in @denov comment. –  J"
            embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
            embed_response = stub.GetEmbeding(embed_request)

        # embed_request = EmbedRequest(content=content_to_embedd, model="microsoft/mpnet-base")
        # embed_response = stub.GetEmbeding(embed_request)
        #
        # embed_request = EmbedRequest(content=content_to_embedd, model="distilbert/distilroberta-base")
        # embed_response = stub.GetEmbeding(embed_request)

        # embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/paraphrase-multilingual-mpnet-base-v2")
        # embed_response = stub.GetEmbeding(embed_request)

        # embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/natural-questions")
        # embed_response = stub.GetEmbeding(embed_request)

        # embed_request = EmbedRequest(content=content_to_embedd, model="sentence-transformers/wikianswers-duplicates")
        # embed_response = stub.GetEmbeding(embed_request)
        
        
        print("GetEmbed Response Received:")
        # print(embed_response.vector)

if __name__ == "__main__":
    run()
