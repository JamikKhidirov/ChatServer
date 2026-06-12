package chathandler

import (
	"io"
	"os"
	"path/filepath"

	chatdomain "ChatServerGolang/internal/domain/chat"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/internal/ws"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type ChatHandler struct {
	chatService service.ChatService
	hub         *ws.Hub
}

func NewChatHandler(chatService service.ChatService, hub *ws.Hub) *ChatHandler {
	return &ChatHandler{chatService: chatService, hub: hub}
}

// CreateChat creates a new chat (private or group)
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body chatdomain.CreateChatRequest true "Chat details"
// @Success 201 {object} chatdomain.ChatResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Router /chats [post]
func (h *ChatHandler) CreateChat(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req chatdomain.CreateChatRequest
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

// GetChat returns chat details
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} chatdomain.ChatResponse
// @Failure 404 {object} response.ErrorResponse "Chat not found"
// @Router /chats/{id} [get]
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

// SearchChats searches the user's chats by name
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search query"
// @Success 200 {array} chatdomain.ChatResponse
// @Failure 400 {object} response.ErrorResponse "Missing query"
// @Router /chats/search [get]
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

// ListChats returns all chats for the authenticated user
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Success 200 {array} chatdomain.ChatResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats [get]
func (h *ChatHandler) ListChats(c *gin.Context) {
	userID, _ := c.Get("userID")

	chats, err := h.chatService.ListChats(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

// DeleteChat deletes a chat (owner only)
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id} [delete]
func (h *ChatHandler) DeleteChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.DeleteChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat deleted"})
}

// AddParticipant adds a user to a group chat
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body chatdomain.AddParticipantRequest true "User ID to add"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/participants [post]
func (h *ChatHandler) AddParticipant(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req chatdomain.AddParticipantRequest
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

// RemoveParticipant removes a user from a group chat
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Param userId path string true "Target user ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/participants/{userId} [delete]
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

// MarkAsRead marks all messages in a chat as read
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/read [post]
func (h *ChatHandler) MarkAsRead(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.MarkAsRead(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "marked as read"})
}

// SetRole changes a participant's role (admin/member)
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param userId path string true "Target user ID"
// @Param request body object{role=string} true 'Role: "admin" or "member"'
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/participants/{userId}/role [put]
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

// LeaveGroup removes the authenticated user from a group chat
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/leave [post]
func (h *ChatHandler) LeaveGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.LeaveGroup(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "left the group"})
}

// HideChat hides a chat from the user's chat list
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/hide [post]
func (h *ChatHandler) HideChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.HideChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat hidden"})
}

// UpdateGroup updates group chat name/description/photo
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body chatdomain.UpdateGroupRequest true "Fields to update"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id} [put]
func (h *ChatHandler) UpdateGroup(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req chatdomain.UpdateGroupRequest
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

// PinChat pins a chat to the top of the list
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/pin [post]
func (h *ChatHandler) PinChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.PinChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat pinned"})
}

// UnpinChat unpins a chat from the top of the list
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/pin [delete]
func (h *ChatHandler) UnpinChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.UnpinChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat unpinned"})
}

// ArchiveChat archives a chat
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/archive [post]
func (h *ChatHandler) ArchiveChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.ArchiveChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat archived"})
}

// UnarchiveChat unarchives a chat
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/unarchive [post]
func (h *ChatHandler) UnarchiveChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	if err := h.chatService.UnarchiveChat(chatID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "chat unarchived"})
}

// ListArchivedChats returns the user's archived chats
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Success 200 {array} chatdomain.ChatResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/archived [get]
func (h *ChatHandler) ListArchivedChats(c *gin.Context) {
	userID, _ := c.Get("userID")

	chats, err := h.chatService.ListArchivedChats(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chats)
}

// TransferOwnership transfers group ownership to another participant
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body object{userId=string} true "New owner's user ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/transfer-ownership [post]
func (h *ChatHandler) TransferOwnership(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.TransferOwnership(chatID, userID.(string), req.UserID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "ownership transferred"})
}

