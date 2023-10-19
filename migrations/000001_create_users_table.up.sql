CREATE TABLE IF NOT EXISTS users (
    id uuid NOT NULL DEFAULT uuid_generate_v4() PRIMARY KEY,
    created_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    updated_at timestamp(0) with time zone NOT NULL DEFAULT NOW(),
    username text UNIQUE NOT NULL,
    email citext UNIQUE NOT NULL,
    avatar path NOT NULL,
    access_token bytea NOT NULL,
    deleted bool NOT NULL DEFAULT false
);
