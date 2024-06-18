-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS chat_message_document_pairs(
    id SERIAL PRIMARY KEY,
    chat_message_id bigint NOT NULL REFERENCES chat_messages(id) ON DELETE CASCADE,
    document_id bigint NOT NULL REFERENCES  documents(id) ON DELETE CASCADE
);
ALTER TABLE chat_messages DROP COLUMN IF EXISTS citations;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROp TABLE IF EXISTS chat_message_document_pairs;
ALTER TABLE chat_messages ADD COLUMN IF NOT EXISTS citations jsonb;
-- +goose StatementEnd
