package posts

import "time"

type PostCreatedEvent struct {
	ID        int       `json:"id"`
	UserID    int       `json:"user_id"`
	Title     string    `json:"title"`
	CreatedAt time.Time `json:"created_at"`
}

type PostViewedEvent struct {
	PostID int `json:"post_id"`
	Views  int `json:"views"`
}

type PostLikedEvent struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
	Likes  int `json:"likes"`
}

type PostUnlikedEvent struct {
	PostID int `json:"post_id"`
	UserID int `json:"user_id"`
	Likes  int `json:"likes"`
}
