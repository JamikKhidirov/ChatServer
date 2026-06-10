package handler

import (
	"ChatServerGolang/internal/domain"
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

func (h *StickerHandler) CreatePack(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.CreateStickerPackRequest
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

func (h *StickerHandler) ListPacks(c *gin.Context) {
	packs, err := h.stickerService.GetPacks()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, packs)
}

func (h *StickerHandler) GetMyPacks(c *gin.Context) {
	userID, _ := c.Get("userID")
	packs, err := h.stickerService.GetMyPacks(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, packs)
}

func (h *StickerHandler) GetPack(c *gin.Context) {
	id := c.Param("id")
	pack, err := h.stickerService.GetPackByID(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}
	response.JSON(c, 200, pack)
}

func (h *StickerHandler) AddSticker(c *gin.Context) {
	userID, _ := c.Get("userID")
	packID := c.Param("id")
	var req domain.AddStickerRequest
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

func (h *StickerHandler) DeletePack(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")
	if err := h.stickerService.DeletePack(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "pack deleted"})
}

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

func (h *StickerHandler) GetLibrary(c *gin.Context) {
	userID, _ := c.Get("userID")
	stickers, err := h.stickerService.GetLibrary(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, stickers)
}
