package handler

import (
	"ChatServerGolang/internal/domain"
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

func (h *DraftHandler) SaveDraft(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.SaveDraftRequest
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

func (h *DraftHandler) DeleteDraft(c *gin.Context) {
	userID, _ := c.Get("userID")
	draftID := c.Param("id")
	if err := h.draftService.DeleteDraft(userID.(string), draftID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "draft deleted"})
}
