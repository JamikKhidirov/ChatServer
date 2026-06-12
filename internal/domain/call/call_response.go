package calldomain

import (
	"time"

	userdomain "ChatServerGolang/internal/domain/user"
)

type CallResponse struct {
	ID        string                  `json:"id"`
	ChatID    string                  `json:"chatId"`
	Caller    *userdomain.UserResponse `json:"caller"`
	Callee    *userdomain.UserResponse `json:"callee"`
	Type      CallType                `json:"type"`
	Status    CallStatus              `json:"status"`
	StartedAt time.Time               `json:"startedAt"`
	EndedAt   *time.Time              `json:"endedAt,omitempty"`
	Duration  int                     `json:"duration,omitempty"`
}
