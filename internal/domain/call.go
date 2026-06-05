package domain

import "time"

type CallStatus string

const (
	CallInitiated  CallStatus = "initiated"
	CallOngoing    CallStatus = "ongoing"
	CallEnded      CallStatus = "ended"
	CallMissed     CallStatus = "missed"
	CallRejected   CallStatus = "rejected"
)

type Call struct {
	ID        string     `json:"id"`
	ChatID    string     `json:"chatId"`
	CallerID  string     `json:"callerId"`
	CalleeID  string     `json:"calleeId"`
	Status    CallStatus `json:"status"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt,omitempty"`
}

type CallOfferData struct {
	CallID string `json:"callId"`
	ChatID string `json:"chatId"`
	SDP    string `json:"sdp"`
}

type CallAnswerData struct {
	CallID string `json:"callId"`
	SDP    string `json:"sdp"`
}

type CallICEData struct {
	CallID      string `json:"callId"`
	Candidate   string `json:"candidate"`
	SDPMLineIdx int    `json:"sdpMLineIdx"`
}

type CallEndData struct {
	CallID string `json:"callId"`
}

type CallResponse struct {
	ID        string     `json:"id"`
	ChatID    string     `json:"chatId"`
	Caller    *UserResponse `json:"caller"`
	Callee    *UserResponse `json:"callee"`
	Status    CallStatus `json:"status"`
	StartedAt time.Time  `json:"startedAt"`
	EndedAt   *time.Time `json:"endedAt,omitempty"`
}
