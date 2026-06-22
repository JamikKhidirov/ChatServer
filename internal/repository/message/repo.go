package messagerepo

import (
	"database/sql"
	"strings"
	"time"

	chatdomain "ChatServerGolang/internal/domain/chat"
	messagedomain "ChatServerGolang/internal/domain/message"
	"ChatServerGolang/internal/repository"
)

type messageRepository struct {
	db *sql.DB
}

func NewMessageRepository(db *sql.DB) repository.MessageRepository {
	return &messageRepository{db: db}
}

func (r *messageRepository) Create(msg *messagedomain.Message) error {
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
		repository.BoolToInt(msg.Pinned),
		msg.CreatedAt.Format(time.RFC3339), msg.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *messageRepository) FindByID(id string) (*messagedomain.Message, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE id = ? AND deleted_at IS NULL`, id,
	)
	return scanMessage(row)
}

func (r *messageRepository) FindByChatID(chatID string, limit, offset int) ([]*messagedomain.Message, error) {
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

	messages := make([]*messagedomain.Message, 0)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *messageRepository) CountByChatID(chatID string) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM messages WHERE chat_id = ? AND deleted_at IS NULL`, chatID,
	).Scan(&count)
	return count, err
}

func (r *messageRepository) CountChatMedia(chatID, mediaType string) (int, error) {
	var count int
	var err error
	if mediaType != "" {
		err = r.db.QueryRow(
			`SELECT COUNT(*) FROM messages WHERE chat_id = ? AND type = ? AND deleted_at IS NULL`,
			chatID, mediaType,
		).Scan(&count)
	} else {
		err = r.db.QueryRow(
			`SELECT COUNT(*) FROM messages WHERE chat_id = ? AND type IN ('image','file','video','audio','gif') AND deleted_at IS NULL`,
			chatID,
		).Scan(&count)
	}
	return count, err
}

