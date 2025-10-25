package auth

import (
	"database/sql"
	"errors"
	_ "errors"
	"mpb/configs"
	model "mpb/internal/user"
	"mpb/pkg/db"
	"mpb/pkg/errors_constant"

	"golang.org/x/crypto/bcrypt"
)

type AuthRepository struct {
	config *configs.Config
	db     *db.Db
}

func NewAuthRepository(config *configs.Config, db *db.Db) *AuthRepository {
	return &AuthRepository{
		config: config,
		db:     db,
	}
}

func (repo *AuthRepository) Register(username, password, email string) error {
	var exists bool
	err := repo.db.Db.Get(&exists, `SELECT EXISTS (SELECT 1 FROM users WHERE username = $1)`, username)
	if err != nil {
		return err
	}
	if exists {
		return errors_constant.UserAlreadyExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	_, err = repo.db.Db.Exec(`INSERT INTO users (username, password_hash, email)
	VALUES ($1, $2, $3)`, username, string(hashed), email)

	return err
}

func (repo *AuthRepository) Login(username, password string) (*model.User, error) {
	user, _ := repo.FindByUsername(username)

	err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return nil, errors_constant.UserAuthenticationFailed
	}

	return user, nil
}

func (repo *AuthRepository) FindByUsername(username string) (*model.User, error) {
	var user *model.User
	err := repo.db.Db.Get(&user, `SELECT * FROM users WHERE username = $1 AND deleted_at IS NULL`, username)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, errors_constant.UserNotFound
		}
	}

	return user, nil
}
