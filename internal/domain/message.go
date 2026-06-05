package domain

import "time"

type MessageType string

const (
	MessageText  MessageType = "text"
	MessageImage MessageType = "image"
	MessageFile  MessageType = "file"
	MessageSystem MessageType = "system"
)

type Message struct {
	ID        string      `json:"id"`
	ChatID    string      `json:"chatId"`
	SenderID  string      `json:"senderId"`
	Content   string      `json:"content"`
	Type      MessageType `json:"type"`
	ReplyToID *string     `json:"replyToId,omitempty"`
	CreatedAt time.Time   `json:"createdAt"`
	UpdatedAt time.Time   `json:"updatedAt"`
	DeletedAt *time.Time  `json:"deletedAt,omitempty"`
}

type SendMessageRequest struct {
	Content   string      `json:"content" binding:"required"`
	Type      MessageType `json:"type" binding:"required,oneof=text image file"`
	ReplyToID *string     `json:"replyToId,omitempty"`
}

type EditMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type MessageResponse struct {
	ID        string          `json:"id"`
	ChatID    string          `json:"chatId"`
	Sender    *UserResponse   `json:"sender"`
	Content   string          `json:"content"`
	Type      MessageType     `json:"type"`
	ReplyTo   *MessageResponse `json:"replyTo,omitempty"`
	CreatedAt time.Time       `json:"createdAt"`
	UpdatedAt time.Time       `json:"updatedAt"`
	Edited    bool            `json:"edited"`
	Deleted   bool            `json:"deleted"`
}
