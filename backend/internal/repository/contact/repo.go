package contactrepo

import (
	"database/sql"
	"time"

	contactdomain "ChatServerGolang/backend/internal/domain/contact"
	userdomain "ChatServerGolang/backend/internal/domain/user"
	"ChatServerGolang/backend/internal/repository"
)

type contactRepo struct {
	db *sql.DB
}

func NewContactRepository(db *sql.DB) repository.ContactRepository {
	return &contactRepo{db: db}
}

func (r *contactRepo) SyncContacts(userID string, contacts []contactdomain.ContactInput) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	_, err = tx.Exec("DELETE FROM contacts WHERE user_id = ?", userID)
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO contacts (user_id, phone, name, user_id_ref) VALUES (?, ?, ?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, c := range contacts {
		// Look up if this phone belongs to a registered user
		var userIDRef string
		r.db.QueryRow("SELECT id FROM users WHERE phone = ? AND deleted = 0", c.Phone).Scan(&userIDRef)
		if _, err := stmt.Exec(userID, c.Phone, c.Name, userIDRef); err != nil {
			return err
		}
	}

	return tx.Commit()
}

func (r *contactRepo) GetContacts(userID string) ([]*contactdomain.ContactResponse, error) {
	rows, err := r.db.Query(`
		SELECT c.phone, c.name, COALESCE(c.photo_url, ''), COALESCE(c.user_id_ref, ''), COALESCE(u.avatar_url, '')
		FROM contacts c
		LEFT JOIN users u ON u.id = c.user_id_ref AND u.deleted = 0
		WHERE c.user_id = ?
		ORDER BY c.name`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*contactdomain.ContactResponse
	for rows.Next() {
		c := &contactdomain.ContactResponse{}
		var avatarURL string
		if err := rows.Scan(&c.Phone, &c.Name, &c.PhotoURL, &c.UserIDRef, &avatarURL); err != nil {
			return nil, err
		}
		if c.PhotoURL == "" && avatarURL != "" {
			c.PhotoURL = avatarURL
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *contactRepo) SearchByPhone(userID, phoneQuery string) ([]*contactdomain.ContactResponse, error) {
	rows, err := r.db.Query(`
		SELECT c.phone, c.name, COALESCE(c.photo_url, ''), COALESCE(c.user_id_ref, ''), COALESCE(u.avatar_url, '')
		FROM contacts c
		LEFT JOIN users u ON u.id = c.user_id_ref AND u.deleted = 0
		WHERE c.user_id = ? AND c.phone LIKE ?
		ORDER BY c.name`, userID, "%"+phoneQuery+"%")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var contacts []*contactdomain.ContactResponse
	for rows.Next() {
		c := &contactdomain.ContactResponse{}
		var avatarURL string
		if err := rows.Scan(&c.Phone, &c.Name, &c.PhotoURL, &c.UserIDRef, &avatarURL); err != nil {
			return nil, err
		}
		if c.PhotoURL == "" && avatarURL != "" {
			c.PhotoURL = avatarURL
		}
		contacts = append(contacts, c)
	}
	return contacts, nil
}

func (r *contactRepo) FindRegisteredByPhone(phones []string) ([]*userdomain.UserResponse, error) {
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

	var users []*userdomain.UserResponse
	for rows.Next() {
		u := &userdomain.UserResponse{}
		var onlineInt int
		var lastSeenStr string
		if err := rows.Scan(&u.ID, &u.Username, &u.DisplayName, &u.AvatarURL, &u.Bio, &u.Phone, &u.Status, &onlineInt, &lastSeenStr); err != nil {
			return nil, err
		}
		u.Online = onlineInt == 1
		u.LastSeen = repository.ParseTime(lastSeenStr)
		users = append(users, u)
	}
	return users, nil
}

func (r *contactRepo) UpdateContactPhoto(userID, phone, photoURL string) error {
	_, err := r.db.Exec(
		`UPDATE contacts SET photo_url = ?, updated_at = ? WHERE user_id = ? AND phone = ?`,
		photoURL, time.Now().Format(time.RFC3339), userID, phone,
	)
	return err
}



