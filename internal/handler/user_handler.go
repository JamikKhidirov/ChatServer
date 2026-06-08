package handler

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	"ChatServerGolang/internal/domain"
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

func (h *UserHandler) GetProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	profile, err := h.userService.GetProfile(userID.(string))
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, profile)
}

func (h *UserHandler) UpdateProfile(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.UpdateProfileRequest
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

func (h *UserHandler) DeleteAccount(c *gin.Context) {
	userID, _ := c.Get("userID")

	if err := h.userService.DeleteAccount(userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "account deleted"})
}

func (h *UserHandler) SearchUsers(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		response.BadRequest(c, "search query required")
		return
	}

	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	users, err := h.userService.SearchUsers(query, limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, users)
}

func (h *UserHandler) GetUserByID(c *gin.Context) {
	userID := c.Param("id")

	user, err := h.userService.GetUserByID(userID)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.JSON(c, 200, user.ToResponse())
}

func (h *UserHandler) GetUserByUsername(c *gin.Context) {
	username := c.Param("username")

	user, err := h.userService.GetByUsername(username)
	if err != nil {
		response.NotFound(c, "user not found")
		return
	}

	response.JSON(c, 200, user)
}

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

func (h *UserHandler) UpdateStatus(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.UpdateStatusRequest
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

func (h *UserHandler) UpdatePushToken(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.UpdatePushTokenRequest
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

func (h *UserHandler) BlockUser(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.BlockUserRequest
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

func (h *UserHandler) UnblockUser(c *gin.Context) {
	userID, _ := c.Get("userID")
	blockedID := c.Param("userId")

	if err := h.userService.UnblockUser(userID.(string), blockedID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "user unblocked"})
}

func (h *UserHandler) GetBlockedUsers(c *gin.Context) {
	userID, _ := c.Get("userID")

	users, err := h.userService.GetBlockedUsers(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, users)
}

func (h *UserHandler) SetNotificationMuted(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req domain.UpdateNotificationSettingRequest
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

	updated, err := h.userService.UpdateProfile(userID.(string), &domain.UpdateProfileRequest{AvatarURL: avatarURL})
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, updated)
}

// Account settings
func (h *UserHandler) GetAccountSetting(c *gin.Context) {
	userID, _ := c.Get("userID")

	setting, err := h.userService.GetAccountSetting(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, setting)
}

func (h *UserHandler) UpdateAccountSetting(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.UpdateAccountSettingRequest
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
