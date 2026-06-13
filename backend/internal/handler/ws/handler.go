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
// @Description Устанавливает WebSocket соединение для real-time коммуникации.
// @Description 
// @Description Формат сообщений:
// @Description   Все сообщения — JSON: {"type": "событие", "payload": { ... }}
// @Description 
// @Description Исходящие события (сервер → клиент):
// @Description   message:new — Новое сообщение. Payload: Message
// @Description   message:edited — Сообщение отредактировано. Payload: Message
// @Description   message:deleted — Сообщение удалено. Payload: {messageId, chatId}
// @Description   message:read — Сообщение прочитано. Payload: {messageId, userId, chatId}
// @Description   message:reaction — Добавлена/удалена реакция. Payload: Message
// @Description   message:pinned — Сообщение закреплено/откреплено. Payload: Message
// @Description   message:starred — Сообщение добавлено в избранное
// @Description   message:forward — Сообщение переслано
// @Description   user:online — Пользователь онлайн. Payload: {userId, online: true}
// @Description   user:offline — Пользователь офлайн. Payload: {userId, online: false}
// @Description   user:typing — Пользователь печатает. Payload: {chatId, userId}
// @Description   user:stop_typing — Пользователь перестал печатать. Payload: {chatId, userId}
// @Description   user:keyboard_opened — Клавиатура открыта. Payload: {chatId, userId}
// @Description   user:keyboard_closed — Клавиатура закрыта. Payload: {chatId, userId}
// @Description   chat:created — Создан новый чат. Payload: Chat
// @Description   chat:updated — Чат обновлён. Payload: Chat
// @Description   chat:deleted — Чат удалён
// @Description   call:offer — Входящий звонок. Payload: {chatId, callId, sdp}
// @Description   call:answer — Ответ на звонок. Payload: {chatId, callId, sdp}
// @Description   call:ice — ICE-кандидат для WebRTC. Payload: {callId, candidate}
// @Description   call:end — Звонок завершён. Payload: {callId, userId}
// @Description   call:missed — Пропущенный звонок
// @Description   call:accept — Звонок принят
// @Description   call:reject — Звонок отклонён. Payload: {callId, userId}
// @Description 
// @Description Входящие события (клиент → сервер):
// @Description   message:send — Отправить сообщение. Payload: {chatId, content, type, replyToId?}
// @Description   message:edit — Редактировать сообщение. Payload: {messageId, content}
// @Description   message:delete — Удалить сообщение. Payload: {messageId, chatId}
// @Description   message:read — Отметить как прочитанное. Payload: {messageId, chatId}
// @Description   message:react — Добавить реакцию. Payload: {messageId, emoji}
// @Description   message:unreact — Удалить реакцию. Payload: {messageId, emoji}
// @Description   message:pin — Закрепить/открепить. Payload: {messageId, pin: bool}
// @Description   message:star — В избранное. Payload: {messageId}
// @Description   message:unstar — Из избранного. Payload: {messageId}
// @Description   message:forward — Переслать сообщение. Payload: {messageId, toChatId}
// @Description   chat:create — Создать чат. Payload: {type, name?, participantIds, description?}
// @Description   chat:update — Обновить чат. Payload: {chatId, name?, description?, avatarUrl?}
// @Description   chat:add_participant — Добавить участника. Payload: {chatId, userId}
// @Description   chat:remove_participant — Удалить участника. Payload: {chatId, userId}
// @Description   chat:leave — Покинуть чат. Payload: {chatId}
// @Description   chat:pin — Закрепить чат. Payload: {chatId}
// @Description   chat:unpin — Открепить чат. Payload: {chatId}
// @Description   chat:archive — Архивировать чат. Payload: {chatId}
// @Description   chat:unarchive — Разархивировать чат. Payload: {chatId}
// @Description   user:typing — Индикатор печатания. Payload: {chatId}
// @Description   user:stop_typing — Прекратил печатать. Payload: {chatId}
// @Description   user:keyboard_opened — Клавиатура открыта. Payload: {chatId}
// @Description   user:keyboard_closed — Клавиатура закрыта. Payload: {chatId}
// @Description   user:block — Заблокировать пользователя. Payload: {userId}
// @Description   user:unblock — Разблокировать пользователя. Payload: {userId}
// @Description   call:offer — WebRTC offer. Payload: {chatId, callId, sdp}
// @Description   call:answer — WebRTC answer. Payload: {chatId, callId, sdp}
// @Description   call:ice — WebRTC ICE candidate. Payload: {callId, candidate}
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


