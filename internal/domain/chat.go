package domain

import "time"

type ChatType string

const (
	ChatPrivate ChatType = "private"
	ChatGroup   ChatType = "group"
)

type Chat struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	AvatarURL   string    `json:"avatarUrl"`
	Type        ChatType  `json:"type"`
	CreatedBy   string    `json:"createdBy"`
	CreatedAt   time.Time `json:"createdAt"`
	UpdatedAt   time.Time `json:"updatedAt"`
}

type ChatParticipant struct {
	ChatID    string    `json:"chatId"`
	UserID    string    `json:"userId"`
	Role      string    `json:"role"`
	JoinedAt  time.Time `json:"joinedAt"`
	LastReadAt time.Time `json:"lastReadAt"`
}

type CreateChatRequest struct {
	Name         string   `json:"name" binding:"max=64"`
	Type         ChatType `json:"type" binding:"required,oneof=private group"`
	ParticipantIDs []string `json:"participantIds" binding:"required,min=1"`
	Description  string   `json:"description,omitempty" binding:"max=512"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name,omitempty" binding:"min=1,max=64"`
	Description string `json:"description,omitempty" binding:"max=512"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}

type ChatResponse struct {
	ID           string          `json:"id"`
	Name         string          `json:"name"`
	Description  string          `json:"description"`
	AvatarURL    string          `json:"avatarUrl"`
	Type         ChatType        `json:"type"`
	CreatedBy    string          `json:"createdBy"`
	Participants []*UserResponse `json:"participants"`
	LastMessage  *MessageResponse `json:"lastMessage,omitempty"`
	UnreadCount  int             `json:"unreadCount"`
	CreatedAt    time.Time       `json:"createdAt"`
}

type AddParticipantRequest struct {
	UserID string `json:"userId" binding:"required"`
}
