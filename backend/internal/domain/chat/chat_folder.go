package chatdomain

import "time"

type ChatFolder struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Name      string    `json:"name"`
	Emoji     string    `json:"emoji,omitempty"`
	Order     int       `json:"order"`
	CreatedAt time.Time `json:"createdAt"`
}

type ChatFolderItem struct {
	FolderID string `json:"folderId"`
	ChatID   string `json:"chatId"`
}

type ChatFolderWithChats struct {
	ChatFolder
	ChatIDs []string      `json:"chatIds,omitempty"`
	Chats   []*ChatResponse `json:"chats,omitempty"`
}

type CreateChatFolderRequest struct {
	Name    string   `json:"name" binding:"required"`
	Emoji   string   `json:"emoji,omitempty"`
	ChatIDs []string `json:"chatIds,omitempty"`
}

type UpdateChatFolderRequest struct {
	Name    string   `json:"name,omitempty"`
	Emoji   string   `json:"emoji,omitempty"`
	Order   int      `json:"order,omitempty"`
	ChatIDs []string `json:"chatIds,omitempty"`
}
