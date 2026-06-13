package authhandler

import (
	"ChatServerGolang/backend/internal/domain/auth"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	authService service.AuthService
}

func NewAuthHandler(authService service.AuthService) *AuthHandler {
	return &AuthHandler{authService: authService}
}

// RegisterAdmin creates a new admin user account
// @Tags Аутентификация
// @Summary Зарегистрировать администратора
// @Accept json
// @Produce json
// @Description Создаёт учётную запись администратора с повышенными правами. Требует указания секретного ключа администратора, а также email, username и password.
// @Param request body authdomain.AdminRegisterRequest true "Данные для регистрации администратора: username (логин, обязательно), email (почта, обязательно), password (пароль, обязательно), secret (секретный ключ администратора, обязательно), display_name (отображаемое имя, опционально)"
// @Success 201 {object} authdomain.AuthResponse "Администратор создан, возвращается JWT-токен и информация о пользователе с правами администратора"
// @Failure 400 {object} response.ErrorResponse "Неверный секретный ключ, email или username уже заняты, или неверные входные данные"
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
// @Tags Аутентификация
// @Summary Зарегистрировать новый аккаунт
// @Accept json
// @Produce json
// @Description Создаёт нового пользователя с указанными email, username и паролем. После успешной регистрации возвращается JWT-токен для авторизации.
// @Param request body authdomain.RegisterRequest true "Данные для регистрации: username (логин, обязательно), email (почта, обязательно), password (пароль, обязательно), display_name (отображаемое имя, опционально)"
// @Success 201 {object} authdomain.AuthResponse "Аккаунт создан, возвращается JWT-токен и информация о пользователе"
// @Failure 400 {object} response.ErrorResponse "Email или username уже зарегистрированы, или неверные поля ввода"
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
// @Tags Аутентификация
// @Summary Войти в аккаунт
// @Accept json
// @Produce json
// @Description Аутентифицирует пользователя по email и паролю. Возвращает JWT-токен для последующих запросов, а также информацию о пользователе.
// @Param request body authdomain.LoginRequest true "Учётные данные: email (почта, обязательно) и password (пароль, обязательно)"
// @Success 200 {object} authdomain.AuthResponse "Успешный вход, возвращается JWT-токен и данные пользователя"
// @Failure 401 {object} response.ErrorResponse "Неверный email или пароль"
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
// @Tags Аутентификация
// @Summary Обновить JWT-токен
// @Security BearerAuth
// @Produce json
// @Description Обновляет срок действия текущего JWT-токена для аутентифицированного пользователя. Требует действительный токен в заголовке Authorization.
// @Success 200 {object} authdomain.RefreshTokenResponse "Новый JWT-токен успешно сгенерирован"
// @Failure 401 {object} response.ErrorResponse "Недействительный или просроченный токен"
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
// @Tags Аутентификация
// @Summary Изменить пароль
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Обновляет пароль аутентифицированного пользователя. Требуется указать старый пароль для проверки и новый пароль.
// @Param request body authdomain.ChangePasswordRequest true "Параметры смены пароля: old_password (старый пароль, обязательно), new_password (новый пароль, обязательно)"
// @Success 200 {object} response.MessageResponse "Пароль успешно изменён"
// @Failure 400 {object} response.ErrorResponse "Неверные данные или неправильный старый пароль"
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
