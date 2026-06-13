package messagedomain

type ExportChatRequest struct {
	ChatID string `json:"chatId" binding:"required"`
}
