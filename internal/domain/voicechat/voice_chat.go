package voicechatdomain

import "time"

type VoiceChatStatus string

const (
	VoiceChatActive    VoiceChatStatus = "active"
	VoiceChatScheduled VoiceChatStatus = "scheduled"
	VoiceChatEnded     VoiceChatStatus = "ended"
)

type VoiceChat struct {
	ID          string          `json:"id"`
	ChatID      string          `json:"chatId"`
	StartedBy   string          `json:"startedBy"`
	Title       string          `json:"title,omitempty"`
	Status      VoiceChatStatus `json:"status"`
	ParticipantCount int        `json:"participantCount"`
	ScheduledAt *time.Time      `json:"scheduledAt,omitempty"`
	StartedAt   *time.Time      `json:"startedAt,omitempty"`
	EndedAt     *time.Time      `json:"endedAt,omitempty"`
	CreatedAt   time.Time       `json:"createdAt"`
}

type VoiceChatParticipant struct {
	VoiceChatID string    `json:"voiceChatId"`
	UserID      string    `json:"userId"`
	JoinedAt    time.Time `json:"joinedAt"`
	LeftAt      *time.Time `json:"leftAt,omitempty"`
	Muted       bool      `json:"muted"`
}

type CreateVoiceChatRequest struct {
	Title       string `json:"title,omitempty"`
	ScheduledInMins int `json:"scheduledInMins,omitempty"`
}

type VoiceChatResponse struct {
	ID              string          `json:"id"`
	ChatID          string          `json:"chatId"`
	StartedBy       string          `json:"startedBy"`
	Title           string          `json:"title,omitempty"`
	Status          VoiceChatStatus `json:"status"`
	ParticipantCount int            `json:"participantCount"`
	Participants    []string        `json:"participants,omitempty"`
	ScheduledAt     *time.Time      `json:"scheduledAt,omitempty"`
	StartedAt       *time.Time      `json:"startedAt,omitempty"`
	EndedAt         *time.Time      `json:"endedAt,omitempty"`
	CreatedAt       time.Time       `json:"createdAt"`
}
