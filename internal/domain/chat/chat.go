package chatdomain

import "time"

type ChatType string

const (
	ChatPrivate  ChatType = "private"
	ChatGroup    ChatType = "group"
	ChatChannel  ChatType = "channel"
)

type Chat struct {
	ID              string    `json:"id"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	AvatarURL       string    `json:"avatarUrl"`
	Type            ChatType  `json:"type"`
	CreatedBy       string    `json:"createdBy"`
	SlowModeSeconds int       `json:"slowModeSeconds,omitempty"`
	CreatedAt       time.Time `json:"createdAt"`
	UpdatedAt       time.Time `json:"updatedAt"`
}
