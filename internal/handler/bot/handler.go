package bothandler

import (
	"ChatServerGolang/internal/domain/bot"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type BotHandler struct {
	botService service.BotService
}

func NewBotHandler(botService service.BotService) *BotHandler {
	return &BotHandler{botService: botService}
}

// CreateBot creates a new bot
// @Tags Bots
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body botdomain.CreateBotRequest true "Bot name and config"
// @Success 201 {object} botdomain.Bot
// @Failure 400 {object} response.ErrorResponse
// @Router /bots [post]
func (h *BotHandler) CreateBot(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req botdomain.CreateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	bot, err := h.botService.CreateBot(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 201, bot)
}

// GetMyBots returns all bots owned by the authenticated user
// @Tags Bots
// @Security BearerAuth
// @Produce json
// @Success 200 {array} botdomain.Bot
// @Failure 400 {object} response.ErrorResponse
// @Router /bots [get]
func (h *BotHandler) GetMyBots(c *gin.Context) {
	userID, _ := c.Get("userID")
	bots, err := h.botService.GetMyBots(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, bots)
}

// UpdateBot modifies a bot's settings
// @Tags Bots
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Bot ID"
// @Param request body botdomain.UpdateBotRequest true "Fields to update"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /bots/{id} [put]
func (h *BotHandler) UpdateBot(c *gin.Context) {
	userID, _ := c.Get("userID")
	botID := c.Param("id")
	var req botdomain.UpdateBotRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.botService.UpdateBot(botID, userID.(string), &req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "bot updated"})
}

// DeleteBot deletes a bot
// @Tags Bots
// @Security BearerAuth
// @Produce json
// @Param id path string true "Bot ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /bots/{id} [delete]
func (h *BotHandler) DeleteBot(c *gin.Context) {
	userID, _ := c.Get("userID")
	botID := c.Param("id")
	if err := h.botService.DeleteBot(botID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "bot deleted"})
}

// RegenerateToken generates a new API token for a bot
// @Tags Bots
// @Security BearerAuth
// @Produce json
// @Param id path string true "Bot ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /bots/{id}/token [post]
func (h *BotHandler) RegenerateToken(c *gin.Context) {
	userID, _ := c.Get("userID")
	botID := c.Param("id")
	if err := h.botService.RegenerateToken(botID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "token regenerated"})
}
