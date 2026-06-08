package domain

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

type InitiateCallRequest struct {
	ChatID string `json:"chatId" binding:"required"`
}

type RespondCallRequest struct {
	Action string `json:"action" binding:"required,oneof=accept reject"`
}
