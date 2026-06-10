package domain

import "time"

type StarredMessage struct {
	UserID    string    `json:"userId"`
	MessageID string    `json:"messageId"`
	ChatID    string    `json:"chatId"`
	CreatedAt time.Time `json:"createdAt"`
}

type StarredMessageResponse struct {
	Message   *MessageResponse `json:"message"`
	Chat      *ChatResponse    `json:"chat"`
	CreatedAt time.Time        `json:"createdAt"`
}
