package posts

import (
	"context"
	"fmt"
	"mpb/pkg/db"
	"strings"
	"time"
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
	return &PostsRepository{db: db}
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
			return fmt.Errorf("failed to scan inserted post: %w", err)
		}
	}

	return nil
}

func (r *PostsRepository) FindByID(ctx context.Context, postID int) (*Post, error) {
	var post Post
	const query = `SELECT * FROM posts WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.Conn.GetContext(ctx, &post, query, postID); err != nil {
		return nil, fmt.Errorf("failed to find post by id: %w", err)
	}
	return &post, nil
}

func (r *PostsRepository) Delete(ctx context.Context, postID int) error {
	const query = `UPDATE posts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`
	res, err := r.db.Conn.ExecContext(ctx, query, postID)
	if err != nil {
		return fmt.Errorf("failed to delete post: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("no post found to delete")
	}

	return nil
}

func (r *PostsRepository) Update(ctx context.Context, post *Post) error {
	const query = `
		UPDATE posts
		SET title = :title,
		    description = :description,
		    tag = :tag,
		    "like" = :like,
		    count_viewers = :count_viewers,
		    updated_at = NOW()
		WHERE id = :id AND deleted_at IS NULL
		RETURNING updated_at
	`

	rows, err := r.db.Conn.NamedQueryContext(ctx, query, post)
	if err != nil {
		return fmt.Errorf("failed to update post: %w", err)
	}
	defer rows.Close()

	if rows.Next() {
		if err := rows.Scan(&post.UpdatedAt); err != nil {
			return fmt.Errorf("failed to scan updated post: %w", err)
		}
	}

	return nil
}

func (r *PostsRepository) List(ctx context.Context, f PostFilter) ([]Post, error) {
	query := `SELECT * FROM posts WHERE 1=1`
	var args []interface{}

	if f.OnlyActive {
		query += ` AND deleted_at IS NULL`
	}

	if f.UserID != nil {
		args = append(args, *f.UserID)
		query += fmt.Sprintf(" AND user_id = $%d", len(args))
	}
	if f.Tag != nil {
		args = append(args, *f.Tag)
		query += fmt.Sprintf(" AND tag = $%d", len(args))
	}
	if f.Title != nil {
		args = append(args, "%"+*f.Title+"%")
		query += fmt.Sprintf(" AND title ILIKE $%d", len(args))
	}
	if f.FromDate != nil {
		args = append(args, *f.FromDate)
		query += fmt.Sprintf(" AND created_at >= $%d", len(args))
	}
	if f.ToDate != nil {
		args = append(args, *f.ToDate)
		query += fmt.Sprintf(" AND created_at <= $%d", len(args))
	}

	orderBy := "created_at DESC"
	if f.OrderBy != "" {
		validOrderColumns := map[string]bool{
			"created_at": true, "updated_at": true, "like": true, "count_viewers": true,
		}
		if validOrderColumns[strings.Split(f.OrderBy, " ")[0]] {
			orderBy = f.OrderBy
		}
	}
	query += fmt.Sprintf(" ORDER BY %s", orderBy)

	if f.Limit > 0 {
		args = append(args, f.Limit)
		query += fmt.Sprintf(" LIMIT $%d", len(args))
	}
	if f.Offset > 0 {
		args = append(args, f.Offset)
		query += fmt.Sprintf(" OFFSET $%d", len(args))
	}

	var posts []Post
	if err := r.db.Conn.SelectContext(ctx, &posts, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list posts: %w", err)
	}

	return posts, nil
}
