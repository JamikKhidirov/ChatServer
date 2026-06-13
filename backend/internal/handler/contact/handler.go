package contacthandler

import (
	"ChatServerGolang/backend/internal/domain/contact"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactService service.ContactService
}

func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

// SyncContacts synchronizes the user's phone contacts
// @Tags Контакты
// @Summary Синхронизировать контакты телефона
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Синхронизирует телефонную книгу пользователя с сервером. Позволяет найти друзей, уже зарегистрированных в приложении.
// @Param request body contactdomain.SyncContactsRequest true "Список контактов для синхронизации: contacts (массив объектов с phone и name, обязательно)"
// @Success 200 {object} response.MessageResponse "Контакты успешно синхронизированы"
// @Failure 400 {object} response.ErrorResponse "Ошибка синхронизации"
// @Router /contacts/sync [post]
func (h *ContactHandler) SyncContacts(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req contactdomain.SyncContactsRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.contactService.SyncContacts(userID.(string), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "contacts synced"})
}

// GetContacts returns all synced contacts for the user
// @Tags Контакты
// @Summary Получить список контактов
// @Security BearerAuth
// @Produce json
// @Description Возвращает все синхронизированные контакты пользователя с информацией о регистрации на платформе.
// @Success 200 {array} contactdomain.ContactResponse "Список контактов пользователя"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения контактов"
// @Router /contacts [get]
func (h *ContactHandler) GetContacts(c *gin.Context) {
	userID, _ := c.Get("userID")

	contacts, err := h.contactService.GetContacts(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, contacts)
}

// SearchByPhone searches contacts by phone number
// @Tags Контакты
// @Summary Найти контакты по номеру телефона
// @Security BearerAuth
// @Produce json
// @Description Ищет контакты пользователя по номеру телефона. Поддерживает частичный поиск.
// @Param q query string true "Номер телефона или его часть для поиска"
// @Success 200 {array} contactdomain.ContactResponse "Найденные контакты"
// @Failure 400 {object} response.ErrorResponse "Отсутствует поисковый запрос"
// @Router /contacts/search [get]
func (h *ContactHandler) SearchByPhone(c *gin.Context) {
	userID, _ := c.Get("userID")
	query := c.Query("q")

	if query == "" {
		response.BadRequest(c, "search query required")
		return
	}

	contacts, err := h.contactService.SearchByPhone(userID.(string), query)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, contacts)
}

// FindRegistered returns contacts that are registered on the platform
// @Tags Контакты
// @Summary Найти зарегистрированных пользователей среди контактов
// @Security BearerAuth
// @Produce json
// @Description Проверяет, какие из контактов телефонной книги пользователя уже зарегистрированы на платформе.
// @Success 200 {array} userdomain.UserResponse "Список зарегистрированных пользователей из контактов"
// @Failure 400 {object} response.ErrorResponse "Ошибка поиска"
// @Router /contacts/registered [get]
func (h *ContactHandler) FindRegistered(c *gin.Context) {
	userID, _ := c.Get("userID")

	users, err := h.contactService.FindRegisteredByPhone(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, users)
}

// UpdateContactPhoto updates the photo associated with a contact
// @Tags Контакты
// @Summary Обновить фото контакта
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Обновляет фотографию, связанную с контактом из телефонной книги.
// @Param request body contactdomain.UpdateContactPhotoRequest true "Данные для обновления: phone (номер телефона, обязательно), photo_url (URL фотографии, обязательно)"
// @Success 200 {object} response.MessageResponse "Фото контакта обновлено"
// @Failure 400 {object} response.ErrorResponse "Ошибка обновления фото"
// @Router /contacts/photo [post]
func (h *ContactHandler) UpdateContactPhoto(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req contactdomain.UpdateContactPhotoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.contactService.UpdateContactPhoto(userID.(string), req.Phone, req.PhotoURL); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "contact photo updated"})
}
