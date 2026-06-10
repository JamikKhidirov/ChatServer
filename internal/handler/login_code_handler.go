package handler

import (
	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginCodeHandler struct {
	verService service.VerificationService
	authService service.AuthService
}

func NewLoginCodeHandler(verService service.VerificationService, authService service.AuthService) *LoginCodeHandler {
	return &LoginCodeHandler{verService: verService, authService: authService}
}

func (h *LoginCodeHandler) SendEmailCode(c *gin.Context) {
	var req domain.LoginByEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if _, err := h.verService.LoginSendEmailCode(req.Email); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "code sent"})
}

func (h *LoginCodeHandler) VerifyEmailCode(c *gin.Context) {
	var req domain.LoginByEmailVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, err := h.verService.LoginVerifyEmailCode(req.Email, req.Code)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	token, err := h.authService.RefreshToken(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"token": token})
}

func (h *LoginCodeHandler) SendPhoneCode(c *gin.Context) {
	var req domain.LoginByPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if _, err := h.verService.LoginSendPhoneCode(req.Phone); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "code sent"})
}

func (h *LoginCodeHandler) VerifyPhoneCode(c *gin.Context) {
	var req domain.LoginByPhoneVerifyRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	userID, err := h.verService.LoginVerifyPhoneCode(req.Phone, req.Code)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	token, err := h.authService.RefreshToken(userID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"token": token})
}
