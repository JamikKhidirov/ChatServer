package stickerhandler

import (
	"ChatServerGolang/internal/domain/sticker"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type StickerHandler struct {
	stickerService service.StickerService
}

func NewStickerHandler(stickerService service.StickerService) *StickerHandler {
	return &StickerHandler{stickerService: stickerService}
}

// CreatePack creates a new sticker pack
// @Tags Stickers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body stickerdomain.CreateStickerPackRequest true "Pack name and stickers"
// @Success 201 {object} stickerdomain.StickerPack
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/packs [post]
func (h *StickerHandler) CreatePack(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req stickerdomain.CreateStickerPackRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	pack, err := h.stickerService.CreatePack(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 201, pack)
}

// ListPacks returns all public sticker packs
// @Tags Stickers
// @Security BearerAuth
// @Produce json
// @Success 200 {array} stickerdomain.StickerPack
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/packs [get]
func (h *StickerHandler) ListPacks(c *gin.Context) {
	packs, err := h.stickerService.GetPacks()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, packs)
}

// GetMyPacks returns the authenticated user's sticker packs
// @Tags Stickers
// @Security BearerAuth
// @Produce json
// @Success 200 {array} stickerdomain.StickerPack
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/packs/my [get]
func (h *StickerHandler) GetMyPacks(c *gin.Context) {
	userID, _ := c.Get("userID")
	packs, err := h.stickerService.GetMyPacks(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, packs)
}

// GetPack returns a sticker pack by ID
// @Tags Stickers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pack ID"
// @Success 200 {object} stickerdomain.StickerPack
// @Failure 404 {object} response.ErrorResponse "Not found"
// @Router /stickers/packs/{id} [get]
func (h *StickerHandler) GetPack(c *gin.Context) {
	id := c.Param("id")
	pack, err := h.stickerService.GetPackByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.JSON(c, 200, pack)
}

// AddSticker adds a sticker to an existing pack
// @Tags Stickers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Pack ID"
// @Param request body stickerdomain.AddStickerRequest true "Sticker details"
// @Success 201 {object} stickerdomain.Sticker
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/packs/{id}/stickers [post]
func (h *StickerHandler) AddSticker(c *gin.Context) {
	userID, _ := c.Get("userID")
	packID := c.Param("id")
	var req stickerdomain.AddStickerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	sticker, err := h.stickerService.AddSticker(packID, userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 201, sticker)
}

// DeletePack deletes a sticker pack (owner only)
// @Tags Stickers
// @Security BearerAuth
// @Produce json
// @Param id path string true "Pack ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/packs/{id} [delete]
func (h *StickerHandler) DeletePack(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	if err := h.stickerService.DeletePack(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "pack deleted"})
}

// AddToLibrary adds a sticker to the user's personal library
// @Tags Stickers
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param request body object{stickerId=string} true "Sticker ID to add"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/library [post]
func (h *StickerHandler) AddToLibrary(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req struct {
		StickerID string `json:"stickerId" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.stickerService.AddToLibrary(userID.(string), req.StickerID); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "sticker added to library"})
}

// GetLibrary returns the user's personal sticker library
// @Tags Stickers
// @Security BearerAuth
// @Produce json
// @Success 200 {array} stickerdomain.Sticker
// @Failure 400 {object} response.ErrorResponse
// @Router /stickers/library [get]
func (h *StickerHandler) GetLibrary(c *gin.Context) {
	userID, _ := c.Get("userID")
	stickers, err := h.stickerService.GetLibrary(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, stickers)
}
