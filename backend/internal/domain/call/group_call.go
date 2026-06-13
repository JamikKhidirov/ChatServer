package calldomain

import "time"

type GroupCall struct {
	ID          string      `json:"id"`
	ChatID      string      `json:"chatId"`
	CallerID    string      `json:"callerId"`
	Type        CallType    `json:"type"`
	Status      CallStatus  `json:"status"`
	StartedAt   time.Time   `json:"startedAt"`
	EndedAt     *time.Time  `json:"endedAt,omitempty"`
}

type GroupCallParticipant struct {
	CallID    string    `json:"callId"`
	UserID    string    `json:"userId"`
	JoinedAt  time.Time `json:"joinedAt"`
	LeftAt    *time.Time `json:"leftAt,omitempty"`
	AudioMuted bool     `json:"audioMuted"`
	VideoMuted bool     `json:"videoMuted"`
}

type GroupCallResponse struct {
	ID           string                  `json:"id"`
	ChatID       string                  `json:"chatId"`
	CallerID     string                  `json:"callerId"`
	Type         CallType                `json:"type"`
	Status       CallStatus              `json:"status"`
	Participants []*GroupCallParticipant `json:"participants"`
	StartedAt    time.Time               `json:"startedAt"`
	EndedAt      *time.Time              `json:"endedAt,omitempty"`
}

type GroupCallInitiateRequest struct {
	ChatID string   `json:"chatId" binding:"required"`
	Type   CallType `json:"type" binding:"required,oneof=audio video"`
}

type GroupCallActionRequest struct {
	CallID     string `json:"callId" binding:"required"`
	Action     string `json:"action" binding:"required,oneof=join leave mute unmute_audio unmute_video mute_video"`
}
