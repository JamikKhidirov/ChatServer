package domain

type UpdatePushTokenRequest struct {
	Token    string `json:"token" binding:"required"`
	Provider string `json:"provider" binding:"required,oneof=fcm apns"`
}
