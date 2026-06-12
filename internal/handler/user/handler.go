package userhandler

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	userdomain "ChatServerGolang/internal/domain/user"
	notificationdomain "ChatServerGolang/internal/domain/notification"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService service.UserService
	pushService service.PushService
}

func NewUserHandler(userService service.UserService, pushService service.PushService) *UserHandler {
	return &UserHandler{userService: userService, pushService: pushService}
}

// GetProfile returns the authenticated user's profile
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} userdomain.UserResponse
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/profile [get]
func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	profile, err := h.userService.GetProfile(userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, profile)
}

// UpdateProfile updates the authenticated user's profile
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userdomain.UpdateProfileRequest true "Profile fields to update"
// @Success 200 {object} userdomain.UserResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Router /users/profile [put]
func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req userdomain.UpdateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	profile, err := h.userService.UpdateProfile(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, profile)
}

// DeleteAccount permanently deletes the authenticated user's account
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Deletion failed"
// @Router /account [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := h.userService.DeleteAccount(userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "account deleted"})
}

// SearchUsers searches for users by query
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param q query string true "Search query"
// @Param limit query int false "Max results (default 50)"
// @Param offset query int false "Pagination offset (default 0)"
// @Success 200 {object} response.APIResponse "Paginated user list"
// @Failure 400 {object} response.ErrorResponse "Missing query"
// @Router /users/search [get]
func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "search query required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, total, err := h.userService.SearchUsers(query, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Paginated(c, 200, users, total, limit, offset)
}

// GetUserByID returns a user by their ID
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} userdomain.UserResponse
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/{id} [get]
func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.JSON(c, 200, user.ToResponse())
}

// GetUserByUsername returns a user by their username
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param username path string true "Username"
// @Success 200 {object} userdomain.UserResponse
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/username/{username} [get]
func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.GetByUsername(username)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.JSON(c, 200, user)
}

