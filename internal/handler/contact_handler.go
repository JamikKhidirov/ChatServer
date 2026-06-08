package handler

import (
	"ChatServerGolang/internal/domain"
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

func (h *ContactHandler) SyncContacts(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req domain.SyncContactsRequest
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

func (h *ContactHandler) GetContacts(c *gin.Context) {
	userID, _ := c.Get("userID")

	contacts, err := h.contactService.GetContacts(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, contacts)
}

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

func (h *ContactHandler) FindRegistered(c *gin.Context) {
	userID, _ := c.Get("userID")

	users, err := h.contactService.FindRegisteredByPhone(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, users)
}
