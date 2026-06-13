package stickerhandler

import (
	"ChatServerGolang/backend/internal/domain/sticker"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type StickerHandler struct {
	stickerService service.StickerService
}

func NewStickerHandler(stickerService service.StickerService) *StickerHandler {
	return &StickerHandler{stickerService: stickerService}
}

// CreatePack creates a new sticker pack
// @Tags Стикеры
// @Summary Создать набор стикеров
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Создаёт новый набор стикеров с указанным названием и списком стикеров.
// @Param request body stickerdomain.CreateStickerPackRequest true "Данные набора: name (название, обязательно), stickers (массив стикеров, обязательно)"
// @Success 201 {object} stickerdomain.StickerPack "Набор стикеров создан"
// @Failure 400 {object} response.ErrorResponse "Ошибка создания набора"
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
// @Tags Стикеры
// @Summary Получить публичные наборы стикеров
// @Security BearerAuth
// @Produce json
// @Description Возвращает список всех общедоступных наборов стикеров на платформе.
// @Success 200 {array} stickerdomain.StickerPack "Список публичных наборов стикеров"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения наборов"
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
// @Tags Стикеры
// @Summary Получить мои наборы стикеров
// @Security BearerAuth
// @Produce json
// @Description Возвращает наборы стикеров, созданные аутентифицированным пользователем.
// @Success 200 {array} stickerdomain.StickerPack "Список ваших наборов стикеров"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения наборов"
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
// @Tags Стикеры
// @Summary Получить набор стикеров по ID
// @Security BearerAuth
// @Produce json
// @Description Возвращает информацию о наборе стикеров и его содержимом по идентификатору.
// @Param id path string true "ID набора"
// @Success 200 {object} stickerdomain.StickerPack "Информация о наборе стикеров"
// @Failure 404 {object} response.ErrorResponse "Набор не найден"
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
// @Tags Стикеры
// @Summary Добавить стикер в набор
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Добавляет новый стикер в существующий набор стикеров. Доступно только владельцу набора.
// @Param id path string true "ID набора"
// @Param request body stickerdomain.AddStickerRequest true "Данные стикера: image_url (URL изображения, обязательно), emoji (связанный эмодзи, опционально)"
// @Success 201 {object} stickerdomain.Sticker "Стикер добавлен"
// @Failure 400 {object} response.ErrorResponse "Ошибка добавления стикера"
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
// @Tags Стикеры
// @Summary Удалить набор стикеров
// @Security BearerAuth
// @Produce json
// @Description Удаляет набор стикеров и все его содержимое. Доступно только владельцу набора.
// @Param id path string true "ID набора"
// @Success 200 {object} response.MessageResponse "Набор стикеров удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка удаления набора"
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
// @Tags Стикеры
// @Summary Добавить стикер в библиотеку
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Добавляет стикер в личную библиотеку пользователя для быстрого доступа.
// @Param request body object{stickerId=string} true "Параметры: stickerId (ID стикера, обязательно)"
// @Success 200 {object} response.MessageResponse "Стикер добавлен в библиотеку"
// @Failure 400 {object} response.ErrorResponse "Ошибка добавления"
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
// @Tags Стикеры
// @Summary Получить библиотеку стикеров
// @Security BearerAuth
// @Produce json
// @Description Возвращает список стикеров, добавленных пользователем в личную библиотеку.
// @Success 200 {array} stickerdomain.Sticker "Библиотека стикеров"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения библиотеки"
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
