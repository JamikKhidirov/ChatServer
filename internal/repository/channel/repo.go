package channelrepo

import (
	"database/sql"
	"time"

	channeldomain "ChatServerGolang/internal/domain/channel"
	"ChatServerGolang/internal/repository"
)

type channelSubscriberRepository struct {
	db *sql.DB
}

func NewChannelSubscriberRepository(db *sql.DB) repository.ChannelSubscriberRepository {
	return &channelSubscriberRepository{db: db}
}

func (r *channelSubscriberRepository) Subscribe(channelID, userID string, role string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO channel_subscribers (channel_id, user_id, role, subscribed_at)
		VALUES (?, ?, ?, ?)`,
		channelID, userID, role, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *channelSubscriberRepository) Unsubscribe(channelID, userID string) error {
	_, err := r.db.Exec(
		`DELETE FROM channel_subscribers WHERE channel_id = ? AND user_id = ?`,
		channelID, userID,
	)
	return err
}

func (r *channelSubscriberRepository) IsSubscribed(channelID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM channel_subscribers WHERE channel_id = ? AND user_id = ?`,
		channelID, userID,
	).Scan(&count)
	return count > 0, err
}

func (r *channelSubscriberRepository) GetSubscribers(channelID string) ([]*channeldomain.ChannelSubscriber, error) {
	rows, err := r.db.Query(
		`SELECT channel_id, user_id, role, subscribed_at
		FROM channel_subscribers WHERE channel_id = ? ORDER BY subscribed_at DESC`,
		channelID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var subscribers []*channeldomain.ChannelSubscriber
	for rows.Next() {
		var s channeldomain.ChannelSubscriber
		var subscribedAt string
		if err := rows.Scan(&s.ChannelID, &s.UserID, &s.Role, &subscribedAt); err != nil {
			return nil, err
		}
		s.SubscribedAt = repository.ParseTime(subscribedAt)
		subscribers = append(subscribers, &s)
	}
	return subscribers, nil
}

func (r *channelSubscriberRepository) GetSubscribedChannels(userID string) ([]string, error) {
	rows, err := r.db.Query(
		`SELECT channel_id FROM channel_subscribers WHERE user_id = ? ORDER BY subscribed_at DESC`,
		userID,
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

func (r *channelSubscriberRepository) SetRole(channelID, userID, role string) error {
	_, err := r.db.Exec(
		`UPDATE channel_subscribers SET role = ? WHERE channel_id = ? AND user_id = ?`,
		role, channelID, userID,
	)
	return err
}

func (r *channelSubscriberRepository) GetRole(channelID, userID string) (string, error) {
	var role string
	err := r.db.QueryRow(
		`SELECT COALESCE(role, '') FROM channel_subscribers WHERE channel_id = ? AND user_id = ?`,
		channelID, userID,
	).Scan(&role)
	if err == sql.ErrNoRows {
		return "", nil
	}
	return role, err
}
