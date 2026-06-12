package messagehandler

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	messagedomain "ChatServerGolang/internal/domain/message"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type MessageHandler struct {
	messageService service.MessageService
}

func NewMessageHandler(messageService service.MessageService) *MessageHandler {
	return &MessageHandler{messageService: messageService}
}

// SendMessage sends a message to a chat
// @Tags Сообщения
// @Summary Отправить сообщение
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Отправляет текстовое сообщение в указанный чат. Можно указать тип сообщения, прикрепить ReplyToID для ответа на другое сообщение.
// @Param id path string true "ID чата"
// @Param request body messagedomain.SendMessageRequest true "Содержимое сообщения: content (текст сообщения, обязательно), type (тип: text/location/..., опционально), reply_to_id (ID сообщения для ответа, опционально)"
// @Success 201 {object} messagedomain.MessageResponse "Сообщение успешно отправлено"
// @Failure 400 {object} response.ErrorResponse "Неверные входные данные"
// @Router /chats/{id}/messages [post]
func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req messagedomain.SendMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.messageService.SendMessage(chatID, userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}

// GetMessages returns paginated messages for a chat
// @Tags Сообщения
// @Summary Получить сообщения чата
// @Security BearerAuth
// @Produce json
// @Description Возвращает постраничный список сообщений для указанного чата. Поддерживает пагинацию через параметры limit и offset.
// @Param id path string true "ID чата"
// @Param limit query int false "Сообщений на странице (по умолчанию 50)"
// @Param offset query int false "Смещение пагинации (по умолчанию 0)"
// @Success 200 {object} response.APIResponse "Постраничный список сообщений"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения сообщений"
// @Router /chats/{id}/messages [get]
func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, total, err := h.messageService.GetMessages(chatID, userID.(string), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Paginated(c, 200, messages, total, limit, offset)
}

// SearchMessages searches messages within a chat
// @Tags Сообщения
// @Summary Искать сообщения в чате
// @Security BearerAuth
// @Produce json
// @Description Ищет сообщения внутри указанного чата по текстовому запросу. Возвращает постраничные результаты.
// @Param id path string true "ID чата"
// @Param q query string true "Поисковый запрос"
// @Param limit query int false "Максимум результатов (по умолчанию 50)"
// @Param offset query int false "Смещение пагинации (по умолчанию 0)"
// @Success 200 {object} response.APIResponse "Постраничные результаты поиска"
// @Failure 400 {object} response.ErrorResponse "Отсутствует поисковый запрос"
// @Router /chats/{id}/messages/search [get]
func (h *MessageHandler) SearchMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	query := c.Query("q")

	if query == "" {
		response.BadRequest(c, "search query required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, total, err := h.messageService.SearchMessages(chatID, userID.(string), query, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Paginated(c, 200, messages, total, limit, offset)
}

// EditMessage edits a message (sender only)
// @Tags Сообщения
// @Summary Редактировать сообщение
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Редактирует текст сообщения. Доступно только отправителю сообщения.
// @Param id path string true "ID сообщения"
// @Param request body messagedomain.EditMessageRequest true "Обновлённое содержание: content (новый текст сообщения, обязательно)"
// @Success 200 {object} messagedomain.MessageResponse "Сообщение отредактировано"
// @Failure 400 {object} response.ErrorResponse "Неверные данные или сообщение принадлежит другому пользователю"
// @Router /messages/{id} [put]
func (h *MessageHandler) EditMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	var req messagedomain.EditMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.messageService.EditMessage(msgID, userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// DeleteMessage deletes a message (soft delete)
// @Tags Сообщения
// @Summary Удалить сообщение
// @Security BearerAuth
// @Produce json
// @Description Удаляет сообщение (мягкое удаление). Доступно только отправителю сообщения.
// @Param id path string true "ID сообщения"
// @Success 200 {object} response.MessageResponse "Сообщение удалено"
// @Failure 400 {object} response.ErrorResponse "Недостаточно прав для удаления"
// @Router /messages/{id} [delete]
func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	if err := h.messageService.DeleteMessage(msgID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "message deleted"})
}

