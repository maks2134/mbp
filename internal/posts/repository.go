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
	db     *db.Db
	config *configs.Config
}

func NewPostsRepository(db *db.Db, config *configs.Config) *PostsRepository {
	return &PostsRepository{
		db:     db,
		config: config,
	}
}

func (repo *PostsRepository) Save(ctx context.Context, post *Post) error {
	query := `INSERT INTO posts (user_id, title, description, tag, "like", count_viewers, created_at, updated_at)
	Values ($1, $2, $3, $4, $5, $6, $7, $8)
	RETURNING id, created_at, updated_at`

	err := repo.db.Conn.QueryRowContext(ctx, query,
		post.UserID,
		post.Title,
		post.Description,
		post.Tag,
		post.Like,
		post.CountViewers,
		post.CreatedAt,
		post.UpdatedAt,
	).Scan(&post.ID, &post.CreatedAt, &post.UpdatedAt)

	if err != nil {
		return fmt.Errorf("failed to create post: %w", err)
	}

	return nil
}
