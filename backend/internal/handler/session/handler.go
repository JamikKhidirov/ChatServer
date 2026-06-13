package sessionhandler

import (
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

// GetSessions returns all active sessions for the authenticated user
// @Tags Сессии
// @Summary Получить активные сессии
// @Security BearerAuth
// @Produce json
// @Description Возвращает список всех активных сессий аутентифицированного пользователя, включая информацию об устройстве и времени входа.
// @Success 200 {array} sessiondomain.Session "Список активных сессий"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения сессий"
// @Router /sessions [get]
func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID, _ := c.Get("userID")

	sessions, err := h.sessionService.GetSessions(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, sessions)
}

// DeleteSession terminates a specific session
// @Tags Сессии
// @Summary Завершить сессию
// @Security BearerAuth
// @Produce json
// @Description Завершает указанную активную сессию пользователя. Пользователь будет вынужден авторизоваться заново на этом устройстве.
// @Param id path string true "ID сессии"
// @Success 200 {object} response.MessageResponse "Сессия завершена"
// @Failure 400 {object} response.ErrorResponse "Ошибка завершения сессии"
// @Router /sessions/{id} [delete]
func (h *SessionHandler) DeleteSession(c *gin.Context) {
	userID, _ := c.Get("userID")
	sessionID := c.Param("id")

	if err := h.sessionService.DeleteSession(userID.(string), sessionID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "session terminated"})
}

// DeleteAllSessions terminates all sessions except current
// @Tags Сессии
// @Summary Завершить все сессии
// @Security BearerAuth
// @Produce json
// @Description Завершает все активные сессии пользователя, кроме текущей. Полезно при подозрении на несанкционированный доступ.
// @Success 200 {object} response.MessageResponse "Все остальные сессии завершены"
// @Failure 400 {object} response.ErrorResponse "Ошибка завершения сессий"
// @Router /sessions [delete]
func (h *SessionHandler) DeleteAllSessions(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := h.sessionService.DeleteAllSessions(userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "all other sessions terminated"})
}
