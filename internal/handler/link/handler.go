package linkhandler

import (
	chatdomain "ChatServerGolang/internal/domain/chat"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type InviteLinkHandler struct {
	linkService service.InviteLinkService
}

func NewInviteLinkHandler(linkService service.InviteLinkService) *InviteLinkHandler {
	return &InviteLinkHandler{linkService: linkService}
}

// CreateInviteLink creates a new invite link for a chat
// @Summary Создать приглашение
// @Description Создаёт пригласительную ссылку для чата с опциональным сроком действия и лимитом использования.
// @Tags InviteLinks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body chatdomain.CreateInviteLinkRequest true "Optional expiration and usage limit"
// @Success 201 {object} chatdomain.InviteLink
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/invite-links [post]
func (h *InviteLinkHandler) CreateInviteLink(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req chatdomain.CreateInviteLinkRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	link, err := h.linkService.CreateInviteLink(chatID, userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, link)
}

// GetInviteLinks returns all invite links for a chat
// @Summary Список приглашений
// @Description Возвращает все пригласительные ссылки для указанного чата.
// @Tags InviteLinks
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {array} chatdomain.InviteLink
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/invite-links [get]
func (h *InviteLinkHandler) GetInviteLinks(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	links, err := h.linkService.GetInviteLinks(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, links)
}

// DeleteInviteLink deletes an invite link
// @Summary Удалить приглашение
// @Description Удаляет пригласительную ссылку. Существующие участники, присоединившиеся по ней, остаются в чате.
// @Tags InviteLinks
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Param linkId path string true "Link ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/{id}/invite-links/{linkId} [delete]
func (h *InviteLinkHandler) DeleteInviteLink(c *gin.Context) {
	userID, _ := c.Get("userID")
	linkID := c.Param("linkId")

	if err := h.linkService.DeleteInviteLink(linkID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "invite link deleted"})
}

// JoinByInviteLink joins a chat using an invite code
// @Summary Присоединиться по ссылке
// @Description Присоединяет пользователя к чату по пригласительному коду/ссылке.
// @Tags InviteLinks
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body chatdomain.JoinByInviteRequest true "Invite code"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /chats/join [post]
func (h *InviteLinkHandler) JoinByInviteLink(c *gin.Context) {
	userID, _ := c.Get("userID")

	var req chatdomain.JoinByInviteRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	if err := h.linkService.JoinByInviteLink(req.Code, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, gin.H{"message": "joined chat successfully"})
}
