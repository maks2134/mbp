package auth

import (
	_ "errors"
	"golang.org/x/crypto/bcrypt"
	"mpb/configs"
	"mpb/pkg/db"
	"mpb/pkg/errors_constant"
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

func (repo *AuthRepository) Login(username, password string) (string, error) {
	var exists bool
	var hashed string

}
