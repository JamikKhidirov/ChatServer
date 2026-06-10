package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type scheduledMessageRepository struct {
	db *sql.DB
}

func NewScheduledMessageRepository(db *sql.DB) ScheduledMessageRepository {
	return &scheduledMessageRepository{db: db}
}

func (r *scheduledMessageRepository) Create(msg *domain.ScheduledMessage) error {
	_, err := r.db.Exec(
		`INSERT INTO scheduled_messages (id, chat_id, sender_id, content, type, reply_to_id, scheduled_at, created_at, sent)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.Type, msg.ReplyToID,
		msg.ScheduledAt, msg.CreatedAt.Format(time.RFC3339), boolToInt(msg.Sent),
	)
	return err
}

func (r *scheduledMessageRepository) FindPending() ([]*domain.ScheduledMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, scheduled_at, created_at, sent
		FROM scheduled_messages WHERE sent = 0 AND scheduled_at <= ?`,
		time.Now().Format(time.RFC3339),
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domain.ScheduledMessage, 0)
	for rows.Next() {
		m, err := scanScheduledMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *scheduledMessageRepository) FindByUserID(userID string) ([]*domain.ScheduledMessage, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, scheduled_at, created_at, sent
		FROM scheduled_messages WHERE sender_id = ? AND sent = 0 ORDER BY scheduled_at ASC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domain.ScheduledMessage, 0)
	for rows.Next() {
		m, err := scanScheduledMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, m)
	}
	return messages, nil
}

func (r *scheduledMessageRepository) MarkAsSent(id string) error {
	_, err := r.db.Exec(`UPDATE scheduled_messages SET sent = 1 WHERE id = ?`, id)
	return err
}

func (r *scheduledMessageRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM scheduled_messages WHERE id = ?`, id)
	return err
}

type scheduledMsgScanner interface {
	Scan(dest ...interface{}) error
}

func scanScheduledMessage(row scheduledMsgScanner) (*domain.ScheduledMessage, error) {
	var m domain.ScheduledMessage
	var replyToID sql.NullString
	var createdAt, scheduledAt string
	var sent int
	if err := row.Scan(&m.ID, &m.ChatID, &m.SenderID, &m.Content, &m.Type, &replyToID, &scheduledAt, &createdAt, &sent); err != nil {
		return nil, err
	}
	m.Sent = sent == 1
	m.ScheduledAt = scheduledAt
	m.CreatedAt = parseTime(createdAt)
	if replyToID.Valid {
		m.ReplyToID = &replyToID.String
	}
	return &m, nil
}
