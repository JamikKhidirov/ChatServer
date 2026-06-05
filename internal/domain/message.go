package domain

import "time"

type MessageType string

const (
	MessageText   MessageType = "text"
	MessageImage  MessageType = "image"
	MessageFile   MessageType = "file"
	MessageSystem MessageType = "system"
)

type Message struct {
	ID         string      `json:"id"`
	ChatID     string      `json:"chatId"`
	SenderID   string      `json:"senderId"`
	Content    string      `json:"content"`
	Type       MessageType `json:"type"`
	ReplyToID  *string     `json:"replyToId,omitempty"`
	ForwardFrom *string    `json:"forwardFrom,omitempty"`
	FileName   string      `json:"fileName,omitempty"`
	FileSize   int64       `json:"fileSize,omitempty"`
	FilePath   string      `json:"-"`
	CreatedAt  time.Time   `json:"createdAt"`
	UpdatedAt  time.Time   `json:"updatedAt"`
	Pinned     bool        `json:"pinned"`
	DeletedAt  *time.Time  `json:"deletedAt,omitempty"`
}

type SendMessageRequest struct {
	Content    string      `json:"content" binding:"required"`
	Type       MessageType `json:"type" binding:"required,oneof=text image file"`
	ReplyToID  *string     `json:"replyToId,omitempty"`
	ForwardMsgID *string   `json:"forwardMsgId,omitempty"`
}

type EditMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type SearchMessagesRequest struct {
	Query  string `json:"query" form:"q"`
	Limit  int    `json:"limit" form:"limit,default=50"`
	Offset int    `json:"offset" form:"offset,default=0"`
}

type MessageResponse struct {
	ID          string           `json:"id"`
	ChatID      string           `json:"chatId"`
	Sender      *UserResponse    `json:"sender"`
	Content     string           `json:"content"`
	Type        MessageType      `json:"type"`
	ReplyTo     *MessageResponse `json:"replyTo,omitempty"`
	ForwardFrom *UserResponse    `json:"forwardFrom,omitempty"`
	FileName    string           `json:"fileName,omitempty"`
	FileSize    int64            `json:"fileSize,omitempty"`
	FileURL     string           `json:"fileUrl,omitempty"`
	Reactions   []*Reaction      `json:"reactions,omitempty"`
	Pinned      bool             `json:"pinned"`
	ReadBy      []*UserResponse  `json:"readBy,omitempty"`
	CreatedAt   time.Time        `json:"createdAt"`
	UpdatedAt   time.Time        `json:"updatedAt"`
	Edited      bool             `json:"edited"`
	Deleted     bool             `json:"deleted"`
}

type Reaction struct {
	MessageID string `json:"messageId"`
	UserID    string `json:"userId"`
	Emoji     string `json:"emoji"`
	CreatedAt time.Time `json:"createdAt"`
	User      *UserResponse `json:"user,omitempty"`
}

type AddReactionRequest struct {
	Emoji string `json:"emoji" binding:"required"`
}

type ReadReceipt struct {
	MessageID string    `json:"messageId"`
	UserID    string    `json:"userId"`
	ReadAt    time.Time `json:"readAt"`
}

type PinMessageRequest struct {
	Pin bool `json:"pin"`
}

type ResendMessageRequest struct {
	MessageID string `json:"messageId" binding:"required"`
}
