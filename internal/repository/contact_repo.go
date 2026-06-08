package repository

import (
	"database/sql"

	"ChatServerGolang/internal/domain"
)

type contactRepo struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) ContactRepository {
	return &contactRepo{db: db}
}

func (r *contactRepo) SyncContacts(userID string, contacts []domain.ContactInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM contacts WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO contacts (user_id, phone, name) VALUES (?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range contacts {
		if _, err := stmt.Exec(userID, c.Phone, c.Name); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *contactRepo) GetContacts(userID string) ([]*domain.ContactResponse, error) {
	rows, err := r.db.Query("SELECT phone, name FROM contacts WHERE user_id = ? ORDER BY name", userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*domain.ContactResponse
	for rows.Next() {
		c := &domain.ContactResponse{}
		if err := rows.Scan(&c.Phone, &c.Name); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *contactRepo) SearchByPhone(userID, phoneQuery string) ([]*domain.ContactResponse, error) {
	rows, err := r.db.Query("SELECT phone, name FROM contacts WHERE user_id = ? AND phone LIKE ? ORDER BY name", userID, "%"+phoneQuery+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*domain.ContactResponse
	for rows.Next() {
		c := &domain.ContactResponse{}
		if err := rows.Scan(&c.Phone, &c.Name); err != nil {
			return nil, err
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *contactRepo) FindRegisteredByPhone(phones []string) ([]*domain.UserResponse, error) {
	if len(phones) == 0 {
		return nil, nil
	}
	query := "SELECT id, username, display_name, avatar_url, bio, phone, user_status, online, last_seen FROM users WHERE phone IN ("
	args := make([]interface{}, len(phones))
	for i, p := range phones {
		if i > 0 {
			query += ","
		}
		query += "?"
		args[i] = p
	}
	query += ") AND deleted = 0"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var users []*domain.UserResponse
	for rows.Next() {
		u := &domain.UserResponse{}
		var onlineInt int
		var lastSeenStr string
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Phone, &u.Status, &onlineInt, &lastSeenStr); err != nil {
			return nil, err
		}
		u.Online = onlineInt == 1
		u.LastSeen = parseTime(lastSeenStr)
		users = append(users, u)
	}
	return users, nil
}
