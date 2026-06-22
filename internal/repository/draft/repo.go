package draftrepo

import (
	"database/sql"
	"time"

	draftdomain "ChatServerGolang/internal/domain/draft"
	"ChatServerGolang/internal/repository"
)

type draftRepository struct {
	db *sql.DB
}

func NewDraftRepository(db *sql.DB) repository.DraftRepository {
	return &draftRepository{db: db}
}

func (r *draftRepository) Save(draft *draftdomain.Draft) error {
	if draft.ID == "" {
		return r.create(draft)
	}
	_, err := r.db.Exec(
		`INSERT OR REPLACE INTO drafts (id, user_id, chat_id, content, reply_to_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		draft.ID, draft.UserID, draft.ChatID, draft.Content, draft.ReplyToID,
		draft.CreatedAt.Format(time.RFC3339), draft.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *draftRepository) create(draft *draftdomain.Draft) error {
	_, err := r.db.Exec(
		`INSERT INTO drafts (id, user_id, chat_id, content, reply_to_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		draft.ID, draft.UserID, draft.ChatID, draft.Content, draft.ReplyToID,
		draft.CreatedAt.Format(time.RFC3339), draft.UpdatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *draftRepository) FindByUserAndChat(userID, chatID string) (*draftdomain.Draft, error) {
	row := r.db.QueryRow(
		`SELECT id, user_id, chat_id, content, reply_to_id, created_at, updated_at
		FROM drafts WHERE user_id = ? AND chat_id = ?`, userID, chatID,
	)
	var d draftdomain.Draft
	var replyToID sql.NullString
	var createdAt, updatedAt string
	if err := row.Scan(&d.ID, &d.UserID, &d.ChatID, &d.Content, &replyToID, &createdAt, &updatedAt); err != nil {
		return nil, err
	}
	if replyToID.Valid {
		d.ReplyToID = &replyToID.String
	}
	d.CreatedAt = repository.ParseTime(createdAt)
	d.UpdatedAt = repository.ParseTime(updatedAt)
	return &d, nil
}

func (r *draftRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM drafts WHERE id = ?`, id)
	return err
}

func (r *draftRepository) DeleteByUserAndChat(userID, chatID string) error {
	_, err := r.db.Exec(`DELETE FROM drafts WHERE user_id = ? AND chat_id = ?`, userID, chatID)
	return err
}



