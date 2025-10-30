package dto

type CreatePostRequest struct {
	Title       string `json:"title" validate:"required,min=3,max=200"`
	Description string `json:"description" validate:"required,min=10"`
	Tag         string `json:"tag" validate:"omitempty,max=50"`
}

type UpdatePostRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=3,max=200"`
	Description *string `json:"description" validate:"omitempty,min=10"`
	Tag         *string `json:"tag" validate:"omitempty,max=50"`
}
