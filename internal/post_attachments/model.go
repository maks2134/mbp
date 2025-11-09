package post_attachments

import "time"

type PostAttachment struct {
	ID        int        `db:"id" json:"id"`
	PostID    int        `db:"post_id" json:"post_id"`
	FileURL   string     `db:"file_url" json:"file_url"`
	FileType  string     `db:"file_type" json:"file_type"`
	FileSize  int64      `db:"file_size" json:"file_size"`
	CreatedAt time.Time  `db:"created_at" json:"created_at"`
	DeletedAt *time.Time `db:"deleted_at" json:"deleted_at,omitempty"`
}
