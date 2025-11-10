-- +goose Up
-- +goose StatementBegin
CREATE TABLE stories (
    id SERIAL PRIMARY KEY,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    file_url TEXT NOT NULL,             -- ссылка на S3
    file_type TEXT NOT NULL,            -- image/jpeg, video/mp4 и т.д.
    views_count INT DEFAULT 0,
    expires_at TIMESTAMP NOT NULL,       -- время истечения
    created_at TIMESTAMP NOT NULL DEFAULT NOW(),
    deleted_at TIMESTAMP NULL
);

CREATE TABLE story_views (
    id SERIAL PRIMARY KEY,
    story_id INT NOT NULL REFERENCES stories(id) ON DELETE CASCADE,
    user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    viewed_at TIMESTAMP NOT NULL DEFAULT NOW(),
    UNIQUE(story_id, user_id)
);

CREATE INDEX idx_stories_user_id ON stories (user_id);
CREATE INDEX idx_stories_expires_at ON stories (expires_at);
CREATE INDEX idx_stories_created_at ON stories (created_at);
CREATE INDEX idx_stories_deleted_at ON stories (deleted_at);
CREATE INDEX idx_story_views_story_id ON story_views (story_id);
CREATE INDEX idx_story_views_user_id ON story_views (user_id);
CREATE INDEX idx_story_views_viewed_at ON story_views (viewed_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS story_views CASCADE;
DROP TABLE IF EXISTS stories CASCADE;
-- +goose StatementEnd

