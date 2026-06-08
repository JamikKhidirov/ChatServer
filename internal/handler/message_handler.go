package handler

import (
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"ChatServerGolang/internal/domain"
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

func (h *MessageHandler) SendMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req domain.SendMessageRequest
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

func (h *MessageHandler) GetMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, err := h.messageService.GetMessages(chatID, userID.(string), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, messages)
}

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

	messages, err := h.messageService.SearchMessages(chatID, userID.(string), query, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, messages)
}

func (h *MessageHandler) EditMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	var req domain.EditMessageRequest
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

func (h *MessageHandler) DeleteMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	if err := h.messageService.DeleteMessage(msgID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "message deleted"})
}

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

func (h *MessageHandler) AddReaction(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	var req domain.AddReactionRequest
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

func (h *MessageHandler) TogglePin(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	var req domain.PinMessageRequest
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

func (h *MessageHandler) MarkMessageRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	msgID := c.Param("id")

	if err := h.messageService.MarkMessageRead(msgID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "marked as read"})
}

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
