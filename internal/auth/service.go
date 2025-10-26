package auth

import (
	"errors"
	"mpb/internal/auth/dto"
	"mpb/pkg/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo     *AuthRepository
	jwtKey   []byte
	tokenTTL time.Duration
}

func NewAuthService(repo *AuthRepository, jwtKey []byte, ttl time.Duration) *AuthService {
	return &AuthService{
		repo:     repo,
		jwtKey:   jwtKey,
		tokenTTL: ttl,
	}
}

func (s *AuthService) Register(req dto.RegisterRequest) error {
	hashed, err := security.HashPassword(req.Password)
	if err != nil {
		return err
	}
	return s.repo.Register(req.Username, hashed, req.Email, req.Name, req.Age)
}

func (s *AuthService) Login(username, password string) (*dto.LoginResponse, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	if !security.CheckPasswordHash(password, user.PasswordHash) {
		return nil, errors.New("invalid username or password")
	}

	claims := jwt.MapClaims{
		"user_id":  user.ID,
		"username": user.Username,
		"exp":      time.Now().Add(s.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(s.jwtKey)
	if err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token: tokenString,
		User:  user,
	}, nil
}
