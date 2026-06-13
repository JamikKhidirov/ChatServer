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
// @Description Устанавливает WebSocket соединение для real-time коммуникации. Токен передаётся в query-параметре.
// @Description 
// @Description ## Формат сообщений
// @Description Все сообщения — JSON: {"type": "событие", "payload": { ... }}
// @Description 
// @Description ## Исходящие события (сервер → клиент)
// @Description | Тип | Payload | Описание |
// @Description |------|---------|---------|
// @Description | `message:new` | `Message` | Новое сообщение (текст, изображение, файл и т.д.) |
// @Description | `message:edited` | `Message` | Сообщение отредактировано |
// @Description | `message:deleted` | `{messageId, chatId}` | Сообщение удалено |
// @Description | `message:read` | `{messageId, userId, chatId}` | Сообщение прочитано |
// @Description | `message:reaction` | `Message` | Добавлена/удалена реакция |
// @Description | `message:pinned` | `Message` | Сообщение закреплено/откреплено |
// @Description | `message:starred` | — | Сообщение добавлено в избранное |
// @Description | `message:forward` | — | Сообщение переслано |
// @Description | `user:online` | `{userId, online: true}` | Пользователь онлайн |
// @Description | `user:offline` | `{userId, online: false}` | Пользователь офлайн |
// @Description | `user:typing` | `{chatId, userId}` | Пользователь печатает |
// @Description | `user:stop_typing` | `{chatId, userId}` | Пользователь перестал печатать |
// @Description | `user:keyboard_opened` | `{chatId, userId}` | Клавиатура открыта |
// @Description | `user:keyboard_closed` | `{chatId, userId}` | Клавиатура закрыта |
// @Description | `chat:created` | `Chat` | Создан новый чат |
// @Description | `chat:updated` | `Chat` | Чат обновлён |
// @Description | `chat:deleted` | — | Чат удалён |
// @Description | `call:offer` | `{chatId, callId, sdp}` | Входящий звонок |
// @Description | `call:answer` | `{chatId, callId, sdp}` | Ответ на звонок |
// @Description | `call:ice` | `{callId, candidate}` | ICE-кандидат для WebRTC |
// @Description | `call:end` | `{callId, userId}` | Звонок завершён |
// @Description | `call:missed` | — | Пропущенный звонок |
// @Description | `call:accept` | — | Звонок принят |
// @Description | `call:reject` | `{callId, userId}` | Звонок отклонён |
// @Description 
// @Description ## Входящие события (клиент → сервер)
// @Description | Тип | Payload | Описание |
// @Description |------|---------|---------|
// @Description | `message:send` | `{chatId, content, type, replyToId?}` | Отправить сообщение |
// @Description | `message:edit` | `{messageId, content}` | Редактировать сообщение |
// @Description | `message:delete` | `{messageId, chatId}` | Удалить сообщение |
// @Description | `message:read` | `{messageId, chatId}` | Отметить как прочитанное |
// @Description | `message:react` | `{messageId, emoji}` | Добавить реакцию |
// @Description | `message:unreact` | `{messageId, emoji}` | Удалить реакцию |
// @Description | `message:pin` | `{messageId, pin: bool}` | Закрепить/открепить |
// @Description | `message:star` | `{messageId}` | В избранное |
// @Description | `message:unstar` | `{messageId}` | Из избранного |
// @Description | `message:forward` | `{messageId, toChatId}` | Переслать |
// @Description | `chat:create` | `{type, name?, participantIds, description?}` | Создать чат |
// @Description | `chat:update` | `{chatId, name?, description?, avatarUrl?}` | Обновить чат |
// @Description | `chat:add_participant` | `{chatId, userId}` | Добавить участника |
// @Description | `chat:remove_participant` | `{chatId, userId}` | Удалить участника |
// @Description | `chat:leave` | `{chatId}` | Покинуть чат |
// @Description | `chat:pin` | `{chatId}` | Закрепить чат |
// @Description | `chat:unpin` | `{chatId}` | Открепить чат |
// @Description | `chat:archive` | `{chatId}` | Архивировать чат |
// @Description | `chat:unarchive` | `{chatId}` | Разархивировать чат |
// @Description | `user:typing` | `{chatId}` | Индикатор печатания |
// @Description | `user:stop_typing` | `{chatId}` | Прекратил печатать |
// @Description | `user:keyboard_opened` | `{chatId}` | Клавиатура открыта |
// @Description | `user:keyboard_closed` | `{chatId}` | Клавиатура закрыта |
// @Description | `user:block` | `{userId}` | Заблокировать пользователя |
// @Description | `user:unblock` | `{userId}` | Разблокировать пользователя |
// @Description | `call:offer` | `{chatId, callId, sdp}` | WebRTC offer |
// @Description | `call:answer` | `{chatId, callId, sdp}` | WebRTC answer |
// @Description | `call:ice` | `{callId, candidate}` | WebRTC ICE candidate |
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