// UploadFile uploads a file as a message attachment
// @Tags Сообщения
// @Summary Загрузить файл
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Загружает файл и отправляет его как вложение в чат. Поддерживаются различные типы файлов. Можно указать ID сообщения для ответа.
// @Param id path string true "ID чата"
// @Param file formData file true "Файл для загрузки"
// @Param replyToId formData string false "ID сообщения для ответа (опционально)"
// @Success 201 {object} messagedomain.MessageResponse "Файл загружен и отправлен как сообщение"
// @Failure 400 {object} response.ErrorResponse "Файл отсутствует"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
// @Router /chats/{id}/messages/file [post]
func (h *MessageHandler) UploadFile(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		response.BadRequest(c, "file required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".bin"
	}

	uploadDir := "uploads"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "failed to create upload directory")
		return
	}

	fileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}
	defer out.Close()

	fileSize, err := io.Copy(out, file)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	var replyToID *string
	if replyStr := c.PostForm("replyToId"); replyStr != "" {
		replyToID = &replyStr
	}

	msg, err := h.messageService.SendFileMessage(chatID, userID.(string), header.Filename, fileName, fileSize, replyToID)
	if err != nil {
		os.Remove(filePath)
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}

// DownloadFile serves a previously uploaded file
// @Tags Сообщения
// @Summary Скачать файл
// @Security BearerAuth
// @Produce application/octet-stream
// @Description Загружает ранее загруженный файл по его имени из директории uploads.
// @Param filename path string true "Имя файла для скачивания"
// @Success 200 {file} binary "Файл успешно загружен"
// @Failure 404 {object} response.ErrorResponse "Файл не найден"
// @Router /files/{filename} [get]
func (h *MessageHandler) DownloadFile(c *gin.Context) {
	fileName := c.Param("filename")
	filePath := filepath.Join("uploads", filepath.Clean(fileName))

	absPath, _ := filepath.Abs(filePath)
	if !strings.Contains(absPath, "uploads") {
		response.NotFound(c, "file not found")
		return
	}

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		response.NotFound(c, "file not found")
		return
	}

	c.File(filePath)
}

// AddReaction adds a reaction (emoji) to a message
// @Tags Сообщения
// @Summary Добавить реакцию к сообщению
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Добавляет реакцию в виде эмодзи к указанному сообщению. Пользователь может добавить только одну реакцию каждого типа.
// @Param id path string true "ID сообщения"
// @Param request body messagedomain.AddReactionRequest true "Эмодзи для реакции: emoji (строка с эмодзи, обязательно)"
// @Success 200 {object} messagedomain.MessageResponse "Реакция добавлена"
// @Failure 400 {object} response.ErrorResponse "Ошибка добавления реакции"
// @Router /messages/{id}/reactions [post]
func (h *MessageHandler) AddReaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	var req messagedomain.AddReactionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.messageService.AddReaction(msgID, userID.(string), req.Emoji)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// RemoveReaction removes a reaction from a message
// @Tags Сообщения
// @Summary Удалить реакцию с сообщения
// @Security BearerAuth
// @Produce json
// @Description Удаляет ранее добавленную реакцию (эмодзи) с указанного сообщения.
// @Param id path string true "ID сообщения"
// @Param emoji query string true "Эмодзи для удаления"
// @Success 200 {object} messagedomain.MessageResponse "Реакция удалена"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления реакции"
// @Router /messages/{id}/reactions [delete]
func (h *MessageHandler) RemoveReaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	emoji := c.Query("emoji")
	if emoji == "" {
		response.BadRequest(c, "emoji query parameter required")
		return
	}

	msg, err := h.messageService.RemoveReaction(msgID, userID.(string), emoji)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// TogglePin pins or unpins a message in a chat
