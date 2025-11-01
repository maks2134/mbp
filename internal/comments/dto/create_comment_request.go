package dto

type CreateCommentRequest struct {
	PostID int    `json:"post_id" validate:"required"`
	Text   string `json:"text" validate:"required,min=1,max=500"`
}
