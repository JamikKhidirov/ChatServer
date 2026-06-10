package domain

type MessageSelfDestruct struct {
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
	DeleteAt  string `json:"delete_at"`
}

type SetSelfDestructRequest struct {
	MessageID    string `json:"message_id" binding:"required"`
	DeleteAfter  int    `json:"delete_after" binding:"required"` // seconds
}
