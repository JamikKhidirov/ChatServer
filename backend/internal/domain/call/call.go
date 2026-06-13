package calldomain

import "time"

type CallType string

const (
	CallAudio CallType = "audio"
	CallVideo CallType = "video"
)

type CallStatus string

const (
	CallInitiated CallStatus = "initiated"
	CallOngoing   CallStatus = "ongoing"
	CallEnded     CallStatus = "ended"
	CallMissed    CallStatus = "missed"
	CallRejected  CallStatus = "rejected"
)

type Call struct {
	ID        string     `json:"id"`
	ChatID    string     `json:"chatId"`
	CallerID  string     `json:"callerId"`
	CalleeID  string     `json:"calleeId"`
	Type      CallType   `json:"type"`
	Status    CallStatus `json:"status"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt,omitempty"`
}
