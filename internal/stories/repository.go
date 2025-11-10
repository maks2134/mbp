package stories

import (
	"context"
	"fmt"
	"mpb/pkg/db"
)

type StoriesRepository struct {
	db *db.Db
}

func NewStoriesRepository(db *db.Db) *StoriesRepository {
	return &StoriesRepository{db: db}
}

func (r *StoriesRepository) Create(ctx context.Context, story *Story) error {
	const query = `
		INSERT INTO stories (user_id, file_url, file_type, expires_at)
		VALUES ($1, $2, $3, $4)
		RETURNING id, created_at`
	if err := r.db.Conn.QueryRowxContext(ctx, query,
		story.UserID, story.FileURL, story.FileType, story.ExpiresAt).
		Scan(&story.ID, &story.CreatedAt); err != nil {
		return fmt.Errorf("failed to insert story: %w", err)
	}
	return nil
}

func (r *StoriesRepository) FindByID(ctx context.Context, id int) (*Story, error) {
	var story Story
	const query = `SELECT * FROM stories WHERE id = $1 AND deleted_at IS NULL AND expires_at > NOW()`
	if err := r.db.Conn.GetContext(ctx, &story, query, id); err != nil {
		return nil, fmt.Errorf("failed to find story by id: %w", err)
	}
	return &story, nil
}

func (r *StoriesRepository) ListByUser(ctx context.Context, userID int) ([]Story, error) {
	const query = `
		SELECT * FROM stories 
		WHERE user_id = $1 AND deleted_at IS NULL AND expires_at > NOW() 
		ORDER BY created_at DESC`
	var stories []Story
	if err := r.db.Conn.SelectContext(ctx, &stories, query, userID); err != nil {
		return nil, fmt.Errorf("failed to list stories: %w", err)
	}
	return stories, nil
}

func (r *StoriesRepository) ListActive(ctx context.Context, excludeUserID *int) ([]Story, error) {
	query := `
		SELECT * FROM stories 
		WHERE deleted_at IS NULL AND expires_at > NOW()`
	var args []interface{}

	if excludeUserID != nil {
		args = append(args, *excludeUserID)
		query += fmt.Sprintf(" AND user_id != $%d", len(args))
	}

	query += " ORDER BY created_at DESC"

	var stories []Story
	if err := r.db.Conn.SelectContext(ctx, &stories, query, args...); err != nil {
		return nil, fmt.Errorf("failed to list active stories: %w", err)
	}
	return stories, nil
}

func (r *StoriesRepository) IncrementViews(ctx context.Context, storyID int) error {
	const query = `UPDATE stories SET views_count = views_count + 1 WHERE id = $1`
	if _, err := r.db.Conn.ExecContext(ctx, query, storyID); err != nil {
		return fmt.Errorf("failed to increment views: %w", err)
	}
	return nil
}

func (r *StoriesRepository) RecordView(ctx context.Context, storyID, userID int) error {
	const query = `
		INSERT INTO story_views (story_id, user_id, viewed_at)
		VALUES ($1, $2, NOW())
		ON CONFLICT (story_id, user_id) DO NOTHING`
	if _, err := r.db.Conn.ExecContext(ctx, query, storyID, userID); err != nil {
		return fmt.Errorf("failed to record view: %w", err)
	}
	return nil
}

func (r *StoriesRepository) HasUserViewed(ctx context.Context, storyID, userID int) (bool, error) {
	var count int
	const query = `SELECT COUNT(*) FROM story_views WHERE story_id = $1 AND user_id = $2`
	if err := r.db.Conn.GetContext(ctx, &count, query, storyID, userID); err != nil {
		return false, fmt.Errorf("failed to check view: %w", err)
	}
	return count > 0, nil
}

func (r *StoriesRepository) DeleteExpired(ctx context.Context) error {
	const query = `UPDATE stories SET deleted_at = NOW() WHERE expires_at <= NOW() AND deleted_at IS NULL`
	if _, err := r.db.Conn.ExecContext(ctx, query); err != nil {
		return fmt.Errorf("failed to delete expired stories: %w", err)
	}
	return nil
}

func (r *StoriesRepository) Delete(ctx context.Context, storyID, userID int) error {
	const query = `UPDATE stories SET deleted_at = NOW() WHERE id = $1 AND user_id = $2 AND deleted_at IS NULL`
	res, err := r.db.Conn.ExecContext(ctx, query, storyID, userID)
	if err != nil {
		return fmt.Errorf("failed to delete story: %w", err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get affected rows: %w", err)
	}
	if rowsAffected == 0 {
		return fmt.Errorf("story not found or already deleted")
	}

	return nil
}
