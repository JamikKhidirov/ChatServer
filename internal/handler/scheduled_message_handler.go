package handler

import (
	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type ScheduledMessageHandler struct {
	schedService service.ScheduledMessageService
}

func NewScheduledMessageHandler(schedService service.ScheduledMessageService) *ScheduledMessageHandler {
	return &ScheduledMessageHandler{schedService: schedService}
}

func (h *ScheduledMessageHandler) Schedule(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.ScheduleMessageRequest
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

func (h *ScheduledMessageHandler) GetScheduled(c *gin.Context) {
	userID, _ := c.Get("userID")
	messages, err := h.schedService.GetScheduled(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, messages)
}

func (h *ScheduledMessageHandler) CancelScheduled(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	if err := h.schedService.CancelScheduled(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "scheduled message cancelled"})
}
