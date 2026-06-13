package gifhandler

import (
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type GifHandler struct {
	gifService service.SavedGifService
}

func NewGifHandler(gifService service.SavedGifService) *GifHandler {
	return &GifHandler{gifService: gifService}
}

// SaveGif saves a GIF URL to the user's collection
// @Tags GIF
// @Summary Сохранить GIF
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Сохраняет URL GIF-изображения в коллекцию пользователя для быстрого доступа при отправке.
// @Param request body object{url=string} true "Параметры: url (URL GIF-изображения, обязательно)"
// @Success 200 {object} response.MessageResponse "GIF сохранён"
// @Failure 400 {object} response.ErrorResponse "Ошибка сохранения GIF"
// @Router /gifs [post]
func (h *GifHandler) SaveGif(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.gifService.SaveGif(userID.(string), req.URL); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "gif saved"})
}

// GetSavedGifs returns all saved GIFs for the user
// @Tags GIF
// @Summary Получить сохранённые GIF
// @Security BearerAuth
// @Produce json
// @Description Возвращает список URL всех GIF-изображений, сохранённых пользователем.
// @Success 200 {array} string "Список URL сохранённых GIF"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения GIF"
// @Router /gifs [get]
func (h *GifHandler) GetSavedGifs(c *gin.Context) {
	userID, _ := c.Get("userID")
	gifs, err := h.gifService.GetSavedGifs(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gifs)
}

// DeleteGif removes a GIF from the user's collection
// @Tags GIF
// @Summary Удалить GIF из коллекции
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Удаляет GIF-изображение из коллекции пользователя по его URL.
// @Param request body object{url=string} true "Параметры: url (URL GIF для удаления, обязательно)"
// @Success 200 {object} response.MessageResponse "GIF удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления GIF"
// @Router /gifs [delete]
func (h *GifHandler) DeleteGif(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		URL string `json:"url" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.gifService.DeleteGif(userID.(string), req.URL); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "gif deleted"})
}
