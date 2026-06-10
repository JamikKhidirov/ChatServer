package repository

import (
	"database/sql"
	"strings"
	"time"

	"ChatServerGolang/internal/domain"
)

type userRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) UserRepository {
	return &userRepository{db: db}
}

func (r *userRepository) Create(user *domain.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth, user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.PasswordHash, user.DisplayName,
		user.AvatarURL, user.Bio, user.Phone, user.Gender, user.DateOfBirth,
		user.Status, user.PushToken, user.PushProvider,
		boolToInt(user.Online), boolToInt(user.Deleted),
		user.LastSeen.Format(time.RFC3339),
		user.CreatedAt.Format(time.RFC3339), user.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *userRepository) FindByID(id string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth,
			user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at
		FROM users WHERE id = ? AND deleted = 0`, id,
	)
	return scanUser(row)
}

func (r *userRepository) FindByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth,
			user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at
		FROM users WHERE email = ? AND deleted = 0`, email,
	)
	return scanUser(row)
}

func (r *userRepository) FindByUsername(username string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth,
			user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at
		FROM users WHERE username = ? AND deleted = 0`, username,
	)
	return scanUser(row)
}

func (r *userRepository) Search(query string, limit, offset int) ([]*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth,
			user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at
		FROM users WHERE (username LIKE ? OR display_name LIKE ?) AND deleted = 0
		ORDER BY username ASC LIMIT ? OFFSET ?`,
		"%"+query+"%", "%"+query+"%", limit, offset,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) Update(user *domain.User) error {
	_, err := r.db.Exec(
		`UPDATE users SET display_name=?, avatar_url=?, bio=?, phone=?, gender=?, date_of_birth=?,
			user_status=?, push_token=?, push_provider=?, online=?, last_seen=?, updated_at=?
		WHERE id=? AND deleted = 0`,
		user.DisplayName, user.AvatarURL, user.Bio,
		user.Phone, user.Gender, user.DateOfBirth,
		user.Status, user.PushToken, user.PushProvider,
		boolToInt(user.Online), user.LastSeen.Format(time.RFC3339),
		user.UpdatedAt.Format(time.RFC3339), user.ID,
	)
	return err
}

func (r *userRepository) UpdatePushToken(userID, token, provider string) error {
	_, err := r.db.Exec(
		`UPDATE users SET push_token=?, push_provider=?, updated_at=? WHERE id=? AND deleted = 0`,
		token, provider, time.Now().Format(time.RFC3339), userID,
	)
	return err
}

func (r *userRepository) SetOnline(userID string, online bool) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE users SET online=?, last_seen=?, updated_at=? WHERE id=? AND deleted = 0`,
		boolToInt(online), now.Format(time.RFC3339), now.Format(time.RFC3339), userID,
	)
	return err
}

func (r *userRepository) SoftDelete(userID string) error {
	_, err := r.db.Exec(
		`UPDATE users SET deleted=1, updated_at=? WHERE id=?`,
		time.Now().Format(time.RFC3339), userID,
	)
	return err
}

func (r *userRepository) UpdatePassword(userID, hash string) error {
	_, err := r.db.Exec(
		`UPDATE users SET password_hash=?, updated_at=? WHERE id=? AND deleted = 0`,
		hash, time.Now().Format(time.RFC3339), userID,
	)
	return err
}

func (r *userRepository) GetParticipantsInChat(chatID string) ([]*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.username, u.email, u.password_hash, u.display_name, u.avatar_url, u.bio, u.phone, u.gender, u.date_of_birth,
			u.user_status, u.push_token, u.push_provider, u.online, u.deleted, u.last_seen, u.created_at, u.updated_at
		FROM users u
		INNER JOIN chat_participants cp ON cp.user_id = u.id
		WHERE cp.chat_id = ? AND u.deleted = 0`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) BlockUser(userID, blockedID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO blocked_users (user_id, blocked_id, created_at) VALUES (?, ?, ?)`,
		userID, blockedID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *userRepository) UnblockUser(userID, blockedID string) error {
	_, err := r.db.Exec(
		`DELETE FROM blocked_users WHERE user_id = ? AND blocked_id = ?`,
		userID, blockedID,
	)
	return err
}

func (r *userRepository) GetBlockedUsers(userID string) ([]*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.username, u.email, u.password_hash, u.display_name, u.avatar_url, u.bio, u.phone, u.gender, u.date_of_birth,
			u.user_status, u.push_token, u.push_provider, u.online, u.deleted, u.last_seen, u.created_at, u.updated_at
		FROM users u
		INNER JOIN blocked_users b ON b.blocked_id = u.id
		WHERE b.user_id = ? AND u.deleted = 0`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	users := make([]*domain.User, 0)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		users = append(users, u)
	}
	return users, nil
}

func (r *userRepository) IsBlocked(userID, blockedID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM blocked_users WHERE user_id = ? AND blocked_id = ?`,
		userID, blockedID,
	).Scan(&count)
	return count > 0, err
}

func (r *userRepository) FindByIDs(ids []string) (map[string]*domain.User, error) {
	if len(ids) == 0 {
		return make(map[string]*domain.User), nil
	}
	placeholders := make([]string, len(ids))
	args := make([]interface{}, len(ids))
	for i, id := range ids {
		placeholders[i] = "?"
		args[i] = id
	}
	query := `SELECT id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth,
		user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at
		FROM users WHERE id IN (` + strings.Join(placeholders, ",") + `) AND deleted = 0`
	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	result := make(map[string]*domain.User)
	for rows.Next() {
		u, err := scanUser(rows)
		if err != nil {
			return nil, err
		}
		result[u.ID] = u
	}
	return result, nil
}

func (r *userRepository) FindByIDIncludeDeleted(id string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, phone, gender, date_of_birth,
			user_status, push_token, push_provider, online, deleted, last_seen, created_at, updated_at
		FROM users WHERE id = ?`, id,
	)
	return scanUser(row)
}

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(row scanner) (*domain.User, error) {
	var (
		u          domain.User
		onlineInt  int
		deletedInt int
		lastSeen   string
		createdAt  string
		updatedAt  string
	)
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.DisplayName, &u.AvatarURL, &u.Bio,
		&u.Phone, &u.Gender, &u.DateOfBirth,
		&u.Status, &u.PushToken, &u.PushProvider, &onlineInt,
		&deletedInt, &lastSeen, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Online = onlineInt == 1
	u.Deleted = deletedInt == 1
	u.LastSeen = parseTime(lastSeen)
	u.CreatedAt = parseTime(createdAt)
	u.UpdatedAt = parseTime(updatedAt)
	return &u, nil
}
