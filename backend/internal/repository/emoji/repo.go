package emojirepo

import (
	"database/sql"
	"time"

	emojidomain "ChatServerGolang/backend/internal/domain/emoji"
	"ChatServerGolang/backend/internal/repository"
)

type customEmojiRepository struct {
	db *sql.DB
}

func NewCustomEmojiRepository(db *sql.DB) repository.CustomEmojiRepository {
	return &customEmojiRepository{db: db}
}

func (r *customEmojiRepository) Create(emoji *emojidomain.CustomEmoji) error {
	_, err := r.db.Exec(
		`INSERT INTO custom_emojis (id, user_id, shortcode, file_url, file_path, animated, created_at) VALUES (?, ?, ?, ?, ?, ?, ?)`,
		emoji.ID, emoji.UserID, emoji.Shortcode, emoji.FileURL, emoji.FilePath, repository.BoolToInt(emoji.Animated), emoji.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *customEmojiRepository) FindByID(id string) (*emojidomain.CustomEmoji, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, shortcode, file_url, file_path, animated, created_at FROM custom_emojis WHERE id = ?`, id,
	)
	return scanEmoji(row)
}

func (r *customEmojiRepository) FindByUserID(userID string) ([]*emojidomain.CustomEmoji, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, shortcode, file_url, file_path, animated, created_at FROM custom_emojis WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emojis []*emojidomain.CustomEmoji
	for rows.Next() {
		e, err := scanEmoji(rows)
		if err != nil {
			return nil, err
		}
		emojis = append(emojis, e)
	}
	return emojis, nil
}

func (r *customEmojiRepository) FindAll() ([]*emojidomain.CustomEmoji, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, shortcode, file_url, file_path, animated, created_at FROM custom_emojis ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var emojis []*emojidomain.CustomEmoji
	for rows.Next() {
		e, err := scanEmoji(rows)
		if err != nil {
			return nil, err
		}
		emojis = append(emojis, e)
	}
	return emojis, nil
}

func (r *customEmojiRepository) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM custom_emojis WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

func scanEmoji(row repository.Scanner) (*emojidomain.CustomEmoji, error) {
	var e emojidomain.CustomEmoji
	var createdAt string
	var animated int
	err := row.Scan(&e.ID, &e.UserID, &e.Shortcode, &e.FileURL, &e.FilePath, &animated, &createdAt)
	if err != nil {
		return nil, err
	}
	e.Animated = animated == 1
	e.CreatedAt = repository.ParseTime(createdAt)
	return &e, nil
}
