package linkrepo

import (
	"database/sql"
	"time"

	chatdomain "ChatServerGolang/internal/domain/chat"
	"ChatServerGolang/internal/repository"
)

type inviteLinkRepo struct {
	db *sql.DB
}

func NewInviteLinkRepository(db *sql.DB) repository.InviteLinkRepository {
	return &inviteLinkRepo{db: db}
}

func (r *inviteLinkRepo) Create(link *chatdomain.InviteLink) error {
	var expiresAt *string
	if link.ExpiresAt != nil {
		s := link.ExpiresAt.Format(time.RFC3339)
		expiresAt = &s
	}
	_, err := r.db.Exec(
		`INSERT INTO invite_links (id, chat_id, creator_id, code, expires_at, usage_limit, usage_count, active, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		link.ID, link.ChatID, link.CreatorID, link.Code,
		expiresAt, link.UsageLimit, link.UsageCount, repository.BoolToInt(link.Active),
		link.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *inviteLinkRepo) FindByCode(code string) (*chatdomain.InviteLink, error) {
	var expiresAt sql.NullString
	link := &chatdomain.InviteLink{}
	err := r.db.QueryRow(
		`SELECT id, chat_id, creator_id, code, expires_at, usage_limit, usage_count, active, created_at
		FROM invite_links WHERE code = ?`, code,
	).Scan(&link.ID, &link.ChatID, &link.CreatorID, &link.Code,
		&expiresAt, &link.UsageLimit, &link.UsageCount, &link.Active, &link.CreatedAt)
	if err != nil {
		return nil, err
	}
	if expiresAt.Valid {
		t := repository.ParseTime(expiresAt.String)
		link.ExpiresAt = &t
	}
	return link, nil
}

func (r *inviteLinkRepo) FindByChatID(chatID string) ([]*chatdomain.InviteLink, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, creator_id, code, expires_at, usage_limit, usage_count, active, created_at
		FROM invite_links WHERE chat_id = ? ORDER BY created_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var links []*chatdomain.InviteLink
	for rows.Next() {
		var expiresAt sql.NullString
		link := &chatdomain.InviteLink{}
		if err := rows.Scan(&link.ID, &link.ChatID, &link.CreatorID, &link.Code,
			&expiresAt, &link.UsageLimit, &link.UsageCount, &link.Active, &link.CreatedAt); err != nil {
			return nil, err
		}
		if expiresAt.Valid {
			t := repository.ParseTime(expiresAt.String)
			link.ExpiresAt = &t
		}
		links = append(links, link)
	}
	return links, nil
}

func (r *inviteLinkRepo) IncrementUsage(id string) error {
	_, err := r.db.Exec(
		`UPDATE invite_links SET usage_count = usage_count + 1 WHERE id = ?`, id,
	)
	return err
}

func (r *inviteLinkRepo) Deactivate(id string) error {
	_, err := r.db.Exec(
		`UPDATE invite_links SET active = 0 WHERE id = ?`, id,
	)
	return err
}

func (r *inviteLinkRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM invite_links WHERE id = ?`, id)
	return err
}



