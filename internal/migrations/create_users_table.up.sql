CREATE TABLE IF NOT EXISTS users (
    id bigserial PRIMARY KEY,
    fullname text NOT NULL,
    email text NOT NULL,
    created_at timestamp(0) WITH TIME ZONE NOT NULL DEFAULT NOW(),
    version integer NOT NULL DEFAULT 1
);