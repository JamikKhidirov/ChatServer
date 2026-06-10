package repository

import (
	"database/sql"
	"strings"
	"time"

	"ChatServerGolang/internal/domain"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *domain.Message) error {
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

func (r *messageRepository) FindByID(id string) (*domain.Message, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE id = ?`, id,
	)
	return scanMessage(row)
}

func (r *messageRepository) FindByChatID(chatID string, limit, offset int) ([]*domain.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND deleted_at IS NULL
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

func (r *messageRepository) Search(chatID, query string, limit, offset int) ([]*domain.Message, error) {
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

func (r *messageRepository) Update(msg *domain.Message) error {
	_, err := r.db.Exec(
		`UPDATE messages SET content=?, updated_at=? WHERE id=? AND deleted_at IS NULL`,
		msg.Content, time.Now().Format(time.RFC3339), msg.ID,
	)
	return err
}

func (r *messageRepository) SoftDelete(id string) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE messages SET deleted_at=?, updated_at=? WHERE id=?`,
		now.Format(time.RFC3339), now.Format(time.RFC3339), id,
	)
	return err
}

func (r *messageRepository) GetLastMessagesByChatIDs(chatIDs []string) (map[string]*domain.Message, error) {
	if len(chatIDs) == 0 {
		return make(map[string]*domain.Message), nil
	}
	placeholders := make([]string, len(chatIDs))
	args := make([]interface{}, len(chatIDs))
	for i, id := range chatIDs {
		placeholders[i] = "?"
		args[i] = id
	}
	rows, err := r.db.Query(
		`SELECT m.id, m.chat_id, m.sender_id, m.content, m.type, m.reply_to_id, m.forward_from, m.file_name, m.file_size, m.file_path, m.pinned, m.created_at, m.updated_at, m.deleted_at
		FROM messages m
		INNER JOIN (
			SELECT chat_id, MAX(created_at) AS max_created
			FROM messages
			WHERE chat_id IN (`+strings.Join(placeholders, ",")+`) AND deleted_at IS NULL
			GROUP BY chat_id
		) latest ON latest.chat_id = m.chat_id AND latest.max_created = m.created_at`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*domain.Message)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		result[msg.ChatID] = msg
	}
	return result, nil
}

func (r *messageRepository) GetLastMessage(chatID string) (*domain.Message, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC LIMIT 1`, chatID,
	)
	return scanMessage(row)
}

func (r *messageRepository) TogglePin(msgID string, pinned bool) error {
	_, err := r.db.Exec(
		`UPDATE messages SET pinned=?, updated_at=? WHERE id=?`,
		boolToInt(pinned), time.Now().Format(time.RFC3339), msgID,
	)
	return err
}

func (r *messageRepository) GetPinned(chatID string) ([]*domain.Message, error) {
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

func (r *messageRepository) AddReaction(msgID, userID, emoji string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO reactions (message_id, user_id, emoji, created_at) VALUES (?, ?, ?, ?)`,
		msgID, userID, emoji, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *messageRepository) RemoveReaction(msgID, userID, emoji string) error {
	_, err := r.db.Exec(
		`DELETE FROM reactions WHERE message_id = ? AND user_id = ? AND emoji = ?`,
		msgID, userID, emoji,
	)
	return err
}

func (r *messageRepository) GetReactions(msgID string) ([]*domain.Reaction, error) {
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
		reaction.CreatedAt = parseTime(createdAt)
		reactions = append(reactions, &reaction)
	}
	return reactions, nil
}

func (r *messageRepository) AddReadReceipt(msgID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO read_receipts (message_id, user_id, read_at) VALUES (?, ?, ?)`,
		msgID, userID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *messageRepository) GetReadReceipts(msgID string) ([]*domain.ReadReceipt, error) {
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
		receipt.ReadAt = parseTime(readAt)
		receipts = append(receipts, &receipt)
	}
	return receipts, nil
}

func (r *messageRepository) FindByIDs(ids []string) (map[string]*domain.Message, error) {
	if len(ids) == 0 {
		return make(map[string]*domain.Message), nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE id IN (`+strings.Join(placeholders, ",")+`)`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*domain.Message)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		result[msg.ID] = msg
	}
	return result, nil
}

