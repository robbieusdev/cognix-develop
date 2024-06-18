-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS test (
   id uuid PRIMARY KEY,
  name text

);
DROP TABLE  test ;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
SELECT 'down SQL query';

-- +goose StatementEnd
