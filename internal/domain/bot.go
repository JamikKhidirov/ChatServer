package domain

import "time"

type Bot struct {
	ID         string    `json:"id"`
	Token      string    `json:"-"`
	OwnerID    string    `json:"ownerId"`
	Name       string    `json:"name"`
	AvatarURL  string    `json:"avatarUrl,omitempty"`
	WebhookURL string    `json:"webhookUrl,omitempty"`
	CreatedAt  time.Time `json:"createdAt"`
	Active     bool      `json:"active"`
}

type CreateBotRequest struct {
	Name       string `json:"name" binding:"required"`
	WebhookURL string `json:"webhookUrl,omitempty"`
}

type UpdateBotRequest struct {
	Name       string `json:"name,omitempty"`
	AvatarURL  string `json:"avatarUrl,omitempty"`
	WebhookURL string `json:"webhookUrl,omitempty"`
}
