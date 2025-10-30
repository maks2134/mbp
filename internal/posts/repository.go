package posts

import (
	"context"
	"fmt"
	"log"
	"mpb/pkg/db"
	"strings"
	"time"

	"github.com/jmoiron/sqlx"
)

type PostFilter struct {
	UserID     *int
	Tag        *string
	Title      *string
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

func NewPostsRepository(db *db.Db) *PostsRepository {
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
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

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
	defer func(rows *sqlx.Rows) {
		err := rows.Close()
		if err != nil {
			log.Printf("failed to close rows: %v", err)
		}
	}(rows)

	if rows.Next() {
		if err := rows.Scan(&post.UpdatedAt); err != nil {
			return fmt.Errorf("failed to update post: %w", err)
		}
	}

	return nil
}

func (r *PostsRepository) List(ctx context.Context, f PostFilter) ([]Post, error) {
	query := `SELECT * FROM posts WHERE deleted_at IS NULL`
	var args []interface{}
	var conditions []string

	if f.UserID != nil {
		conditions = append(conditions, fmt.Sprintf("user_id = $%d", len(args)+1))
		args = append(args, *f.UserID)
	}
	if f.Tag != nil {
		conditions = append(conditions, fmt.Sprintf("tag = $%d", len(args)+1))
		args = append(args, *f.Tag)
	}
	if f.Title != nil {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", len(args)+1))
		args = append(args, "%"+*f.Title+"%")
	}
	if f.FromDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at >= $%d", len(args)+1))
		args = append(args, *f.FromDate)
	}
	if f.ToDate != nil {
		conditions = append(conditions, fmt.Sprintf("created_at <= $%d", len(args)+1))
		args = append(args, *f.ToDate)
	}
	if f.OnlyActive {
		conditions = append(conditions, "deleted_at IS NULL")
	}

	if len(conditions) > 0 {
		query += " AND " + strings.Join(conditions, " AND ")
	}

	orderBy := "created_at DESC"
	if f.OrderBy != "" {
		orderBy = f.OrderBy
	}
	query += fmt.Sprintf("ORDER BY %s", orderBy)

	if f.Limit > 0 {
		query += fmt.Sprintf(" LIMIT %d", f.Limit)
	}
	if f.Offset > 0 {
		query += fmt.Sprintf(" OFFSET %d", f.Offset)
	}

	var posts []Post
	if err := r.db.Conn.SelectContext(ctx, &posts, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}