func (r *messageRepository) GetReactionsByMessageIDs(ids []string) (map[string][]*domain.Reaction, error) {
	if len(ids) == 0 {
		return make(map[string][]*domain.Reaction), nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	rows, err := r.db.Query(
		`SELECT message_id, user_id, emoji, created_at FROM reactions WHERE message_id IN (`+strings.Join(placeholders, ",")+`)`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*domain.Reaction)
	for rows.Next() {
		var (
			reaction  domain.Reaction
			createdAt string
		)
		if err := rows.Scan(&reaction.MessageID, &reaction.UserID, &reaction.Emoji, &createdAt); err != nil {
			return nil, err
		}
		reaction.CreatedAt = parseTime(createdAt)
		result[reaction.MessageID] = append(result[reaction.MessageID], &reaction)
	}
	return result, nil
}

func (r *messageRepository) SearchByUser(userID, query string, limit, offset int) ([]*domain.Message, error) {
	rows, err := r.db.Query(
		`SELECT m.id, m.chat_id, m.sender_id, m.content, m.type, m.reply_to_id, m.forward_from, m.file_name, m.file_size, m.file_path, m.pinned, m.created_at, m.updated_at, m.deleted_at
		FROM messages m
		INNER JOIN chat_participants cp ON cp.chat_id = m.chat_id AND cp.user_id = ?
		WHERE m.content LIKE ? AND m.deleted_at IS NULL
		AND m.id NOT IN (SELECT dfm.message_id FROM deleted_messages dfm WHERE dfm.user_id = ?)
		ORDER BY m.created_at DESC LIMIT ? OFFSET ?`,
		userID, "%"+query+"%", userID, limit, offset,
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

func (r *messageRepository) StarMessage(userID, messageID, chatID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO starred_messages (user_id, message_id, chat_id, created_at) VALUES (?, ?, ?, ?)`,
		userID, messageID, chatID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *messageRepository) UnstarMessage(userID, messageID string) error {
	_, err := r.db.Exec(
		`DELETE FROM starred_messages WHERE user_id = ? AND message_id = ?`,
		userID, messageID,
	)
	return err
}

func (r *messageRepository) GetStarredMessages(userID string) ([]*domain.StarredMessage, error) {
	rows, err := r.db.Query(
		`SELECT user_id, message_id, chat_id, created_at FROM starred_messages WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*domain.StarredMessage, 0)
	for rows.Next() {
		var (
			sm        domain.StarredMessage
			createdAt string
		)
		if err := rows.Scan(&sm.UserID, &sm.MessageID, &sm.ChatID, &createdAt); err != nil {
			return nil, err
		}
		sm.CreatedAt = parseTime(createdAt)
		messages = append(messages, &sm)
	}
	return messages, nil
}

func (r *messageRepository) DeleteMessageForMe(userID, messageID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO deleted_messages (user_id, message_id, deleted_at) VALUES (?, ?, ?)`,
		userID, messageID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *messageRepository) FindDeletedForMe(userID string, messageIDs []string) (map[string]bool, error) {
	if len(messageIDs) == 0 {
		return make(map[string]bool), nil
	}
	placeholders := make([]string, len(messageIDs))
	args := make([]interface{}, len(messageIDs)+1)
	args[0] = userID
	for i, id := range messageIDs {
		placeholders[i] = "?"
		args[i+1] = id
	}
	rows, err := r.db.Query(
		`SELECT message_id FROM deleted_messages WHERE user_id = ? AND message_id IN (`+strings.Join(placeholders, ",")+`)`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]bool)
	for rows.Next() {
		var mid string
		if err := rows.Scan(&mid); err != nil {
			return nil, err
		}
		result[mid] = true
	}
	return result, nil
}

func (r *messageRepository) GetReadReceiptsByMessageIDs(ids []string) (map[string][]*domain.ReadReceipt, error) {
	if len(ids) == 0 {
		return make(map[string][]*domain.ReadReceipt), nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	rows, err := r.db.Query(
		`SELECT message_id, user_id, read_at FROM read_receipts WHERE message_id IN (`+strings.Join(placeholders, ",")+`)`, args...,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string][]*domain.ReadReceipt)
	for rows.Next() {
		var (
			receipt domain.ReadReceipt
			readAt  string
		)
		if err := rows.Scan(&receipt.MessageID, &receipt.UserID, &readAt); err != nil {
			return nil, err
		}
		receipt.ReadAt = parseTime(readAt)
		result[receipt.MessageID] = append(result[receipt.MessageID], &receipt)
	}
	return result, nil
}

func (r *messageRepository) SaveMention(messageID, userID, username string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO mentions (message_id, user_id, username) VALUES (?, ?, ?)`,
		messageID, userID, username,
	)
	return err
}

func (r *messageRepository) GetMentionsByMessageID(messageID string) ([]*domain.Mention, error) {
	rows, err := r.db.Query(
		`SELECT message_id, user_id, username FROM mentions WHERE message_id = ?`, messageID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentions := make([]*domain.Mention, 0)
	for rows.Next() {
		var m domain.Mention
		if err := rows.Scan(&m.MessageID, &m.UserID, &m.Username); err != nil {
			return nil, err
		}
		mentions = append(mentions, &m)
	}
	return mentions, nil
}

func (r *messageRepository) FindMediaByChatID(chatID string, mediaType string, limit, offset int) ([]*domain.Message, error) {
	var rows *sql.Rows
	var err error
	if mediaType != "" {
		rows, err = r.db.Query(
			`SELECT m.id, m.chat_id, m.sender_id, m.content, m.type, m.reply_to_id, m.forward_from, m.file_name, m.file_size, m.file_path, m.pinned, m.created_at, m.updated_at, m.deleted_at
			FROM messages m
			WHERE m.chat_id = ? AND m.type = ? AND m.deleted_at IS NULL
			ORDER BY m.created_at DESC LIMIT ? OFFSET ?`,
			chatID, mediaType, limit, offset,
		)
	} else {
		rows, err = r.db.Query(
			`SELECT m.id, m.chat_id, m.sender_id, m.content, m.type, m.reply_to_id, m.forward_from, m.file_name, m.file_size, m.file_path, m.pinned, m.created_at, m.updated_at, m.deleted_at
			FROM messages m
			WHERE m.chat_id = ? AND m.type IN ('image','file','video','audio','gif') AND m.deleted_at IS NULL
			ORDER BY m.created_at DESC LIMIT ? OFFSET ?`,
			chatID, limit, offset,
		)
	}
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

type messageScanner interface {
	Scan(dest ...interface{}) error
}

func scanMessage(row messageScanner) (*domain.Message, error) {
	var (
		msg         domain.Message
		replyToID   sql.NullString
		forwardFrom sql.NullString
		pinnedInt   int
		createdAt   string
		updatedAt   string
		deletedAt   sql.NullString
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
		t := parseTime(deletedAt.String)
		msg.DeletedAt = &t
	}
	msg.CreatedAt = parseTime(createdAt)
	msg.UpdatedAt = parseTime(updatedAt)
	return &msg, nil
}
