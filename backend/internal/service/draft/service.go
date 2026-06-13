package draftservice

import (
	"errors"
	"time"

	"ChatServerGolang/backend/internal/domain/draft"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type draftService struct {
	draftRepo repository.DraftRepository
}

func NewDraftService(draftRepo repository.DraftRepository) service.DraftService {
	return &draftService{draftRepo: draftRepo}
}

func (s *draftService) SaveDraft(userID string, req *draftdomain.SaveDraftRequest) (*draftdomain.Draft, error) {
	existing, _ := s.draftRepo.FindByUserAndChat(userID, req.ChatID)
	now := time.Now()

	if existing != nil {
		existing.Content = req.Content
		existing.ReplyToID = req.ReplyToID
		existing.UpdatedAt = now
		if err := s.draftRepo.Save(existing); err != nil {
			return nil, err
		}
		return existing, nil
	}

	draft := &draftdomain.Draft{
		ID:        uuid.New().String(),
		UserID:    userID,
		ChatID:    req.ChatID,
		Content:   req.Content,
		ReplyToID: req.ReplyToID,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.draftRepo.Save(draft); err != nil {
		return nil, err
	}
	return draft, nil
}

func (s *draftService) GetDraft(userID, chatID string) (*draftdomain.Draft, error) {
	draft, err := s.draftRepo.FindByUserAndChat(userID, chatID)
	if err != nil {
		return nil, errors.New("draft not found")
	}
	return draft, nil
}

func (s *draftService) DeleteDraft(userID, draftID string) error {
	return s.draftRepo.Delete(draftID)
}


