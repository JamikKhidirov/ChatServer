package userdomain

import "time"

type Block struct {
	UserID    string    `json:"userId"`
	BlockedID string    `json:"blockedId"`
	CreatedAt time.Time `json:"createdAt"`
}

type BlockUserRequest struct {
	BlockedID string `json:"blockedId" binding:"required"`
}
