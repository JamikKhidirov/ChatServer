package handler

import (
	"strconv"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService *service.UserService
	pushService *service.PushService
}

func NewUserHandler(userService *service.UserService, pushService *service.PushService) *UserHandler {
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
		response.InternalError(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "push token updated"})
}
