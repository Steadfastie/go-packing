CREATE TABLE IF NOT EXISTS pack_configurations (
    id SMALLINT PRIMARY KEY,
    version BIGINT NOT NULL,
    pack_sizes INTEGER[] NOT NULL,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