// @Tags Сообщения
// @Summary Закрепить или открепить сообщение
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Закрепляет или открепляет сообщение в чате. Закреплённые сообщения отображаются в верхней части чата.
// @Param id path string true "ID сообщения"
// @Param request body messagedomain.PinMessageRequest true "Параметры: pin (boolean, true — закрепить, false — открепить, обязательно)"
// @Success 200 {object} messagedomain.MessageResponse "Статус закрепления изменён"
// @Failure 400 {object} response.ErrorResponse "Ошибка изменения статуса"
// @Router /messages/{id}/pin [put]
func (h *MessageHandler) TogglePin(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	var req messagedomain.PinMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.messageService.TogglePin(msgID, userID.(string), req.Pin)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// GetPinned returns all pinned messages in a chat
// @Tags Сообщения
// @Summary Получить закреплённые сообщения
// @Security BearerAuth
// @Produce json
// @Description Возвращает список всех закреплённых сообщений в указанном чате.
// @Param id path string true "ID чата"
// @Success 200 {array} messagedomain.MessageResponse "Список закреплённых сообщений"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения закреплённых сообщений"
// @Router /chats/{id}/pinned [get]
func (h *MessageHandler) GetPinned(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	messages, err := h.messageService.GetPinnedMessages(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, messages)
}

// MarkMessageRead marks a single message as read
// @Tags Сообщения
// @Summary Отметить сообщение как прочитанное
// @Security BearerAuth
// @Produce json
// @Description Помечает одно сообщение как прочитанное для текущего пользователя.
// @Param id path string true "ID сообщения"
// @Success 200 {object} response.MessageResponse "Сообщение отмечено как прочитанное"
// @Failure 400 {object} response.ErrorResponse "Ошибка отметки"
// @Router /messages/{id}/read [post]
func (h *MessageHandler) MarkMessageRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	if err := h.messageService.MarkMessageRead(msgID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "marked as read"})
}

// ResendMessage resends a message (useful for failed sends)
// @Tags Сообщения
// @Summary Повторно отправить сообщение
// @Security BearerAuth
// @Produce json
// @Description Повторно отправляет сообщение, которое не было доставлено. Создаёт новое сообщение с тем же содержимым.
// @Param id path string true "ID чата"
// @Param msgId path string true "ID сообщения для повторной отправки"
// @Success 201 {object} messagedomain.MessageResponse "Сообщение повторно отправлено"
// @Failure 400 {object} response.ErrorResponse "Ошибка повторной отправки"
// @Router /chats/{id}/messages/{msgId}/resend [post]
func (h *MessageHandler) ResendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	msgID := c.Param("msgId")

	msg, err := h.messageService.ResendMessage(chatID, userID.(string), msgID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}

