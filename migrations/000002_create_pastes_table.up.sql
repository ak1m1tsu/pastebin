CREATE TABLE IF NOT EXISTS pastes (
    hash varchar(8) NOT NULL PRIMARY KEY,
    user_id uuid REFERENCES users(id) DEFAULT NULL,
    title varchar(255) NOT NULL DEFAULT 'Untitled',
    format varchar(255) NOT NULL,
    password_hash bytea,
    expires_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP + INTERVAL '2' YEAR,
    created_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at timestamp(0) with time zone NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS pastes_user_id_idx ON pastes (user_id);
CREATE INDEX IF NOT EXISTS pastes_title_idx ON pastes USING GIN (to_tsvector('simple', title));
