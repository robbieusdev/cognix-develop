-- +goose Up
-- +goose StatementBegin
CREATE TABLE IF NOT EXISTS tenants (
    id uuid PRIMARY KEY,
    name varchar(255),
    configuration jsonb not null default '{}'::jsonb
);
CREATE TABLE IF NOT EXISTS users (
    id uuid PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    user_name varchar(255) UNIQUE NOT NULL,
    first_name varchar(255),
    last_name varchar(255),
    external_id text,
    roles text[] NOT NULL DEFAULT '{}'::text[]
);
CREATE TABLE IF NOT EXISTS llm (
    id SERIAL PRIMARY KEY,
    tenant_id uuid NOT NULL REFERENCES tenants(id),
    name varchar(255) NOT NULL,
    model_id varchar(255) NOT NULL,
    url varchar NOT NULL,
    api_key varchar,
    endpoint varchar,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS embedding_models (
    id SERIAL PRIMARY KEY,
    tenant_id uuid NOT NULL,
    model_id varchar NOT NULL,
    model_name varchar NOT NULL,
    model_dim integer NOT NULL,
    normalize boolean NOT NULL,
    query_prefix varchar NOT NULL,
    passage_prefix varchar NOT NULL,
    index_name varchar NOT NULL,
    "url" varchar,
    is_active boolean NOT NULL DEFAULT false,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS personas (
    id SERIAL PRIMARY KEY,
    name varchar NOT NULL,
    llm_id bigint references llm(id),
    default_persona boolean NOT NULL,
    description varchar NOT NULL,
    tenant_id uuid NOT NULL references tenants(id),
    is_visible boolean NOT NULL,
    display_priority integer,
    starter_messages jsonb NOT NULL DEFAULT '{}'::jsonb,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS prompts (
    id SERIAL PRIMARY KEY,
    persona_id bigint NOT NULL REFERENCES personas(id),
    user_id uuid NOT NULL references users(id),
    name varchar NOT NULL,
    description varchar NOT NULL,
    system_prompt text NOT NULL,
    task_prompt text NOT NULL,
    include_citations boolean NOT NULL,
    datetime_aware boolean NOT NULL,
    default_prompt boolean NOT NULL,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (now()),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS credentials (
    id SERIAL PRIMARY KEY,
    credential_json jsonb NOT NULL DEFAULT '{}'::jsonb,
    user_id uuid NOT NULL references users(id),
    tenant_id uuid NOT NULL references tenants(id),
    source varchar(50) NOT NULL,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (now()),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE,
    shared boolean NOT NULL
);


CREATE TABLE IF NOT EXISTS connectors (
    id SERIAL PRIMARY KEY,
    credential_id bigint NULL references credentials(id),
    name varchar NOT NULL,
    source varchar(50) NOT NULL,
    input_type varchar(10),
    connector_specific_config jsonb NOT NULL,
    refresh_freq integer,
    user_id uuid NOT NULL references users(id),
    tenant_id uuid NOT NULL references tenants(id),
    shared boolean NOT NULL,
    disabled boolean NOT NULL,
    last_successful_index_time timestamp WITHOUT TIME ZONE,
    last_attempt_status varchar,
    total_docs_indexed integer NOT NULL,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (now()),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS chat_sessions (
    id SERIAL PRIMARY KEY,
    user_id uuid NOT NULL references users(id),
    description text NOT NULL,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (now()),
    deleted_date timestamp WITHOUT TIME ZONE,
    persona_id integer NOT NULL references personas(id),
    one_shot boolean NOT NULL
);

CREATE TABLE IF NOT EXISTS chat_messages (
    id SERIAL PRIMARY KEY,
    chat_session_id bigint NOT NULL references chat_sessions(id),
    message text NOT NULL,
    message_type varchar(9) NOT NULL,
    time_sent timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (now()),
    token_count integer NOT NULL,
    parent_message integer,
    latest_child_message integer,
    rephrased_query text,
    citations jsonb NOT NULL default '{}'::jsonb,
    error text
);

CREATE TABLE chat_message_feedbacks (
    id              SERIAL NOT NULL PRIMARY KEY,
    chat_message_id BIGINT NOT NULL REFERENCES chat_messages,
    user_id         UUID NOT NULL REFERENCES users,
    up_votes        BOOLEAN NOT NULL,
    feedback        VARCHAR NOT NULL DEFAULT ''
);

CREATE TABLE IF NOT EXISTS documents (
    id serial PRIMARY KEY NOT NULL,
    document_id varchar NOT NULL ,
    connector_id bigint NOT NULL references connectors(id),
    boost integer NOT NULL,
    hidden boolean NOT NULL,
    semantic_id varchar NOT NULL,
    link varchar,
    from_ingestion_api boolean,
    signature text,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT now(),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

create table document_feedbacks(
    id              serial        primary key,
    document_id     bigint  references documents(id) ,
    user_id         uuid references users(id),
    document_rank   integer not null,
    up_votes        boolean not null,
    feedback        varchar not null default ''
);


CREATE TABLE IF NOT EXISTS document_sets (
    id SERIAL PRIMARY KEY,
    user_id uuid references users(id),
    name varchar NOT NULL,
    description varchar NOT NULL,
    is_up_to_date boolean NOT NULL,
    created_date timestamp WITHOUT TIME ZONE NOT NULL DEFAULT (now()),
    updated_date timestamp WITHOUT TIME ZONE,
    deleted_date timestamp WITHOUT TIME ZONE
);

CREATE TABLE IF NOT EXISTS document_set_connector_pairs (
    document_set_id integer NOT NULL references document_sets(id),
    connector_id integer NOT NULL references connectors(id),
    is_current boolean NOT NULL,
    primary key (document_set_id,connector_id)
);

-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS embedding_models;
DROP TABLE IF EXISTS chat_message_feedback;
DROP TABLE IF EXISTS chat_messages;
DROP TABLE IF EXISTS chat_sessions;
DROP TABLE IF EXISTS prompts;
DROP TABLE IF EXISTS personas;
DROP TABLE IF EXISTS llm;
DROP TABLE IF EXISTS document_feedbacks;
DROP TABLE IF EXISTS documents;
DROP TABLE IF EXISTS document_set_connector_pairs ;
DROP TABLE IF EXISTS document_sets;
DROP TABLE IF EXISTS connectors;
DROP TABLE IF EXISTS credentials;
DROP TABLE IF EXISTS users;
DROP TABLE IF EXISTS tenants;

-- +goose StatementEnd
