package posts

import (
	"context"
	"fmt"
	"mpb/configs"
	"mpb/pkg/db"
	"time"
)

type PostFilter struct {
	UserID     int
	Tag        string
	Title      string
	FromDate   *time.Time
	ToDate     *time.Time
	OnlyActive bool
	Limit      int
	Offset     int
	OrderBy    string
}

type PostsRepository struct {
	db *db.Db
}

func NewPostsRepository(db *db.Db, config *configs.Config) *PostsRepository {
	return &PostsRepository{
		db: db,
	}
}

func (r *PostsRepository) Save(ctx context.Context, post *Post) error {
	const query = `
		INSERT INTO posts (user_id, title, description, tag, "like", count_viewers)
		VALUES (:user_id, :title, :description, :tag, :like, :count_viewers)
		RETURNING id, created_at, updated_at
	`

	rows, err := r.db.Conn.NamedQueryContext(ctx, query, post)
	if err != nil {
		return fmt.Errorf("failed to insert post: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan post: %w", err)
		}
	}

	return nil
}

func (r *PostsRepository) FindByID(ctx context.Context, postID int) (*Post, error) {
	var post Post
	const query = `SELECT * FROM posts WHERE id = $1`
	if err := r.db.Conn.GetContext(ctx, &post, query, postID); err != nil {
		return nil, fmt.Errorf("failed to find post by id: %w", err)
	}

	return &post, nil
}

func (r *PostsRepository) Delete(ctx context.Context, postID int) error {
	query := `UPDATE posts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Conn.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if rowsAffected == 0 || err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	return nil
}

func (r *PostsRepository) Update(ctx context.Context, post *Post) error {
	const query = `UPDATE posts SET title = :title, description = :description,
    "like" = :like, count_viewers = :count_viewers
	WHERE id = :id AND deleted_at IS NULL
	RETURNING updated_at`

	rows, err := r.db.Conn.NamedQueryContext(ctx, query, post)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&post.UpdatedAt); err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}
	}

	return nil
}
