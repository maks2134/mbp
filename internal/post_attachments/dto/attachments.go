package dto

type UploadPostAttachmentRequest struct {
	PostID int `params:"id" validate:"required,min=1"`
}

type DeletePostAttachmentRequest struct {
	ID int `params:"id" validate:"required,min=1"`
}
