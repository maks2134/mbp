package post_attachments

import (
	"context"
	"fmt"
	"mpb/pkg/db"
)

type PostAttacmentsRepository struct {
	db *db.Db
}

func NewPostAttacmentsRepository(db *db.Db) *PostAttacmentsRepository {
	return &PostAttacmentsRepository{
		db: db,
	}
}

func (r *PostAttacmentsRepository) Create(ctx context.Context, att *PostAttachment) error {
	const query = `
		INSERT INTO post_attachments (post_id, file_url, file_type, file_size)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	if err := r.db.Conn.QueryRowxContext(ctx, query,
		att.PostID, att.FileURL, att.FileType, att.FileSize).
		Scan(&att.ID, &att.CreatedAt); err != nil {
		return fmt.Errorf("failed to insert attachment: %w", err)
	}
	return nil
}

func (r *PostAttacmentsRepository) ListByPost(ctx context.Context, postID int) ([]PostAttachment, error) {
	const query = `SELECT * FROM post_attachments WHERE post_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	var result []PostAttachment
	if err := r.db.Conn.SelectContext(ctx, &result, query, postID); err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	return result, nil
}

func (r *PostAttacmentsRepository) Delete(ctx context.Context, id int) error {
	const query = `UPDATE post_attachments SET deleted_at = NOW() WHERE id = $1`
	if _, err := r.db.Conn.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}
	return nil
}
