package savedmsghandler

import (
	"strconv"

	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type SavedMessageHandler struct {
	savedMsgService service.SavedMessageService
}

func NewSavedMessageHandler(savedMsgService service.SavedMessageService) *SavedMessageHandler {
	return &SavedMessageHandler{savedMsgService: savedMsgService}
}

// SaveMessage сохраняет указанное сообщение в личную коллекцию "Избранное" пользователя
// @Tags Сохранённые сообщения
// @Summary Сохранить сообщение в избранное по его ID
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Добавляет сообщение в избранное пользователя по ID сообщения и ID чата. Позволяет быстро находить важные сообщения позже через GET /saved-messages. Одно и то же сообщение нельзя сохранить дважды.
// @Param id path string true "ID сообщения, которое нужно сохранить"
// @Param chatId query string true "ID чата, в котором находится сообщение"
// @Success 200 {object} response.MessageResponse "Сообщение успешно сохранено, возвращает объект SavedMessageResponse с данными сообщения и чата"
// @Failure 400 {object} response.ErrorResponse "Ошибка: сообщение уже сохранено или неверный запрос"
// @Router /messages/{id}/save [post]
func (h *SavedMessageHandler) SaveMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	messageID := c.Param("id")
	chatID := c.Query("chatId")
	if chatID == "" {
		response.BadRequest(c, "chatId обязателен")
		return
	}

	result, err := h.savedMsgService.SaveMessage(userID.(string), messageID, chatID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, result)
}

// GetSavedMessages возвращает постраничный список всех сохранённых сообщений текущего пользователя
// @Tags Сохранённые сообщения
// @Summary Получить список сохранённых сообщений с пагинацией
// @Security BearerAuth
// @Produce json
// @Description Возвращает список сохранённых сообщений с пагинацией. Каждый элемент содержит полную информацию о сообщении (текст, тип, отправитель) и о чате, в котором оно находится. Упорядочено по дате сохранения (новые сверху).
// @Param limit query int false "Количество записей на странице (по умолчанию 50, максимум 100)"
// @Param offset query int false "Смещение от начала списка для пагинации"
// @Success 200 {object} response.APIResponse "Массив SavedMessageResponse с meta (total, limit, offset)"
// @Failure 400 {object} response.ErrorResponse "Ошибка при получении списка"
// @Router /saved-messages [get]
func (h *SavedMessageHandler) GetSavedMessages(c *gin.Context) {
	userID, _ := c.Get("userID")
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset, _ := strconv.Atoi(c.DefaultQuery("offset", "0"))

	messages, total, err := h.savedMsgService.GetSavedMessages(userID.(string), limit, offset)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.Paginated(c, 200, messages, total, limit, offset)
}

// DeleteSavedMessage удаляет сообщение из коллекции "Избранное" по его ID
// @Tags Сохранённые сообщения
// @Summary Удалить сообщение из избранного по ID записи
// @Security BearerAuth
// @Produce json
// @Description Удаляет ранее сохранённое сообщение из избранного. Только владелец может удалить своё сохранённое сообщение. ID можно получить из GET /saved-messages.
// @Param id path string true "ID записи сохранённого сообщения (не ID самого сообщения, а ID записи в избранном)"
// @Success 200 {object} response.MessageResponse "Сообщение успешно удалено из избранного"
// @Failure 400 {object} response.ErrorResponse "Ошибка: запись не найдена или нет прав"
// @Router /saved-messages/{id} [delete]
func (h *SavedMessageHandler) DeleteSavedMessage(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	if err := h.savedMsgService.DeleteSavedMessage(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, map[string]interface{}{"message": "сообщение удалено из избранного"})
}
