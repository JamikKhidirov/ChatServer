package handler

import (
	"ChatServerGolang/internal/domain"
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

func (h *PollHandler) CreatePoll(c *gin.Context) {
	userID, _ := c.Get("userID")
	var req domain.CreatePollRequest
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

func (h *PollHandler) Vote(c *gin.Context) {
	userID, _ := c.Get("userID")
	pollID := c.Param("pollId")
	var req domain.VotePollRequest
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

func (h *PollHandler) ClosePoll(c *gin.Context) {
	userID, _ := c.Get("userID")
	pollID := c.Param("pollId")
	if err := h.pollService.ClosePoll(pollID, userID.(string)); err != nil {
		response.BadRequest(c, err.Error())
		return
	}
	response.JSON(c, 200, gin.H{"message": "poll closed"})
}
