package domain

type AdminUser struct {
	UserID    string `json:"user_id"`
	Role      string `json:"role"`
	CreatedAt string `json:"created_at"`
}

type AdminLog struct {
	ID         string `json:"id"`
	AdminID    string `json:"admin_id"`
	Action     string `json:"action"`
	TargetType string `json:"target_type"`
	TargetID   string `json:"target_id"`
	Details    string `json:"details"`
	CreatedAt  string `json:"created_at"`
}

type AppSetting struct {
	Key       string `json:"key"`
	Value     string `json:"value"`
	UpdatedAt string `json:"updated_at"`
}

type AdminDashboard struct {
	TotalUsers      int64 `json:"total_users"`
	TotalChats      int64 `json:"total_chats"`
	TotalMessages   int64 `json:"total_messages"`
	ActiveToday     int64 `json:"active_today"`
	PendingReports  int64 `json:"pending_reports"`
	BlockedIPs      int64 `json:"blocked_ips"`
}

type AdminUserResponse struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	DisplayName  string `json:"display_name"`
	Phone        string `json:"phone"`
	Online       int    `json:"online"`
	LastSeen     string `json:"last_seen"`
	CreatedAt    string `json:"created_at"`
	IsAdmin      bool   `json:"is_admin"`
	MessageCount int64  `json:"message_count"`
}

type AdminMessageResponse struct {
	ID         string      `json:"id"`
	ChatID     string      `json:"chat_id"`
	SenderID   string      `json:"sender_id"`
	SenderName string      `json:"sender_name"`
	Content    string      `json:"content"`
	Type       string      `json:"type"`
	CreatedAt  string      `json:"created_at"`
	Decrypted  string      `json:"decrypted,omitempty"`
}

type AdminBanRequest struct {
	UserID string `json:"user_id" binding:"required"`
	Reason string `json:"reason" binding:"required"`
}

type AdminUpdateSettingRequest struct {
	Key   string `json:"key" binding:"required"`
	Value string `json:"value" binding:"required"`
}
