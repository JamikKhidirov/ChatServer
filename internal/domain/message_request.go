package domain

type SendMessageRequest struct {
	Content      string      `json:"content" binding:"required"`
	Type         MessageType `json:"type" binding:"required,oneof=text image file gif voice video"`
	ReplyToID    *string     `json:"replyToId,omitempty"`
	ForwardMsgID *string     `json:"forwardMsgId,omitempty"`
}

type EditMessageRequest struct {
	Content string `json:"content" binding:"required"`
}

type SearchMessagesRequest struct {
	Query  string `json:"query" form:"q"`
	Limit  int    `json:"limit" form:"limit,default=50"`
	Offset int    `json:"offset" form:"offset,default=0"`
}

type PinMessageRequest struct {
	Pin bool `json:"pin"`
}

type ResendMessageRequest struct {
	MessageID string `json:"messageId" binding:"required"`
}
