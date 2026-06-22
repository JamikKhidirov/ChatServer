package messagedomain

type MessageEditHistory struct {
	ID         string `json:"id"`
	MessageID  string `json:"message_id"`
	OldContent string `json:"old_content"`
	EditedAt   string `json:"edited_at"`
}
