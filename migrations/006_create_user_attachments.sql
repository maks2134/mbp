-- +goose Up
-- +goose StatementBegin
CREATE TABLE user_attachments (
       id SERIAL PRIMARY KEY,
       user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
       file_url TEXT NOT NULL,             -- ссылка на S3
       file_type TEXT NOT NULL,            -- image/jpeg, video/mp4 и т.д.
       file_size BIGINT,                   -- в байтах
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_user_attachments_user_id ON user_attachments (user_id);
CREATE INDEX idx_user_attachments_file_type ON user_attachments (file_type);
CREATE INDEX idx_user_attachments_created_at ON user_attachments (created_at);
CREATE INDEX idx_user_attachments_deleted_at ON user_attachments (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS user_attachments CASCADE;
-- +goose StatementEnd

