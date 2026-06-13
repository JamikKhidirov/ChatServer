package userhandler

import (
	"io"
	"os"
	"path/filepath"
	"strconv"

	userdomain "ChatServerGolang/backend/internal/domain/user"
	notificationdomain "ChatServerGolang/backend/internal/domain/notification"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

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
// @Tags Профиль
// @Summary Получить профиль пользователя
// @Security BearerAuth
// @Produce json
// @Description Возвращает полную информацию о профиле аутентифицированного пользователя, включая имя, email, аватар и настройки.
// @Success 200 {object} userdomain.UserResponse "Профиль пользователя успешно получен"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
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
// @Tags Профиль
// @Summary Обновить профиль пользователя
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Обновляет поля профиля аутентифицированного пользователя. Можно изменить отображаемое имя, биографию, аватар и другие настройки профиля.
// @Param request body userdomain.UpdateProfileRequest true "Обновляемые поля профиля: display_name (отображаемое имя, опционально), bio (биография, опционально), avatar_url (URL аватара, опционально)"
// @Success 200 {object} userdomain.UserResponse "Профиль успешно обновлён"
// @Failure 400 {object} response.ErrorResponse "Неверные входные данные"
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
// @Tags Профиль
// @Summary Удалить аккаунт
// @Security BearerAuth
// @Produce json
// @Description Полностью удаляет учётную запись пользователя и все связанные данные. Это действие необратимо.
// @Success 200 {object} response.MessageResponse "Аккаунт успешно удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка при удалении аккаунта"
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
// @Tags Профиль
// @Summary Найти пользователей
// @Security BearerAuth
// @Produce json
// @Description Ищет пользователей по имени, username или email. Возвращает постраничный список результатов.
// @Param q query string true "Поисковый запрос (имя, username или email)"
// @Param limit query int false "Максимум результатов (по умолчанию 50)"
// @Param offset query int false "Смещение пагинации (по умолчанию 0)"
// @Success 200 {object} response.APIResponse "Постраничный список пользователей"
// @Failure 400 {object} response.ErrorResponse "Отсутствует поисковый запрос"
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
// @Tags Профиль
// @Summary Получить пользователя по ID
// @Security BearerAuth
// @Produce json
// @Description Возвращает информацию о пользователе по его уникальному идентификатору.
// @Param id path string true "ID пользователя"
// @Success 200 {object} userdomain.UserResponse "Информация о пользователе"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
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
// @Tags Профиль
// @Summary Получить пользователя по username
// @Security BearerAuth
// @Produce json
// @Description Возвращает информацию о пользователе по его уникальному имени пользователя (username).
// @Param username path string true "Username пользователя"
// @Success 200 {object} userdomain.UserResponse "Информация о пользователе"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
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
// @Tags Профиль
// @Summary Отправить тестовое push-уведомление
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Отправляет тестовое push-уведомление на устройство пользователя для проверки настроек уведомлений.
// @Param request body object{title=string,body=string} true "Параметры уведомления: title (заголовок, обязательно), body (текст, обязательно)"
// @Success 200 {object} response.MessageResponse "Тестовое уведомление отправлено"
// @Failure 400 {object} response.ErrorResponse "Неверные входные данные"
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
// @Tags Профиль
// @Summary Обновить статус пользователя
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Изменяет статус активности пользователя (online, offline, busy, away).
// @Param request body userdomain.UpdateStatusRequest true "Новый статус: status (строка, статус: online/offline/busy/away, обязательно)"
// @Success 200 {object} userdomain.UserResponse "Статус успешно обновлён"
// @Failure 400 {object} response.ErrorResponse "Неверный статус"
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
// @Tags Профиль
// @Summary Обновить push-токен
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Обновляет токен push-уведомлений для устройства пользователя. Необходимо для получения уведомлений на мобильные устройства.
// @Param request body userdomain.UpdatePushTokenRequest true "Данные push-токена: token (строка, токен устройства, обязательно), provider (строка, провайдер: fcm/apns, обязательно)"
// @Success 200 {object} response.MessageResponse "Push-токен успешно обновлён"
// @Failure 400 {object} response.ErrorResponse "Неверный токен"
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
// @Tags Профиль
// @Summary Заблокировать пользователя
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Блокирует указанного пользователя. Заблокированный пользователь не сможет отправлять вам сообщения или звонить.
// @Param request body userdomain.BlockUserRequest true "Данные для блокировки: blocked_id (ID пользователя для блокировки, обязательно)"
// @Success 200 {object} response.MessageResponse "Пользователь заблокирован"
// @Failure 400 {object} response.ErrorResponse "Ошибка блокировки"
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
// @Tags Профиль
// @Summary Разблокировать пользователя
// @Security BearerAuth
// @Produce json
// @Description Снимает блокировку с указанного пользователя, восстанавливая возможность общения.
// @Param userId path string true "ID пользователя для разблокировки"
// @Success 200 {object} response.MessageResponse "Пользователь разблокирован"
// @Failure 400 {object} response.ErrorResponse "Ошибка разблокировки"
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
// @Tags Профиль
// @Summary Получить список заблокированных пользователей
// @Security BearerAuth
// @Produce json
// @Description Возвращает список всех пользователей, заблокированных аутентифицированным пользователем.
// @Success 200 {array} userdomain.UserResponse "Список заблокированных пользователей"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения списка"
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
// @Tags Профиль
// @Summary Включить или отключить звук уведомлений чата
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Позволяет отключить или включить звук уведомлений для указанной беседы.
// @Param id path string true "ID чата"
// @Param request body notificationdomain.UpdateNotificationSettingRequest true "Настройки звука: muted (boolean, true — отключить звук, false — включить, обязательно)"
// @Success 200 {object} response.MessageResponse "Настройки уведомлений обновлены"
// @Failure 400 {object} response.ErrorResponse "Ошибка обновления настроек"
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
// @Tags Профиль
// @Summary Проверить статус уведомлений чата
// @Security BearerAuth
// @Produce json
// @Description Проверяет, отключён ли звук уведомлений для указанного чата.
// @Param id path string true "ID чата"
// @Success 200 {object} object{muted=boolean} "Статус уведомлений: muted (true — звук отключён, false — включён)"
// @Failure 400 {object} response.ErrorResponse "Ошибка проверки статуса"
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
// @Tags Профиль
// @Summary Загрузить аватар
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Загружает новое изображение аватара для профиля пользователя. Поддерживаются изображения в форматах JPG, PNG и WEBP.
// @Param avatar formData file true "Файл изображения аватара (JPG, PNG, WEBP)"
// @Success 200 {object} userdomain.UserResponse "Аватар успешно загружен и профиль обновлён"
// @Failure 400 {object} response.ErrorResponse "Файл отсутствует или неверный формат"
// @Failure 500 {object} response.ErrorResponse "Внутренняя ошибка сервера"
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
// @Tags Профиль
// @Summary Получить настройки аккаунта
// @Security BearerAuth
// @Produce json
// @Description Возвращает текущие настройки учётной записи пользователя, включая параметры приватности и уведомлений.
// @Success 200 {object} userdomain.AccountSetting "Настройки аккаунта"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения настроек"
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
// @Tags Профиль
// @Summary Получить время последнего посещения
// @Security BearerAuth
// @Produce json
// @Description Возвращает информацию о последнем посещении пользователя: статус онлайн и временную метку.
// @Param id path string true "ID пользователя"
// @Success 200 {object} object{userId=string,online=boolean,lastSeen=string} "Информация о последнем посещении: userId (ID), online (онлайн ли), lastSeen (время последнего визита)"
// @Failure 404 {object} response.ErrorResponse "Пользователь не найден"
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
// @Tags Профиль
// @Summary Изменить username
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Изменяет имя пользователя (username). Длина username должна быть от 3 до 32 символов.
// @Param request body object{username=string} true "Новый username (строка, от 3 до 32 символов, обязательно)"
// @Success 200 {object} response.MessageResponse "Username успешно изменён"
// @Failure 400 {object} response.ErrorResponse "Ошибка: неверный формат или username уже занят"
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
// @Tags Профиль
// @Summary Изменить email
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Изменяет адрес электронной почты пользователя. Новый email должен быть валидным и не занятым другим пользователем.
// @Param request body object{email=string} true "Новый email (строка, валидный email, обязательно)"
// @Success 200 {object} response.MessageResponse "Email успешно изменён"
// @Failure 400 {object} response.ErrorResponse "Ошибка: неверный формат email или адрес уже занят"
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
// @Tags Профиль
// @Summary Обновить настройки аккаунта
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Обновляет настройки учётной записи пользователя, такие как приватность, язык интерфейса и параметры уведомлений.
// @Param request body userdomain.UpdateAccountSettingRequest true "Обновляемые настройки: language (язык, опционально), privacy (настройки приватности, опционально)"
// @Success 200 {object} userdomain.AccountSetting "Настройки успешно обновлены"
// @Failure 400 {object} response.ErrorResponse "Неверные настройки"
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
