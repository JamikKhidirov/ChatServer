package domain

type ExportChatRequest struct {
	ChatID string `json:"chatId" binding:"required"`
}
