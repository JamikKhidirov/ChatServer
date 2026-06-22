package userdomain

import "time"

type UserResponse struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	AvatarURL   string    `json:"avatarUrl"`
	Bio         string    `json:"bio"`
	Phone       string    `json:"phone,omitempty"`
	Gender      string    `json:"gender,omitempty"`
	DateOfBirth string    `json:"dateOfBirth,omitempty"`
	Status      string    `json:"status"`
	Online      bool      `json:"online"`
	LastSeen    time.Time `json:"lastSeen"`
	IsAdmin     bool      `json:"isAdmin"`
}
