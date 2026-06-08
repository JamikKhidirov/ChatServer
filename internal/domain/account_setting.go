package domain

type AccountSetting struct {
	UserID        string `json:"userId"`
	Language      string `json:"language"`
	Theme         string `json:"theme"`
	Notifications bool   `json:"notifications"`
	SoundEnabled  bool   `json:"soundEnabled"`
	LastSeenMode  string `json:"lastSeenMode"` // everyone, nobody, contacts
	UpdatedAt     string `json:"updatedAt"`
}

type UpdateAccountSettingRequest struct {
	Language      *string `json:"language,omitempty"`
	Theme         *string `json:"theme,omitempty"`
	Notifications *bool   `json:"notifications,omitempty"`
	SoundEnabled  *bool   `json:"soundEnabled,omitempty"`
	LastSeenMode  *string `json:"lastSeenMode,omitempty"`
}
