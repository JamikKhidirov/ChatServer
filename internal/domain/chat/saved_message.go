package chatdomain

import (
	"time"
	messagedomain "ChatServerGolang/internal/domain/message"
)

type SavedMessage struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	MessageID string    `json:"messageId"`
	ChatID    string    `json:"chatId"`
	CreatedAt time.Time `json:"createdAt"`
}

type SavedMessageResponse struct {
	ID        string                         `json:"id"`
	Message   *messagedomain.MessageResponse `json:"message"`
	Chat      *ChatResponse                  `json:"chat"`
	CreatedAt time.Time                      `json:"createdAt"`
}
