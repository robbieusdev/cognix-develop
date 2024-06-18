
### Milvus schema 
```
| name        | type      | description                            |
| ---         | ---       | ---                                    |          
| id          | int64     |  primary key |  
| document_id | int64     |   id of document in cockroach database |    
| content     | json      | text content expected format {"content":""}|
| vector      | []float32 |  vector array |

```
?? ***Need to define dimension fo vector array*** 
for testing with open-ai i used 1536 

??? Milvus schemas for user and tenant can be created 
- when user signin first time. (golang)
- when first document will be processed in embedder service (python)

```shell
cockroach table documents
(
    id           bigint  default unique_rowid() not null  primary key,
    document_id  varchar not null,
    connector_id bigint not null references public.connectors,
    link         varchar,
    signature    text,
    created_date timestamp default now()  not null,
    updated_date timestamp,
    deleted_date timestamp,
    status       varchar(100) default 'new'::STRING  not null
)

supported document statuses 

    StatusPending    = "pending"	
	StatusChunking   = "chunking"
	StatusEmbedding  = "embedding"
	StatusComplete   = "complete"

```

### Document Statuses in  connector flow 

#### golang connector service 
- Set status ***pending*** for all documents stored in database when connector is started.
- Read document one by one. 
  - if document was not modified (check hash) set status ***complete*** and do not send message for chunking 
  - if document new or was changed set status ***chunking*** and send message for chunking 
- WWhen all documents were read from source delete ( or mark for delete need to discuss ) all documents with status ***pending***

#### python services 
- chunking service set status ***embedding*** and send message to embedding service 
- embedding service set status ***complete*** and store content in milvus. 
