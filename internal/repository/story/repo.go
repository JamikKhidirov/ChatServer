package storyrepo

import (
	"database/sql"
	"time"

	storydomain "ChatServerGolang/internal/domain/story"
	"ChatServerGolang/internal/repository"
)

type storyRepository struct {
	db *sql.DB
}

func NewStoryRepository(db *sql.DB) repository.StoryRepository {
	return &storyRepository{db: db}
}

func (r *storyRepository) Create(story *storydomain.Story) error {
	_, err := r.db.Exec(
		`INSERT INTO stories (id, user_id, file_path, file_url, type, caption, expires_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		story.ID, story.UserID, story.FilePath, story.FileURL, story.Type, story.Caption,
		story.ExpiresAt.Format(time.RFC3339), story.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *storyRepository) FindByID(id string) (*storydomain.Story, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, file_path, file_url, type, caption, expires_at, created_at
		FROM stories WHERE id = ? AND expires_at > ?`,
		id, time.Now().Format(time.RFC3339),
	)
	return scanStory(row)
}

func (r *storyRepository) FindActiveByUserID(userID string) ([]*storydomain.Story, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, file_path, file_url, type, caption, expires_at, created_at
		FROM stories WHERE user_id = ? AND expires_at > ?
		ORDER BY created_at DESC`, userID, time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []*storydomain.Story
	for rows.Next() {
		s, err := scanStory(rows)
		if err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, nil
}

func (r *storyRepository) FindActiveByFollowing(userIDs []string) ([]*storydomain.Story, error) {
	if len(userIDs) == 0 {
		return []*storydomain.Story{}, nil
	}
	placeholders := make([]string, len(userIDs))
	args := make([]interface{}, len(userIDs)+1)
	args[0] = time.Now().Format(time.RFC3339)
	for i, uid := range userIDs {
		placeholders[i] = "?"
		args[i+1] = uid
	}
	query := `SELECT id, user_id, file_path, file_url, type, caption, expires_at, created_at
		FROM stories WHERE user_id IN (` + joinStrings(placeholders, ",") + `) AND expires_at > ?
		ORDER BY created_at DESC`
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var stories []*storydomain.Story
	for rows.Next() {
		s, err := scanStory(rows)
		if err != nil {
			return nil, err
		}
		stories = append(stories, s)
	}
	return stories, nil
}

func (r *storyRepository) MarkExpired() error {
	_, err := r.db.Exec(`DELETE FROM stories WHERE expires_at <= ?`, time.Now().Format(time.RFC3339))
	return err
}

func (r *storyRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM stories WHERE id = ?`, id)
	return err
}

func (r *storyRepository) AddView(storyID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO story_views (story_id, user_id, viewed_at) VALUES (?, ?, ?)`,
		storyID, userID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *storyRepository) GetViewCount(storyID string) (int64, error) {
	var count int64
	err := r.db.QueryRow(`SELECT COUNT(*) FROM story_views WHERE story_id = ?`, storyID).Scan(&count)
	return count, err
}

func (r *storyRepository) GetViews(storyID string) ([]*storydomain.StoryView, error) {
	rows, err := r.db.Query(
		`SELECT story_id, user_id, viewed_at FROM story_views WHERE story_id = ? ORDER BY viewed_at DESC`,
		storyID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var views []*storydomain.StoryView
	for rows.Next() {
		var v storydomain.StoryView
		var viewedAt string
		if err := rows.Scan(&v.StoryID, &v.UserID, &viewedAt); err != nil {
			return nil, err
		}
		v.ViewedAt = repository.ParseTime(viewedAt)
		views = append(views, &v)
	}
	return views, nil
}

func (r *storyRepository) HasViewed(storyID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM story_views WHERE story_id = ? AND user_id = ?`, storyID, userID).Scan(&count)
	return count > 0, err
}

func scanStory(row repository.Scanner) (*storydomain.Story, error) {
	var s storydomain.Story
	var expiresAt, createdAt string
	err := row.Scan(&s.ID, &s.UserID, &s.FilePath, &s.FileURL, &s.Type, &s.Caption, &expiresAt, &createdAt)
	if err != nil {
		return nil, err
	}
	s.ExpiresAt = repository.ParseTime(expiresAt)
	s.CreatedAt = repository.ParseTime(createdAt)
	return &s, nil
}

func joinStrings(strs []string, sep string) string {
	result := ""
	for i, s := range strs {
		if i > 0 {
			result += sep
		}
		result += s
	}
	return result
}