// SetSlowMode sets slow mode interval for a group chat
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body object{seconds=integer} true "Slow mode interval in seconds (0-3600, 0=disabled)"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Invalid value"
// @Router /chats/{id}/slow-mode [put]
func (h *ChatHandler) SetSlowMode(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		Seconds int `json:"seconds" binding:"min=0,max=3600"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetSlowMode(chatID, userID.(string), req.Seconds); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "slow mode updated"})
}

// PromoteToAdmin promotes a participant to admin role
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body object{userId=string} true "Target user ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/promote [post]
func (h *ChatHandler) PromoteToAdmin(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetRole(chatID, req.UserID, userID.(string), "admin"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user promoted to admin"})
}

// DemoteFromAdmin demotes a participant from admin to member
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body object{userId=string} true "Target user ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/demote [post]
func (h *ChatHandler) DemoteFromAdmin(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req struct {
		UserID string `json:"userId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.chatService.SetRole(chatID, req.UserID, userID.(string), "member"); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user demoted to member"})
}

// UploadChatPhoto uploads a photo for a group chat
// @Tags Chats
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Chat ID"
// @Param photo formData file true "Chat photo image"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/photo [post]
func (h *ChatHandler) UploadChatPhoto(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("photo")
	if err != nil {
		response.BadRequest(c, "photo file required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	uploadDir := "uploads/chat_photos"
	os.MkdirAll(uploadDir, 0755)

	fileName := chatID + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save file")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		response.InternalError(c, "failed to save file")
		return
	}

	photoURL := "/uploads/chat_photos/" + fileName
	if err := h.chatService.UpdateGroup(chatID, userID.(string), &chatdomain.UpdateGroupRequest{AvatarURL: photoURL}); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"photoUrl": photoURL})
}

// GetOnlineMembers returns online members of a chat
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} object{userIds=[]string}
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/online [get]
func (h *ChatHandler) GetOnlineMembers(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	chat, err := h.chatService.GetChat(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	var onlineIDs []string
	for _, p := range chat.Participants {
		if h.hub.IsOnline(p.ID) {
			onlineIDs = append(onlineIDs, p.ID)
		}
	}

	response.JSON(c, 200, gin.H{"userIds": onlineIDs})
}

// SetChatPermissions sets group permissions (who can send messages, add members, etc.)
// @Tags Chats
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body object{whoCanSend=string} true 'Permissions: "everyone", "admins"'
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/permissions [put]
func (h *ChatHandler) SetChatPermissions(c *gin.Context) {
	chatID := c.Param("id")

	var req struct {
		WhoCanSend string `json:"whoCanSend" binding:"required,oneof=everyone admins"`
		WhoCanAdd  string `json:"whoCanAdd" binding:"required,oneof=everyone admins"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"chatId": chatID, "whoCanSend": req.WhoCanSend, "whoCanAdd": req.WhoCanAdd})
}

// SetChatWallpaper sets the wallpaper for a chat
// @Tags Chats
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param id path string true "Chat ID"
// @Param wallpaper formData file true "Wallpaper image"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/wallpaper [post]
func (h *ChatHandler) SetChatWallpaper(c *gin.Context) {
	chatID := c.Param("id")

	file, header, err := c.Request.FormFile("wallpaper")
	if err != nil {
		response.BadRequest(c, "wallpaper file required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	uploadDir := "uploads/wallpapers"
	os.MkdirAll(uploadDir, 0755)

	fileName := chatID + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "failed to save wallpaper")
		return
	}
	defer out.Close()

	io.Copy(out, file)

	wallpaperURL := "/uploads/wallpapers/" + fileName
	response.JSON(c, 200, gin.H{"wallpaperUrl": wallpaperURL})
}

// StartPrivateChat finds or creates a private chat with another user
// @Tags Chats
// @Security BearerAuth
// @Produce json
// @Param userId path string true "Target user ID to chat with"
// @Success 200 {object} chatdomain.ChatResponse "Existing private chat"
// @Success 201 {object} chatdomain.ChatResponse "New private chat created"
// @Failure 400 {object} response.ErrorResponse "Invalid user ID"
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /chats/start/{userId} [post]
func (h *ChatHandler) StartPrivateChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	targetUserID := c.Param("userId")

	if targetUserID == userID.(string) {
		response.BadRequest(c, "cannot start chat with yourself")
		return
	}

	req := &chatdomain.CreateChatRequest{
		Type:           chatdomain.ChatPrivate,
		ParticipantIDs: []string{targetUserID},
	}

	chat, err := h.chatService.CreateChat(userID.(string), req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, chat)
}
