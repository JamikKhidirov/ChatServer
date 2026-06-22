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
// @Tags Боты
// @Summary Создать бота
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Создаёт нового бота с указанным именем и конфигурацией. Владелец бота может управлять его настройками.
// @Param request body botdomain.CreateBotRequest true "Данные бота: name (имя, обязательно), description (описание, опционально), avatar_url (URL аватара, опционально)"
// @Success 201 {object} botdomain.Bot "Бот создан"
// @Failure 400 {object} response.ErrorResponse "Ошибка создания бота"
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
// @Tags Боты
// @Summary Получить моих ботов
// @Security BearerAuth
// @Produce json
// @Description Возвращает список всех ботов, созданных аутентифицированным пользователем.
// @Success 200 {array} botdomain.Bot "Список ботов"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения ботов"
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
// @Tags Боты
// @Summary Обновить настройки бота
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Изменяет настройки существующего бота, такие как имя, описание или аватар. Доступно только владельцу.
// @Param id path string true "ID бота"
// @Param request body botdomain.UpdateBotRequest true "Обновляемые поля: name (имя, опционально), description (описание, опционально)"
// @Success 200 {object} response.MessageResponse "Бот обновлён"
// @Failure 400 {object} response.ErrorResponse "Ошибка обновления бота"
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
// @Tags Боты
// @Summary Удалить бота
// @Security BearerAuth
// @Produce json
// @Description Удаляет бота и все связанные с ним данные. Доступно только владельцу бота.
// @Param id path string true "ID бота"
// @Success 200 {object} response.MessageResponse "Бот удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления бота"
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
// @Tags Боты
// @Summary Перегенерировать токен бота
// @Security BearerAuth
// @Produce json
// @Description Создаёт новый API-токен для бота, аннулируя старый. Используется при компрометации текущего токена.
// @Param id path string true "ID бота"
// @Success 200 {object} response.MessageResponse "Токен перегенерирован"
// @Failure 400 {object} response.ErrorResponse "Ошибка генерации токена"
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
