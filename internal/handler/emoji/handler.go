package emojihandler

import (
	"io"
	"os"
	"path/filepath"

	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type EmojiHandler struct {
	emojiService service.CustomEmojiService
}

func NewEmojiHandler(emojiService service.CustomEmojiService) *EmojiHandler {
	return &EmojiHandler{emojiService: emojiService}
}

// CreateEmoji загружает новый кастомный эмодзи на сервер
// @Tags Кастомные эмодзи
// @Security BearerAuth
// @Accept multipart/form-data
// @Produce json
// @Description Позволяет загрузить собственное изображение эмодзи и задать ему короткое имя (shortcode). После загрузки эмодзи можно использовать в чатах. Поддерживаются форматы PNG, JPG, GIF (включая анимированные). Файл сохраняется в директории uploads/emojis/.
// @Param shortcode formData string true "Короткое имя для эмодзи (например 'party', 'thumbsup', 'cat'). Будет использоваться в клиенте для вызова эмодзи."
// @Param emoji formData file true "Файл изображения эмодзи (PNG, JPG, GIF). Рекомендуемый размер: 128x128 пикселей."
// @Success 201 {object} emojidomain.CustomEmojiResponse "Эмодзи успешно создан, возвращает ID, shortcode и URL файла"
// @Failure 400 {object} response.ErrorResponse "Ошибка: файл не загружен, shortcode пустой или формат не поддерживается"
// @Router /emojis [post]
func (h *EmojiHandler) CreateEmoji(c *gin.Context) {
	userID, _ := c.Get("userID")
	shortcode := c.PostForm("shortcode")
	if shortcode == "" {
		response.BadRequest(c, "shortcode обязателен")
		return
	}

	file, header, err := c.Request.FormFile("emoji")
	if err != nil {
		response.BadRequest(c, "файл эмодзи обязателен")
		return
	}
	defer file.Close()

	uploadDir := "uploads/emojis"
	os.MkdirAll(uploadDir, 0755)

	ext := filepath.Ext(header.Filename)
	fileName := uuid.New().String() + ext
	filePath := filepath.Join(uploadDir, fileName)

	out, err := os.Create(filePath)
	if err != nil {
		response.InternalError(c, "не удалось сохранить эмодзи")
		return
	}
	defer out.Close()

	if _, err := io.Copy(out, file); err != nil {
		os.Remove(filePath)
		response.InternalError(c, "не удалось сохранить эмодзи")
		return
	}

	fileURL := "/uploads/emojis/" + fileName

	result, err := h.emojiService.CreateEmoji(userID.(string), shortcode, filePath, fileURL)
	if err != nil {
		os.Remove(filePath)
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, result)
}

// GetMyEmojis возвращает список кастомных эмодзи, загруженных текущим пользователем
// @Tags Кастомные эмодзи
// @Security BearerAuth
// @Produce json
// @Description Возвращает только те эмодзи, которые были загружены текущим пользователем. Для просмотра всех публичных эмодзи используйте GET /emojis.
// @Success 200 {array} emojidomain.CustomEmojiResponse "Массив кастомных эмодзи пользователя с ID, shortcode и URL"
// @Failure 400 {object} response.ErrorResponse
// @Router /emojis/my [get]
func (h *EmojiHandler) GetMyEmojis(c *gin.Context) {
	userID, _ := c.Get("userID")
	emojis, err := h.emojiService.GetMyEmojis(userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, emojis)
}

// GetAllEmojis возвращает список всех кастомных эмодзи на сервере
// @Tags Кастомные эмодзи
// @Security BearerAuth
// @Produce json
// @Description Возвращает все кастомные эмодзи, загруженные всеми пользователями. Подходит для общей галереи эмодзи в клиенте.
// @Success 200 {array} emojidomain.CustomEmojiResponse "Массив всех кастомных эмодзи на сервере"
// @Failure 400 {object} response.ErrorResponse
// @Router /emojis [get]
func (h *EmojiHandler) GetAllEmojis(c *gin.Context) {
	emojis, err := h.emojiService.GetAllEmojis()
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, emojis)
}

// DeleteEmoji удаляет кастомный эмодзи по его ID
// @Tags Кастомные эмодзи
// @Security BearerAuth
// @Produce json
// @Description Удаляет эмодзи из базы данных и файл изображения с диска. Только владелец эмодзи может его удалить. После удаления эмодзи перестанет отображаться в галерее.
// @Param id path string true "ID эмодзи для удаления"
// @Success 200 {object} response.MessageResponse "Эмодзи успешно удалён"
// @Failure 400 {object} response.ErrorResponse "Ошибка: эмодзи не найден или нет прав на удаление"
// @Router /emojis/{id} [delete]
func (h *EmojiHandler) DeleteEmoji(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	if err := h.emojiService.DeleteEmoji(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, map[string]interface{}{"message": "эмодзи удалён"})
}
