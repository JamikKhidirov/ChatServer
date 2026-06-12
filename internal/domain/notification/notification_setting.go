package notificationdomain

type NotificationSetting struct {
	UserID string `json:"userId"`
	ChatID string `json:"chatId"`
	Muted  bool   `json:"muted"`
}

type UpdateNotificationSettingRequest struct {
	Muted bool `json:"muted"`
}
