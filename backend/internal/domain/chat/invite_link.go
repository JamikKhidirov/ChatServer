package chatdomain

import "time"

type InviteLink struct {
	ID          string     `json:"id"`
	ChatID      string     `json:"chatId"`
	CreatorID   string     `json:"creatorId"`
	Code        string     `json:"code"`
	ExpiresAt   *time.Time `json:"expiresAt,omitempty"`
	UsageLimit  int        `json:"usageLimit,omitempty"`
	UsageCount  int        `json:"usageCount,omitempty"`
	Active      bool       `json:"active"`
	CreatedAt   time.Time  `json:"createdAt"`
}

type CreateInviteLinkRequest struct {
	ExpiresInMins int `json:"expiresInMins,omitempty"`
	UsageLimit    int `json:"usageLimit,omitempty"`
}

type JoinByInviteRequest struct {
	Code string `json:"code" binding:"required"`
}
