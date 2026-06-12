package loginhandler

import (
	"ChatServerGolang/internal/domain/auth"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type LoginCodeHandler struct {
	verService  service.VerificationService
	authService service.AuthService
}

func NewLoginCodeHandler(verService service.VerificationService, authService service.AuthService) *LoginCodeHandler {
	return &LoginCodeHandler{verService: verService, authService: authService}
}

// SendEmailCode sends login code to email
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.LoginByEmailRequest true "Email address"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/login/email [post]
func (h *LoginCodeHandler) SendEmailCode(c *gin.Context) {
	var req authdomain.LoginByEmailRequest
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

// VerifyEmailCode verifies email login code and returns token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.LoginByEmailVerifyRequest true "Email and code"
// @Success 200 {object} object{token=string}
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/login/email/verify [post]
func (h *LoginCodeHandler) VerifyEmailCode(c *gin.Context) {
	var req authdomain.LoginByEmailVerifyRequest
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

// SendPhoneCode sends login code via SMS
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.LoginByPhoneRequest true "Phone number"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/login/phone [post]
func (h *LoginCodeHandler) SendPhoneCode(c *gin.Context) {
	var req authdomain.LoginByPhoneRequest
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

// VerifyPhoneCode verifies phone login code and returns token
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.LoginByPhoneVerifyRequest true "Phone and code"
// @Success 200 {object} object{token=string}
// @Failure 400 {object} response.ErrorResponse
// @Router /auth/login/phone/verify [post]
func (h *LoginCodeHandler) VerifyPhoneCode(c *gin.Context) {
	var req authdomain.LoginByPhoneVerifyRequest
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
