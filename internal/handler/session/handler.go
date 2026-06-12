package sessionhandler

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

// GetSessions returns all active sessions for the authenticated user
// @Tags Sessions
// @Security BearerAuth
// @Produce json
// @Success 200 {array} sessiondomain.Session
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Sessions
// @Security BearerAuth
// @Produce json
// @Param id path string true "Session ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Sessions
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /sessions [delete]
func (h *SessionHandler) DeleteAllSessions(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := h.sessionService.DeleteAllSessions(userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "all other sessions terminated"})
}
