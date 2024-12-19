CREATE TABLE size (
id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
name STRING NOT NULL,
created_at TIMESTAMP NOT NULL DEFAULT NOW(),
updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
deleted_at TIMESTAMPTZ DEFAULT NULL
);