package auth

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"mpb/configs"
	model "mpb/internal/user"
	"mpb/pkg/db"
	"mpb/pkg/errors_constant"
	"time"

	"github.com/redis/go-redis/v9"
)

type AuthRepository struct {
	config *configs.Config
	db     *db.Db
	redis  *redis.Client
}

func NewAuthRepository(config *configs.Config, database *db.Db, red *redis.Client) *AuthRepository {
	return &AuthRepository{
		config: config,
		db:     database,
		redis:  red,
	}
}

func (repo *AuthRepository) SetRefreshToken(userId int, token string, ttl time.Duration) error {
	key := fmt.Sprintf("refresh_token:%d", userId)
	return repo.redis.Set(context.Background(), key, token, ttl).Err()
}

func (repo *AuthRepository) GetRefreshToken(userId int) (string, error) {
	key := fmt.Sprintf("refresh_token:%d", userId)
	return repo.redis.Get(context.Background(), key).Result()
}

func (repo *AuthRepository) DeleteRefreshToken(userId int) error {
	key := fmt.Sprintf("refresh_token:%d", userId)
	return repo.redis.Del(context.Background(), key).Err()
}

func (repo *AuthRepository) Register(username, passwordHash, email, name string, age int) error {
	var exists bool
	err := repo.db.Conn.Get(&exists, `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`, username)
	if err != nil {
		return err
	}
	if exists {
		return errors_constant.UserAlreadyExists
	}

	_, err = repo.db.Conn.Exec(
		`INSERT INTO users (username, password_hash, email, name, age, is_active) VALUES ($1,$2,$3,$4,$5, TRUE)`,
		username, passwordHash, email, name, age,
	)
	return err
}

func (repo *AuthRepository) FindByUsername(username string) (*model.User, error) {
	var user model.User
	err := repo.db.Conn.Get(&user, `SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors_constant.UserNotFound
		}
		return nil, err
	}
	return &user, nil
}
