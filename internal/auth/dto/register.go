package dto

type RegisterRequest struct {
	Name     string `json:"name" validate:"required,min=2,max=100"`
	Username string `json:"username" validate:"required,min=3,max=50"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
	Age      int    `json:"age" validate:"gte=0,lte=123"`
}
