CREATE TABLE users (
                       id SERIAL PRIMARY KEY,
                       name TEXT NOT NULL,
                       age INT NOT NULL,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_users_name ON users (name);
CREATE INDEX idx_users_created_at ON users (created_at);
CREATE INDEX idx_users_updated_at ON users (updated_at);
CREATE INDEX idx_users_deleted_at ON users (deleted_at);

CREATE OR REPLACE FUNCTION set_updated_at_timestamp()
    RETURNS TRIGGER AS $BODY$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$BODY$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_updated_at
    BEFORE UPDATE ON users
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();

ALTER TABLE users
    ADD COLUMN username TEXT UNIQUE NOT NULL DEFAULT '',
    ADD COLUMN password_hash TEXT NOT NULL DEFAULT '',
    ADD COLUMN email TEXT UNIQUE,
    ADD COLUMN is_active BOOLEAN NOT NULL DEFAULT TRUE;
