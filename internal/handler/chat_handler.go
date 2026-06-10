package handler

import (
	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
}

func NewChatHandler(chatService service.ChatService) *ChatHandler {
	return &ChatHandler{chatService: chatService}
}

func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.CreateChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	chat, err := h.chatService.CreateChat(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	c.Set("chatResponse", chat)
	response.JSON(c, 201, chat)
}

func (h *ChatHandler) GetChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	chat, err := h.chatService.GetChat(chatID, userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, chat)
}

func (h *ChatHandler) SearchChats(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "query parameter q is required")
		return
	}

	chats, err := h.chatService.SearchChats(userID.(string), query)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

func (h *ChatHandler) ListChats(c *gin.Context) {
	userID, _ := c.Get("userID")

	chats, err := h.chatService.ListChats(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

func (h *ChatHandler) DeleteChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.DeleteChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat deleted"})
}

func (h *ChatHandler) AddParticipant(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req domain.AddParticipantRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.AddParticipant(chatID, req.UserID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "participant added"})
}

func (h *ChatHandler) RemoveParticipant(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	targetUserID := c.Param("userId")

	if err := h.chatService.RemoveParticipant(chatID, targetUserID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "participant removed"})
}

func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.MarkAsRead(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "marked as read"})
}

func (h *ChatHandler) SetRole(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	targetUserID := c.Param("userId")

	var req struct {
		Role string `json:"role" binding:"required,oneof=admin member"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetRole(chatID, targetUserID, userID.(string), req.Role); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "role updated to " + req.Role})
}

func (h *ChatHandler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.LeaveGroup(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "left the group"})
}

func (h *ChatHandler) HideChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.HideChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat hidden"})
}

func (h *ChatHandler) UpdateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req domain.UpdateGroupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.UpdateGroup(chatID, userID.(string), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "group updated"})
}
