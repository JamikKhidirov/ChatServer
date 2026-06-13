package voicechathandler

import (
	voicechatdomain "ChatServerGolang/backend/internal/domain/voicechat"
	"ChatServerGolang/backend/internal/service"
	"ChatServerGolang/backend/pkg/response"

	"github.com/gin-gonic/gin"
)

type VoiceChatHandler struct {
	voiceChatService service.VoiceChatService
}

func NewVoiceChatHandler(voiceChatService service.VoiceChatService) *VoiceChatHandler {
	return &VoiceChatHandler{voiceChatService: voiceChatService}
}

// CreateVoiceChat создаёт новую голосовую комнату в указанной группе
// @Tags Голосовые чаты
// @Summary Создать голосовую комнату в группе
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Создаёт постоянную голосовую комнату в групповом чате. В отличие от группового звонка, голосовая комната может быть активна длительное время и не требует принятия вызова. Можно указать название комнаты и время отложенного старта. Создатель автоматически присоединяется к комнате.
// @Param id path string true "ID группы (chatId), в которой создаётся голосовая комната. Группа должна существовать и иметь тип 'group'."
// @Param request body voicechatdomain.CreateVoiceChatRequest true "Параметры создания: title (название комнаты, опционально), scheduledInMins (отложенный старт через N минут, опционально)"
// @Success 201 {object} voicechatdomain.VoiceChatResponse "Голосовая комната создана, возвращает ID, статус и количество участников"
// @Failure 400 {object} response.ErrorResponse "Ошибка: чат не найден, это не группа или неверные параметры"
// @Router /chats/{id}/voice-chat [post]
func (h *VoiceChatHandler) CreateVoiceChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")

	var req voicechatdomain.CreateVoiceChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	result, err := h.voiceChatService.CreateVoiceChat(chatID, userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 201, result)
}

// GetVoiceChat возвращает детальную информацию о голосовой комнате
// @Tags Голосовые чаты
// @Summary Получить информацию о голосовой комнате по ID
// @Security BearerAuth
// @Produce json
// @Description Возвращает полную информацию о голосовой комнате: статус (active/scheduled/ended), количество участников, время начала и окончания.
// @Param id path string true "ID голосовой комнаты"
// @Success 200 {object} voicechatdomain.VoiceChatResponse "Информация о голосовой комнате"
// @Failure 400 {object} response.ErrorResponse "Ошибка: голосовая комната не найдена"
// @Router /voice-chats/{id} [get]
func (h *VoiceChatHandler) GetVoiceChat(c *gin.Context) {
	id := c.Param("id")

	result, err := h.voiceChatService.GetVoiceChat(id)
	if err != nil {
		response.NotFound(c, err.Error())
		return
	}

	response.JSON(c, 200, result)
}

// GetActiveVoiceChats возвращает список активных голосовых комнат в указанном чате
// @Tags Голосовые чаты
// @Summary Получить активные голосовые комнаты в чате
// @Security BearerAuth
// @Produce json
// @Description Показывает все активные (не завершённые) голосовые комнаты в группе. Если комната завершена, она не отображается в этом списке — используйте GET /chats/{id}/voice-chats/history.
// @Param id path string true "ID группы (chatId)"
// @Success 200 {array} voicechatdomain.VoiceChatResponse "Список активных голосовых комнат"
// @Failure 400 {object} response.ErrorResponse "Ошибка при получении списка"
// @Router /chats/{id}/voice-chats/active [get]
func (h *VoiceChatHandler) GetActiveVoiceChats(c *gin.Context) {
	chatID := c.Param("id")

	results, err := h.voiceChatService.GetActiveVoiceChats(chatID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, results)
}

// GetVoiceChatHistory возвращает историю всех голосовых комнат в чате (включая завершённые)
// @Tags Голосовые чаты
// @Summary Получить историю всех голосовых комнат в чате
// @Security BearerAuth
// @Produce json
// @Description Возвращает все голосовые комнаты, созданные в группе, включая активные, запланированные и завершённые. Упорядочено по дате создания (новые сверху).
// @Param id path string true "ID группы (chatId)"
// @Success 200 {array} voicechatdomain.VoiceChatResponse "Полная история голосовых комнат"
// @Failure 400 {object} response.ErrorResponse "Ошибка при получении истории"
// @Router /chats/{id}/voice-chats/history [get]
func (h *VoiceChatHandler) GetVoiceChatHistory(c *gin.Context) {
	chatID := c.Param("id")

	results, err := h.voiceChatService.GetVoiceChatHistory(chatID)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, results)
}

