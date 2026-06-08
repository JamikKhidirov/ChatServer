package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type chatRepository struct {
	db *sql.DB
}

func NewChatRepository(db *sql.DB) ChatRepository {
	return &chatRepository{db: db}
}

func (r *chatRepository) Create(chat *domain.Chat) error {
	_, err := r.db.Exec(
		`INSERT INTO chats (id, name, description, avatar_url, type, created_by, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		chat.ID, chat.Name, chat.Description, chat.AvatarURL, chat.Type, chat.CreatedBy,
		chat.CreatedAt.Format(time.RFC3339), chat.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *chatRepository) FindByID(id string) (*domain.Chat, error) {
	row := r.db.QueryRow(
		`SELECT id, name, COALESCE(description,''), avatar_url, type, created_by, created_at, updated_at
		FROM chats WHERE id = ?`, id,
	)
	return scanChat(row)
}

func (r *chatRepository) FindByUserID(userID string) ([]*domain.Chat, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.name, COALESCE(c.description,''), c.avatar_url, c.type, c.created_by, c.created_at, c.updated_at
		FROM chats c
		INNER JOIN chat_participants cp ON cp.chat_id = c.id
		WHERE cp.user_id = ?
		ORDER BY c.updated_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*domain.Chat, 0)
	for rows.Next() {
		c, err := scanChat(rows)
		if err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	return chats, nil
}

func (r *chatRepository) Update(chat *domain.Chat) error {
	_, err := r.db.Exec(
		`UPDATE chats SET name=?, description=?, avatar_url=?, updated_at=? WHERE id=?`,
		chat.Name, chat.Description, chat.AvatarURL, time.Now().Format(time.RFC3339), chat.ID,
	)
	return err
}

func (r *chatRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM chat_participants WHERE chat_id = ?`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE messages SET deleted_at=? WHERE chat_id=?`, time.Now().Format(time.RFC3339), id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM chats WHERE id = ?`, id)
	return err
}

func (r *chatRepository) AddParticipant(chatID, userID, role string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO chat_participants (chat_id, user_id, role, joined_at, last_read_at)
		VALUES (?, ?, ?, ?, ?)`,
		chatID, userID, role, time.Now().Format(time.RFC3339), time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *chatRepository) RemoveParticipant(chatID, userID string) error {
	_, err := r.db.Exec(
		`DELETE FROM chat_participants WHERE chat_id = ? AND user_id = ?`,
		chatID, userID,
	)
	return err
}

func (r *chatRepository) GetParticipants(chatID string) ([]*domain.ChatParticipant, error) {
	rows, err := r.db.Query(
		`SELECT chat_id, user_id, role, joined_at, last_read_at
		FROM chat_participants WHERE chat_id = ?`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	participants := make([]*domain.ChatParticipant, 0)
	for rows.Next() {
		var (
			p        domain.ChatParticipant
			joinedAt string
			lastRead string
		)
		if err := rows.Scan(&p.ChatID, &p.UserID, &p.Role, &joinedAt, &lastRead); err != nil {
			return nil, err
		}
		p.JoinedAt = parseTime(joinedAt)
		p.LastReadAt = parseTime(lastRead)
		participants = append(participants, &p)
	}
	return participants, nil
}

func (r *chatRepository) IsParticipant(chatID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM chat_participants WHERE chat_id = ? AND user_id = ?`,
		chatID, userID,
	).Scan(&count)
	return count > 0, err
}

func (r *chatRepository) GetPrivateChat(user1ID, user2ID string) (*domain.Chat, error) {
	row := r.db.QueryRow(
		`SELECT c.id, c.name, COALESCE(c.description,''), c.avatar_url, c.type, c.created_by, c.created_at, c.updated_at
		FROM chats c
		INNER JOIN chat_participants cp1 ON cp1.chat_id = c.id AND cp1.user_id = ?
		INNER JOIN chat_participants cp2 ON cp2.chat_id = c.id AND cp2.user_id = ?
		WHERE c.type = 'private'
		LIMIT 1`, user1ID, user2ID,
	)
	return scanChat(row)
}

func (r *chatRepository) SetRole(chatID, userID, role string) error {
	_, err := r.db.Exec(
		`UPDATE chat_participants SET role = ? WHERE chat_id = ? AND user_id = ?`,
		role, chatID, userID,
	)
	return err
}

func (r *chatRepository) UpdateLastRead(chatID, userID string) error {
	_, err := r.db.Exec(
		`UPDATE chat_participants SET last_read_at = ? WHERE chat_id = ? AND user_id = ?`,
		time.Now().Format(time.RFC3339), chatID, userID,
	)
	return err
}

func (r *chatRepository) GetUnreadCount(chatID, userID string) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM messages m
		WHERE m.chat_id = ? AND m.sender_id != ?
		AND (m.deleted_at IS NULL)
		AND m.created_at > (
			SELECT COALESCE(last_read_at, '1970-01-01') FROM chat_participants
			WHERE chat_id = ? AND user_id = ?
		)`,
		chatID, userID, chatID, userID,
	).Scan(&count)
	return count, err
}

func (r *chatRepository) SetNotificationMuted(userID, chatID string, muted bool) error {
	_, err := r.db.Exec(
		`INSERT OR REPLACE INTO notification_settings (user_id, chat_id, muted) VALUES (?, ?, ?)`,
		userID, chatID, boolToInt(muted),
	)
	return err
}

func (r *chatRepository) IsNotificationMuted(userID, chatID string) (bool, error) {
	var mutedInt int
	err := r.db.QueryRow(
		`SELECT COALESCE(muted, 0) FROM notification_settings WHERE user_id = ? AND chat_id = ?`,
		userID, chatID,
	).Scan(&mutedInt)
	if err == sql.ErrNoRows {
		return false, nil
	}
	return mutedInt == 1, err
}

func (r *chatRepository) HideChat(userID, chatID string) error {
	_, err := r.db.Exec(
		`INSERT OR REPLACE INTO hidden_chats (user_id, chat_id, hidden_at) VALUES (?, ?, ?)`,
		userID, chatID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *chatRepository) IsHidden(userID, chatID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM hidden_chats WHERE user_id = ? AND chat_id = ?`,
		userID, chatID,
	).Scan(&count)
	return count > 0, err
}

func (r *chatRepository) FindByUserIDExcludeHidden(userID string) ([]*domain.Chat, error) {
	rows, err := r.db.Query(
		`SELECT c.id, c.name, COALESCE(c.description,''), c.avatar_url, c.type, c.created_by, c.created_at, c.updated_at
		FROM chats c
		INNER JOIN chat_participants cp ON cp.chat_id = c.id
		LEFT JOIN hidden_chats hc ON hc.chat_id = c.id AND hc.user_id = ?
		WHERE cp.user_id = ? AND hc.chat_id IS NULL
		ORDER BY c.updated_at DESC`, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	chats := make([]*domain.Chat, 0)
	for rows.Next() {
		c, err := scanChat(rows)
		if err != nil {
			return nil, err
		}
		chats = append(chats, c)
	}
	return chats, nil
}

func scanChat(row scanner) (*domain.Chat, error) {
	var (
		c         domain.Chat
		createdAt string
		updatedAt string
	)
	err := row.Scan(&c.ID, &c.Name, &c.Description, &c.AvatarURL, &c.Type, &c.CreatedBy, &createdAt, &updatedAt)
	if err != nil {
		return nil, err
	}
	c.CreatedAt = parseTime(createdAt)
	c.UpdatedAt = parseTime(updatedAt)
	return &c, nil
}
