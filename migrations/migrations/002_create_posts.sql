CREATE TABLE posts (
                       id SERIAL PRIMARY KEY,
                       user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
                       title TEXT NOT NULL,
                       description TEXT NOT NULL,
                       tag TEXT NOT NULL,
                       "like" INT DEFAULT 0,
                       count_viewers INT DEFAULT 0,
                       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       updated_at TIMESTAMP NOT NULL DEFAULT NOW(),
                       deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_posts_user_id ON posts (user_id);
CREATE INDEX idx_posts_title ON posts (title);
CREATE INDEX idx_posts_tag ON posts (tag);
CREATE INDEX idx_posts_created_at ON posts (created_at);
CREATE INDEX idx_posts_updated_at ON posts (updated_at);
CREATE INDEX idx_posts_deleted_at ON posts (deleted_at);

CREATE TRIGGER trigger_set_updated_at_posts
    BEFORE UPDATE ON posts
    FOR EACH ROW
EXECUTE FUNCTION set_updated_at_timestamp();
