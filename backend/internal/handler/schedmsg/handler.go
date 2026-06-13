package schedmsghandler

import (
	draftdomain "ChatServerGolang/backend/internal/domain/draft"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ScheduledMessageHandler struct {
	schedService service.ScheduledMessageService
}

func NewScheduledMessageHandler(schedService service.ScheduledMessageService) *ScheduledMessageHandler {
	return &ScheduledMessageHandler{schedService: schedService}
}

// Schedule schedules a message for later delivery
// @Summary Запланировать сообщение
// @Description Планирует отправку сообщения на указанное время. Сообщение будет доставлено автоматически.
// @Tags ScheduledMessages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body draftdomain.ScheduleMessageRequest true "Chat, content, and send time"
// @Success 201 {object} draftdomain.ScheduledMessage
// @Failure 400 {object} response.ErrorResponse
// @Router /messages/schedule [post]
func (h *ScheduledMessageHandler) Schedule(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req draftdomain.ScheduleMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	msg, err := h.schedService.Schedule(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 201, msg)
}

// GetScheduled returns all scheduled messages for the user
// @Summary Запланированные сообщения
// @Description Возвращает список всех запланированных сообщений текущего пользователя.
// @Tags ScheduledMessages
// @Security BearerAuth
// @Produce json
// @Success 200 {array} draftdomain.ScheduledMessage
// @Failure 400 {object} response.ErrorResponse
// @Router /messages/scheduled [get]
func (h *ScheduledMessageHandler) GetScheduled(c *gin.Context) {
	userID, _ := c.Get("userID")
	messages, err := h.schedService.GetScheduled(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, messages)
}

// CancelScheduled cancels a scheduled message
// @Summary Отменить запланированное
// @Description Отменяет запланированную отправку сообщения. Сообщение не будет доставлено.
// @Tags ScheduledMessages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Scheduled message ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /messages/scheduled/{id} [delete]
func (h *ScheduledMessageHandler) CancelScheduled(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	if err := h.schedService.CancelScheduled(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "scheduled message cancelled"})
}
