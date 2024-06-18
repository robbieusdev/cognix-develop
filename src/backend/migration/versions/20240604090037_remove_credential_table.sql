-- +goose Up
-- +goose StatementBegin
alter table connectors drop column if exists credential_id;
alter table connectors drop column if exists disabled;
alter table connectors rename column last_successful_index_date to last_successful_analyzed;
alter table connectors rename column total_docs_indexed to total_docs_analyzed;
alter table connectors rename column last_attempt_status to status;
drop table if exists credentials;
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
create table if not exists credentials
(
    id              bigint default unique_rowid() not null        primary key,
    credential_json jsonb  default '{}'::JSONB    not null,
    user_id         uuid                          not null        references public.users,
    tenant_id       uuid                          not null        references public.tenants,
    source          varchar(50)                   not null,
    creation_date   timestamp                     not null,
    last_update     timestamp,
    deleted_date    timestamp,
    shared          boolean                       not null
);
alter table connectors add column if not exists credential_id bigint references credentials(id);
alter table connectors rename column status to last_attempt_status;
alter table connectors add column if not exists disabled bool ;
alter table connectors rename column last_successful_analyzed to last_successful_index_date;
alter table connectors rename column total_docs_analyzed to total_docs_indexed;
-- +goose StatementEnd
