package handler

import (
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type CallHandler struct {
	callService service.CallService
}

func NewCallHandler(callService service.CallService) *CallHandler {
	return &CallHandler{callService: callService}
}

type InitiateCallRequest struct {
	ChatID string `json:"chatId" binding:"required"`
}

type RespondCallRequest struct {
	Action string `json:"action" binding:"required,oneof=accept reject"`
}

func (h *CallHandler) InitiateCall(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req InitiateCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	call, err := h.callService.InitiateCall(req.ChatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Set("callResponse", call)
	response.JSON(c, 201, call)
}

func (h *CallHandler) RespondCall(c *gin.Context) {
	userID, _ := c.Get("userID")
	callID := c.Param("id")

	var req RespondCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var err error
	switch req.Action {
	case "accept":
		err = h.callService.AcceptCall(callID, userID.(string))
	case "reject":
		err = h.callService.RejectCall(callID, userID.(string))
	}

	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "call " + req.Action + "ed"})
}

func (h *CallHandler) EndCall(c *gin.Context) {
	userID, _ := c.Get("userID")
	callID := c.Param("id")

	if err := h.callService.EndCall(callID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "call ended"})
}

func (h *CallHandler) GetCall(c *gin.Context) {
	callID := c.Param("id")

	call, err := h.callService.GetCallByID(callID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, call)
}

func (h *CallHandler) GetCallHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("chatId")

	calls, err := h.callService.GetCallHistory(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, calls)
}
