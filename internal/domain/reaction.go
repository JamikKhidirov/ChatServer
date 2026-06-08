package domain

import "time"

type Reaction struct {
	MessageID string       `json:"messageId"`
	UserID    string       `json:"userId"`
	Emoji     string       `json:"emoji"`
	CreatedAt time.Time    `json:"createdAt"`
	User      *UserResponse `json:"user,omitempty"`
}

type AddReactionRequest struct {
	Emoji string `json:"emoji" binding:"required"`
}
