package domain

import "time"

type ChatResponse struct {
	ID           string           `json:"id"`
	Name         string           `json:"name"`
	Description  string           `json:"description"`
	AvatarURL    string           `json:"avatarUrl"`
	Type         ChatType         `json:"type"`
	CreatedBy    string           `json:"createdBy"`
	Participants []*UserResponse  `json:"participants"`
	LastMessage  *MessageResponse `json:"lastMessage,omitempty"`
	UnreadCount  int              `json:"unreadCount"`
	CreatedAt    time.Time        `json:"createdAt"`
}
