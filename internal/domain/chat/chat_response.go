package chatdomain

import (
	"time"

	userdomain "ChatServerGolang/internal/domain/user"
	messagedomain "ChatServerGolang/internal/domain/message"
)

type ChatResponse struct {
	ID           string                      `json:"id"`
	Name         string                      `json:"name"`
	Description  string                      `json:"description"`
	AvatarURL    string                      `json:"avatarUrl"`
	Type         ChatType                    `json:"type"`
	CreatedBy    string                      `json:"createdBy"`
	Participants []*userdomain.UserResponse  `json:"participants"`
	LastMessage  *messagedomain.MessageResponse `json:"lastMessage,omitempty"`
	UnreadCount  int                         `json:"unreadCount"`
	CreatedAt    time.Time                   `json:"createdAt"`
}
