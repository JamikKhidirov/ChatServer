package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (r *UserRepository) Create(user *domain.User) error {
	_, err := r.db.Exec(
		`INSERT INTO users (id, username, email, password_hash, display_name, avatar_url, bio, user_status, push_token, push_provider, online, last_seen, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		user.ID, user.Username, user.Email, user.PasswordHash, user.DisplayName,
		user.AvatarURL, user.Bio, user.Status, user.PushToken, user.PushProvider,
		boolToInt(user.Online), user.LastSeen.Format(time.RFC3339),
		user.CreatedAt.Format(time.RFC3339), user.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *UserRepository) FindByID(id string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, user_status,
			push_token, push_provider, online, last_seen, created_at, updated_at
		FROM users WHERE id = ?`, id,
	)
	return scanUser(row)
}

func (r *UserRepository) FindByEmail(email string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, user_status,
			push_token, push_provider, online, last_seen, created_at, updated_at
		FROM users WHERE email = ?`, email,
	)
	return scanUser(row)
}

func (r *UserRepository) FindByUsername(username string) (*domain.User, error) {
	row := r.db.QueryRow(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, user_status,
			push_token, push_provider, online, last_seen, created_at, updated_at
		FROM users WHERE username = ?`, username,
	)
	return scanUser(row)
}

func (r *UserRepository) Search(query string, limit, offset int) ([]*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT id, username, email, password_hash, display_name, avatar_url, bio, user_status,
			push_token, push_provider, online, last_seen, created_at, updated_at
		FROM users WHERE username LIKE ? OR display_name LIKE ?
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

func (r *UserRepository) Update(user *domain.User) error {
	_, err := r.db.Exec(
		`UPDATE users SET display_name=?, avatar_url=?, bio=?, user_status=?, push_token=?, push_provider=?,
			online=?, last_seen=?, updated_at=?
		WHERE id=?`,
		user.DisplayName, user.AvatarURL, user.Bio, user.Status, user.PushToken, user.PushProvider,
		boolToInt(user.Online), user.LastSeen.Format(time.RFC3339),
		user.UpdatedAt.Format(time.RFC3339), user.ID,
	)
	return err
}

func (r *UserRepository) UpdatePushToken(userID, token, provider string) error {
	_, err := r.db.Exec(
		`UPDATE users SET push_token=?, push_provider=?, updated_at=? WHERE id=?`,
		token, provider, time.Now().Format(time.RFC3339), userID,
	)
	return err
}

func (r *UserRepository) SetOnline(userID string, online bool) error {
	now := time.Now()
	_, err := r.db.Exec(
		`UPDATE users SET online=?, last_seen=?, updated_at=? WHERE id=?`,
		boolToInt(online), now.Format(time.RFC3339), now.Format(time.RFC3339), userID,
	)
	return err
}

func (r *UserRepository) GetParticipantsInChat(chatID string) ([]*domain.User, error) {
	rows, err := r.db.Query(
		`SELECT u.id, u.username, u.email, u.password_hash, u.display_name, u.avatar_url, u.bio, u.user_status,
			u.push_token, u.push_provider, u.online, u.last_seen, u.created_at, u.updated_at
		FROM users u
		INNER JOIN chat_participants cp ON cp.user_id = u.id
		WHERE cp.chat_id = ?`, chatID,
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

type scanner interface {
	Scan(dest ...interface{}) error
}

func scanUser(row scanner) (*domain.User, error) {
	var (
		u         domain.User
		onlineInt int
		lastSeen  string
		createdAt string
		updatedAt string
	)
	err := row.Scan(
		&u.ID, &u.Username, &u.Email, &u.PasswordHash,
		&u.DisplayName, &u.AvatarURL, &u.Bio, &u.Status,
		&u.PushToken, &u.PushProvider, &onlineInt,
		&lastSeen, &createdAt, &updatedAt,
	)
	if err != nil {
		return nil, err
	}
	u.Online = onlineInt == 1
	u.LastSeen, _ = time.Parse(time.RFC3339, lastSeen)
	u.CreatedAt, _ = time.Parse(time.RFC3339, createdAt)
	u.UpdatedAt, _ = time.Parse(time.RFC3339, updatedAt)
	return &u, nil
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}
