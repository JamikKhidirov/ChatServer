package domain

import "time"

type CallResponse struct {
	ID        string       `json:"id"`
	ChatID    string       `json:"chatId"`
	Caller    *UserResponse `json:"caller"`
	Callee    *UserResponse `json:"callee"`
	Status    CallStatus   `json:"status"`
	StartedAt time.Time    `json:"startedAt"`
	EndedAt   *time.Time   `json:"endedAt,omitempty"`
}
