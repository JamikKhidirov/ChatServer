package groupcallhandler

import (
	calldomain "ChatServerGolang/backend/internal/domain/call"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type GroupCallHandler struct {
	groupCallService service.GroupCallService
}

func NewGroupCallHandler(groupCallService service.GroupCallService) *GroupCallHandler {
	return &GroupCallHandler{groupCallService: groupCallService}
}

// InitiateGroupCall starts a group voice/video call
// @Tags Звонки
// @Summary Начать групповой звонок
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Инициирует групповой голосовой или видеозвонок в указанном чате. Все участники чата получают уведомление.
// @Param request body calldomain.GroupCallInitiateRequest true "Параметры группы: chat_id (ID чата, обязательно), type (тип: voice/video, обязательно)"
// @Success 201 {object} calldomain.GroupCallResponse "Групповой звонок начат"
// @Failure 400 {object} response.ErrorResponse "Ошибка инициализации группового звонка"
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
// @Tags Звонки
// @Summary Присоединиться к групповому звонку
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Позволяет присоединиться к активному групповому звонку, выйти из него, а также управлять микрофоном и камерой.
// @Param request body calldomain.GroupCallActionRequest true "Действие: call_id (ID звонка, обязательно), action (join/leave/mute/unmute_audio/mute_video/unmute_video, обязательно)"
// @Success 200 {object} response.MessageResponse "Действие выполнено"
// @Failure 400 {object} response.ErrorResponse "Ошибка выполнения действия"
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
// @Tags Звонки
// @Summary Завершить групповой звонок
// @Security BearerAuth
// @Produce json
// @Description Завершает активный групповой звонок. Доступно только инициатору звонка.
// @Param id path string true "ID звонка"
// @Success 200 {object} response.MessageResponse "Групповой звонок завершён"
// @Failure 400 {object} response.ErrorResponse "Ошибка завершения звонка"
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
// @Tags Звонки
// @Summary Получить детали группового звонка
// @Security BearerAuth
// @Produce json
// @Description Возвращает информацию о групповом звонке: участники, статус аудио/видео, длительность.
// @Param id path string true "ID звонка"
// @Success 200 {object} calldomain.GroupCallResponse "Детали группового звонка"
// @Failure 404 {object} response.ErrorResponse "Звонок не найден"
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
// @Tags Звонки
// @Summary Получить активные групповые звонки
// @Security BearerAuth
// @Produce json
// @Description Возвращает список активных групповых звонков в указанном чате.
// @Param id path string true "ID чата"
// @Success 200 {array} calldomain.GroupCallResponse "Список активных групповых звонков"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения звонков"
// @Router /chats/{chatId}/active-calls [get]
func (h *GroupCallHandler) GetActiveGroupCalls(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	calls, err := h.groupCallService.GetActiveGroupCalls(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, calls)
}
