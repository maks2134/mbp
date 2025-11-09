package users

import (
	"context"
	"fmt"
	"mpb/internal/user"
	"mpb/pkg/db"
	"strings"
)

type UsersRepository struct {
	db *db.Db
}

func NewUsersRepository(db *db.Db) *UsersRepository {
	return &UsersRepository{db: db}
}

func (r *UsersRepository) FindByID(ctx context.Context, userID int) (*user.User, error) {
	var u user.User
	const query = `SELECT * FROM users WHERE id = $1 AND deleted_at IS NULL`
	if err := r.db.Conn.GetContext(ctx, &u, query, userID); err != nil {
		return nil, fmt.Errorf("failed to find user by id: %w", err)
	}
	return &u, nil
}

func (r *UsersRepository) FindByUsername(ctx context.Context, username string) (*user.User, error) {
	var u user.User
	const query = `SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL`
	if err := r.db.Conn.GetContext(ctx, &u, query, username); err != nil {
		return nil, fmt.Errorf("failed to find user by username: %w", err)
	}
	return &u, nil
}

func (r *UsersRepository) List(ctx context.Context, f UserFilter) ([]user.User, error) {
	query := `SELECT * FROM users WHERE deleted_at IS NULL`
	var args []interface{}

	if f.Username != nil {
		args = append(args, "%"+*f.Username+"%")
		query += fmt.Sprintf(" AND username ILIKE $%d", len(args))
	}

	if f.IsActive != nil {
		args = append(args, *f.IsActive)
		query += fmt.Sprintf(" AND is_active = $%d", len(args))
	}

	orderBy := "created_at DESC"
	if f.OrderBy != "" {
		validOrderColumns := map[string]bool{
			"created_at": true, "updated_at": true, "username": true,
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

	var users []user.User
	if err := r.db.Conn.SelectContext(ctx, &users, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list users: %w", err)
	}

	return users, nil
}

func (r *UsersRepository) GetPostsCount(ctx context.Context, userID int) (int, error) {
	var count int
	const query = `SELECT COUNT(*) FROM posts WHERE user_id = $1 AND deleted_at IS NULL`
	if err := r.db.Conn.GetContext(ctx, &count, query, userID); err != nil {
		return 0, fmt.Errorf("failed to get posts count: %w", err)
	}
	return count, nil
}

func (r *UsersRepository) GetAttachmentsCount(ctx context.Context, userID int) (int, error) {
	var count int
	const query = `SELECT COUNT(*) FROM user_attachments WHERE user_id = $1 AND deleted_at IS NULL`
	if err := r.db.Conn.GetContext(ctx, &count, query, userID); err != nil {
		return 0, fmt.Errorf("failed to get attachments count: %w", err)
	}
	return count, nil
}
