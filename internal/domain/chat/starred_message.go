package chatdomain

import (
	"time"

	messagedomain "ChatServerGolang/internal/domain/message"
)

type StarredMessage struct {
	UserID    string    `json:"userId"`
	MessageID string    `json:"messageId"`
	ChatID    string    `json:"chatId"`
	CreatedAt time.Time `json:"createdAt"`
}

type StarredMessageResponse struct {
	Message   *messagedomain.MessageResponse `json:"message"`
	Chat      *ChatResponse                  `json:"chat"`
	CreatedAt time.Time                      `json:"createdAt"`
}
