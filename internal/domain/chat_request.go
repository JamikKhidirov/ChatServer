package domain

type CreateChatRequest struct {
	Name           string   `json:"name" binding:"max=64"`
	Type           ChatType `json:"type" binding:"required,oneof=private group"`
	ParticipantIDs []string `json:"participantIds" binding:"required,min=1"`
	Description    string   `json:"description,omitempty" binding:"max=512"`
}

type UpdateGroupRequest struct {
	Name        string `json:"name,omitempty" binding:"min=1,max=64"`
	Description string `json:"description,omitempty" binding:"max=512"`
	AvatarURL   string `json:"avatarUrl,omitempty"`
}

type AddParticipantRequest struct {
	UserID string `json:"userId" binding:"required"`
}
