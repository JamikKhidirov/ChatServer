package wshandler

import (
	"net/http"

	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/ws"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WSHandler struct {
	hub         *ws.Hub
	authService service.AuthService
	userRepo    repository.UserRepository
	chatRepo    repository.ChatRepository
}

func NewWSHandler(hub *ws.Hub, authService service.AuthService, userRepo repository.UserRepository, chatRepo repository.ChatRepository) *WSHandler {
	return &WSHandler{hub: hub, authService: authService, userRepo: userRepo, chatRepo: chatRepo}
}

// HandleWebSocket WebSocket endpoint for real-time communication
// @Summary WebSocket соединение
// @Description Устанавливает WebSocket соединение для получения событий в реальном времени. Токен передаётся в query-параметре. После подключения клиент отправляет JSON-команды. События: newMessage, editMessage, deleteMessage, reaction, readMessage, pinMessage, callOffer, callAccept, callEnd, onlineStatus.
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param token query string true "JWT token"
// @Success 101 {string} string "Upgraded to WebSocket"
// @Failure 401 {object} response.ErrorResponse
// @Router /ws [get]
func (h *WSHandler) HandleWebSocket(c *gin.Context) {
	token := c.Query("token")
	if token == "" {
		c.JSON(401, gin.H{"error": "token required"})
		return
	}

	userID, err := h.authService.ValidateToken(token)
	if err != nil {
		c.JSON(401, gin.H{"error": "invalid token"})
		return
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return
	}

	go func() {
		h.userRepo.SetOnline(userID, true)
		h.broadcastOnlineStatus(userID, true)
	}()

	client := ws.NewClient(h.hub, conn, userID, func() {
		h.userRepo.SetOnline(userID, false)
		h.broadcastOnlineStatus(userID, false)
	})
	client.Start()
}

func (h *WSHandler) broadcastOnlineStatus(userID string, online bool) {
	chats, err := h.chatRepo.FindByUserID(userID)
	if err != nil {
		return
	}

	for _, chat := range chats {
		participants, _ := h.chatRepo.GetParticipants(chat.ID)
		for _, p := range participants {
			if p.UserID != userID {
				h.hub.SendToUser(p.UserID, ws.WSOutgoingMessage{
					Type: map[bool]ws.MessageType{true: ws.MsgOnline, false: ws.MsgOffline}[online],
					Payload: map[string]interface{}{
						"userId": userID,
						"online": online,
					},
				})
			}
		}
	}
}


