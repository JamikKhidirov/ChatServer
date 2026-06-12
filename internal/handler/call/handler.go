package callhandler

import (
	"ChatServerGolang/internal/domain/call"
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

// InitiateCall initiates a voice or video call in a chat
// @Tags Calls
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body calldomain.InitiateCallRequest true "Chat ID and call type (voice/video)"
// @Success 201 {object} calldomain.Call
// @Failure 400 {object} response.ErrorResponse
// @Router /calls [post]
func (h *CallHandler) InitiateCall(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req calldomain.InitiateCallRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	call, err := h.callService.InitiateCall(req.ChatID, userID.(string), req.Type)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Set("callResponse", call)
	response.JSON(c, 201, call)
}

// RespondCall accepts or rejects an incoming call
// @Tags Calls
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Call ID"
// @Param request body calldomain.RespondCallRequest true "Action: accept or reject"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /calls/{id}/respond [post]
func (h *CallHandler) RespondCall(c *gin.Context) {
	userID, _ := c.Get("userID")
	callID := c.Param("id")

	var req calldomain.RespondCallRequest
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

// EndCall terminates an active call
// @Tags Calls
// @Security BearerAuth
// @Produce json
// @Param id path string true "Call ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /calls/{id}/end [post]
func (h *CallHandler) EndCall(c *gin.Context) {
	userID, _ := c.Get("userID")
	callID := c.Param("id")

	if err := h.callService.EndCall(callID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "call ended"})
}

// GetCall returns details of a specific call
// @Tags Calls
// @Security BearerAuth
// @Produce json
// @Param id path string true "Call ID"
// @Success 200 {object} calldomain.Call
// @Failure 404 {object} response.ErrorResponse "Not found"
// @Router /calls/{id} [get]
func (h *CallHandler) GetCall(c *gin.Context) {
	callID := c.Param("id")

	call, err := h.callService.GetCallByID(callID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, call)
}

// GetCallHistory returns call history for a chat
// @Tags Calls
// @Security BearerAuth
// @Produce json
// @Param chatId path string true "Chat ID"
// @Success 200 {array} calldomain.Call
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{chatId}/calls [get]
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
