package auth

import (
	"errors"
	"mpb/internal/auth/dto"
	"mpb/pkg/security"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type AuthService struct {
	repo       *AuthRepository
	jwtKey     []byte
	tokenTTL   time.Duration
	refreshTTL time.Duration
}

func NewAuthService(repo *AuthRepository, jwtKey []byte, ttl time.Duration) *AuthService {
	return &AuthService{
		repo:       repo,
		jwtKey:     jwtKey,
		tokenTTL:   ttl,
		refreshTTL: 7 * 24 * time.Hour,
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

	accessToken, err := s.generateAccessToken(user.ID, user.Username)
	if err != nil {
		return nil, err
	}

	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SetRefreshToken(user.ID, refreshToken, s.refreshTTL); err != nil {
		return nil, err
	}

	return &dto.LoginResponse{
		Token:        accessToken,
		RefreshToken: refreshToken,
		User:         user,
	}, nil
}

func (s *AuthService) Refresh(refreshToken string) (*dto.RefreshResponse, error) {
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		return s.jwtKey, nil
	})
	if err != nil || !token.Valid {
		return nil, errors.New("invalid refresh token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return nil, errors.New("invalid token payload")
	}
	userID := int(userIDFloat)

	storedToken, err := s.repo.GetRefreshToken(userID)
	if err != nil || storedToken != refreshToken {
		return nil, errors.New("refresh token not found or expired")
	}

	newAccess, err := s.generateAccessToken(userID, claims["username"].(string))
	if err != nil {
		return nil, err
	}

	newRefresh, err := s.generateRefreshToken(userID)
	if err != nil {
		return nil, err
	}

	if err := s.repo.SetRefreshToken(userID, newRefresh, s.refreshTTL); err != nil {
		return nil, err
	}

	return &dto.RefreshResponse{
		AccessToken:  newAccess,
		RefreshToken: newRefresh,
	}, nil
}

func (s *AuthService) generateAccessToken(userID int, username string) (string, error) {
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      time.Now().Add(s.tokenTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}

func (s *AuthService) generateRefreshToken(userID int) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.refreshTTL).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(s.jwtKey)
}
