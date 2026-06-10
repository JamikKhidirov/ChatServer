package domain

import "time"

type Draft struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	ChatID     string    `json:"chatId"`
	Content    string    `json:"content"`
	ReplyToID  *string   `json:"replyToId,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	UpdatedAt  time.Time `json:"updatedAt"`
}

type SaveDraftRequest struct {
	ChatID    string  `json:"chatId" binding:"required"`
	Content   string  `json:"content"`
	ReplyToID *string `json:"replyToId,omitempty"`
}

type ScheduledMessage struct {
	ID          string      `json:"id"`
	ChatID      string      `json:"chatId"`
	SenderID    string      `json:"senderId"`
	Content     string      `json:"content"`
	Type        MessageType `json:"type"`
	ReplyToID   *string     `json:"replyToId,omitempty"`
	ScheduledAt string      `json:"scheduledAt"`
	CreatedAt   time.Time   `json:"createdAt"`
	Sent        bool        `json:"sent"`
}

type ScheduleMessageRequest struct {
	ChatID      string      `json:"chatId" binding:"required"`
	Content     string      `json:"content" binding:"required"`
	Type        MessageType `json:"type" binding:"required,oneof=text image file gif voice video"`
	ScheduledAt string      `json:"scheduledAt" binding:"required"`
	ReplyToID   *string     `json:"replyToId,omitempty"`
}
