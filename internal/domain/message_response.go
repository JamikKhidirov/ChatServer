package domain

import "time"

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