func (r *messageRepository) Search(chatID, query string, limit, offset int) ([]*messagedomain.Message, error) {
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

	messages := make([]*messagedomain.Message, 0)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

func (r *messageRepository) Update(msg *messagedomain.Message) error {
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

func (r *messageRepository) GetLastMessagesByChatIDs(chatIDs []string) (map[string]*messagedomain.Message, error) {
	if len(chatIDs) == 0 {
		return make(map[string]*messagedomain.Message), nil
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

	result := make(map[string]*messagedomain.Message)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		result[msg.ChatID] = msg
	}
	return result, nil
}

func (r *messageRepository) GetLastMessage(chatID string) (*messagedomain.Message, error) {
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
		repository.BoolToInt(pinned), time.Now().Format(time.RFC3339), msgID,
	)
	return err
}

func (r *messageRepository) GetPinned(chatID string) ([]*messagedomain.Message, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, sender_id, content, type, reply_to_id, forward_from, file_name, file_size, file_path, pinned, created_at, updated_at, deleted_at
		FROM messages WHERE chat_id = ? AND pinned = 1 AND deleted_at IS NULL
		ORDER BY created_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*messagedomain.Message, 0)
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

func (r *messageRepository) GetReactions(msgID string) ([]*messagedomain.Reaction, error) {
	rows, err := r.db.Query(
		`SELECT message_id, user_id, emoji, created_at FROM reactions WHERE message_id = ?`,
		msgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	reactions := make([]*messagedomain.Reaction, 0)
	for rows.Next() {
		var (
			reaction  messagedomain.Reaction
			createdAt string
		)
		if err := rows.Scan(&reaction.MessageID, &reaction.UserID, &reaction.Emoji, &createdAt); err != nil {
			return nil, err
		}
		reaction.CreatedAt = repository.ParseTime(createdAt)
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

func (r *messageRepository) GetReadReceipts(msgID string) ([]*messagedomain.ReadReceipt, error) {
	rows, err := r.db.Query(
		`SELECT message_id, user_id, read_at FROM read_receipts WHERE message_id = ?`,
		msgID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	receipts := make([]*messagedomain.ReadReceipt, 0)
	for rows.Next() {
		var (
			receipt messagedomain.ReadReceipt
			readAt  string
		)
		if err := rows.Scan(&receipt.MessageID, &receipt.UserID, &readAt); err != nil {
			return nil, err
		}
		receipt.ReadAt = repository.ParseTime(readAt)
		receipts = append(receipts, &receipt)
	}
	return receipts, nil
}

func (r *messageRepository) FindByIDs(ids []string) (map[string]*messagedomain.Message, error) {
	if len(ids) == 0 {
		return make(map[string]*messagedomain.Message), nil
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

	result := make(map[string]*messagedomain.Message)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		result[msg.ID] = msg
	}
	return result, nil
}

func (r *messageRepository) GetReactionsByMessageIDs(ids []string) (map[string][]*messagedomain.Reaction, error) {
	if len(ids) == 0 {
		return make(map[string][]*messagedomain.Reaction), nil
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

	result := make(map[string][]*messagedomain.Reaction)
	for rows.Next() {
		var (
			reaction  messagedomain.Reaction
			createdAt string
		)
		if err := rows.Scan(&reaction.MessageID, &reaction.UserID, &reaction.Emoji, &createdAt); err != nil {
			return nil, err
		}
		reaction.CreatedAt = repository.ParseTime(createdAt)
		result[reaction.MessageID] = append(result[reaction.MessageID], &reaction)
	}
	return result, nil
}

func (r *messageRepository) SearchByUser(userID, query string, limit, offset int) ([]*messagedomain.Message, error) {
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

	messages := make([]*messagedomain.Message, 0)
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

func (r *messageRepository) GetStarredMessages(userID string) ([]*chatdomain.StarredMessage, error) {
	rows, err := r.db.Query(
		`SELECT user_id, message_id, chat_id, created_at FROM starred_messages WHERE user_id = ? ORDER BY created_at DESC`,
		userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	messages := make([]*chatdomain.StarredMessage, 0)
	for rows.Next() {
		var (
			sm        chatdomain.StarredMessage
			createdAt string
		)
		if err := rows.Scan(&sm.UserID, &sm.MessageID, &sm.ChatID, &createdAt); err != nil {
			return nil, err
		}
		sm.CreatedAt = repository.ParseTime(createdAt)
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

func (r *messageRepository) GetReadReceiptsByMessageIDs(ids []string) (map[string][]*messagedomain.ReadReceipt, error) {
	if len(ids) == 0 {
		return make(map[string][]*messagedomain.ReadReceipt), nil
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

	result := make(map[string][]*messagedomain.ReadReceipt)
	for rows.Next() {
		var (
			receipt messagedomain.ReadReceipt
			readAt  string
		)
		if err := rows.Scan(&receipt.MessageID, &receipt.UserID, &readAt); err != nil {
			return nil, err
		}
		receipt.ReadAt = repository.ParseTime(readAt)
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

func (r *messageRepository) GetMentionsByMessageID(messageID string) ([]*messagedomain.Mention, error) {
	rows, err := r.db.Query(
		`SELECT message_id, user_id, username FROM mentions WHERE message_id = ?`, messageID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	mentions := make([]*messagedomain.Mention, 0)
	for rows.Next() {
		var m messagedomain.Mention
		if err := rows.Scan(&m.MessageID, &m.UserID, &m.Username); err != nil {
			return nil, err
		}
		mentions = append(mentions, &m)
	}
	return mentions, nil
}

func (r *messageRepository) FindMediaByChatID(chatID string, mediaType string, limit, offset int) ([]*messagedomain.Message, error) {
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

	messages := make([]*messagedomain.Message, 0)
	for rows.Next() {
		msg, err := scanMessage(rows)
		if err != nil {
			return nil, err
		}
		messages = append(messages, msg)
	}
	return messages, nil
}

type messageScannerLocal interface {
	Scan(dest ...interface{}) error
}

func scanMessage(row repository.MessageScanner) (*messagedomain.Message, error) {
	var (
		msg         messagedomain.Message
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
		t := repository.ParseTime(deletedAt.String)
		msg.DeletedAt = &t
	}
	msg.CreatedAt = repository.ParseTime(createdAt)
	msg.UpdatedAt = repository.ParseTime(updatedAt)
	return &msg, nil
}

func (r *messageRepository) SetSelfDestruct(msgID, chatID string, deleteAt time.Time) error {
	_, err := r.db.Exec(
		`INSERT OR REPLACE INTO message_self_destruct (message_id, chat_id, delete_at) VALUES (?, ?, ?)`,
		msgID, chatID, deleteAt.UTC().Format(time.RFC3339),
	)
	return err
}

func (r *messageRepository) GetExpiredSelfDestruct() ([]messagedomain.MessageSelfDestruct, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	rows, err := r.db.Query(
		`SELECT message_id, chat_id, delete_at FROM message_self_destruct WHERE delete_at <= ?`,
		now,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var results []messagedomain.MessageSelfDestruct
	for rows.Next() {
		var sd messagedomain.MessageSelfDestruct
		if err := rows.Scan(&sd.MessageID, &sd.ChatID, &sd.DeleteAt); err != nil {
			return nil, err
		}
		results = append(results, sd)
	}
	return results, rows.Err()
}

func (r *messageRepository) DeleteSelfDestructByMessageID(messageID string) error {
	_, err := r.db.Exec(`DELETE FROM message_self_destruct WHERE message_id = ?`, messageID)
	return err
}


