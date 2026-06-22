package emojidomain

import "time"

type CustomEmoji struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	Shortcode string    `json:"shortcode"`
	FileURL   string    `json:"fileUrl"`
	FilePath  string    `json:"-"`
	Animated  bool      `json:"animated"`
	CreatedAt time.Time `json:"createdAt"`
}

type CustomEmojiResponse struct {
	ID        string    `json:"id"`
	Shortcode string    `json:"shortcode"`
	FileURL   string    `json:"fileUrl"`
	Animated  bool      `json:"animated"`
	CreatedAt time.Time `json:"createdAt"`
}

type CreateEmojiRequest struct {
	Shortcode string `json:"shortcode" binding:"required,min=1,max=32"`
}
