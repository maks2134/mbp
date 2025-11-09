package dto

import "time"

type UserProfileResponse struct {
	ID               int       `json:"id"`
	Name             string    `json:"name"`
	Username         string    `json:"username"`
	Email            *string   `json:"email,omitempty"`
	Age              int       `json:"age"`
	IsActive         bool      `json:"is_active"`
	PostsCount       int       `json:"posts_count"`
	AttachmentsCount int       `json:"attachments_count"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}
