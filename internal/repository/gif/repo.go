package gifrepo

import (
	"database/sql"
	"ChatServerGolang/internal/repository"
)

type savedGifRepository struct {
	db *sql.DB
}

func NewSavedGifRepository(db *sql.DB) repository.SavedGifRepository {
	return &savedGifRepository{db: db}
}

func (r *savedGifRepository) Save(userID, gifURL string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO saved_gifs (user_id, gif_url, created_at) VALUES (?, ?, CURRENT_TIMESTAMP)`,
		userID, gifURL,
	)
	return err
}

func (r *savedGifRepository) FindByUserID(userID string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT gif_url FROM saved_gifs WHERE user_id = ? ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	urls := make([]string, 0)
	for rows.Next() {
		var url string
		if err := rows.Scan(&url); err != nil {
			return nil, err
		}
		urls = append(urls, url)
	}
	return urls, nil
}

func (r *savedGifRepository) Delete(userID, gifURL string) error {
	_, err := r.db.Exec(`DELETE FROM saved_gifs WHERE user_id = ? AND gif_url = ?`, userID, gifURL)
	return err
}




