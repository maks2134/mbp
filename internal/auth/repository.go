package auth

import (
	"mpb/configs"
	"mpb/pkg/db"
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
