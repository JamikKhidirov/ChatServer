package chatdomain

import "time"

type ChatParticipant struct {
	ChatID     string    `json:"chatId"`
	UserID     string    `json:"userId"`
	Role       string    `json:"role"`
	JoinedAt   time.Time `json:"joinedAt"`
	LastReadAt time.Time `json:"lastReadAt"`
}
