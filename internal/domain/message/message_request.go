package messagedomain

type SendMessageRequest struct {
	Content       string      `json:"content" binding:"required"`
	Type          MessageType `json:"type" binding:"required,oneof=text image file gif voice video audio system location"`
	ReplyToID     *string     `json:"replyToId,omitempty"`
	ForwardMsgID  *string     `json:"forwardMsgId,omitempty"`
	Latitude      float64     `json:"latitude,omitempty"`
	Longitude     float64     `json:"longitude,omitempty"`
	LocationTitle string      `json:"locationTitle,omitempty"`
	Effect        string      `json:"effect,omitempty" binding:"oneof=confetti fireworks hearts balloons stars ''"`
}

type SendLocationRequest struct {
	Latitude      float64 `json:"latitude" binding:"required"`
	Longitude     float64 `json:"longitude" binding:"required"`
	Title         string  `json:"title,omitempty"`
	ReplyToID     *string `json:"replyToId,omitempty"`
	Effect        string  `json:"effect,omitempty" binding:"oneof=confetti fireworks hearts balloons stars ''"`
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

type ForwardMessageRequest struct {
	MessageID  string `json:"messageId" binding:"required"`
	FromChatID string `json:"fromChatId" binding:"required"`
	ToChatID   string `json:"toChatId" binding:"required"`
}
