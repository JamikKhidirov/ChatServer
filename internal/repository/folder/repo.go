package folderrepo

import (
	"database/sql"
	"time"

	chatdomain "ChatServerGolang/internal/domain/chat"
	"ChatServerGolang/internal/repository"
)

type chatFolderRepo struct {
	db *sql.DB
}

func NewChatFolderRepository(db *sql.DB) repository.ChatFolderRepository {
	return &chatFolderRepo{db: db}
}

func (r *chatFolderRepo) Create(folder *chatdomain.ChatFolder) error {
	_, err := r.db.Exec(
		`INSERT INTO chat_folders (id, user_id, name, emoji, folder_order, created_at)
		VALUES (?, ?, ?, ?, ?, ?)`,
		folder.ID, folder.UserID, folder.Name, folder.Emoji, folder.Order, folder.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *chatFolderRepo) FindByUserID(userID string) ([]*chatdomain.ChatFolder, error) {
	rows, err := r.db.Query(
		`SELECT id, user_id, name, emoji, folder_order, created_at
		FROM chat_folders WHERE user_id = ? ORDER BY folder_order ASC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var folders []*chatdomain.ChatFolder
	for rows.Next() {
		f := &chatdomain.ChatFolder{}
		if err := rows.Scan(&f.ID, &f.UserID, &f.Name, &f.Emoji, &f.Order, &f.CreatedAt); err != nil {
			return nil, err
		}
		folders = append(folders, f)
	}
	return folders, nil
}

func (r *chatFolderRepo) FindByID(id string) (*chatdomain.ChatFolder, error) {
	f := &chatdomain.ChatFolder{}
	err := r.db.QueryRow(
		`SELECT id, user_id, name, emoji, folder_order, created_at
		FROM chat_folders WHERE id = ?`, id,
	).Scan(&f.ID, &f.UserID, &f.Name, &f.Emoji, &f.Order, &f.CreatedAt)
	return f, err
}

func (r *chatFolderRepo) Update(folder *chatdomain.ChatFolder) error {
	_, err := r.db.Exec(
		`UPDATE chat_folders SET name = ?, emoji = ?, folder_order = ? WHERE id = ?`,
		folder.Name, folder.Emoji, folder.Order, folder.ID,
	)
	return err
}

func (r *chatFolderRepo) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM chat_folder_items WHERE folder_id = ?`, id)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`DELETE FROM chat_folders WHERE id = ?`, id)
	return err
}

func (r *chatFolderRepo) AddChatToFolder(folderID, chatID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO chat_folder_items (folder_id, chat_id, added_at) VALUES (?, ?, ?)`,
		folderID, chatID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *chatFolderRepo) RemoveChatFromFolder(folderID, chatID string) error {
	_, err := r.db.Exec(
		`DELETE FROM chat_folder_items WHERE folder_id = ? AND chat_id = ?`,
		folderID, chatID,
	)
	return err
}

func (r *chatFolderRepo) GetChatIDsByFolder(folderID string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT chat_id FROM chat_folder_items WHERE folder_id = ? ORDER BY added_at DESC`, folderID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if err := rows.Scan(&id); err != nil {
			return nil, err
		}
		ids = append(ids, id)
	}
	return ids, nil
}

func (r *chatFolderRepo) SetChatsForFolder(folderID string, chatIDs []string) error {
	tx, err := r.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()

	if _, err := tx.Exec(`DELETE FROM chat_folder_items WHERE folder_id = ?`, folderID); err != nil {
		return err
	}

	stmt, err := tx.Prepare(`INSERT INTO chat_folder_items (folder_id, chat_id, added_at) VALUES (?, ?, ?)`)
	if err != nil {
		return err
	}
	defer stmt.Close()

	now := time.Now().Format(time.RFC3339)
	for _, chatID := range chatIDs {
		if _, err := stmt.Exec(folderID, chatID, now); err != nil {
			return err
		}
	}

	return tx.Commit()
}



