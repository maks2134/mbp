package auth

import (
	model "mpb/internal/user"
)

type AuthService struct {
	repo *AuthRepository
}

func NewAuthService(repo *AuthRepository) *AuthService {
	return &AuthService{
		repo: repo,
	}
}

func (s *AuthService) Register(username, password, email string) error {
	return s.repo.Register(username, password, email)
}

func (s *AuthService) Login(username, password string) (*model.User, error) {
	return s.repo.Login(username, password)
}
