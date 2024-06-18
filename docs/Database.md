in folder backend/ migration /version will created new file like a  20240606162757_new_script.sql
in this file

`-- +goose Up
-- +goose StatementBegin
 write queries for upgrade dataabse 
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
write queries for downgrade database 
-- +goose StatementEnd`


go orm https://github.com/go-pg/pg