create table tenants
(
    id            uuid                      not null        primary key,
    name          varchar(255),
    configuration jsonb default '{}'::JSONB not null
);

create table users
(
    id          uuid                          not null        primary key,
    tenant_id   uuid                          not null        references tenants,
    user_name   varchar(255)                  not null        unique,
    first_name  varchar(255),
    last_name   varchar(255),
    external_id text,
    roles       text[] default '{}'::STRING[] not null
);

create table llms
(
    id            bigint default unique_rowid() not null        primary key,
    tenant_id     uuid                          not null        references tenants,
    name          varchar(255)                  not null,
    model_id      varchar(255)                  not null,
    url           varchar                       not null,
    api_key       varchar,
    endpoint      varchar,
    creation_date timestamp                     not null,
    last_update   timestamp,
    deleted_date  timestamp
);

create table embedding_models
(
    id            bigint  default unique_rowid() not null        primary key,
    tenant_id     uuid                           not null,
    model_id      varchar                        not null,
    model_name    varchar                        not null,
    model_dim     bigint                         not null,
    url           varchar,
    is_active     boolean default false          not null,
    creation_date timestamp                      not null,
    last_update   timestamp,
    deleted_date  timestamp
);

create table personas
(
    id               bigint default unique_rowid() not null        primary key,
    name             varchar                       not null,
    llm_id           bigint        references llm,
    default_persona  boolean                       not null,
    description      varchar                       not null,
    tenant_id        uuid                          not null        references tenants,
    is_visible       boolean                       not null,
    display_priority bigint,
    starter_messages jsonb  default '{}'::JSONB    not null,
    creation_date    timestamp                     not null,
    last_update      timestamp,
    deleted_date     timestamp
);

create table prompts
(
    id            bigint default unique_rowid() not null        primary key,
    persona_id    bigint                        not null        references personas,
    user_id       uuid                          not null        references users,
    name          varchar                       not null,
    description   varchar                       not null,
    system_prompt text                          not null,
    task_prompt   text                          not null,
    creation_date timestamp                     not null,
    last_update   timestamp,
    deleted_date  timestamp
);

create table connectors
(
    id                        bigint default unique_rowid() not null        primary key,
    name                      varchar                       not null,
    type                      varchar(50)                   not null,
    connector_specific_config jsonb                         not null,
    refresh_freq              bigint,
    user_id                   uuid                          not null        references public.users,
    tenant_id                 uuid        references public.tenants,
    last_successful_analyzed  timestamp,
    status                    varchar,
    total_docs_analyzed       bigint                        not null,
    creation_date             timestamp                     not null,
    last_update               timestamp,
    deleted_date              timestamp
);

create table chat_sessions
(
    id            bigint default unique_rowid() not null        primary key,
    user_id       uuid                          not null        references users,
    description   text                          not null,
    creation_date timestamp                     not null,
    deleted_date  timestamp,
    persona_id    bigint                        not null        references personas,
    one_shot      boolean                       not null
);

create table chat_messages
(
    id                   bigint default unique_rowid() not null        primary key,
    chat_session_id      bigint                        not null        references chat_sessions,
    message              text                          not null,
    message_type         varchar(9)                    not null,
    time_sent            timestamp                     not null,
    token_count          bigint                        not null,
    parent_message       bigint,
    latest_child_message bigint,
    rephrased_query      text,
    error                text
);

create table chat_message_feedbacks
(
    id              bigint  default unique_rowid() not null        primary key,
    chat_message_id bigint                         not null        references chat_messages,
    user_id         uuid                           not null        references users,
    up_votes        boolean                        not null,
    feedback        varchar default ''::STRING     not null
);

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

create table chat_message_document_pairs
(
    id              bigint default unique_rowid() not null primary key,
    chat_message_id bigint not null references public.chat_messages on delete cascade,
    document_id     bigint not null references public.documents on delete cascade
);
