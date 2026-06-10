package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type stickerRepository struct {
	db *sql.DB
}

func NewStickerRepository(db *sql.DB) StickerRepository {
	return &stickerRepository{db: db}
}

func (r *stickerRepository) CreatePack(pack *domain.StickerPack) error {
	_, err := r.db.Exec(
		`INSERT INTO sticker_packs (id, name, creator_id, animated, created_at) VALUES (?, ?, ?, ?, ?)`,
		pack.ID, pack.Name, pack.CreatorID, boolToInt(pack.Animated), pack.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *stickerRepository) GetPackByID(id string) (*domain.StickerPack, error) {
	row := r.db.QueryRow(
		`SELECT id, name, creator_id, animated, created_at FROM sticker_packs WHERE id = ?`, id,
	)
	var p domain.StickerPack
	var createdAt string
	var animated int
	if err := row.Scan(&p.ID, &p.Name, &p.CreatorID, &animated, &createdAt); err != nil {
		return nil, err
	}
	p.Animated = animated == 1
	p.CreatedAt = parseTime(createdAt)
	return &p, nil
}

func (r *stickerRepository) GetPacksByUserID(userID string) ([]*domain.StickerPack, error) {
	rows, err := r.db.Query(
		`SELECT id, name, creator_id, animated, created_at FROM sticker_packs WHERE creator_id = ? ORDER BY created_at DESC`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packs := make([]*domain.StickerPack, 0)
	for rows.Next() {
		var p domain.StickerPack
		var createdAt string
		var animated int
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatorID, &animated, &createdAt); err != nil {
			return nil, err
		}
		p.Animated = animated == 1
		p.CreatedAt = parseTime(createdAt)
		packs = append(packs, &p)
	}
	return packs, nil
}

func (r *stickerRepository) ListPacks() ([]*domain.StickerPack, error) {
	rows, err := r.db.Query(
		`SELECT id, name, creator_id, animated, created_at FROM sticker_packs ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	packs := make([]*domain.StickerPack, 0)
	for rows.Next() {
		var p domain.StickerPack
		var createdAt string
		var animated int
		if err := rows.Scan(&p.ID, &p.Name, &p.CreatorID, &animated, &createdAt); err != nil {
			return nil, err
		}
		p.Animated = animated == 1
		p.CreatedAt = parseTime(createdAt)
		packs = append(packs, &p)
	}
	return packs, nil
}

func (r *stickerRepository) AddSticker(sticker *domain.Sticker) error {
	_, err := r.db.Exec(
		`INSERT INTO stickers (id, pack_id, emoji, image_url, file_path, created_at) VALUES (?, ?, ?, ?, ?, ?)`,
		sticker.ID, sticker.PackID, sticker.Emoji, sticker.ImageURL, sticker.FilePath, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *stickerRepository) GetStickersByPackID(packID string) ([]*domain.Sticker, error) {
	rows, err := r.db.Query(
		`SELECT id, pack_id, emoji, image_url, file_path FROM stickers WHERE pack_id = ? ORDER BY created_at`, packID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stickers := make([]*domain.Sticker, 0)
	for rows.Next() {
		var s domain.Sticker
		if err := rows.Scan(&s.ID, &s.PackID, &s.Emoji, &s.ImageURL, &s.FilePath); err != nil {
			return nil, err
		}
		stickers = append(stickers, &s)
	}
	return stickers, nil
}

func (r *stickerRepository) DeletePack(id string) error {
	r.db.Exec(`DELETE FROM stickers WHERE pack_id = ?`, id)
	_, err := r.db.Exec(`DELETE FROM sticker_packs WHERE id = ?`, id)
	return err
}

func (r *stickerRepository) DeleteSticker(id string) error {
	_, err := r.db.Exec(`DELETE FROM stickers WHERE id = ?`, id)
	return err
}

func (r *stickerRepository) AddToUserLibrary(userID, stickerID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO user_stickers (user_id, sticker_id) VALUES (?, ?)`, userID, stickerID,
	)
	return err
}

func (r *stickerRepository) GetUserLibrary(userID string) ([]*domain.Sticker, error) {
	rows, err := r.db.Query(
		`SELECT s.id, s.pack_id, s.emoji, s.image_url, s.file_path
		FROM stickers s INNER JOIN user_stickers us ON us.sticker_id = s.id
		WHERE us.user_id = ?`, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	stickers := make([]*domain.Sticker, 0)
	for rows.Next() {
		var s domain.Sticker
		if err := rows.Scan(&s.ID, &s.PackID, &s.Emoji, &s.ImageURL, &s.FilePath); err != nil {
			return nil, err
		}
		stickers = append(stickers, &s)
	}
	return stickers, nil
}
