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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body messagedomain.SendMessageRequest true "Message content"
// @Success 201 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Param limit query int false "Messages per page (default 50)"
// @Param offset query int false "Pagination offset (default 0)"
// @Success 200 {object} response.APIResponse "Paginated messages"
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Param q query string true "Search query"
// @Param limit query int false "Max results (default 50)"
// @Param offset query int false "Pagination offset (default 0)"
// @Success 200 {object} response.APIResponse "Paginated results"
// @Failure 400 {object} response.ErrorResponse "Missing query"
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param request body messagedomain.EditMessageRequest true "Updated content"
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input or not owner"
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Not owner"
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
// @Tags Messages
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Chat ID"
// @Param file formData file true "File to upload"
// @Param replyToId formData string false "Optional message ID to reply to"
// @Success 201 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Missing file"
// @Failure 500 {object} response.ErrorResponse "Server error"
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
// @Tags Messages
// @Security BearerAuth
// @Produce application/octet-stream
// @Param filename path string true "Filename to download"
// @Success 200 {file} binary
// @Failure 404 {object} response.ErrorResponse "File not found"
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param request body messagedomain.AddReactionRequest true "Emoji to add"
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Param emoji query string true "Emoji to remove"
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param request body messagedomain.PinMessageRequest true 'pin: true to pin, false to unpin'
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {array} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Param msgId path string true "Message ID to resend"
// @Success 201 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{messageIds=[]string,chatId=string} true "Message IDs to mark as read"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{messageIds=[]string} true "Message IDs to delete"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Chat ID"
// @Param voice formData file true "Voice recording (opus/ogg)"
// @Success 201 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Max results (default 50)"
// @Param offset query int false "Pagination offset (default 0)"
// @Success 200 {object} response.APIResponse "Paginated results"
// @Failure 400 {object} response.ErrorResponse "Missing query"
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body messagedomain.ForwardMessageRequest true "Forward details"
// @Success 201 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Message ID"
// @Param request body object{reason=string} true "Report reason"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Success 200 {array} messagedomain.StarredMessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Param type query string false "Media type filter (photo, video, audio, document)"
// @Param limit query int false "Max results (default 50)"
// @Param offset query int false "Pagination offset (default 0)"
// @Success 200 {object} response.APIResponse "Paginated media"
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Message ID"
// @Success 200 {object} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Messages
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {array} messagedomain.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
