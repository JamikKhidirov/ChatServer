package wshandler

import (
	"net/http"

	"ChatServerGolang/internal/service"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/ws"

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
// @Description # WebSocket API
// @Description
// @Description ## Быстрый старт
// @Description 1. Зарегистрируйся: `POST /api/auth/register` → получи JWT token
// @Description 2. Подключись: `ws://localhost:8080/ws?token={JWT_TOKEN}`
// @Description 3. Отправляй и получай события в формате `{ "type": "...", "payload": {...} }`
// @Description
// @# Description # Технические детали
// @Description - **Пинг/Pong**: сервер шлёт ping каждые 54с, клиент отвечает pong (встроено в ws-библиотеку)
// @Description - **Таймаут**: если pong не получен 60с — соединение закрывается
// @Description - **Макс. размер сообщения**: 64 КБ
// @Description - **Ограничение**: 100 запросов/мин на API (REST + WebSocket суммарно)
// @Description - **Реконнект**: при обрыве клиент должен переподключаться с тем же токеном
// @Description
// @Description ---
// @Description
// @Description ## 1. Сервер → Клиент (входящие события)
// @Description Сервер автоматически рассылает эти события участникам чата в реальном времени.
// @Description
// @Description ### 1.1 Сообщения
// @Description
// @Description **`message:new`** — Новое сообщение.
// @Description ```json
// @Description { "type": "message:new", "payload": { "id": "uuid", "chatId": "uuid", "senderId": "uuid", "content": "Привет!", "type": "text", "createdAt": "2026-01-01T00:00:00Z" } }
// @Description ```
// @Description Поля payload: id, chatId, senderId, content, type (text|image|file|gif|voice|video|audio|location), createdAt, replyTo (если ответ), forwardedFrom (если переслано).
// @Description
// @Description **`message:edited`** — Сообщение отредактировано. Payload: полный объект Message с обновлённым content + editHistory.
// @Description
// @Description **`message:deleted`** — Сообщение удалено.
// @Description ```json
// @Description { "type": "message:deleted", "payload": { "messageId": "uuid", "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`message:read`** — Сообщение прочитано другим участником.
// @Description ```json
// @Description { "type": "message:read", "payload": { "messageId": "uuid", "userId": "uuid", "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`message:reaction`** — Изменение реакций. Payload: полный объект Message с обновлённым полем `reactions`.
// @Description
// @Description **`message:pinned`** — Сообщение закреплено/откреплено. Payload: полный объект Message.
// @Description
// @Description **`message:starred`** — Сообщение добавлено в избранное. Payload: { "messageId", "userId" }
// @Description
// @Description **`message:forward`** — Сообщение переслано. Payload: полный объект Message.
// @Description
// @Description ### 1.2 Пользователи
// @Description
// @Description | Тип | Описание | Payload |
// @Description |-----|----------|---------|
// @Description | `user:online` | Пользователь стал онлайн | `{ "userId": "uuid", "online": true }` |
// @Description | `user:offline` | Пользователь стал офлайн | `{ "userId": "uuid", "online": false }` |
// @Description | `user:typing` | Пользователь печатает | `{ "chatId": "uuid", "userId": "uuid" }` |
// @Description | `user:stop_typing` | Перестал печатать | `{ "chatId": "uuid", "userId": "uuid" }` |
// @Description | `user:keyboard_opened` | Клавиатура открыта (моб.) | `{ "chatId": "uuid", "userId": "uuid" }` |
// @Description | `user:keyboard_closed` | Клавиатура закрыта | `{ "chatId": "uuid", "userId": "uuid" }` |
// @Description
// @Description ### 1.3 Чаты
// @Description
// @Description **`chat:created`** — Создан новый чат. Payload: полный объект Chat (id, name, type, participants, avatarUrl, description).
// @Description
// @Description **`chat:updated`** — Чат обновлён (название, аватар, участники). Payload: полный объект Chat.
// @Description
// @Description **`chat:deleted`** — Чат удалён.
// @Description ```json
// @Description { "type": "chat:deleted", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description ### 1.4 Звонки (WebRTC)
// @Description
// @Description | Тип | Описание | Payload |
// @Description |-----|----------|---------|
// @Description | `call:offer` | Входящий звонок (WebRTC offer) | `{ "chatId", "callId", "callerId", "type" }` |
// @Description | `call:answer` | Ответ на звонок | `{ "callId", "userId" }` |
// @Description | `call:ice` | ICE-кандидат | `{ "callId", "candidate" }` |
// @Description | `call:end` | Звонок завершён | `{ "callId", "userId" }` |
// @Description | `call:reject` | Звонок отклонён | `{ "callId", "userId" }` |
// @Description | `call:accept` | Звонок принят | `{ "callId", "userId" }` |
// @Description | `call:missed` | Пропущенный звонок | `{ "callId", "userId" }` |
// @Description
// @Description ---
// @Description
// @Description ## 2. Клиент → Сервер (исходящие события)
// @Description Отправляй эти события через WebSocket, чтобы выполнить действия.
// @Description
// @Description ### 2.1 Отправка и управление сообщениями
// @Description
// @Description **`message:send`** — Отправить сообщение в чат.
// @Description ```json
// @Description { "type": "message:send", "payload": { "chatId": "uuid", "content": "Привет!", "type": "text", "replyToId": "uuid" } }
// @Description ```
// @Description
// @Description | Поле | Обязательное | Описание |
// @Description |------|-------------|----------|
// @Description | chatId | ✅ | ID чата |
// @Description | content | ✅ | Текст / URL файла / base64 |
// @Description | type | ✅ | text, image, file, gif, voice, video, audio, location |
// @Description | replyToId | ❌ | ID сообщения, на которое отвечаем |
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
// @Description **`message:react`** — Добавить реакцию. emoji: "👍", "❤️", "😆", "😮", "😢", "🙏"
// @Description ```json
// @Description { "type": "message:react", "payload": { "messageId": "uuid", "emoji": "👍" } }
// @Description ```
// @Description
// @Description **`message:unreact`** — Удалить реакцию.
// @Description ```json
// @Description { "type": "message:unreact", "payload": { "messageId": "uuid", "emoji": "👍" } }
// @Description ```
// @Description
// @Description **`message:pin`** — Закрепить (pin: true) / открепить (pin: false) сообщение.
// @Description ```json
// @Description { "type": "message:pin", "payload": { "messageId": "uuid", "pin": true } }
// @Description ```
// @Description
// @Description **`message:star`** — Добавить в избранное.
// @Description ```json
// @Description { "type": "message:star", "payload": { "messageId": "uuid" } }
// @Description ```
// @Description
// @Description **`message:unstar`** — Удалить из избранного.
// @Description ```json
// @Description { "type": "message:unstar", "payload": { "messageId": "uuid" } }
// @Description ```
// @Description
// @Description **`message:forward`** — Переслать сообщение в другой чат.
// @Description ```json
// @Description { "type": "message:forward", "payload": { "messageId": "uuid", "toChatId": "uuid" } }
// @Description ```
// @Description
// @Description ### 2.2 Управление чатами
// @Description
// @Description **`chat:create`** — Создать новый чат.
// @Description ```json
// @Description { "type": "chat:create", "payload": { "type": "group", "name": "Friends", "participantIds": ["uuid1","uuid2"], "description": "Чат для друзей" } }
// @Description ```
// @Description | Поле | Обязательное | Описание |
// @Description |------|-------------|----------|
// @Description | type | ✅ | "private", "group", "channel" |
// @Description | participantIds | ✅ | Список ID участников |
// @Description | name | ❌ | Название (обяз. для group/channel) |
// @Description | description | ❌ | Описание чата |
// @Description | avatarUrl | ❌ | Ссылка на аватарку |
// @Description
// @Description **`chat:update`** — Обновить название/аватар/описание чата.
// @Description ```json
// @Description { "type": "chat:update", "payload": { "chatId": "uuid", "name": "Новое название", "description": "Описание", "avatarUrl": "https://..." } }
// @Description ```
// @Description
// @Description **`chat:add_participant`** — Добавить участника.
// @Description ```json
// @Description { "type": "chat:add_participant", "payload": { "chatId": "uuid", "userId": "uuid" } }
// @Description ```
// @Description
// @Description **`chat:remove_participant`** — Удалить участника.
// @Description ```json
// @Description { "type": "chat:remove_participant", "payload": { "chatId": "uuid", "userId": "uuid" } }
// @Description ```
// @Description
// @Description **`chat:leave`** — Покинуть чат.
// @Description ```json
// @Description { "type": "chat:leave", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`chat:pin`** — Закрепить чат в списке.
// @Description ```json
// @Description { "type": "chat:pin", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`chat:unpin`** — Открепить чат.
// @Description ```json
// @Description { "type": "chat:unpin", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`chat:archive`** — Архивировать чат.
// @Description ```json
// @Description { "type": "chat:archive", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`chat:unarchive`** — Разархивировать чат.
// @Description ```json
// @Description { "type": "chat:unarchive", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description ### 2.3 Статус пользователя
// @Description
// @Description **`user:typing`** — Индикатор печатания.
// @Description ```json
// @Description { "type": "user:typing", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`user:stop_typing`** — Остановить индикатор.
// @Description ```json
// @Description { "type": "user:stop_typing", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`user:keyboard_opened`** — Клавиатура открыта (моб.).
// @Description ```json
// @Description { "type": "user:keyboard_opened", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`user:keyboard_closed`** — Клавиатура закрыта.
// @Description ```json
// @Description { "type": "user:keyboard_closed", "payload": { "chatId": "uuid" } }
// @Description ```
// @Description
// @Description **`user:block`** — Заблокировать пользователя.
// @Description ```json
// @Description { "type": "user:block", "payload": { "userId": "uuid" } }
// @Description ```
// @Description
// @Description **`user:unblock`** — Разблокировать пользователя.
// @Description ```json
// @Description { "type": "user:unblock", "payload": { "userId": "uuid" } }
// @Description ```
// @Description
// @Description ### 2.4 WebRTC звонки (сигналинг)
// @Description
// @Description **`call:offer`** — Отправить WebRTC offer.
// @Description ```json
// @Description { "type": "call:offer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "offer_sdp_string" } }
// @Description ```
// @Description
// @Description **`call:answer`** — Ответить на звонок (WebRTC answer).
// @Description ```json
// @Description { "type": "call:answer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "answer_sdp_string" } }
// @Description ```
// @Description
// @Description **`call:ice`** — Отправить ICE-кандидат.
// @Description ```json
// @Description { "type": "call:ice", "payload": { "callId": "uuid", "candidate": "ice_candidate_string" } }
// @Description ```
// @Description
// @Description **`call:reject`** — Отклонить звонок.
// @Description ```json
// @Description { "type": "call:reject", "payload": { "callId": "uuid" } }
// @Description ```
// @Description
// @Description ---
// @Description
// @Description ## 3. Полный пример: отправка сообщения
// @Description
// @Description 1. Клиент А отправляет через WebSocket:
// @Description    ```json
// @Description    { "type": "message:send", "payload": { "chatId": "123", "content": "Привет!", "type": "text" } }
// @Description    ```
// @Description 2. Сервер сохраняет в БД и рассылает всем участникам чата:
// @Description    ```json
// @Description    { "type": "message:new", "payload": { "id": "msg-uuid", "chatId": "123", "senderId": "userA-uuid", "content": "Привет!", "type": "text", "createdAt": "2026-01-01T00:00:00Z" } }
// @Description    ```
// @Description 3. Клиент Б получает `message:new` и отображает сообщение в реальном времени.
// @Description
// @Description ## 4. Полный пример: звонок
// @Description
// @Description 1. Клиент А отправляет `call:offer` → сервер шлёт `call:offer` всем участникам чата
// @Description 2. Клиент Б отвечает `call:answer` → сервер шлёт `call:accept` всем
// @Description 3. Клиенты обмениваются `call:ice` (ICE candidates) через сервер
// @Description 4. Любая сторона шлёт `call:end` → сервер шлёт `call:end` всем
// @Description
// @Tags WebSocket
// @Accept json
// @Produce json
// @Param token query string true "JWT token (получить через POST /api/auth/login)"
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


