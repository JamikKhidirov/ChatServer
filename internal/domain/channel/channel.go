package channeldomain

import "time"

type ChannelSubscriber struct {
	ChannelID   string    `json:"channelId"`
	UserID      string    `json:"userId"`
	Role        string    `json:"role"`
	SubscribedAt time.Time `json:"subscribedAt"`
}

type SubscribeRequest struct {
	ChannelID string `json:"channelId" binding:"required"`
}
