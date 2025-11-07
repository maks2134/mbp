package dto

type LoginRequest struct {
	Username     string `json:"username" validate:"required,min=3"`
	Password     string `json:"password" validate:"required,min=6"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
