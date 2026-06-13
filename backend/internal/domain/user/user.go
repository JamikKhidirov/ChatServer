package userdomain

import "time"

type User struct {
	ID            string    `json:"id"`
	Username      string    `json:"username"`
	Email         string    `json:"email"`
	PasswordHash  string    `json:"-"`
	DisplayName   string    `json:"displayName"`
	AvatarURL     string    `json:"avatarUrl"`
	Bio           string    `json:"bio"`
	Phone         string    `json:"phone,omitempty"`
	Gender        string    `json:"gender,omitempty"`
	DateOfBirth   string    `json:"dateOfBirth,omitempty"`
	Status        string    `json:"status"`
	PushToken     string    `json:"-"`
	PushProvider  string    `json:"-"`
	Online        bool      `json:"online"`
	LastSeen      time.Time `json:"lastSeen"`
	Deleted       bool      `json:"-"`
	CreatedAt     time.Time `json:"createdAt"`
	UpdatedAt     time.Time `json:"updatedAt"`
	IsAdmin       bool      `json:"isAdmin"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		Bio:         u.Bio,
		Phone:       u.Phone,
		Gender:      u.Gender,
		DateOfBirth: u.DateOfBirth,
		Status:      u.Status,
		Online:      u.Online,
		LastSeen:    u.LastSeen,
		IsAdmin:     u.IsAdmin,
	}
}
