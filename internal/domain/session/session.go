package sessiondomain

import "time"

type Session struct {
	ID         string    `json:"id"`
	UserID     string    `json:"userId"`
	DeviceName string    `json:"deviceName"`
	IPAddress  string    `json:"ipAddress"`
	LastActive time.Time `json:"lastActive"`
	CreatedAt  time.Time `json:"createdAt"`
}

type CreateSessionRequest struct {
	DeviceName string `json:"deviceName" binding:"required"`
}
