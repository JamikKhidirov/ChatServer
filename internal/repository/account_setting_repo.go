package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type accountSettingRepository struct {
	db *sql.DB
}

func NewAccountSettingRepository(db *sql.DB) AccountSettingRepository {
	return &accountSettingRepository{db: db}
}

func (r *accountSettingRepository) GetByUserID(userID string) (*domain.AccountSetting, error) {
	row := r.db.QueryRow(
		`SELECT user_id, language, theme, notifications, sound_enabled, last_seen_mode, updated_at
		FROM account_settings WHERE user_id = ?`, userID,
	)
	var (
		s         domain.AccountSetting
		updatedAt string
	)
	err := row.Scan(&s.UserID, &s.Language, &s.Theme, &s.Notifications, &s.SoundEnabled, &s.LastSeenMode, &updatedAt)
	if err != nil {
		return nil, err
	}
	s.UpdatedAt = updatedAt
	return &s, nil
}

func (r *accountSettingRepository) Upsert(setting *domain.AccountSetting) error {
	_, err := r.db.Exec(
		`INSERT INTO account_settings (user_id, language, theme, notifications, sound_enabled, last_seen_mode, updated_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)
		ON CONFLICT(user_id) DO UPDATE SET
			language = COALESCE(NULLIF(?, ''), language),
			theme = COALESCE(NULLIF(?, ''), theme),
			notifications = COALESCE(?, notifications),
			sound_enabled = COALESCE(?, sound_enabled),
			last_seen_mode = COALESCE(NULLIF(?, ''), last_seen_mode),
			updated_at = ?`,
		setting.UserID, setting.Language, setting.Theme, boolToInt(setting.Notifications), boolToInt(setting.SoundEnabled), setting.LastSeenMode, time.Now().Format(time.RFC3339),
		setting.Language, setting.Theme, boolToInt(setting.Notifications), boolToInt(setting.SoundEnabled), setting.LastSeenMode, time.Now().Format(time.RFC3339),
	)
	return err
}
