package domain

type UpdateProfileRequest struct {
	DisplayName string `json:"displayName,omitempty" binding:"omitempty,min=1,max=64"`
	Bio         string `json:"bio,omitempty" binding:"max=256"`
	Phone       string `json:"phone,omitempty" binding:"max=20"`
	Gender      string `json:"gender,omitempty" binding:"omitempty,oneof=male female other"`
	DateOfBirth string `json:"dateOfBirth,omitempty" binding:"max=10"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}
