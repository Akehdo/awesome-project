CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TABLE meetings (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    original_filename TEXT NOT NULL CHECK (btrim(original_filename) <> ''),
    object_key TEXT NOT NULL UNIQUE CHECK (btrim(object_key) <> ''),
    content_type TEXT NOT NULL CHECK (btrim(content_type) <> ''),
    size_bytes BIGINT NOT NULL CHECK (size_bytes > 0),
    created_at TIMESTAMPTZ NOT NULL DEFAULT now()
);
