-- +goose Up
-- +goose StatementBegin
ALTER TABLE documents DROP COLUMN IF EXISTS boost;
ALTER TABLE documents DROP COLUMN IF EXISTS hidden;
ALTER TABLE documents DROP COLUMN IF EXISTS semantic_id;
ALTER TABLE documents DROP COLUMN IF EXISTS from_ingestion_api;
ALTER TABLE connectors DROP COLUMN IF EXISTS input_type;
ALTER TABLE embedding_models DROP COLUMN IF EXISTS normalize;
ALTER TABLE embedding_models DROP COLUMN IF EXISTS query_prefix;
ALTER TABLE embedding_models DROP COLUMN IF EXISTS passage_prefix;
ALTER TABLE embedding_models DROP COLUMN IF EXISTS index_name;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
-- +goose StatementEnd
