package dto

type UpdateCommentRequest struct {
	Text *string `json:"text" validate:"omitempty,min=1,max=500"`
}
