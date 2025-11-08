package dto

import (
	"mpb/internal/user"
)

type LoginResponse struct {
	Token        string     `json:"token"`
	User         *user.User `json:"user"`
	RefreshToken string     `json:"refresh_token"`
}
