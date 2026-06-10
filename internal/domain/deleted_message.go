package domain

import "time"

type DeletedMessageForMe struct {
	UserID    string    `json:"userId"`
	MessageID string    `json:"messageId"`
	DeletedAt time.Time `json:"deletedAt"`
}
