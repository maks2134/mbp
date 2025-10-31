-- +goose Up
-- +goose StatementBegin
CREATE TABLE comments (
    id SERIAL PRIMARY KEY,
    post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    text TEXT NOT NULL,
    "like" INT DEFAULT 0,
    blocked BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_comments_post_id ON comments (post_id);
CREATE INDEX idx_comments_user_id ON comments (user_id);
CREATE INDEX idx_comments_created_at ON comments (created_at);
CREATE INDEX idx_comments_updated_at ON comments (updated_at);
CREATE INDEX idx_comments_deleted_at ON comments (deleted_at);

CREATE TRIGGER trigger_set_updated_at_comments
    BEFORE UPDATE ON comments
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TRIGGER IF EXISTS trigger_set_updated_at_comments ON comments;
DROP TABLE IF EXISTS comments CASCADE;
-- +goose StatementEnd
