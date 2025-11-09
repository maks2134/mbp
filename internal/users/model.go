package users

import (
	"mpb/internal/user"
)

type UserProfile struct {
	user.User
	PostsCount       int `json:"posts_count"`
	AttachmentsCount int `json:"attachments_count"`
}

type UserFilter struct {
	Username *string
	IsActive *bool
	Limit    int
	Offset   int
	OrderBy  string
}
