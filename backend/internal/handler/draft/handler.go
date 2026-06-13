package drafthandler

import (
	"ChatServerGolang/backend/internal/domain/draft"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type DraftHandler struct {
	draftService service.DraftService
}

func NewDraftHandler(draftService service.DraftService) *DraftHandler {
	return &DraftHandler{draftService: draftService}
}

// SaveDraft saves a message draft for a chat
// @Tags Черновики
// @Summary Сохранить черновик сообщения
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Сохраняет черновик сообщения для указанного чата. Черновик можно будет позже восстановить и отредактировать.
// @Param request body draftdomain.SaveDraftRequest true "Данные черновика: chat_id (ID чата, обязательно), content (текст сообщения, обязательно), reply_to_id (ID сообщения для ответа, опционально)"
// @Success 200 {object} draftdomain.Draft "Черновик сохранён"
// @Failure 400 {object} response.ErrorResponse "Ошибка сохранения черновика"
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
// @Tags Черновики
// @Summary Получить черновик чата
// @Security BearerAuth
// @Produce json
// @Description Возвращает сохранённый черновик сообщения для указанного чата.
// @Param chatId query string true "ID чата"
// @Success 200 {object} draftdomain.Draft "Черновик сообщения"
// @Failure 404 {object} response.ErrorResponse "Черновик не найден"
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
// @Tags Черновики
// @Summary Удалить черновик
// @Security BearerAuth
// @Produce json
// @Description Удаляет сохранённый черновик сообщения по его идентификатору.
// @Param id path string true "ID черновика"
// @Success 200 {object} response.MessageResponse "Черновик удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления черновика"
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
