package authhandler

import (
	"ChatServerGolang/internal/domain/auth"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterAdmin creates a new admin user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.AdminRegisterRequest true "Admin registration details with secret key"
// @Success 201 {object} authdomain.AuthResponse "Returns JWT token + user with isAdmin:true"
// @Failure 400 {object} response.ErrorResponse "Invalid admin secret, email/username taken, or invalid input"
// @Router /auth/admin/register [post]
func (h *AuthHandler) RegisterAdmin(c *gin.Context) {
	var req authdomain.AdminRegisterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	resp, err := h.authService.RegisterAdmin(&req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, resp)
}

// Register creates a new user account
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.RegisterRequest true "Registration details (username, email, password, display_name)"
// @Success 201 {object} authdomain.AuthResponse "Returns JWT token + user object"
// @Failure 400 {object} response.ErrorResponse "Email or username already registered, or invalid input fields"
// @Router /auth/register [post]
func (h *AuthHandler) Register(c *gin.Context) {
	var req authdomain.RegisterRequest
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

// Login authenticates user credentials
// @Tags Auth
// @Accept json
// @Produce json
// @Param request body authdomain.LoginRequest true "Login credentials"
// @Success 200 {object} authdomain.AuthResponse
// @Failure 401 {object} response.ErrorResponse "Invalid credentials"
// @Router /auth/login [post]
func (h *AuthHandler) Login(c *gin.Context) {
	var req authdomain.LoginRequest
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

// RefreshToken returns a new JWT token for the authenticated user
// @Tags Auth
// @Security BearerAuth
// @Produce json
// @Success 200 {object} authdomain.RefreshTokenResponse
// @Failure 401 {object} response.ErrorResponse "Invalid or expired token"
// @Router /auth/refresh [get]
func (h *AuthHandler) RefreshToken(c *gin.Context) {
	userID, _ := c.Get("userID")

	token, err := h.authService.RefreshToken(userID.(string))
	if err != nil {
		response.Unauthorized(c, err.Error())
		return
	}

	response.JSON(c, 200, authdomain.RefreshTokenResponse{Token: token})
}

// ChangePassword updates the authenticated user's password
// @Tags Auth
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body authdomain.ChangePasswordRequest true "Old and new passwords"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Invalid input or wrong password"
// @Router /auth/change-password [put]
func (h *AuthHandler) ChangePassword(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req authdomain.ChangePasswordRequest
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
