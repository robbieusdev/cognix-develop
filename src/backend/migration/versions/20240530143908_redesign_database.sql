-- +goose Up
-- +goose StatementBegin
DROP TABLE IF EXISTS document_feedbacks;
DROP TABLE IF EXISTS document_set_connector_pairs;
DROP TABLE IF EXISTS document_sets;
DROP TABLE IF EXISTS documents;

CREATE TABLE documents (
    id SERIAL PRIMARY KEY NOT NULL,
    parent_id bigint REFERENCES documents(id), -- Allows nulls, used for URLs
    connector_id bigint NOT NULL REFERENCES connectors(id),
    source_id text NOT NULL, -- unique id from source url for web, id for other services
    url text, -- url for web connector, link (minio:bucket:file) for file in minio
    signature text,
    chunking_session uuid, -- Allows nulls
    analyzed bool NOT NULL DEFAULT FALSE,  -- default false, true when semantic created the embeddings in the vector db
    creation_date timestamp WITHOUT TIME ZONE NOT NULL, --datetime utc IMPORTANT now() will not get the utc date!!!!
    last_update timestamp WITHOUT TIME ZONE
);

ALTER TABLE connectors RENAME COLUMN source TO type;
ALTER TABLE connectors RENAME COLUMN created_date TO creation_date;
ALTER TABLE connectors ALTER COLUMN creation_date DROP DEFAULT;
ALTER TABLE connectors RENAME COLUMN updated_date TO last_update;
ALTER TABLE connectors RENAME COLUMN last_successful_index_time TO last_successful_index_date;
ALTER TABLE connectors ALTER COLUMN tenant_id DROP NOT NULL;
ALTER TABLE connectors DROP COLUMN IF EXISTS shared;

ALTER TABLE chat_messages ALTER COLUMN time_sent DROP DEFAULT;

ALTER TABLE chat_sessions RENAME COLUMN created_date TO creation_date;
ALTER TABLE chat_sessions ALTER COLUMN creation_date DROP DEFAULT;

ALTER TABLE credentials RENAME COLUMN created_date TO creation_date;
ALTER TABLE credentials ALTER COLUMN creation_date DROP DEFAULT;
ALTER TABLE credentials RENAME COLUMN updated_date TO last_update;

ALTER TABLE embedding_models RENAME COLUMN created_date TO creation_date;
ALTER TABLE embedding_models ALTER COLUMN creation_date DROP DEFAULT;
ALTER TABLE embedding_models RENAME COLUMN updated_date TO last_update;

ALTER TABLE llm RENAME COLUMN created_date TO creation_date;
ALTER TABLE llm ALTER COLUMN creation_date DROP DEFAULT;
ALTER TABLE llm RENAME COLUMN updated_date TO last_update;

ALTER TABLE personas RENAME COLUMN created_date TO creation_date;
ALTER TABLE personas ALTER COLUMN creation_date DROP DEFAULT;
ALTER TABLE personas RENAME COLUMN updated_date TO last_update;

ALTER TABLE prompts RENAME COLUMN created_date TO creation_date;
ALTER TABLE prompts ALTER COLUMN creation_date DROP DEFAULT;
ALTER TABLE prompts RENAME COLUMN updated_date TO last_update;

ALTER TABLE prompts  DROP COLUMN IF EXISTS include_citations;
ALTER TABLE prompts  DROP COLUMN IF EXISTS default_prompt;
ALTER TABLE prompts  DROP COLUMN IF EXISTS datetime_aware;

ALTER TABLE llm RENAME TO llms;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';
-- +goose StatementEnd
