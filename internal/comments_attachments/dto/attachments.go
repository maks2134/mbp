package dto

type UploadCommentAttachmentRequest struct {
	CommentID int `params:"id" validate:"required,min=1"`
}

type DeleteCommentAttachmentRequest struct {
	ID int `params:"id" validate:"required,min=1"`
}
