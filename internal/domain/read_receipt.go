package domain

import "time"

type ReadReceipt struct {
	MessageID string    `json:"messageId"`
	UserID    string    `json:"userId"`
	ReadAt    time.Time `json:"readAt"`
}
