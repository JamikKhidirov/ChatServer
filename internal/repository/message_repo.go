package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type MessageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) *MessageRepository {
	return &MessageRepository{db: db}
}

func (r *MessageRepository) Create(msg *domain.Message) error {
	var replyToID *string
	if msg.ReplyToID != nil && *msg.ReplyToID != "" {
		replyToID = msg.ReplyToID
	}
	_, err := r.db.Exec(
		`INSERT INTO messages (id, chat_id, sender_id, content, type, reply_to_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.Type, replyToID,
		msg.CreatedAt.Format(time.RFC3339), msg.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *MessageRepository) FindByID(id string) (*domain.Message, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, created_at, updated_at, deleted_at
		FROM messages WHERE id = ?`, id,
	)
	return scanMessage(row)
}

func (r *MessageRepository) FindByChatID(chatID string, limit, offset int) ([]*domain.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ?
		ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		chatID, limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domain.Message, 0)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *MessageRepository) Update(msg *domain.Message) error {
	_, err := r.db.Exec(
		`UPDATE messages SET content=?, updated_at=? WHERE id=? AND deleted_at IS NULL`,
		msg.Content, time.Now().Format(time.RFC3339), msg.ID,
	)
	return err
}

func (r *MessageRepository) SoftDelete(id string) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE messages SET deleted_at=?, updated_at=? WHERE id=?`,
		now.Format(time.RFC3339), now.Format(time.RFC3339), id,
	)
	return err
}

func (r *MessageRepository) GetLastMessage(chatID string) (*domain.Message, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1`, chatID,
	)
	return scanMessage(row)
}

type messageScanner interface {
	Scan(dest ...interface{}) error
}

func scanMessage(row messageScanner) (*domain.Message, error) {
	var (
		msg       domain.Message
		replyToID sql.NullString
		createdAt string
		updatedAt string
		deletedAt sql.NullString
	)
	err := row.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.Type,
		&replyToID, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if replyToID.Valid {
		msg.ReplyToID = &replyToID.String
	}
	if deletedAt.Valid {
		t, _ := time.Parse(time.RFC3339, deletedAt.String)
		msg.DeletedAt = &t
	}
	msg.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	msg.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &msg, nil
}
