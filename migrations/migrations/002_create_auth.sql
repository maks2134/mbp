CREATE TABLE auth (
                      id SERIAL PRIMARY KEY,
                      user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                      email TEXT UNIQUE NOT NULL,
                      password_hash TEXT NOT NULL,
                      last_login_at TIMESTAMP NULL,
                      failed_attempts INT NOT NULL DEFAULT 0,
                      locked_until TIMESTAMP NULL,
                      created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                      updated_at TIMESTAMP NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_auth_user_id ON auth (user_id);
CREATE INDEX idx_auth_email ON auth (email);
CREATE INDEX idx_auth_created_at ON auth (created_at);
CREATE INDEX idx_auth_updated_at ON auth (updated_at);

CREATE OR REPLACE FUNCTION set_auth_updated_at()
    RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at := NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trigger_set_auth_updated_at
    BEFORE UPDATE ON auth
    FOR EACH ROW
EXECUTE FUNCTION set_auth_updated_at();

CREATE TABLE auth_tokens (
                             id SERIAL PRIMARY KEY,
                             user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                             token TEXT NOT NULL,
                             expires_at TIMESTAMP NOT NULL,
                             created_at TIMESTAMP DEFAULT NOW()
);
CREATE INDEX idx_auth_tokens_user_id ON auth_tokens (user_id);
CREATE INDEX idx_auth_tokens_expires_at ON auth_tokens (expires_at);
