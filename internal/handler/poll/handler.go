package pollhandler

import (
	"ChatServerGolang/internal/domain/poll"
	"ChatServerGolang/internal/service"
	"ChatServerGolang/pkg/response"

	"github.com/gin-gonic/gin"
)

type PollHandler struct {
	pollService service.PollService
}

func NewPollHandler(pollService service.PollService) *PollHandler {
	return &PollHandler{pollService: pollService}
}

// CreatePoll creates a poll in a chat
// @Tags Голосований
// @Summary Создать голосование
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Создаёт новое голосование в указанном чате. Можно задать вопрос и несколько вариантов ответа.
// @Param id path string true "ID чата"
// @Param request body polldomain.CreatePollRequest true "Параметры голосования: question (вопрос, обязательно), options (варианты ответа, обязательно), is_anonymous (анонимное, опционально)"
// @Success 201 {object} polldomain.PollWithResults "Голосование создано"
// @Failure 400 {object} response.ErrorResponse "Ошибка создания голосования"
// @Router /chats/{id}/polls [post]
func (h *PollHandler) CreatePoll(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req polldomain.CreatePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	poll, err := h.pollService.CreatePoll(userID.(string), &req)
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 201, poll)
}

// GetPolls returns all polls in a chat
// @Tags Голосований
// @Summary Получить список голосований
// @Security BearerAuth
// @Produce json
// @Description Возвращает все голосования, созданные в указанном чате.
// @Param id path string true "ID чата"
// @Success 200 {array} polldomain.PollWithResults "Список голосований"
// @Failure 400 {object} response.ErrorResponse "Ошибка получения голосований"
// @Router /chats/{id}/polls [get]
func (h *PollHandler) GetPolls(c *gin.Context) {
	userID, _ := c.Get("userID")
	chatID := c.Param("id")
	polls, err := h.pollService.GetPollsByChatID(chatID, userID.(string))
	if err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, polls)
}

// Vote casts a vote in a poll
// @Tags Голосований
// @Summary Проголосовать
// @Security BearerAuth
// @Accept json
// @Produce json
// @Description Отдаёт голос за один из вариантов ответа в указанном голосовании.
// @Param pollId path string true "ID голосования"
// @Param request body polldomain.VotePollRequest true "Выбранный вариант: option_index (индекс варианта, обязательно)"
// @Success 200 {object} response.MessageResponse "Голос учтён"
// @Failure 400 {object} response.ErrorResponse "Уже голосовали или неверный вариант"
// @Router /polls/{pollId}/vote [post]
func (h *PollHandler) Vote(c *gin.Context) {
	userID, _ := c.Get("userID")
	pollID := c.Param("pollId")
	var req polldomain.VotePollRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	if err := h.pollService.Vote(pollID, userID.(string), req.OptionIndex); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "vote recorded"})
}

// ClosePoll closes a poll (creator only)
// @Tags Голосований
// @Summary Закрыть голосование
// @Security BearerAuth
// @Produce json
// @Description Закрывает голосование для дальнейшего голосования. Доступно только создателю голосования.
// @Param pollId path string true "ID голосования"
// @Success 200 {object} response.MessageResponse "Голосование закрыто"
// @Failure 400 {object} response.ErrorResponse "Недостаточно прав для закрытия"
// @Router /polls/{pollId}/close [post]
func (h *PollHandler) ClosePoll(c *gin.Context) {
	userID, _ := c.Get("userID")
	pollID := c.Param("pollId")
	if err := h.pollService.ClosePoll(pollID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "poll closed"})
}
