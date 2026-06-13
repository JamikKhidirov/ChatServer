package messagedomain

type Bookmark struct {
	ID        string `json:"id"`
	UserID    string `json:"user_id"`
	MessageID string `json:"message_id"`
	ChatID    string `json:"chat_id"`
	CreatedAt string `json:"created_at"`
}

type BookmarkMessageRequest struct {
	MessageID string `json:"message_id" binding:"required"`
	ChatID    string `json:"chat_id" binding:"required"`
}
