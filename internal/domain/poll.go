package domain

import "time"

type Poll struct {
	ID             string    `json:"id"`
	ChatID         string    `json:"chatId"`
	CreatorID      string    `json:"creatorId"`
	Question       string    `json:"question"`
	Options        string    `json:"options"`
	IsAnonymous    bool      `json:"isAnonymous"`
	MultipleChoice bool      `json:"multipleChoice"`
	ExpiresAt      *string   `json:"expiresAt,omitempty"`
	CreatedAt      time.Time `json:"createdAt"`
	Closed         bool      `json:"closed"`
}

type PollVote struct {
	PollID      string    `json:"pollId"`
	UserID      string    `json:"userId"`
	OptionIndex int       `json:"optionIndex"`
	VotedAt     time.Time `json:"votedAt"`
}

type PollWithResults struct {
	Poll
	OptionsList []PollOption `json:"options"`
	TotalVotes  int          `json:"totalVotes"`
	VotedOption *int         `json:"votedOption,omitempty"`
}

type PollOption struct {
	Text  string `json:"text"`
	Votes int    `json:"votes"`
}

type CreatePollRequest struct {
	ChatID         string   `json:"chatId" binding:"required"`
	Question       string   `json:"question" binding:"required"`
	Options        []string `json:"options" binding:"required,min=2,max=10"`
	IsAnonymous    bool     `json:"isAnonymous"`
	MultipleChoice bool     `json:"multipleChoice"`
	ExpiresInMins  *int     `json:"expiresInMins,omitempty"`
}

type VotePollRequest struct {
	OptionIndex int `json:"optionIndex" binding:"required"`
}
