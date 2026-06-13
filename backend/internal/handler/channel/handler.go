package channelhandler

import (
	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

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
// @Tags Чаты
// @Summary Подписаться на канал
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Подписывает пользователя на указанный канал для получения обновлений и новых постов.
// @Param request body object{channelId=string} true "Параметры: channelId (ID канала, обязательно)"
// @Success 200 {object} response.MessageResponse "Подписка оформлена"
// @Failure 400 {object} response.ErrorResponse "Ошибка подписки"
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
// @Tags Чаты
// @Summary Отписаться от канала
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Отменяет подписку пользователя на указанный канал. Пользователь перестанет получать обновления.
// @Param request body object{channelId=string} true "Параметры: channelId (ID канала, обязательно)"
// @Success 200 {object} response.MessageResponse "Подписка отменена"
// @Failure 400 {object} response.ErrorResponse "Ошибка отмены подписки"
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
// @Tags Чаты
// @Summary Получить подписчиков канала
// @Security BearerAuth
// @Produce json
// @Description Возвращает список подписчиков канала. Доступно только администраторам канала.
// @Param id path string true "ID канала"
// @Success 200 {array} channeldomain.ChannelSubscriber "Список подписчиков"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения подписчиков"
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
// @Tags Чаты
// @Summary Получить мои каналы
// @Security BearerAuth
// @Produce json
// @Description Возвращает список каналов, созданных пользователем или на которые он подписан.
// @Success 200 {array} chatdomain.ChatResponse "Список каналов"
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
// @Tags Чаты
// @Summary Изменить роль подписчика канала
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Изменяет роль подписчика канала: admin (администратор) или subscriber (подписчик). Доступно администраторам.
// @Param id path string true "ID канала"
// @Param userId path string true "ID целевого пользователя"
// @Param request body object{role=string} true "Новая роль: role (admin — администратор, subscriber — подписчик, обязательно)"
// @Success 200 {object} response.MessageResponse "Роль обновлена"
// @Failure 400 {object} response.ErrorResponse "Ошибка изменения роли"
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
// @Tags Чаты
// @Summary Проверить подписку на канал
// @Security BearerAuth
// @Produce json
// @Description Проверяет, подписан ли текущий пользователь на указанный канал.
// @Param id path string true "ID канала"
// @Success 200 {object} response.MessageResponse "Статус подписки: subscribed (boolean)"
// @Router /channels/{id}/subscribed [get]
func (h *ChannelHandler) IsSubscribed(c *gin.Context) {
	userID, _ := c.Get("userID")
	channelID := c.Param("id")

	subscribed, _ := h.channelService.IsSubscribed(channelID, userID.(string))
	response.JSON(c, 200, gin.H{"subscribed": subscribed})
}
