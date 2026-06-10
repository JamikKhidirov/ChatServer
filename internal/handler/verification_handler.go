package handler

import (
	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	verService service.VerificationService
}

func NewVerificationHandler(verService service.VerificationService) *VerificationHandler {
	return &VerificationHandler{verService: verService}
}

func (h *VerificationHandler) SendEmail(c *gin.Context) {
	userID := c.GetString("userID")
	var req domain.SendEmailVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.verService.SendEmailVerification(userID, req.Email); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "code sent"})
}

func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	userID := c.GetString("userID")
	var req domain.VerifyEmailRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.verService.VerifyEmail(userID, req.Code); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "email verified"})
}

func (h *VerificationHandler) SendPhone(c *gin.Context) {
	userID := c.GetString("userID")
	var req domain.SendPhoneVerificationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.verService.SendPhoneVerification(userID, req.Phone); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "code sent"})
}

func (h *VerificationHandler) VerifyPhone(c *gin.Context) {
	userID := c.GetString("userID")
	var req domain.VerifyPhoneRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.verService.VerifyPhone(userID, req.Code); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "phone verified"})
}