// JoinVoiceChat присоединяет текущего пользователя к активной голосовой комнате
// @Tags Голосовые чаты
// @Summary Присоединиться к голосовой комнате
// @Security BearerAuth
// @Produce json
// @Description Позволяет пользователю присоединиться к голосовой комнате. Пользователь становится участником и будет отображаться в списке участников комнаты. Не требует подтверждения от создателя.
// @Param id path string true "ID голосовой комнаты"
// @Success 200 {object} response.MessageResponse "Пользователь успешно присоединился к голосовой комнате"
// @Failure 400 {object} response.ErrorResponse "Ошибка: комната не найдена или уже завершена"
// @Router /voice-chats/{id}/join [post]
func (h *VoiceChatHandler) JoinVoiceChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	if err := h.voiceChatService.JoinVoiceChat(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, map[string]interface{}{"message": "вы присоединились к голосовому чату"})
}

// LeaveVoiceChat покидает голосовую комнату
// @Tags Голосовые чаты
// @Summary Покинуть голосовую комнату
// @Security BearerAuth
// @Produce json
// @Description Отсоединяет пользователя от голосовой комнаты. Пользователь удаляется из списка участников, и счётчик участников обновляется.
// @Param id path string true "ID голосовой комнаты"
// @Success 200 {object} response.MessageResponse "Пользователь покинул голосовую комнату"
// @Failure 400 {object} response.ErrorResponse "Ошибка: пользователь не является участником комнаты"
// @Router /voice-chats/{id}/leave [post]
func (h *VoiceChatHandler) LeaveVoiceChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	if err := h.voiceChatService.LeaveVoiceChat(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, map[string]interface{}{"message": "вы покинули голосовой чат"})
}

// EndVoiceChat завершает голосовую комнату (только для создателя или администратора)
// @Tags Голосовые чаты
// @Summary Завершить голосовую комнату (создатель/админ)
// @Security BearerAuth
// @Produce json
// @Description Завершает голосовую комнату. После завершения новые участники не могут присоединиться. Статус меняется на "ended". Все текущие участники отключаются.
// @Param id path string true "ID голосовой комнаты"
// @Success 200 {object} response.MessageResponse "Голосовая комната завершена"
// @Failure 400 {object} response.ErrorResponse "Ошибка: нет прав на завершение"
// @Router /voice-chats/{id}/end [post]
func (h *VoiceChatHandler) EndVoiceChat(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	if err := h.voiceChatService.EndVoiceChat(id, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	response.JSON(c, 200, map[string]interface{}{"message": "голосовой чат завершён"})
}

// MuteParticipant включает или выключает микрофон участника в голосовой комнате
// @Tags Голосовые чаты
// @Summary Включить/выключить микрофон участника
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Позволяет участнику включить или выключить свой микрофон в голосовой комнате. Администраторы могут отключить микрофон другим участникам.
// @Param id path string true "ID голосовой комнаты"
// @Param request body voicechatdomain.MuteParticipantRequest true "Параметры: muted (boolean) — true выключает микрофон, false включает"
// @Success 200 {object} response.MessageResponse "Статус микрофона изменён (muted/unmuted)"
// @Failure 400 {object} response.ErrorResponse "Ошибка при изменении статуса"
// @Router /voice-chats/{id}/mute [post]
func (h *VoiceChatHandler) MuteParticipant(c *gin.Context) {
	userID, _ := c.Get("userID")
	id := c.Param("id")

	var req struct {
		Muted bool `json:"muted"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		response.ValidationError(c, err.Error())
		return
	}

	if err := h.voiceChatService.MuteParticipant(id, userID.(string), req.Muted); err != nil {
		response.BadRequest(c, err.Error())
		return
	}

	action := "микрофон включён"
	if req.Muted {
		action = "микрофон выключен"
	}
	response.JSON(c, 200, map[string]interface{}{"message": action})
}
