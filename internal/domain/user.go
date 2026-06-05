package domain

import "time"

type User struct {
	ID           string    `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	DisplayName  string    `json:"displayName"`
	AvatarURL    string    `json:"avatarUrl"`
	Bio          string    `json:"bio"`
	Status       string    `json:"status"`
	PushToken    string    `json:"-"`
	PushProvider string    `json:"-"`
	Online       bool      `json:"online"`
	LastSeen     time.Time `json:"lastSeen"`
	Deleted      bool      `json:"-"`
	CreatedAt    time.Time `json:"createdAt"`
	UpdatedAt    time.Time `json:"updatedAt"`
}

type RegisterRequest struct {
	Username        string `json:"username" binding:"required,min=3,max=32"`
	Email           string `json:"email" binding:"required,email"`
	Password        string `json:"password" binding:"required,min=6"`
	DisplayName     string `json:"displayName" binding:"required,min=1,max=64"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required"`
}

type AuthResponse struct {
	Token string `json:"token"`
	User  *User  `json:"user"`
}

type RefreshTokenResponse struct {
	Token string `json:"token"`
}

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName" binding:"min=1,max=64"`
	Bio         string `json:"bio" binding:"max=256"`
	AvatarURL   string `json:"avatarUrl"`
}

type UpdatePushTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Provider string `json:"provider" binding:"required,oneof=fcm apns"`
}

type UpdateStatusRequest struct {
	Status string `json:"status" binding:"max=100"`
}

type ChangePasswordRequest struct {
	OldPassword string `json:"oldPassword" binding:"required"`
	NewPassword string `json:"newPassword" binding:"required,min=6"`
}

type UserResponse struct {
	ID          string    `json:"id"`
	Username    string    `json:"username"`
	DisplayName string    `json:"displayName"`
	AvatarURL   string    `json:"avatarUrl"`
	Bio         string    `json:"bio"`
	Status      string    `json:"status"`
	Online      bool      `json:"online"`
	LastSeen    time.Time `json:"lastSeen"`
}

type Block struct {
	UserID    string    `json:"userId"`
	BlockedID string    `json:"blockedId"`
	CreatedAt time.Time `json:"createdAt"`
}

type BlockUserRequest struct {
	BlockedID string `json:"blockedId" binding:"required"`
}

type NotificationSetting struct {
	UserID string `json:"userId"`
	ChatID string `json:"chatId"`
	Muted  bool   `json:"muted"`
}

type UpdateNotificationSettingRequest struct {
	Muted bool `json:"muted"`
}

func (u *User) ToResponse() *UserResponse {
	return &UserResponse{
		ID:          u.ID,
		Username:    u.Username,
		DisplayName: u.DisplayName,
		AvatarURL:   u.AvatarURL,
		Bio:         u.Bio,
		Status:      u.Status,
		Online:      u.Online,
		LastSeen:    u.LastSeen,
	}
}
