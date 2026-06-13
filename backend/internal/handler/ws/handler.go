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
// @Description ## Подключение
// @Description ws://localhost:8080/ws?token={JWT_TOKEN}
// @Description Токен можно получить через POST /api/auth/register или POST /api/auth/login.
// @Description 
// @Description ## Формат сообщений
// @Description ```json
// @Description { "type": "название_события", "payload": { ... поля ... } }
// @Description ```
// @Description type — строка, указывающая тип события. payload — объект с данными события.
// @Description 
// @Description ## Направление: Сервер → Клиент
// @Description Сервер автоматически отправляет эти события всем подключённым участникам чата.
// @Description 
// @Description ### Сообщения (типы)
// @Description **`message:new`** — Новое сообщение отправлено в чат.
// @Description ```json
// @Description { "type": "message:new", "payload": { "id": "uuid", "chatId": "uuid", "senderId": "uuid", "content": "Привет!", "type": "text", "createdAt": "2026-01-01T00:00:00Z" } }
// @Description ```
// @Description 
// @Description **`message:edited`** — Сообщение отредактировано. Payload: объект Message.
// @Description 
// @Description **`message:deleted`** — Сообщение удалено.
// @Description ```json
// @Description { "type": "message:deleted", "payload": { "messageId": "uuid", "chatId": "uuid" } }
// @Description ```
// @Description 
// @Description **`message:read`** — Сообщение прочитано.
// @Description ```json
// @Description { "type": "message:read", "payload": { "messageId": "uuid", "userId": "uuid", "chatId": "uuid" } }
// @Description ```
// @Description 
// @Description **`message:reaction`** — Добавлена/удалена реакция. Payload: обновлённый объект Message с полем reactions.
// @Description 
// @Description **`message:pinned`** — Сообщение закреплено/откреплено. Payload: Message.
// @Description 
// @Description ### Пользователи (типы user:*)
// @Description **`user:online`** — Пользователь стал онлайн.
// @Description ```json
// @Description { "type": "user:online", "payload": { "userId": "uuid", "online": true } }
// @Description ```
// @Description 
// @Description **`user:offline`** — Пользователь стал офлайн.
// @Description ```json
// @Description { "type": "user:offline", "payload": { "userId": "uuid", "online": false } }
// @Description ```
// @Description 
// @Description **`user:typing`** — Пользователь печатает.
// @Description ```json
// @Description { "type": "user:typing", "payload": { "chatId": "uuid", "userId": "uuid" } }
// @Description ```
// @Description 
// @Description **`user:stop_typing`** — Перестал печатать.
// @Description ```json
// @Description { "type": "user:stop_typing", "payload": { "chatId": "uuid", "userId": "uuid" } }
// @Description ```
// @Description 
// @Description **`user:keyboard_opened`** — Клавиатура открыта (мобильные устройства).
// @Description **`user:keyboard_closed`** — Клавиатура закрыта.
// @Description 
// @Description ### Чаты (типы chat:*)
// @Description **`chat:created`** — Создан новый чат. Payload: объект Chat (id, name, type, participants и т.д.).
// @Description **`chat:updated`** — Чат обновлён. Payload: Chat.
// @Description **`chat:deleted`** — Чат удалён.
// @Description 
// @Description ### Звонки (типы call:*)
// @Description **`call:offer`** — Входящий WebRTC звонок.
// @Description ```json
// @Description { "type": "call:offer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "offer_sdp_string" } }
// @Description ```
// @Description 
// @Description **`call:answer`** — Ответ на звонок.
// @Description ```json
// @Description { "type": "call:answer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "answer_sdp_string" } }
// @Description ```
// @Description 
// @Description **`call:ice`** — ICE-кандидат для WebRTC.
// @Description ```json
// @Description { "type": "call:ice", "payload": { "callId": "uuid", "candidate": "ice_candidate_string" } }
// @Description ```
// @Description 
// @Description **`call:end`** — Звонок завершён.
// @Description ```json
// @Description { "type": "call:end", "payload": { "callId": "uuid", "userId": "uuid" } }
// @Description ```
// @Description 
// @Description **`call:reject`** — Звонок отклонён.
// @Description ```json
// @Description { "type": "call:reject", "payload": { "callId": "uuid", "userId": "uuid" } }
// @Description ```
// @Description 
// @Description **`call:missed`** — Пропущенный звонок. **`call:accept`** — Звонок принят.
// @Description 
// @Description ## Направление: Клиент → Сервер
// @Description Клиент отправляет эти события, чтобы выполнить действие на сервере.
// @Description 
// @Description ### Отправка и управление сообщениями
// @Description **`message:send`** — Отправить сообщение в чат.
// @Description ```json
// @Description { "type": "message:send", "payload": { "chatId": "uuid", "content": "Привет!", "type": "text", "replyToId": "uuid" } }
// @Description ```
// @Description Поля: chatId (обяз.), content (обяз.), type (обяз.): text|image|file|gif|voice|video|audio|location|system, replyToId (опц.) — ID сообщения, на который отвечаем.
// @Description 
// @Description **`message:edit`** — Редактировать сообщение.
// @Description ```json
// @Description { "type": "message:edit", "payload": { "messageId": "uuid", "content": "Новый текст" } }
// @Description ```
// @Description 
// @Description **`message:delete`** — Удалить сообщение.
// @Description ```json
// @Description { "type": "message:delete", "payload": { "messageId": "uuid", "chatId": "uuid" } }
// @Description ```
// @Description 
// @Description **`message:read`** — Отметить сообщение как прочитанное.
// @Description ```json
// @Description { "type": "message:read", "payload": { "messageId": "uuid", "chatId": "uuid" } }
// @Description ```
// @Description 
// @Description **`message:react`** — Добавить реакцию.
// @Description ```json
// @Description { "type": "message:react", "payload": { "messageId": "uuid", "emoji": "👍" } }
// @Description ```
// @Description emoji: "👍", "❤️", "😆", "😮", "😢", "🙏"
// @Description 
// @Description **`message:unreact`** — Удалить реакцию.
// @Description ```json
// @Description { "type": "message:unreact", "payload": { "messageId": "uuid", "emoji": "👍" } }
// @Description ```
// @Description 
// @Description **`message:pin`** — Закрепить/открепить сообщение.
// @Description ```json
// @Description { "type": "message:pin", "payload": { "messageId": "uuid", "pin": true } }
// @Description ```
// @Description pin: true — закрепить, false — открепить.
// @Description 
// @Description **`message:star`** — Добавить в избранное. Payload: { "messageId": "uuid" }
// @Description **`message:unstar`** — Удалить из избранного. Payload: { "messageId": "uuid" }
// @Description **`message:forward`** — Переслать сообщение. Payload: { "messageId": "uuid", "toChatId": "uuid" }
// @Description 
// @Description ### Управление чатами (клиент → сервер)
// @Description **`chat:create`** — Создать новый чат.
// @Description ```json
// @Description { "type": "chat:create", "payload": { "type": "group", "name": "Friends", "participantIds": ["uuid1","uuid2"], "description": "Чат для друзей" } }
// @Description ```
// @Description type: "private" | "group" | "channel". participantIds — список ID участников (обяз.). name — название (опц. для private).
// @Description 
// @Description **`chat:update`** — Обновить название/аватар/описание чата.
// @Description ```json
// @Description { "type": "chat:update", "payload": { "chatId": "uuid", "name": "Новое название", "description": "Описание", "avatarUrl": "https://..." } }
// @Description ```
// @Description 
// @Description **`chat:add_participant`** — Добавить участника. Payload: { "chatId": "uuid", "userId": "uuid" }
// @Description **`chat:remove_participant`** — Удалить участника. Payload: { "chatId": "uuid", "userId": "uuid" }
// @Description **`chat:leave`** — Покинуть чат. Payload: { "chatId": "uuid" }
// @Description **`chat:pin`** / **`chat:unpin`** — Закрепить/открепить чат в списке. Payload: { "chatId": "uuid" }
// @Description **`chat:archive`** / **`chat:unarchive`** — Архивировать/разархивировать чат. Payload: { "chatId": "uuid" }
// @Description 
// @Description ### Статус пользователя (клиент → сервер)
// @Description **`user:typing`** — Отправить индикатор печатания.
// @Description ```json
// @Description { "type": "user:typing", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description **`user:stop_typing`** — Остановить индикатор.
// @Description ```json
// @Description { "type": "user:stop_typing", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description **`user:keyboard_opened`** / **`user:keyboard_closed`** — Клавиатура открыта/закрыта. Payload: { "chatId": "uuid" }
// @Description **`user:block`** — Заблокировать пользователя. Payload: { "userId": "uuid" }
// @Description **`user:unblock`** — Разблокировать пользователя. Payload: { "userId": "uuid" }
// @Description 
// @Description ### WebRTC звонки (клиент → сервер)
// @Description **`call:offer`** — Отправить WebRTC offer. Payload: { "chatId": "uuid", "callId": "uuid", "sdp": "..." }
// @Description **`call:answer`** — WebRTC answer. Payload: { "chatId": "uuid", "callId": "uuid", "sdp": "..." }
// @Description **`call:ice`** — ICE candidate. Payload: { "callId": "uuid", "candidate": "..." }
// @Description 
// @Description ## Пример обмена сообщениями
// @Description 1. Клиент А отправляет: { "type": "message:send", "payload": { "chatId": "123", "content": "Привет!", "type": "text" } }
// @Description 2. Сервер принимает сообщение, сохраняет в БД, рассылает всем участникам чата:
// @Description    { "type": "message:new", "payload": { "id": "msg-uuid", "chatId": "123", "senderId": "userA-uuid", "content": "Привет!", "type": "text", "createdAt": "..." } }
// @Description 3. Клиент Б получает `message:new` и отображает сообщение в реальном времени.
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


