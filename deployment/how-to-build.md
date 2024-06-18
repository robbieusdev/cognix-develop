Create a .env file in the same directory as your docker-compose.yml:

DATA_PATH=/users/gp/Developer/cognix-services/data
CONFIG_PATH=/absolute/path/to/config
BACKEND_PATH=/absolute/path/to/backend
MIGRATION_PATH=/absolute/path/to/migration

sudo docker-compose -f deployment/docker-compose-test.yaml up/down

sudo docker-compose -f deployment/docker-compose-test.yaml --build

from the docker folder
docker build -f Dockerfile-embedder -t embedder:latest ..

CONNECT INSIDE THE CONTAINER
docker exec -it <container_name_or_id> /bin/bash
bin/pulsar-admin namespaces set-schema-compatibility-strategy --compatibility AUTO_CONSUME <my-tenant/my-namespace>

inside the container "df" command shows all volumes mounted



Change to Project Root:
and then rund this command to
docker build -f docker/Dockerfile-embedder -t embedder:latest .
sudo docker build -f docker/Dockerfile-embedder -t embedder:latest --build-arg COGNIX_PATH=. .

docker run embedder

docker compose -f docker-compose-embedder.yml build

docker compose -f docker-compose-embedder.yml up -d
docker compose -f docker/docker-compose-services-ai.yaml up -d

docker compose -f docker-compose-embedder.yml logs -f embedder