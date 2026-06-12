package groupcallhandler

import (
	calldomain "ChatServerGolang/internal/domain/call"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type GroupCallHandler struct {
	groupCallService service.GroupCallService
}

func NewGroupCallHandler(groupCallService service.GroupCallService) *GroupCallHandler {
	return &GroupCallHandler{groupCallService: groupCallService}
}

// InitiateGroupCall starts a group voice/video call
// @Tags Group Calls
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body calldomain.GroupCallInitiateRequest true "Chat ID and call type"
// @Success 201 {object} calldomain.GroupCallResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /calls/group/initiate [post]
func (h *GroupCallHandler) InitiateGroupCall(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req calldomain.GroupCallInitiateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	call, err := h.groupCallService.InitiateGroupCall(req.ChatID, userID.(string), req.Type)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, call)
}

// JoinGroupCall joins an active group call
// @Tags Group Calls
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body calldomain.GroupCallActionRequest true "Call ID and action"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /calls/group/respond [post]
func (h *GroupCallHandler) JoinGroupCall(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req calldomain.GroupCallActionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	switch req.Action {
	case "join":
		if err := h.groupCallService.JoinGroupCall(req.CallID, userID.(string)); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	case "leave":
		if err := h.groupCallService.LeaveGroupCall(req.CallID, userID.(string)); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	case "mute":
		if err := h.groupCallService.MuteParticipant(req.CallID, userID.(string), true, false); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	case "unmute_audio":
		if err := h.groupCallService.MuteParticipant(req.CallID, userID.(string), false, false); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	case "mute_video":
		if err := h.groupCallService.MuteParticipant(req.CallID, userID.(string), false, true); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	case "unmute_video":
		if err := h.groupCallService.MuteParticipant(req.CallID, userID.(string), false, false); err != nil {
			response.BadRequest(c, err.Error())
			return
		}
	}

	response.JSON(c, 200, gin.H{"message": "action " + req.Action + " completed"})
}

// EndGroupCall ends an active group call (caller only)
// @Tags Group Calls
// @Security BearerAuth
// @Produce json
// @Param id path string true "Call ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /calls/group/{id}/end [post]
func (h *GroupCallHandler) EndGroupCall(c *gin.Context) {
	userID, _ := c.Get("userID")
	callID := c.Param("id")

	if err := h.groupCallService.EndGroupCall(callID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "call ended"})
}

// GetGroupCall returns group call details
// @Tags Group Calls
// @Security BearerAuth
// @Produce json
// @Param id path string true "Call ID"
// @Success 200 {object} calldomain.GroupCallResponse
// @Failure 404 {object} response.ErrorResponse
// @Router /calls/group/{id} [get]
func (h *GroupCallHandler) GetGroupCall(c *gin.Context) {
	callID := c.Param("id")

	call, err := h.groupCallService.GetGroupCallByID(callID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, call)
}

// GetActiveGroupCalls returns active calls in a chat
// @Tags Group Calls
// @Security BearerAuth
// @Produce json
// @Param chatId path string true "Chat ID"
// @Success 200 {array} calldomain.GroupCallResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{chatId}/active-calls [get]
func (h *GroupCallHandler) GetActiveGroupCalls(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("chatId")

	calls, err := h.groupCallService.GetActiveGroupCalls(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, calls)
}