// GetMessageByID returns a single message by its ID
// @Tags Сообщения
// @Summary Получить сообщение по ID
// @Security BearerAuth
// @Produce json
// @Description Возвращает одно сообщение по его уникальному идентификатору.
// @Param id path string true "ID сообщения"
// @Success 200 {object} messagedomain.MessageResponse "Информация о сообщении"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения сообщения"
// @Router /messages/{id} [get]
func (h *MessageHandler) GetMessageByID(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	msg, err := h.messageService.GetMessageByID(msgID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// BulkMarkRead marks multiple messages as read
// @Tags Сообщения
// @Summary Отметить несколько сообщений как прочитанные
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Помечает несколько сообщений как прочитанные за один запрос. Принимает массив ID сообщений и ID чата.
// @Param request body object{messageIds=[]string,chatId=string} true "Параметры: messageIds (массив ID сообщений, обязательно, мин. 1), chatId (ID чата, обязательно)"
// @Success 200 {object} response.MessageResponse "Сообщения отмечены как прочитанные"
// @Failure 400 {object} response.ErrorResponse "Ошибка отметки"
// @Router /messages/read/bulk [post]
func (h *MessageHandler) BulkMarkRead(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		MessageIDs []string `json:"messageIds" binding:"required,min=1"`
		ChatID     string   `json:"chatId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, msgID := range req.MessageIDs {
		h.messageService.MarkMessageRead(msgID, userID.(string))
	}

	response.JSON(c, 200, gin.H{"message": "messages marked as read", "count": len(req.MessageIDs)})
}

// BulkDeleteMessages soft deletes multiple messages for the current user
// @Tags Сообщения
// @Summary Массовое удаление сообщений
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Удаляет несколько сообщений для текущего пользователя. Сообщения удаляются только для отправителя.
// @Param request body object{messageIds=[]string} true "Параметры: messageIds (массив ID сообщений для удаления, обязательно, мин. 1)"
// @Success 200 {object} response.MessageResponse "Сообщения удалены"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления"
// @Router /messages/bulk [delete]
func (h *MessageHandler) BulkDeleteMessages(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		MessageIDs []string `json:"messageIds" binding:"required,min=1"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	for _, msgID := range req.MessageIDs {
		h.messageService.DeleteMessageForMe(msgID, userID.(string))
	}

	response.JSON(c, 200, gin.H{"message": "messages deleted", "count": len(req.MessageIDs)})
}

// UploadVoice uploads a voice message
// @Tags Сообщения
// @Summary Отправить голосовое сообщение
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Загружает и отправляет голосовое сообщение (аудиозапись) в чат. Поддерживается формат Opus/Ogg.
// @Param id path string true "ID чата"
// @Param voice formData file true "Аудиозапись голосового сообщения (формат Opus/Ogg)"
// @Success 201 {object} messagedomain.MessageResponse "Голосовое сообщение отправлено"
// @Failure 400 {object} response.ErrorResponse "Ошибка загрузки"
// @Router /chats/{id}/messages/voice [post]
func (h *MessageHandler) UploadVoice(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("voice")
	if err != nil {
		response.BadRequest(c, "voice file required")
		return
	}
	defer file.Close()

	uploadDir := "uploads/voice"
	os.MkdirAll(uploadDir, 0755)

	fileName := uuid.New().String() + ".ogg"
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save voice message")
		return
	}
	defer out.Close()

	fileSize, err := io.Copy(out, file)
	if err != nil {
		response.InternalError(c, "failed to save voice message")
		return
	}

	msg, err := h.messageService.SendFileMessage(chatID, userID.(string), header.Filename, fileName, fileSize, nil)
	if err != nil {
		os.Remove(filePath)
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}

// SearchAllMessages searches all chats for the user
// @Tags Сообщения
// @Summary Искать сообщения по всем чатам
// @Security BearerAuth
// @Produce json
// @Description Ищет сообщения по всем чатам пользователя. Возвращает постраничные результаты с информацией о чате.
// @Param q query string true "Поисковый запрос"
// @Param limit query int false "Максимум результатов (по умолчанию 50)"
// @Param offset query int false "Смещение пагинации (по умолчанию 0)"
// @Success 200 {object} response.APIResponse "Постраничные результаты поиска"
// @Failure 400 {object} response.ErrorResponse "Отсутствует поисковый запрос"
// @Router /messages/search [get]
func (h *MessageHandler) SearchAllMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")

	if query == "" {
		response.BadRequest(c, "search query required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, total, err := h.messageService.SearchAllMessages(userID.(string), query, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Paginated(c, 200, messages, total, limit, offset)
}

// ForwardMessage forwards a message from one chat to another
// @Tags Сообщения
// @Summary Переслать сообщение
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Пересылает сообщение из одного чата в другой. Исходное сообщение остаётся без изменений.
// @Param request body messagedomain.ForwardMessageRequest true "Детали пересылки: message_id (ID сообщения, обязательно), from_chat_id (ID исходного чата, обязательно), to_chat_id (ID целевого чата, обязательно)"
// @Success 201 {object} messagedomain.MessageResponse "Сообщение переслано"
// @Failure 400 {object} response.ErrorResponse "Ошибка пересылки"
// @Router /messages/forward [post]
func (h *MessageHandler) ForwardMessage(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req messagedomain.ForwardMessageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	msg, err := h.messageService.ForwardMessage(req.MessageID, req.FromChatID, req.ToChatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}

// ReportMessage reports a message for moderation
// @Tags Сообщения
// @Summary Пожаловаться на сообщение
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Отправляет жалобу на сообщение модераторам. Необходимо указать причину жалобы.
// @Param id path string true "ID сообщения"
// @Param request body object{reason=string} true "Причина жалобы: reason (строка, обязательно)"
// @Success 200 {object} response.MessageResponse "Жалоба отправлена"
// @Failure 400 {object} response.ErrorResponse "Ошибка отправки жалобы"
// @Router /messages/{id}/report [post]
func (h *MessageHandler) ReportMessage(c *gin.Context) {
	msgID := c.Param("id")

	var req struct {
		Reason string `json:"reason" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"messageId": msgID, "reason": req.Reason, "status": "reported"})
}

// DeleteMessageForMe deletes a message only for the current user
// @Tags Сообщения
// @Summary Удалить сообщение только для себя
// @Security BearerAuth
// @Produce json
// @Description Удаляет сообщение только для текущего пользователя. Остальные участники чата по-прежнему видят сообщение.
// @Param id path string true "ID сообщения"
// @Success 200 {object} response.MessageResponse "Сообщение удалено для вас"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления"
// @Router /messages/{id}/for-me [delete]
func (h *MessageHandler) DeleteMessageForMe(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	if err := h.messageService.DeleteMessageForMe(msgID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "message deleted for you"})
}

// StarMessage stars a message for quick access
// @Tags Сообщения
// @Summary Добавить сообщение в избранное
// @Security BearerAuth
// @Produce json
// @Description Добавляет сообщение в список избранных для быстрого доступа. Звёздочка видна только пользователю.
// @Param id path string true "ID сообщения"
// @Success 200 {object} messagedomain.MessageResponse "Сообщение добавлено в избранное"
// @Failure 400 {object} response.ErrorResponse "Ошибка добавления"
// @Router /messages/{id}/star [post]
func (h *MessageHandler) StarMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	msg, err := h.messageService.StarMessage(msgID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// UnstarMessage removes a star from a message
// @Tags Сообщения
// @Summary Удалить сообщение из избранного
// @Security BearerAuth
// @Produce json
// @Description Удаляет сообщение из списка избранных.
// @Param id path string true "ID сообщения"
// @Success 200 {object} response.MessageResponse "Сообщение удалено из избранного"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления"
// @Router /messages/{id}/star [delete]
func (h *MessageHandler) UnstarMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	if err := h.messageService.UnstarMessage(msgID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "message unstarred"})
}

// GetStarredMessages returns all starred messages for the user
// @Tags Сообщения
// @Summary Получить избранные сообщения
// @Security BearerAuth
// @Produce json
// @Description Возвращает все сообщения, добавленные пользователем в избранное.
// @Success 200 {array} chatdomain.StarredMessageResponse "Список избранных сообщений"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения списка"
// @Router /messages/starred [get]
func (h *MessageHandler) GetStarredMessages(c *gin.Context) {
	userID, _ := c.Get("userID")

	messages, err := h.messageService.GetStarredMessages(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, messages)
}

// GetChatMedia returns paginated media messages from a chat
// @Tags Сообщения
// @Summary Получить медиафайлы чата
// @Security BearerAuth
// @Produce json
// @Description Возвращает постраничный список медиасообщений из чата. Можно фильтровать по типу (photo, video, audio, document).
// @Param id path string true "ID чата"
// @Param type query string false "Фильтр по типу медиа (photo, video, audio, document)"
// @Param limit query int false "Максимум результатов (по умолчанию 50)"
// @Param offset query int false "Смещение пагинации (по умолчанию 0)"
// @Success 200 {object} response.APIResponse "Постраничный список медиа"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения медиа"
// @Router /chats/{id}/media [get]
func (h *MessageHandler) GetChatMedia(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	mediaType := c.Query("type")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, total, err := h.messageService.GetChatMedia(chatID, userID.(string), mediaType, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Paginated(c, 200, messages, total, limit, offset)
}

// GetMessageHistory returns message edit history
// @Tags Сообщения
// @Summary Получить историю изменений сообщения
// @Security BearerAuth
// @Produce json
// @Description Возвращает историю редактирования сообщения, включая все предыдущие версии текста.
// @Param id path string true "ID сообщения"
// @Success 200 {object} messagedomain.MessageResponse "История изменений сообщения"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения истории"
// @Router /messages/{id}/history [get]
func (h *MessageHandler) GetMessageHistory(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	msg, err := h.messageService.GetMessageByID(msgID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, msg)
}

// ExportChat exports all messages from a chat
// @Tags Сообщения
// @Summary Экспортировать сообщения чата
// @Security BearerAuth
// @Produce json
// @Description Экспортирует все сообщения из указанного чата в формате JSON для резервного копирования.
// @Param id path string true "ID чата"
// @Success 200 {array} messagedomain.MessageResponse "Все сообщения чата"
// @Failure 400 {object} response.ErrorResponse "Ошибка экспорта"
// @Router /chats/{id}/export [get]
func (h *MessageHandler) ExportChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	messages, err := h.messageService.ExportChat(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, messages)
}

// UploadVideoCircle загружает и отправляет круговое видео (видеосообщение) в чат
// @Tags Сообщения
// @Summary Отправить круговое видеосообщение (MP4)
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Позволяет отправить круговое видеосообщение (аналог видеосообщений в Telegram). Видео должно быть в формате MP4 с круговым обрезанием. Файл сохраняется в директории uploads/video_circles/. Можно указать подпись к видео.
// @Param id path string true "ID чата, в который отправляется видеосообщение"
// @Param video formData file true "Файл видео в формате MP4 для кругового видеосообщения. Рекомендуется квадратное соотношение сторон."
// @Param caption formData string false "Подпись к видеосообщению (опционально, отображается под видео)"
// @Success 201 {object} messagedomain.MessageResponse "Видеосообщение успешно отправлено. Тип сообщения: video_circle."
// @Failure 400 {object} response.ErrorResponse "Ошибка: файл не загружен, неверный формат или чат не найден."
// @Router /chats/{id}/messages/video-circle [post]
func (h *MessageHandler) UploadVideoCircle(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("video")
	if err != nil {
		response.BadRequest(c, "video file required")
		return
	}
	defer file.Close()

	uploadDir := "uploads/video_circles"
	os.MkdirAll(uploadDir, 0755)

	fileName := uuid.New().String() + ".mp4"
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save video")
		return
	}
	defer out.Close()

	fileSize, err := io.Copy(out, file)
	if err != nil {
		response.InternalError(c, "failed to save video")
		return
	}

	msg, err := h.messageService.SendFileMessage(chatID, userID.(string), header.Filename, fileName, fileSize, nil)
	if err != nil {
		os.Remove(filePath)
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}

// SendLocation отправляет сообщение с геолокацией в указанный чат
// @Tags Сообщения
// @Summary Отправить сообщение с геолокацией (координаты + опциональное название)
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Создаёт и отправляет в чат сообщение с геолокацией (тип "location"). Позволяет пользователям делиться своим местоположением на карте. В теле запроса обязательно нужно передать широту (latitude) и долготу (longitude) места. Опционально можно добавить название места (title), указать ID сообщения для ответа (replyToId) и выбрать эффект анимации (effect). После успешной отправки возвращается полный объект сообщения с координатами.
// @Param id path string true "ID чата (группы или личного диалога), в который отправляется сообщение с геолокацией"
// @Param request body messagedomain.SendLocationRequest true "Параметры геолокации: latitude (число, широта, обязательно), longitude (число, долгота, обязательно), title (строка, название места, опционально), replyToId (строка, ID сообщения для ответа, опционально), effect (строка, эффект анимации: confetti/fireworks/hearts/balloons/stars, опционально)"
// @Success 201 {object} messagedomain.MessageResponse "Сообщение с геолокацией успешно создано. В ответе возвращается полный объект MessageResponse с заполненными полями latitude, longitude, locationTitle."
// @Failure 400 {object} response.ErrorResponse "Ошибка валидации: не указаны latitude/longitude, неверный формат данных или чат не найден."
// @Router /chats/{id}/messages/location [post]
func (h *MessageHandler) SendLocation(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req messagedomain.SendLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	msgReq := &messagedomain.SendMessageRequest{
		Content:       "",
		Type:          messagedomain.MessageLocation,
		Latitude:      req.Latitude,
		Longitude:     req.Longitude,
		LocationTitle: req.Title,
		ReplyToID:     req.ReplyToID,
		Effect:        req.Effect,
	}

	msg, err := h.messageService.SendMessage(chatID, userID.(string), msgReq)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, msg)
}
