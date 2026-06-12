package drafthandler

import (
	"ChatServerGolang/internal/domain/draft"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type DraftHandler struct {
	draftService service.DraftService
}

func NewDraftHandler(draftService service.DraftService) *DraftHandler {
	return &DraftHandler{draftService: draftService}
}

// SaveDraft saves a message draft for a chat
// @Tags Drafts
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body draftdomain.SaveDraftRequest true "Chat ID and draft content"
// @Success 200 {object} draftdomain.Draft
// @Failure 400 {object} response.ErrorResponse
// @Router /drafts [post]
func (h *DraftHandler) SaveDraft(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req draftdomain.SaveDraftRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	draft, err := h.draftService.SaveDraft(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, draft)
}

// GetDraft returns the draft for a specific chat
// @Tags Drafts
// @Security BearerAuth
// @Produce json
// @Param chatId query string true "Chat ID"
// @Success 200 {object} draftdomain.Draft
// @Failure 404 {object} response.ErrorResponse "No draft"
// @Router /drafts [get]
func (h *DraftHandler) GetDraft(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Query("chatId")
	if chatID == "" {
		response.BadRequest(c, "chatId query parameter required")
		return
	}
	draft, err := h.draftService.GetDraft(userID.(string), chatID)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.JSON(c, 200, draft)
}

// DeleteDraft deletes a saved draft
// @Tags Drafts
// @Security BearerAuth
// @Produce json
// @Param id path string true "Draft ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /drafts/{id} [delete]
func (h *DraftHandler) DeleteDraft(c *gin.Context) {
	userID, _ := c.Get("userID")
	draftID := c.Param("id")
	if err := h.draftService.DeleteDraft(userID.(string), draftID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "draft deleted"})
}