// TestPush sends a test push notification to the authenticated user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{title=string,body=string} true "Push notification fields"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input"
// @Router /users/push-test [post]
func (h *UserHandler) TestPush(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Title string `json:"title" binding:"required"`
		Body  string `json:"body" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	h.pushService.SendTestPush(userID.(string), req.Title, req.Body)
	response.JSON(c, 200, gin.H{"message": "test push sent"})
}

// UpdateStatus updates the user's online status
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userdomain.UpdateStatusRequest true "New status"
// @Success 200 {object} userdomain.UserResponse
// @Failure 400 {object} response.ErrorResponse "Invalid status"
// @Router /users/status [put]
func (h *UserHandler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req userdomain.UpdateStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	profile, err := h.userService.UpdateStatus(userID.(string), req.Status)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, profile)
}

// UpdatePushToken updates the push notification token
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userdomain.UpdatePushTokenRequest true "Push token details"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Invalid token"
// @Router /users/push-token [put]
func (h *UserHandler) UpdatePushToken(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req userdomain.UpdatePushTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.UpdatePushToken(userID.(string), req.Token, req.Provider); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "push token updated"})
}

// BlockUser blocks a user
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userdomain.BlockUserRequest true "User ID to block"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Block failed"
// @Router /users/block [post]
func (h *UserHandler) BlockUser(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req userdomain.BlockUserRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.BlockUser(userID.(string), req.BlockedID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user blocked"})
}

// UnblockUser unblocks a user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param userId path string true "User ID to unblock"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Unblock failed"
// @Router /users/block/{userId} [delete]
func (h *UserHandler) UnblockUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	blockedID := c.Param("userId")

	if err := h.userService.UnblockUser(userID.(string), blockedID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user unblocked"})
}

// GetBlockedUsers returns the list of blocked users
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {array} userdomain.UserResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /users/blocked [get]
func (h *UserHandler) GetBlockedUsers(c *gin.Context) {
	userID, _ := c.Get("userID")

	users, err := h.userService.GetBlockedUsers(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, users)
}

// SetNotificationMuted mutes or unmutes notifications for a chat
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body notificationdomain.UpdateNotificationSettingRequest true "Mute settings"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/notifications [put]
func (h *UserHandler) SetNotificationMuted(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req notificationdomain.UpdateNotificationSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.userService.SetNotificationMuted(userID.(string), chatID, req.Muted); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	action := "unmuted"
	if req.Muted {
		action = "muted"
	}
	response.JSON(c, 200, gin.H{"message": "notifications " + action})
}

// IsNotificationMuted checks if notifications are muted for a chat
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {object} object{muted=boolean}
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/notifications [get]
func (h *UserHandler) IsNotificationMuted(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	muted, err := h.userService.IsNotificationMuted(userID.(string), chatID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"muted": muted})
}

// UploadAvatar uploads a new avatar image
// @Tags Users
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Param avatar formData file true "Avatar image file"
// @Success 200 {object} userdomain.UserResponse
// @Failure 400 {object} response.ErrorResponse "Missing or invalid file"
// @Failure 500 {object} response.ErrorResponse "Server error"
// @Router /users/avatar [post]
func (h *UserHandler) UploadAvatar(c *gin.Context) {
	userID, _ := c.Get("userID")

	file, header, err := c.Request.FormFile("avatar")
	if err != nil {
		response.BadRequest(c, "avatar file required")
		return
	}
	defer file.Close()

	ext := filepath.Ext(header.Filename)
	if ext == "" {
		ext = ".jpg"
	}

	uploadDir := "uploads/avatars"
	if err := os.MkdirAll(uploadDir, 0755); err != nil {
		response.InternalError(c, "failed to create upload directory")
		return
	}

	fileName := userID.(string) + ext
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

	avatarURL := "/uploads/avatars/" + fileName

	updated, err := h.userService.UpdateProfile(userID.(string), &userdomain.UpdateProfileRequest{AvatarURL: avatarURL})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, updated)
}

// GetAccountSetting returns the user's account settings
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Success 200 {object} userdomain.AccountSetting
// @Failure 400 {object} response.ErrorResponse
// @Router /account/settings [get]
func (h *UserHandler) GetAccountSetting(c *gin.Context) {
	userID, _ := c.Get("userID")

	setting, err := h.userService.GetAccountSetting(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, setting)
}

// GetLastSeen returns the last seen timestamp for a user
// @Tags Users
// @Security BearerAuth
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} object{userId=string,online=boolean,lastSeen=string}
// @Failure 404 {object} response.ErrorResponse "User not found"
// @Router /users/{id}/last-seen [get]
func (h *UserHandler) GetLastSeen(c *gin.Context) {
	targetID := c.Param("id")

	user, err := h.userService.GetUserByID(targetID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.JSON(c, 200, gin.H{
		"userId":   targetID,
		"username": user.Username,
		"online":   user.Online,
		"lastSeen": user.LastSeen,
	})
}

// ChangeUsername changes the user's username
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{username=string} true "New username"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /users/username [put]
func (h *UserHandler) ChangeUsername(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Username string `json:"username" binding:"required,min=3,max=32"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	profile, err := h.userService.UpdateProfile(userID.(string), &userdomain.UpdateProfileRequest{Username: req.Username})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, profile)
}

// ChangeEmail changes the user's email address
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{email=string} true "New email"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /users/email [put]
func (h *UserHandler) ChangeEmail(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req struct {
		Email string `json:"email" binding:"required,email"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	profile, err := h.userService.UpdateProfile(userID.(string), &userdomain.UpdateProfileRequest{Email: req.Email})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, profile)
}

// UpdateAccountSetting updates the user's account settings
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body userdomain.UpdateAccountSettingRequest true "Settings to update"
// @Success 200 {object} userdomain.AccountSetting
// @Failure 400 {object} response.ErrorResponse "Invalid settings"
// @Router /account/settings [put]
func (h *UserHandler) UpdateAccountSetting(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req userdomain.UpdateAccountSettingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	setting, err := h.userService.UpdateAccountSetting(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, setting)
}
