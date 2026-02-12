SELECT 'CREATE DATABASE packing'
WHERE NOT EXISTS (SELECT FROM pg_database WHERE datname = 'packing')
\gexec

\connect packing

CREATE TABLE IF NOT EXISTS pack_configs (
    id SMALLINT PRIMARY KEY CHECK (id = 1),
    pack_sizes INTEGER[] NOT NULL,
    version BIGINT NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
