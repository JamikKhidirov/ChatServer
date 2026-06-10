package domain

type Mention struct {
	MessageID string `json:"messageId"`
	UserID    string `json:"userId"`
	Username  string `json:"username"`
}

type MentionedUser struct {
	UserID   string `json:"userId"`
	Username string `json:"username"`
	Offset   int    `json:"offset"`
	Length   int    `json:"length"`
}
