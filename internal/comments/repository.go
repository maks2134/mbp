package comments

import (
	"context"
	"database/sql"
	"fmt"
	"mpb/pkg/db"
)

type CommentsRepository struct {
	db *db.Db
}

func NewCommentsRepository(db *db.Db) *CommentsRepository {
	return &CommentsRepository{db: db}
}

func (r *CommentsRepository) Create(ctx context.Context, c *Comment) error {
	const query = `
		INSERT INTO comments (post_id, user_id, text, blocked, "like")
		VALUES ($1, $2, $3, $4, $5)
		RETURNING id, created_at, updated_at
	`

	return r.db.Conn.QueryRowContext(ctx, query,
		c.PostID, c.UserID, c.Text, c.Blocked, c.Like,
	).Scan(&c.ID, &c.CreatedAt, &c.UpdatedAt)
}

func (r *CommentsRepository) Update(ctx context.Context, c *Comment) error {
	const query = `
		UPDATE comments
		SET text = $1,
		    blocked = $2,
		    "like" = $3,
		    updated_at = NOW()
		WHERE id = $4 AND deleted_at IS NULL
		RETURNING updated_at
	`

	return r.db.Conn.QueryRowContext(ctx, query,
		c.Text, c.Blocked, c.Like, c.ID,
	).Scan(&c.UpdatedAt)
}

func (r *CommentsRepository) Delete(ctx context.Context, commentID int) error {
	const query = `
		UPDATE comments
		SET deleted_at = NOW()
		WHERE id = $1 AND deleted_at IS NULL
	`
	res, err := r.db.Conn.ExecContext(ctx, query, commentID)
	if err != nil {
		return fmt.Errorf("failed to delete comment: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return sql.ErrNoRows
	}

	return nil
}

func (r *CommentsRepository) List(ctx context.Context, postID int) ([]Comment, error) {
	query := `
		SELECT * FROM comments
		WHERE deleted_at IS NULL AND post_id = $1
		ORDER BY created_at DESC
	`

	var comments []Comment
	if err := r.db.Conn.SelectContext(ctx, &comments, query, postID); err != nil {
		return nil, fmt.Errorf("failed to list comments: %w", err)
	}

	return comments, nil
}
