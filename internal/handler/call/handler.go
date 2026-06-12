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
// @Tags Звонки
// @Summary Начать звонок
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Инициирует голосовой или видеозвонок в указанном чате. Тип звонка указывается в параметрах.
// @Param request body calldomain.InitiateCallRequest true "Параметры звонка: chat_id (ID чата, обязательно), type (тип: voice/video, обязательно)"
// @Success 201 {object} calldomain.Call "Звонок начат"
// @Failure 400 {object} response.ErrorResponse "Ошибка инициализации звонка"
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
// @Tags Звонки
// @Summary Ответить на звонок
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Принимает или отклоняет входящий звонок. Действие указывается в параметрах: accept (принять) или reject (отклонить).
// @Param id path string true "ID звонка"
// @Param request body calldomain.RespondCallRequest true "Действие: action (строка: accept — принять, reject — отклонить, обязательно)"
// @Success 200 {object} response.MessageResponse "Ответ на звонок отправлен"
// @Failure 400 {object} response.ErrorResponse "Ошибка ответа на звонок"
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
// @Tags Звонки
// @Summary Завершить звонок
// @Security BearerAuth
// @Produce json
// @Description Завершает активный звонок. Любой участник звонка может его завершить.
// @Param id path string true "ID звонка"
// @Success 200 {object} response.MessageResponse "Звонок завершён"
// @Failure 400 {object} response.ErrorResponse "Ошибка завершения звонка"
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
// @Tags Звонки
// @Summary Получить детали звонка
// @Security BearerAuth
// @Produce json
// @Description Возвращает информацию о звонке по его идентификатору: участники, статус, длительность.
// @Param id path string true "ID звонка"
// @Success 200 {object} calldomain.Call "Детали звонка"
// @Failure 404 {object} response.ErrorResponse "Звонок не найден"
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
// @Tags Звонки
// @Summary Получить историю звонков
// @Security BearerAuth
// @Produce json
// @Description Возвращает историю всех звонков, совершённых в указанном чате.
// @Param chatId path string true "ID чата"
// @Success 200 {array} calldomain.Call "История звонков"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения истории"
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
