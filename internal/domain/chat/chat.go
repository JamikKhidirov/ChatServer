package chatdomain

import "time"

type ChatType string

const (
	ChatPrivate  ChatType = "private"
	ChatGroup    ChatType = "group"
	ChatChannel  ChatType = "channel"
)

type Role string

const (
	RoleOwner      Role = "owner"
	RoleAdmin      Role = "admin"
	RoleModerator  Role = "moderator"
	RoleEditor     Role = "editor"
	RoleMember     Role = "member"
	RoleReadOnly   Role = "readonly"
	RoleSubscriber Role = "subscriber"
)

func IsAdminRole(role string) bool {
	return role == string(RoleOwner) || role == string(RoleAdmin)
}

func CanSendMessage(role string) bool {
	return role != string(RoleReadOnly) && role != string(RoleSubscriber)
}

func CanManageMessages(role string) bool {
	return role == string(RoleOwner) || role == string(RoleAdmin) || role == string(RoleModerator)
}

func CanManageMembers(role string) bool {
	return role == string(RoleOwner) || role == string(RoleAdmin)
}

func CanPinMessages(role string) bool {
	return role == string(RoleOwner) || role == string(RoleAdmin) || role == string(RoleEditor)
}

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
