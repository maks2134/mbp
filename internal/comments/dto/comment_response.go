package dto

import "time"

type CommentResponse struct {
	ID        int        `json:"id"`
	PostID    int        `json:"post_id"`
	UserID    int        `json:"user_id"`
	Text      string     `json:"text"`
	Like      int        `json:"like"`
	Blocked   bool       `json:"blocked"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}
