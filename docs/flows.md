The orchestrator should play an intelligent role in managing the flow of data between different components, rather than acting as a simple scheduler. Here are my thoughts on how to structure the orchestrator to meet these goals:


[Orchestrator](#orchestrator)
 - [General Behavior](#general-behavior)
 - [Multiple Orchestrator Instances](#multiple-orchestrator-instances)
 - [Connector Status](#connector-status)
     - [Statuses](#statuses)
     - [Important Considerations](#important-considerations)
- [Rules for Re-scanning Connectors](#rules-for-re-scanning-connectors)
  - [URL](#url)
- [File](#file)
- [OneDrive, Google Drive, and Other Cloud Drives](#onedrive-google-drive-and-other-cloud-drives)
- [MS Teams and Slack](#ms-teams-and-slack)
- [Refresh Frequency](#refresh-frequency)
[Connectors](#connectors)
 - [General Behavior](#general-behavior-1)
 - [URL and File](#url-and-file)
 - [OneDrive, Google Drive, and Other Cloud Drives](#onedrive-google-drive-and-other-cloud-drives-1)
[Chunker](#chunker)
[Document Table Schema](#document-table-schema)
  - [document](#document)
  - [connector](#connector)
  - [chunking_data](#chunking_data)



The orchestrator should play an intelligent role in managing the flow of data between different components, rather than acting as a simple scheduler. Here are my thoughts on how to structure the orchestrator to meet these goals:

1. **Intelligent Orchestrator**: The orchestrator should have the capability to decide whether a scan request should go directly to the chunker or through the connector, based on the type of source and its current status.

2. **Direct Handling for Certain Sources**: For sources like files and URLs, it makes sense for the orchestrator to send requests directly to the chunker, bypassing the connector to reduce unnecessary traffic and workload. This can significantly optimize the system's efficiency, especially at scale.

3. **Connector's Role for Complex Sources**: Connectors should handle more complex sources, such as OneDrive, Google Drive, MS Teams, and Slack, where there is a need to scan directories, verify credentials, or handle nested structures. The connector can then process these sources and decide what needs to be sent to the chunker.

Here's a revised structure based on these principles:

# Orchestrator
The orchestrator is responsible for monitoring the connector table in the relational database and intelligently deciding when it's time for a new analysis of the single connector. A single connector row represents a knowledge source the user decided to have analyzed by CogniX and stored in the Vector database.

## General Behavior
- Monitor the connector table for changes.
- Decide when a new scan or analysis is needed.
- Send analysis requests intelligently to the Connector or Chunker, depending on the specific use case.
- Ensure no other processes are running before starting a new scan for that particular source.

## Multiple Orchestrator Instances
TBD: Define how to handle multiple orchestrator instances.

## Connector Status
The status of the connector is currently determined by a field in the connector table. In the future, the status will be retrieved from NATS.

### Statuses:
- **Active**: Connector just created.
- **Pending Scan**: Orchestrator sent a message to Connector or Chunker to scan the source.
- **Working**: Set by Connector and Orchestrator, indicating that the Connector or Chunker is working on it.
- **Scan Completed Successfully**: Set by Chunker.
- **Scan Completed with Errors**: Set by Chunker.
- **Disabled**: No further scans will be started for this connector.
- **Unable to Process**: Set by Chunker, indicating a problem with the source. Only one email will be sent to the user reporting the issue.

### Important Considerations:
- Determine actions for rows in the Pending or Working status for long time, which means the process has crashed without being able to update the status.
- Investigate if NATS can notify when a new item is added to the dead letter queue. If so, the orchestrator or another service should subscribe to this message and set the status to Unable to Process.


###  Rules for Re-scanning Connectors

- The connector is in one of the following statuses:
  - "Active"
  - "Scan Completed Successfully"
  - "Scan Completed with Errors"
- The connector is **not** in the following statuses:
  - "Pending Scan"
  - "Working"
  - "Disabled"
  - "Unable to Process"
- The resulting date after adding the refresh frequency (refresh_freq) to the last update date is less than the current UTC time (now_utc()).
- One-Time Scan Connectors
    - If the connector type is "YT" (YouTube) or "File", the connector should not be triggered again once it reaches an ending status (either "Scan Completed Successfully" or "Scan Completed with Errors").
- Statuses that Prevent Scanning
    - The orchestrator should never trigger a scan if the connector is in "Disabled" or "Unable to Process" status.


## URL
The user can connect a URL as a knowledge source, providing:
- URL (mandatory)
- Sitemap URL (optional)
- Option to scan all links on the page (optional)
- Option to search for a sitemap if not provided

The Orchestrator will forward the request directly to Chunker for this file type. No file is stored in MinIO.
If the user deletes the source, the API will soft delete the row from the relational database and hard delete related entities from the vector database.

## File
The user can upload a file as a knowledge source, providing:
- The file to be analyzed (mandatory)

Files are uploaded to MinIO and scanned only once. Once the status is set to scan completed (with or without errors), no further analysis requests are issued. The Orchestrator will send a request directly to Chunker as soon as the file is uploaded correctly.

## OneDrive, Google Drive, and Other Cloud Drives
The user can connect a cloud drive as a knowledge source, providing:
- Path to be analyzed (mandatory)
- Option to scan the path only or all subfolders
- Necessary credentials to access the path

The Orchestrator will forward the request to Connector for this file type.

## MS Teams and Slack
The user can upload a file as a knowledge source, providing:
- The file to be analyzed (mandatory)

TBD: Define the full flow and how to handle the documents table.

## Refresh Frequency
- **URL**: One week
- **OneDrive, Google Drive, and other cloud drives**: One week
- **MS Teams and Slack**: Daily, only new messages


## Sequence diagram
- Orchestrator monitors the connector table.
- Depending on the connector status, it sends a scan request to the connector or chunker.
- Connector processes the request for complex sources and interacts with chunker as needed.
- Chunker processes the data and updates the status.

![plot](https://github.com/gen-mind/cognix/blob/feature/chunking-3/docs/media/orchestrator_sequence_diagram.png)

## Componenet diagram
![plot](https://github.com/gen-mind/cognix/blob/feature/chunking-3/docs/media/orchestrator_componenet_diagram.png)



# Connectors
The connector does not perform any actions on the vector database. It scans and downloads files from sources, storing them in MinIO, and determines if files need to be reloaded and re-chunked.

## General Behavior
- Verify if files need to be reloaded and re-chunked.
- Send messages to Chunker only for modified or new data.

## URL and File
The Connector does not handle these directly; requests are sent directly to Chunker by the Orchestrator.

## OneDrive - Google Drive and other cloud drives
When the Connector receives a new message from the orchestrator it will
- Creates a GUID “chunking_session” that will be sent to each Chunking message sent by this operation. This way the Chunker will be able to understand when to set this process as completed
- Set the connector status from Pending scan to  Working by Connector. 
- It will scan the drive (given the rules from the orchestrator, all sub-folder or not) and get a list of path/file
- for each path/file item will check in the documents table if the item shall be sent to Chunker, depending on the hash comparison between database and file actually scanned.
    - Update chunking_session with the new chunking_session
    - if the item needs to be scanned (because is new or updated) sets the "analyzed" field to false
    - if the item does not needs to be scanned sets the "analyzed" field to true
    - (it is important to update the database for all the itmes before sending messages to NATS to avoid concurrency)
- delete (physically) all the documents in the database that are not anymore present in the original source  
- Iterage again the list (after DB is updated for all the rows)
    - if the item needs to be scanned (because is new or updated) send a message to chunker

### Sequence diagram

1. **Orchestrator to Connector:** The Orchestrator sends a new message to the Connector to start the process.
2. **Connector:** The Connector creates a GUID named "chunking_session" and sets the status to "Working."
3. **Connector to Drive:** The Connector scans the drive based on the given rules and retrieves a list of files and paths.
4. **Loop for each file/path:**
   - **Connector to Database:** The Connector checks in the database if the item should be sent to the Chunker by comparing hashes.
   - **Database to Connector:** The database returns the hash comparison result.
   - **Alt (Condition):** 
     - If the item needs to be scanned, the Connector updates the chunking_session and sets the "analyzed" field to false.
     - If the item does not need to be scanned, the Connector sets the "analyzed" field to true.
5. **Connector to Database:** The Connector deletes documents from the database that are no longer present in the original source.
6. **Loop for each file/path:**
   - **Alt (Condition):** 
     - If the item needs to be scanned, the Connector sends a message to the Chunker.
![plot](https://github.com/gen-mind/cognix/blob/feature/chunking-3/docs/media/connecor_cloud_drive_sequence_diagram.png)


# Chunker
Chunker processes data from various sources, creating embeddings and storing them in the vector database.

## URL Chunker
Becaus ein golang is pretty time consuming creating a crawler able to analyze the content of a 
This specific chunker has two 
guid “chunking_session”


## OneDrive, Google Drive, and Other Cloud Drives

# Document Table Schema
All the tables in the database shall have the following fields:
```sql
  "creation_date" timestamp NOT NULL DEFAULT (now()),
  "last_update" timestamp
```
TBD if to add also ubdated_by (the user or the service that performed the last update)


## document
```sql
CREATE TABLE "document" (
  "id" varchar PRIMARY KEY NOT NULL,
  "parent_id" integer, -- Allows nulls, used for URLs
  "connector_id" integer NOT NULL,
  "link" varchar,
  "signature" text,
  "chunking_session" guid, -- Allows nulls
  "analyzed" bit -- default false, true when chunker created the embeddings in the vector db
  "creation_date" timestamp NOT NULL DEFAULT (now()), --datetime utc IMPORTANT now() will not get the utc date!!!!
  "last_update" timestamp --datetime utc
);
```
The parent_id is used for URL. 
The orchestrator does not know how may related URL will be found in the starting URL, to store all the data we need to use the parent id
The first document, is the URL the user provided.
URL Chunker will then create as many documents, with the parent_id set to the id of the document the user stored in the connector

## connector
```sql
CREATE TABLE "connectors" (
  "id" SERIAL PRIMARY KEY,
  "credential_id" integer NOT NULL, -- remove this and related tabe
  "name" varchar NOT NULL,
  "type" varchar(50) NOT NULL, -- PDF, URL etc
  "status" varchar,
  "connector_specific_config" jsonb NOT NULL,
  "refresh_freq" integer,
  "user_id" uuid NOT NULL,
  "tenant_id" uuid NOT NULL,
  "last_successful_analyzed" timestamp, --datetime utc
   "total_docs_analyzed" integer, 
  "creation_date" timestamp NOT NULL DEFAULT (now()), --datetime utc
  "last_update" timestamp --datetime utc
);
```
status in the connector will have the following values:
```sql
	ConnectorStatusReadyToProcessed = "READY_TO_BE_PROCESSED"
	ConnectorStatusPending          = "PENDING"
	ConnectorStatusWorking          = "PROCESSING"
	ConnectorStatusSuccess          = "COMPLETED_SUCCESSFULLY"
	ConnectorStatusError            = "COMPLETED_WITH_ERRORS"
	ConnectorStatusDisabled         = "DISABLED"
	ConnectorStatusUnableProcess    = "UNABLE_TO_PROCESS"
```

Removed the shared field. If a connector has the tennat_id it means it is shared, if it has the user_id it means is from the user

## chunking_data
```proto
syntax = "proto3";

package com.embedd;
option go_package = "backend/core/proto;proto";


enum FileType {
  UNKNOWN = 0;
  URL = 1;
  PDF = 2;
  RTF = 3;
  DOC = 4;
  XLS = 5;
  PPT = 6;
  TXT = 7;
  MD  = 8;
  YT = 9;
  // add all supported file that in another document
  // check what with Google docs
};

message ChunkingData {
  // This is the url where the file is located.
  // Based on the chunking type it will be a WEB URL (HTML type)
  // Will be an S3/MINIO link with a proper authentication in case of a file
  string url = 1;
  string site_map = 2;
  bool search_for_sitemap = 3;
  int64  document_id = 4;
  FileType file_type = 5;
  string collection_name = 6;
  string model_name = 7;
  int32 model_dimension = 8;
  uuid chunking_session = 9;
  bool is_internal = 10; // used by chunker if it needs to send message to itself
}
```

- URL (mandatory)
- Sitemap URL (optional)
- Option to scan all links on the page (optional)
- Option to search for a sitemap if not provided
