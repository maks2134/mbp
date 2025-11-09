-- +goose Up
-- +goose StatementBegin
CREATE TABLE post_attachments (
       id SERIAL PRIMARY KEY,
       post_id INT NOT NULL REFERENCES posts(id) ON DELETE CASCADE,
       file_url TEXT NOT NULL,             -- ссылка на S3
       file_type TEXT NOT NULL,            -- image/jpeg, video/mp4 и т.д.
       file_size BIGINT,                   -- в байтах
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_post_attachments_post_id ON post_attachments (post_id);
CREATE INDEX idx_post_attachments_file_type ON post_attachments (file_type);
CREATE INDEX idx_post_attachments_created_at ON post_attachments (created_at);
CREATE INDEX idx_post_attachments_deleted_at ON post_attachments (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS post_attachments CASCADE;
-- +goose StatementEnd
