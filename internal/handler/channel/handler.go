package channelhandler

import (
	chatdomain "ChatServerGolang/internal/domain/chat"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChannelHandler struct {
	channelService service.ChannelService
	chatService    service.ChatService
}

func NewChannelHandler(channelService service.ChannelService, chatService service.ChatService) *ChannelHandler {
	return &ChannelHandler{
		channelService: channelService,
		chatService:    chatService,
	}
}

// Subscribe subscribes to a broadcast channel
// @Tags Channels
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{channelId=string} true "Channel ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /channels/subscribe [post]
func (h *ChannelHandler) Subscribe(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		ChannelID string `json:"channelId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.channelService.Subscribe(req.ChannelID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "subscribed"})
}

// Unsubscribe unsubscribes from a broadcast channel
// @Tags Channels
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{channelId=string} true "Channel ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /channels/unsubscribe [post]
func (h *ChannelHandler) Unsubscribe(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		ChannelID string `json:"channelId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.channelService.Unsubscribe(req.ChannelID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "unsubscribed"})
}

// GetSubscribers returns subscriber list (admin only)
// @Tags Channels
// @Security BearerAuth
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {array} channeldomain.ChannelSubscriber
// @Failure 400 {object} response.ErrorResponse
// @Router /channels/{id}/subscribers [get]
func (h *ChannelHandler) GetSubscribers(c *gin.Context) {
	userID, _ := c.Get("userID")
	channelID := c.Param("id")

	subscribers, err := h.channelService.GetSubscribers(channelID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, subscribers)
}

// GetMyChannels returns channels the user has created or is subscribed to
// @Tags Channels
// @Security BearerAuth
// @Produce json
// @Success 200 {array} chatdomain.ChatResponse
// @Router /channels [get]
func (h *ChannelHandler) GetMyChannels(c *gin.Context) {
	userID, _ := c.Get("userID")

	// Get subscribed channels
	subscribed, err := h.channelService.GetSubscribedChannels(userID.(string))
	if err != nil {
		subscribed = []*chatdomain.ChatResponse{}
	}

	// Get owned channels (from chat list)
	allChats, err := h.chatService.ListChats(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	ownedChannels := make([]*chatdomain.ChatResponse, 0)
	subMap := make(map[string]bool)
	for _, ch := range subscribed {
		subMap[ch.ID] = true
	}
	for _, chat := range allChats {
		if chat.Type == chatdomain.ChatChannel {
			if !subMap[chat.ID] {
				ownedChannels = append(ownedChannels, chat)
			}
		}
	}

	result := append(ownedChannels, subscribed...)
	response.JSON(c, 200, result)
}

// SetSubscriberRole changes a subscriber's role (admin only)
// @Tags Channels
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Channel ID"
// @Param userId path string true "Target user ID"
// @Param request body object{role=string} true "Role: admin or subscriber"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /channels/{id}/subscribers/{userId}/role [put]
func (h *ChannelHandler) SetSubscriberRole(c *gin.Context) {
	userID, _ := c.Get("userID")
	channelID := c.Param("id")
	targetUserID := c.Param("userId")

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin subscriber"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.channelService.SetSubscriberRole(channelID, targetUserID, userID.(string), req.Role); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "role updated"})
}

// IsSubscribed checks if user is subscribed to a channel
// @Tags Channels
// @Security BearerAuth
// @Produce json
// @Param id path string true "Channel ID"
// @Success 200 {object} response.MessageResponse
// @Router /channels/{id}/subscribed [get]
func (h *ChannelHandler) IsSubscribed(c *gin.Context) {
	userID, _ := c.Get("userID")
	channelID := c.Param("id")

	subscribed, _ := h.channelService.IsSubscribed(channelID, userID.(string))
	response.JSON(c, 200, gin.H{"subscribed": subscribed})
}
