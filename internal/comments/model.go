package comments

import "time"

type Comment struct {
	ID        int        `db:"id"`
	PostID    int        `db:"post_id"`
	UserID    int        `db:"user_id"`
	Text      string     `db:"text"`
	Like      int        `db:"like"`
	Blocked   bool       `db:"blocked"`
	CreatedAt time.Time  `db:"created_at"`
	UpdatedAt time.Time  `db:"updated_at"`
	DeletedAt *time.Time `db:"deleted_at"`
}
