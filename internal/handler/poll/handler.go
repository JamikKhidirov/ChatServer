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
// @Tags Polls
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param id path string true "Chat ID"
// @Param request body polldomain.CreatePollRequest true "Poll question and options"
// @Success 201 {object} polldomain.PollWithResults
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Polls
// @Security BearerAuth
// @Produce json
// @Param id path string true "Chat ID"
// @Success 200 {array} polldomain.PollWithResults
// @Failure 400 {object} response.ErrorResponse
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
// @Tags Polls
// @Security BearerAuth
// @Accept json
// @Produce json
// @Param pollId path string true "Poll ID"
// @Param request body polldomain.VotePollRequest true "Selected option index"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Already voted or invalid option"
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
// @Tags Polls
// @Security BearerAuth
// @Produce json
// @Param pollId path string true "Poll ID"
// @Success 200 {object} response.MessageResponse
// @Failure 400 {object} response.ErrorResponse "Not the creator"
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
