package user_attachments

import (
	"context"
	"fmt"
	"mpb/pkg/db"
)

type UserAttachmentsRepository struct {
	db *db.Db
}

func NewUserAttachmentsRepository(db *db.Db) *UserAttachmentsRepository {
	return &UserAttachmentsRepository{db: db}
}

func (r *UserAttachmentsRepository) Create(ctx context.Context, att *UserAttachment) error {
	const query = `
		INSERT INTO user_attachments (user_id, file_url, file_type, file_size)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	if err := r.db.Conn.QueryRowxContext(ctx, query,
		att.UserID, att.FileURL, att.FileType, att.FileSize).
		Scan(&att.ID, &att.CreatedAt); err != nil {
		return fmt.Errorf("failed to insert attachment: %w", err)
	}
	return nil
}

func (r *UserAttachmentsRepository) ListByUser(ctx context.Context, userID int) ([]UserAttachment, error) {
	const query = `SELECT * FROM user_attachments WHERE user_id = $1 AND deleted_at IS NULL ORDER BY created_at DESC`
	var result []UserAttachment
	if err := r.db.Conn.SelectContext(ctx, &result, query, userID); err != nil {
		return nil, fmt.Errorf("failed to list attachments: %w", err)
	}
	return result, nil
}

func (r *UserAttachmentsRepository) Delete(ctx context.Context, id int) error {
	const query = `UPDATE user_attachments SET deleted_at = NOW() WHERE id = $1`
	if _, err := r.db.Conn.ExecContext(ctx, query, id); err != nil {
		return fmt.Errorf("failed to delete attachment: %w", err)
	}
	return nil
}

func (r *UserAttachmentsRepository) FindByID(ctx context.Context, id int) (*UserAttachment, error) {
	var att UserAttachment
	const query = `SELECT * FROM user_attachments WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.Conn.GetContext(ctx, &att, query, id); err != nil {
		return nil, fmt.Errorf("failed to find attachment by id: %w", err)
	}
	return &att, nil
}
