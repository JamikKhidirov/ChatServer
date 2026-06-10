package service

import (
	"encoding/json"
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
)

type pollService struct {
	pollRepo      repository.PollRepository
	chatRepo      repository.ChatRepository
	messageRepo   repository.MessageRepository
	sysMsgService SystemMessageService
}

func NewPollService(pollRepo repository.PollRepository, chatRepo repository.ChatRepository, messageRepo repository.MessageRepository, sysMsgService SystemMessageService) PollService {
	return &pollService{
		pollRepo:      pollRepo,
		chatRepo:      chatRepo,
		messageRepo:   messageRepo,
		sysMsgService: sysMsgService,
	}
}

func (s *pollService) CreatePoll(userID string, req *domain.CreatePollRequest) (*domain.PollWithResults, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(req.ChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	optionsJSON, _ := json.Marshal(req.Options)
	now := time.Now()

	var expiresAt *string
	if req.ExpiresInMins != nil && *req.ExpiresInMins > 0 {
		t := now.Add(time.Duration(*req.ExpiresInMins) * time.Minute).Format(time.RFC3339)
		expiresAt = &t
	}

	poll := &domain.Poll{
		ID:             uuid.New().String(),
		ChatID:         req.ChatID,
		CreatorID:      userID,
		Question:       req.Question,
		Options:        string(optionsJSON),
		IsAnonymous:    req.IsAnonymous,
		MultipleChoice: req.MultipleChoice,
		ExpiresAt:      expiresAt,
		CreatedAt:      now,
	}

	if err := s.pollRepo.Create(poll); err != nil {
		return nil, err
	}

	return s.buildPollResponse(poll, userID)
}

func (s *pollService) GetPollsByChatID(chatID, userID string) ([]*domain.PollWithResults, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	polls, err := s.pollRepo.FindByChatID(chatID)
	if err != nil {
		return nil, err
	}

	results := make([]*domain.PollWithResults, 0, len(polls))
	for _, p := range polls {
		r, err := s.buildPollResponse(p, userID)
		if err != nil {
			continue
		}
		results = append(results, r)
	}
	return results, nil
}

func (s *pollService) Vote(pollID, userID string, optionIndex int) error {
	poll, err := s.pollRepo.FindByID(pollID)
	if err != nil {
		return errors.New("poll not found")
	}

	if poll.Closed {
		return errors.New("poll is closed")
	}

	if poll.ExpiresAt != nil {
		t, err := time.Parse(time.RFC3339, *poll.ExpiresAt)
		if err == nil && time.Now().After(t) {
			return errors.New("poll has expired")
		}
	}

	hasVoted, _ := s.pollRepo.HasVoted(pollID, userID)
	if hasVoted && !poll.MultipleChoice {
		return errors.New("already voted")
	}

	var options []string
	json.Unmarshal([]byte(poll.Options), &options)
	if optionIndex < 0 || optionIndex >= len(options) {
		return errors.New("invalid option")
	}

	vote := &domain.PollVote{
		PollID:      pollID,
		UserID:      userID,
		OptionIndex: optionIndex,
		VotedAt:     time.Now(),
	}
	return s.pollRepo.AddVote(vote)
}

func (s *pollService) ClosePoll(pollID, userID string) error {
	poll, err := s.pollRepo.FindByID(pollID)
	if err != nil {
		return errors.New("poll not found")
	}

	if poll.CreatorID != userID {
		return errors.New("only the creator can close the poll")
	}

	poll.Closed = true
	return s.pollRepo.Update(poll)
}

func (s *pollService) buildPollResponse(poll *domain.Poll, userID string) (*domain.PollWithResults, error) {
	var options []string
	json.Unmarshal([]byte(poll.Options), &options)

	totalVotes, _ := s.pollRepo.GetTotalVotes(poll.ID)

	optionsList := make([]domain.PollOption, len(options))
	for i, opt := range options {
		votes, _ := s.pollRepo.GetVoteCount(poll.ID, i)
		optionsList[i] = domain.PollOption{
			Text:  opt,
			Votes: votes,
		}
	}

	var votedOption *int
	userVote, err := s.pollRepo.GetUserVote(poll.ID, userID)
	if err == nil && userVote != nil {
		votedOption = &userVote.OptionIndex
	}

	expiresAt := poll.ExpiresAt
	if expiresAt != nil {
		t, err := time.Parse(time.RFC3339, *expiresAt)
		if err == nil && time.Now().After(t) {
			poll.Closed = true
		}
	}

	return &domain.PollWithResults{
		Poll:        *poll,
		OptionsList: optionsList,
		TotalVotes:  totalVotes,
		VotedOption: votedOption,
	}, nil
}
