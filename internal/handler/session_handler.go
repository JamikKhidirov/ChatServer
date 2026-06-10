package handler

import (
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type SessionHandler struct {
	sessionService service.SessionService
}

func NewSessionHandler(sessionService service.SessionService) *SessionHandler {
	return &SessionHandler{sessionService: sessionService}
}

func (h *SessionHandler) GetSessions(c *gin.Context) {
	userID, _ := c.Get("userID")
	sessions, err := h.sessionService.GetSessions(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, sessions)
}

func (h *SessionHandler) DeleteSession(c *gin.Context) {
	userID, _ := c.Get("userID")
	sessionID := c.Param("id")
	if err := h.sessionService.DeleteSession(sessionID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "session deleted"})
}

func (h *SessionHandler) DeleteAllSessions(c *gin.Context) {
	userID, _ := c.Get("userID")
	if err := h.sessionService.DeleteAllSessions(userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "all sessions deleted"})
}
