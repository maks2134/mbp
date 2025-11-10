package dto

import "time"

type StoryResponse struct {
	ID         int       `json:"id"`
	UserID     int       `json:"user_id"`
	FileURL    string    `json:"file_url"`
	FileType   string    `json:"file_type"`
	ViewsCount int       `json:"views_count"`
	ExpiresAt  time.Time `json:"expires_at"`
	CreatedAt  time.Time `json:"created_at"`
	IsViewed   bool      `json:"is_viewed,omitempty"`
}

type StoryCreateRequest struct {
	FileURL  string `json:"file_url"`
	FileType string `json:"file_type"`
}
