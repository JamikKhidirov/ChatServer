package contacthandler

import (
	"ChatServerGolang/internal/domain/contact"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type ContactHandler struct {
	contactService service.ContactService
}

func NewContactHandler(contactService service.ContactService) *ContactHandler {
	return &ContactHandler{contactService: contactService}
}

// SyncContacts synchronizes the user's phone contacts
// @Tags Contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body contactdomain.SyncContactsRequest true "Contact list to sync"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Contacts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} contactdomain.ContactResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Contacts
// @Security BearerAuth
// @Produce json
// @Param q query string true "Phone number or partial to search"
// @Success 200 {array} contactdomain.ContactResponse
// @Failure 400 {object} response.ErrorResponse "Missing query"
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
// @Tags Contacts
// @Security BearerAuth
// @Produce json
// @Success 200 {array} userdomain.UserResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Contacts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body contactdomain.UpdateContactPhotoRequest true "Phone and photo URL"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
