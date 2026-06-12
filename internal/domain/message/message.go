package messagedomain

import "time"

type MessageType string

const (
	MessageText   MessageType = "text"
	MessageImage  MessageType = "image"
	MessageFile   MessageType = "file"
	MessageGif    MessageType = "gif"
	MessageVoice      MessageType = "voice"
	MessageVideo      MessageType = "video"
	MessageVideoCircle MessageType = "video_circle"
	MessageAudio      MessageType = "audio"
	MessageSystem     MessageType = "system"
)

type Message struct {
	ID          string      `json:"id"`
	ChatID      string      `json:"chatId"`
	SenderID    string      `json:"senderId"`
	Content     string      `json:"content"`
	Type        MessageType `json:"type"`
	ReplyToID   *string     `json:"replyToId,omitempty"`
	ForwardFrom *string     `json:"forwardFrom,omitempty"`
	FileName    string      `json:"fileName,omitempty"`
	FileSize    int64       `json:"fileSize,omitempty"`
	FilePath    string      `json:"-"`
	Caption     string      `json:"caption,omitempty"`
	MimeType    string      `json:"mimeType,omitempty"`
	Duration    int         `json:"duration,omitempty"`
	Width       int         `json:"width,omitempty"`
	Height      int         `json:"height,omitempty"`
	CreatedAt   time.Time   `json:"createdAt"`
	UpdatedAt   time.Time   `json:"updatedAt"`
	Pinned      bool        `json:"pinned"`
	DeletedAt   *time.Time  `json:"deletedAt,omitempty"`
}
