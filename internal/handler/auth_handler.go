package handler

import (
	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService *service.AuthService
}

func NewAuthHandler(authService *service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

func (h *AuthHandler) Register(c *gin.Context) {
	var req domain.RegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.Register(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, resp)
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req domain.LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.Login(&req)
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.JSON(c, 200, resp)
}

func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, _ := c.Get("userID")

	token, err := h.authService.RefreshToken(userID.(string))
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.JSON(c, 200, domain.RefreshTokenResponse{Token: token})
}

func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.ChangePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.authService.ChangePassword(userID.(string), req.OldPassword, req.NewPassword); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "password changed"})
}
