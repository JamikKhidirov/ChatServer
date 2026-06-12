package sessionrepo

import (
	"database/sql"
	"time"

	sessiondomain "ChatServerGolang/internal/domain/session"
	"ChatServerGolang/internal/repository"
)

type sessionRepository struct {
	db *sql.DB
}

func NewSessionRepository(db *sql.DB) repository.SessionRepository {
	return &sessionRepository{db: db}
}

func (r *sessionRepository) Create(session *sessiondomain.Session) error {
	_, err := r.db.Exec(
		`INSERT INTO sessions (id, user_id, device_name, ip_address, last_active, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		session.ID, session.UserID, session.DeviceName, session.IPAddress,
		session.LastActive.Format(time.RFC3339), session.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *sessionRepository) FindByID(id string) (*sessiondomain.Session, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, device_name, ip_address, last_active, created_at FROM sessions WHERE id = ?`, id,
	)
	var s sessiondomain.Session
	var lastActive, createdAt string
	if err := row.Scan(&s.ID, &s.UserID, &s.DeviceName, &s.IPAddress, &lastActive, &createdAt); err != nil {
		return nil, err
	}
	s.LastActive = repository.ParseTime(lastActive)
	s.CreatedAt = repository.ParseTime(createdAt)
	return &s, nil
}

func (r *sessionRepository) FindByUserID(userID string) ([]*sessiondomain.Session, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, device_name, ip_address, last_active, created_at
		FROM sessions WHERE user_id = ? ORDER BY last_active DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	sessions := make([]*sessiondomain.Session, 0)
	for rows.Next() {
		var s sessiondomain.Session
		var lastActive, createdAt string
		if err := rows.Scan(&s.ID, &s.UserID, &s.DeviceName, &s.IPAddress, &lastActive, &createdAt); err != nil {
			return nil, err
		}
		s.LastActive = repository.ParseTime(lastActive)
		s.CreatedAt = repository.ParseTime(createdAt)
		sessions = append(sessions, &s)
	}
	return sessions, nil
}

func (r *sessionRepository) UpdateLastActive(id string) error {
	_, err := r.db.Exec(
		`UPDATE sessions SET last_active = ? WHERE id = ?`,
		time.Now().Format(time.RFC3339), id,
	)
	return err
}

func (r *sessionRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE id = ?`, id)
	return err
}

func (r *sessionRepository) DeleteByUserID(userID string) error {
	_, err := r.db.Exec(`DELETE FROM sessions WHERE user_id = ?`, userID)
	return err
}



