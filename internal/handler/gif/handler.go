package gifhandler

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

// SaveGif saves a GIF URL to the user's collection
// @Tags Gifs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{url=string} true "GIF URL"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Gifs
// @Security BearerAuth
// @Produce json
// @Success 200 {array} string
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Gifs
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{url=string} true "GIF URL to remove"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
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
