-- +goose Up
-- +goose StatementBegin
CREATE TABLE comment_attachments (
       id SERIAL PRIMARY KEY,
       comment_id INT NOT NULL REFERENCES comments(id) ON DELETE CASCADE,
       file_url TEXT NOT NULL,
       file_type TEXT NOT NULL,
       file_size BIGINT,
       created_at TIMESTAMP NOT NULL DEFAULT NOW(),
       deleted_at TIMESTAMP NULL
);

CREATE INDEX idx_comment_attachments_comment_id ON comment_attachments (comment_id);
CREATE INDEX idx_comment_attachments_created_at ON comment_attachments (created_at);
CREATE INDEX idx_comment_attachments_deleted_at ON comment_attachments (deleted_at);
-- +goose StatementEnd

-- +goose Down
-- +goose StatementBegin
DROP TABLE IF EXISTS comment_attachments CASCADE;
-- +goose StatementEnd
