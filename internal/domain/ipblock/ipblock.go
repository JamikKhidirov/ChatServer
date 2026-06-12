package ipblockdomain

type IPBlock struct {
	IPAddress string `json:"ip_address"`
	Reason    string `json:"reason"`
	BlockedAt string `json:"blocked_at"`
	ExpiresAt string `json:"expires_at,omitempty"`
	Attempts  int    `json:"attempts"`
}

type LoginAttempt struct {
	IPAddress   string `json:"ip_address"`
	AttemptedAt string `json:"attempted_at"`
	Email       string `json:"email"`
	Success     int    `json:"success"`
}
