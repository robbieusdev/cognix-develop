# Standards

**Containers

| image names         | deployment names    | service names       | pod names           |
| ------------------- | ------------------- | ------------------- | ------------------- |
| cognix-embedder     | cognix-embedder     | cognix-embedder     | cognix-embedder     |
| cognix-semantic     | cognix-semantic     | cognix-semantic     | cognix-semantic     |
| cognix-web          | cognix-web          | cognix-web          | cognix-web          |
| cognix-api          | cognix-api          | cognix-api          | cognix-api          |
| cognix-orchestrator | cognix-orchestrator | cognix-orchestrator | cognix-orchestrator |
| cognix-connector    | cognix-connector    | cognix-connector    | cognix-connector    |
| cognix-migration    | cognix-migration    | cognix-migration    | cognix-migration    |

**Configmaps (src/config)

must end with -cli for client related parameters
must end with -srv for server related parameters
env-configmap for the global parameters (.env in docker)

| name             |
| ---------------- |
| api-srv          |
| cockroach-cli    |
| connector-srv    |
| embedder-cli     |
| embedder-srv     |
| env-configmap    |
| milvus-cli       |
| minio-cli        |
| nats-cli         |
| oauth-cli        |
| orchestrator-srv |
| semantic-srv     |
| web-srv          |

**Github Actions workflows 


| workflow file name            | **workflow name**         |
| ----------------------------- | ------------------------- |
| build-cognix-embedder.yml     | Build Cognix Embedder     |
| build-cognix-semantic.yml     | Build Cognix Semantic     |
| build-cognix-web.yml          | Build Cognix Web          |
| build-cognix-api.yml          | Build Cognix API          |
| build-cognix-orchestrator.yml | Build Cognix Orchestrator |
| build-cognix-connector.yml    | Build Cognix Connector    |
| build-cognix-migration.yml    | Build Cognix Migration    |


**Versioning, based on semantic versioning

| image names         | tag    |
| ------------------- | ------ |
| cognix-embedder     | v0.0.1 |
| cognix-semantic     | v0.0.1 |
| cognix-web          | v0.0.1 |
| cognix-api          | v0.0.1 |
| cognix-orchestrator | v0.0.1 |
| cognix-connector    | v0.0.1 |
| cognix-migration    | v0.0.1 |

Volumes

| service names       | volume name           | pvc name               | directory | shared | size |
| ------------------- | --------------------- | ---------------------- | --------- | ------ | ---- |
| cognix-embedder     | models-volume         | embedder-models-volume | /models   | yes    | ??   |
| cognix-semantic     | use ephemeral storage | N/A                    | /temp     | no     | N/A  |
| cognix-web          | N/A                   | N/A                    | N/A       | N/A    | N/A  |
| cognix-api          | N/A                   | N/A                    | N/A       | N/A    | N/A  |
| cognix-orchestrator | N/A                   | N/A                    | N/A       | N/A    | N/A  |
| cognix-connector    | N/A                   | N/A                    | N/A       | N/A    | N/A  |
| cognix-migration    | migration-volume      | migration-volume       | /versions | no     | 20GB |

Folders

| service names       | folder                   |
| ------------------- | ------------------------ |
| cognix-embedder     | src/backend/embedder     |
| cognix-semantic     | src/backend/semantic     |
| cognix-web          | src/web                  |
| cognix-api          | src/backend/api          |
| cognix-orchestrator | src/backend/orchestrator |
| cognix-connector    | src/backend/connector    |
| cognix-migration    | src/backend/migration    |


Ports

| service names       | Port  |
| ------------------- | ----- |
| cognix-embedder     | 50051 |
| cognix-semantic     | 8080  |
| cognix-web          | 8080  |
| cognix-api          | 8080  |
| cognix-orchestrator | none  |
| cognix-connector    | none  |
| cognix-migration    | none  |