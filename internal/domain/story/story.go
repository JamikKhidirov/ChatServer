package storydomain

import "time"

type StoryType string

const (
	StoryPhoto StoryType = "photo"
	StoryVideo StoryType = "video"
)

type Story struct {
	ID        string    `json:"id"`
	UserID    string    `json:"userId"`
	FilePath  string    `json:"-"`
	FileURL   string    `json:"fileUrl"`
	Type      StoryType `json:"type"`
	Caption   string    `json:"caption,omitempty"`
	ExpiresAt time.Time `json:"expiresAt"`
	CreatedAt time.Time `json:"createdAt"`
}

type StoryView struct {
	StoryID   string    `json:"storyId"`
	UserID    string    `json:"userId"`
	ViewedAt  time.Time `json:"viewedAt"`
}

type StoryWithViews struct {
	Story
	Views     int64             `json:"views"`
	Viewers   []*StoryView      `json:"viewers,omitempty"`
}

type StoryResponse struct {
	ID        string     `json:"id"`
	UserID    string     `json:"userId"`
	FileURL   string     `json:"fileUrl"`
	Type      StoryType  `json:"type"`
	Caption   string     `json:"caption,omitempty"`
	ExpiresAt time.Time  `json:"expiresAt"`
	CreatedAt time.Time  `json:"createdAt"`
	Views     int64      `json:"views"`
	Viewed    bool       `json:"viewed"`
}

type CreateStoryRequest struct {
	Type    StoryType `json:"type" binding:"required,oneof=photo video"`
	Caption string    `json:"caption,omitempty"`
}
