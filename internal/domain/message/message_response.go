package messagedomain

import (
	"time"

	userdomain "ChatServerGolang/internal/domain/user"
)

type MessageResponse struct {
	ID          string                    `json:"id"`
	ChatID      string                    `json:"chatId"`
	Sender      *userdomain.UserResponse  `json:"sender"`
	Content     string                    `json:"content"`
	Type        MessageType               `json:"type"`
	ReplyTo     *MessageResponse          `json:"replyTo,omitempty"`
	ForwardFrom *userdomain.UserResponse  `json:"forwardFrom,omitempty"`
	FileName    string                    `json:"fileName,omitempty"`
	FileSize    int64                     `json:"fileSize,omitempty"`
	FileURL     string                    `json:"fileUrl,omitempty"`
	Caption     string                    `json:"caption,omitempty"`
	MimeType    string                    `json:"mimeType,omitempty"`
	Duration    int                       `json:"duration,omitempty"`
	Width       int                       `json:"width,omitempty"`
	Height      int                       `json:"height,omitempty"`
	Reactions   []*Reaction               `json:"reactions,omitempty"`
	Pinned      bool                      `json:"pinned"`
	ReadBy      []*userdomain.UserResponse `json:"readBy,omitempty"`
	CreatedAt   time.Time                 `json:"createdAt"`
	UpdatedAt   time.Time                 `json:"updatedAt"`
	Edited      bool                      `json:"edited"`
	Deleted     bool                      `json:"deleted"`
}
