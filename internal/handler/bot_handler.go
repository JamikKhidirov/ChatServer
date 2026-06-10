package handler

import (
	"ChatServerGolang/internal/domain"
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

func (h *BotHandler) CreateBot(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.CreateBotRequest
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

func (h *BotHandler) GetMyBots(c *gin.Context) {
	userID, _ := c.Get("userID")
	bots, err := h.botService.GetMyBots(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, bots)
}

func (h *BotHandler) UpdateBot(c *gin.Context) {
	userID, _ := c.Get("userID")
	botID := c.Param("id")
	var req domain.UpdateBotRequest
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

func (h *BotHandler) DeleteBot(c *gin.Context) {
	userID, _ := c.Get("userID")
	botID := c.Param("id")
	if err := h.botService.DeleteBot(botID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "bot deleted"})
}

func (h *BotHandler) RegenerateToken(c *gin.Context) {
	userID, _ := c.Get("userID")
	botID := c.Param("id")
	if err := h.botService.RegenerateToken(botID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "token regenerated"})
}
