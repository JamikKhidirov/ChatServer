package domain

import "time"

type StickerPack struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	CreatorID string    `json:"creatorId"`
	Animated  bool      `json:"animated"`
	CreatedAt time.Time `json:"createdAt"`
}

type Sticker struct {
	ID       string `json:"id"`
	PackID   string `json:"packId"`
	Emoji    string `json:"emoji"`
	ImageURL string `json:"imageUrl,omitempty"`
	FilePath string `json:"filePath,omitempty"`
}

type StickerPackWithStickers struct {
	StickerPack
	Stickers []*Sticker `json:"stickers"`
}

type CreateStickerPackRequest struct {
	Name     string `json:"name" binding:"required"`
	Animated bool   `json:"animated"`
}

type AddStickerRequest struct {
	Emoji    string `json:"emoji" binding:"required"`
	ImageURL string `json:"imageUrl,omitempty"`
}
