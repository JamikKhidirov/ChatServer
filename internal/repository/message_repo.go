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
	var replyToID, forwardFrom *string
	if msg.ReplyToID != nil && *msg.ReplyToID != "" {
		replyToID = msg.ReplyToID
	}
	if msg.ForwardFrom != nil && *msg.ForwardFrom != "" {
		forwardFrom = msg.ForwardFrom
	}
	_, err := r.db.Exec(
		`INSERT INTO messages (id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		msg.ID, msg.ChatID, msg.SenderID, msg.Content, msg.Type,
		replyToID, forwardFrom, msg.FileName, msg.FileSize, msg.FilePath,
		boolToInt(msg.Pinned),
		msg.CreatedAt.Format(time.RFC3339), msg.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *MessageRepository) FindByID(id string) (*domain.Message, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE id = ?`, id,
	)
	return scanMessage(row)
}

func (r *MessageRepository) FindByChatID(chatID string, limit, offset int) ([]*domain.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
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

func (r *MessageRepository) Search(chatID, query string, limit, offset int) ([]*domain.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND content LIKE ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT ? OFFSET ?`,
		chatID, "%"+query+"%", limit, offset,
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
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1`, chatID,
	)
	return scanMessage(row)
}

func (r *MessageRepository) TogglePin(msgID string, pinned bool) error {
	_, err := r.db.Exec(
		`UPDATE messages SET pinned=?, updated_at=? WHERE id=?`,
		boolToInt(pinned), time.Now().Format(time.RFC3339), msgID,
	)
	return err
}

func (r *MessageRepository) GetPinned(chatID string) ([]*domain.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND pinned = 1 AND deleted_at IS NULL
		ORDER BY created_at DESC`, chatID,
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

// Reactions
func (r *MessageRepository) AddReaction(msgID, userID, emoji string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO reactions (message_id, user_id, emoji, created_at) VALUES (?, ?, ?, ?)`,
		msgID, userID, emoji, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *MessageRepository) RemoveReaction(msgID, userID, emoji string) error {
	_, err := r.db.Exec(
		`DELETE FROM reactions WHERE message_id = ? AND user_id = ? AND emoji = ?`,
		msgID, userID, emoji,
	)
	return err
}

func (r *MessageRepository) GetReactions(msgID string) ([]*domain.Reaction, error) {
	rows, err := r.db.Query(
		`SELECT message_id, user_id, emoji, created_at FROM reactions WHERE message_id = ?`,
		msgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reactions := make([]*domain.Reaction, 0)
	for rows.Next() {
		var (
			reaction  domain.Reaction
			createdAt string
		)
		if err := rows.Scan(&reaction.MessageID, &reaction.UserID, &reaction.Emoji, &createdAt); err != nil {
			return nil, err
		}
		reaction.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
		reactions = append(reactions, &reaction)
	}
	return reactions, nil
}

// Read receipts
func (r *MessageRepository) AddReadReceipt(msgID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO read_receipts (message_id, user_id, read_at) VALUES (?, ?, ?)`,
		msgID, userID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *MessageRepository) GetReadReceipts(msgID string) ([]*domain.ReadReceipt, error) {
	rows, err := r.db.Query(
		`SELECT message_id, user_id, read_at FROM read_receipts WHERE message_id = ?`,
		msgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := make([]*domain.ReadReceipt, 0)
	for rows.Next() {
		var (
			receipt domain.ReadReceipt
			readAt  string
		)
		if err := rows.Scan(&receipt.MessageID, &receipt.UserID, &readAt); err != nil {
			return nil, err
		}
		receipt.ReadAt, _ = time.Parse(time.RFC3339, readAt)
		receipts = append(receipts, &receipt)
	}
	return receipts, nil
}

type messageScanner interface {
	Scan(dest ...interface{}) error
}

func scanMessage(row messageScanner) (*domain.Message, error) {
	var (
		msg        domain.Message
		replyToID  sql.NullString
		forwardFrom sql.NullString
		pinnedInt  int
		createdAt  string
		updatedAt  string
		deletedAt  sql.NullString
	)
	err := row.Scan(&msg.ID, &msg.ChatID, &msg.SenderID, &msg.Content, &msg.Type,
		&replyToID, &forwardFrom, &msg.FileName, &msg.FileSize, &msg.FilePath,
		&pinnedInt, &createdAt, &updatedAt, &deletedAt)
	if err != nil {
		return nil, err
	}
	if replyToID.Valid {
		msg.ReplyToID = &replyToID.String
	}
	if forwardFrom.Valid {
		msg.ForwardFrom = &forwardFrom.String
	}
	msg.Pinned = pinnedInt == 1
	if deletedAt.Valid {
		t, _ := time.Parse(time.RFC3339, deletedAt.String)
		msg.DeletedAt = &t
	}
	msg.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	msg.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &msg, nil
}
