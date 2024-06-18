we have tow services:
connector
connector orchestrator (we will change this name)
the orcherstrator cheks the connector table and decides if a connector needs to do the work for one line inside connector table
if yes it sends a NATS message so that only one connector service will take the message
Important check how to configure nats for retries and dead letter queue


max 3 retryes and then it shal go in error status
CONNECTOR SERVCES
shall inherit from a base connector class
it shall do the job f crarling
splitting by a fixed size with some char of the previous iteration
it will call the python service to create the embeddings
it will store the embedding in milvus
it will update the status of the connectr
we neeed to store:
    embedding
    original text
reference to find back the original document which shall be a link exposed to the client
milvus store with file index (i'll explain later)
we create a collection for each organization and a collection for each user
in milvus we create a collection for each tennant and one collection for each user foe the private connectors
collection name is fixed by a rule like tennant_{ID}. user_{ID}
pass opentelemetry context to nats for distributed tracing


### file compatibility

 - Office all
   - doc, docx
   - xls, xlsx
   - rtf
   - ppt, pptx
 - pdf
 - google docs all
 - txt
 - md

#### optional 
 - ODT
 - ODS
 - ODP
 - ODG


### connectors 
teams chats
teams teams
sharepoint
slack
email (office/gmail)
google drive
one drive
dropbox
upload files (supported formats)
