package verhandler

import (
	"net/http"

	"ChatServerGolang/backend/internal/domain/verification"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type VerificationHandler struct {
	verService service.VerificationService
}

func NewVerificationHandler(verService service.VerificationService) *VerificationHandler {
	return &VerificationHandler{verService: verService}
}

// SendEmail sends email verification code
// @Summary Отправить код верификации email
// @Description Отправляет код подтверждения на email для верификации адреса электронной почты.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body verificationdomain.SendEmailVerificationRequest true "Email address"
// @Success 200 {object} response.MessageResponse
// @Failure 500 {object} response.ErrorResponse "Failed to send"
// @Router /verification/email [post]
func (h *VerificationHandler) SendEmail(c *gin.Context) {
	userID := c.GetString("userID")
	var req verificationdomain.SendEmailVerificationRequest
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

// VerifyEmail verifies email with code
// @Summary Подтвердить email
// @Description Подтверждает email-адрес с помощью кода, отправленного на почту.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body verificationdomain.VerifyEmailRequest true "Verification code"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Wrong code"
// @Router /verification/email/verify [post]
func (h *VerificationHandler) VerifyEmail(c *gin.Context) {
	userID := c.GetString("userID")
	var req verificationdomain.VerifyEmailRequest
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

// SendPhone sends phone verification SMS
// @Summary Отправить код верификации телефона
// @Description Отправляет SMS с кодом подтверждения на указанный номер телефона.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body verificationdomain.SendPhoneVerificationRequest true "Phone number"
// @Success 200 {object} response.MessageResponse
// @Failure 500 {object} response.ErrorResponse "Failed to send"
// @Router /verification/phone [post]
func (h *VerificationHandler) SendPhone(c *gin.Context) {
	userID := c.GetString("userID")
	var req verificationdomain.SendPhoneVerificationRequest
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

// VerifyPhone verifies phone with code
// @Summary Подтвердить телефон
// @Description Подтверждает номер телефона с помощью кода из SMS.
// @Tags Users
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body verificationdomain.VerifyPhoneRequest true "Verification code"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Wrong code"
// @Router /verification/phone/verify [post]
func (h *VerificationHandler) VerifyPhone(c *gin.Context) {
	userID := c.GetString("userID")
	var req verificationdomain.VerifyPhoneRequest
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
