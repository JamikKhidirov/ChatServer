package savedmsgrepo

import (
	"database/sql"
	"time"

	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	"ChatServerGolang/backend/internal/repository"
)

type savedMessageRepository struct {
	db *sql.DB
}

func NewSavedMessageRepository(db *sql.DB) repository.SavedMessageRepository {
	return &savedMessageRepository{db: db}
}

func (r *savedMessageRepository) Save(msg *chatdomain.SavedMessage) error {
	_, err := r.db.Exec(
		`INSERT INTO saved_messages (id, user_id, message_id, chat_id, created_at) VALUES (?, ?, ?, ?, ?)`,
		msg.ID, msg.UserID, msg.MessageID, msg.ChatID, msg.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *savedMessageRepository) FindByUserID(userID string, limit, offset int) ([]*chatdomain.SavedMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, message_id, chat_id, created_at FROM saved_messages
		WHERE user_id = ? ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		userID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var msgs []*chatdomain.SavedMessage
	for rows.Next() {
		var m chatdomain.SavedMessage
		var createdAt string
		if err := rows.Scan(&m.ID, &m.UserID, &m.MessageID, &m.ChatID, &createdAt); err != nil {
			return nil, err
		}
		m.CreatedAt = repository.ParseTime(createdAt)
		msgs = append(msgs, &m)
	}
	return msgs, nil
}

func (r *savedMessageRepository) CountByUserID(userID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM saved_messages WHERE user_id = ?`, userID).Scan(&count)
	return count, err
}

func (r *savedMessageRepository) Delete(id, userID string) error {
	_, err := r.db.Exec(`DELETE FROM saved_messages WHERE id = ? AND user_id = ?`, id, userID)
	return err
}

func (r *savedMessageRepository) Exists(userID, messageID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM saved_messages WHERE user_id = ? AND message_id = ?`, userID, messageID).Scan(&count)
	return count > 0, err
}
