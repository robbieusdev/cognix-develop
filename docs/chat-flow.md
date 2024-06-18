The user can interact with CogniX, start a new chat in two different ways (in the future will be expanded)

1. Chat with your documents (requires_llm = true)
2. Search only (requires requires_llm=false)


The user asks a question (query): "How do I do get attention in Collaboard" 

Server receives the chat

following are business logic methods:

method name strem_response
argument query (string)
returns the stream

strem_response calls vector_search
    method name vector_search
    argument query (string)
    returns a list of documents

    It performs a query against Milvus, with userid and tennantid
    The query result will contian all the matching documents inside Milvus.

if equires_llm == true
    It retrives the information regarding Persona, prompt and llm

    Top X documents retrieved from Milvus query will be sent to the LLM (LLM Endpoint + eventually api_key) embedded into the prompt

    LLM will stream the answer which will be streamed back to the client
if equires_llm == false
    will send back only the list of documents

Chat and response is saved into chat history