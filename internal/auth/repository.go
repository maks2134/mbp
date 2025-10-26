package auth

import (
	"database/sql"
	"errors"
	"mpb/configs"
	model "mpb/internal/user"
	"mpb/pkg/db"
	"mpb/pkg/errors_constant"
)

type AuthRepository struct {
	config *configs.Config
	db     *db.Db
}

func NewAuthRepository(config *configs.Config, database *db.Db) *AuthRepository {
	return &AuthRepository{
		config: config,
		db:     database,
	}
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
