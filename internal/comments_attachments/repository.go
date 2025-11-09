package comments_attachments

import (
	"context"
	"fmt"
	"mpb/pkg/db"
)

type CommentAttachmentsRepository struct {
	db *db.Db
}

func NewCommentAttachmentsRepository(db *db.Db) *CommentAttachmentsRepository {
	return &CommentAttachmentsRepository{db: db}
}

func (r *CommentAttachmentsRepository) Create(ctx context.Context, att *CommentAttachment) error {
	const query = `
		INSERT INTO comment_attachments (comment_id, file_url, file_type, file_size)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	if err := r.db.Conn.QueryRowxContext(ctx, query, att.CommentID, att.FileURL, att.FileType, att.FileSize).
		Scan(&att.ID, &att.CreatedAt); err != nil {
		return fmt.Errorf("failed to insert comment attachment: %w", err)
	}
	return nil
}

func (r *CommentAttachmentsRepository) ListByComment(ctx context.Context, commentID int) ([]CommentAttachment, error) {
	const query = `SELECT * FROM comment_attachments WHERE comment_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	var result []CommentAttachment
	if err := r.db.Conn.SelectContext(ctx, &result, query, commentID); err != nil {
		return nil, fmt.Errorf("failed to list comment attachments: %w", err)
	}
	return result, nil
}

func (r *CommentAttachmentsRepository) Delete(ctx context.Context, id int) error {
	const query = `UPDATE comment_attachments SET deleted_at = NOW() WHERE id = $1`
	if _, err := r.db.Conn.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete comment attachment: %w", err)
	}
	return nil
}
