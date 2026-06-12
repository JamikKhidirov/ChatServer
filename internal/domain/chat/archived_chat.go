package chatdomain

import "time"

type ArchivedChat struct {
	UserID     string    `json:"userId"`
	ChatID     string    `json:"chatId"`
	ArchivedAt time.Time `json:"archivedAt"`
}
