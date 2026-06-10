package handler

import (
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type GifHandler struct {
	gifService service.SavedGifService
}

func NewGifHandler(gifService service.SavedGifService) *GifHandler {
	return &GifHandler{gifService: gifService}
}

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

func (h *GifHandler) GetSavedGifs(c *gin.Context) {
	userID, _ := c.Get("userID")
	gifs, err := h.gifService.GetSavedGifs(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gifs)
}

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
