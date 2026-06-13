package linkservice

import (
	"errors"
	"time"

	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type inviteLinkService struct {
	linkRepo repository.InviteLinkRepository
	chatRepo repository.ChatRepository
}

func NewInviteLinkService(linkRepo repository.InviteLinkRepository, chatRepo repository.ChatRepository) service.InviteLinkService {
	return &inviteLinkService{linkRepo: linkRepo, chatRepo: chatRepo}
}

func (s *inviteLinkService) CreateInviteLink(chatID, userID string, req *chatdomain.CreateInviteLinkRequest) (*chatdomain.InviteLink, error) {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}
	if chat.CreatedBy != userID {
		return nil, errors.New("only chat creator can create invite links")
	}

	var expiresAt *time.Time
	if req.ExpiresInMins > 0 {
		t := time.Now().Add(time.Duration(req.ExpiresInMins) * time.Minute)
		expiresAt = &t
	}

	link := &chatdomain.InviteLink{
		ID:         uuid.New().String(),
		ChatID:     chatID,
		CreatorID:  userID,
		Code:       uuid.New().String()[:12],
		ExpiresAt:  expiresAt,
		UsageLimit: req.UsageLimit,
		Active:     true,
		CreatedAt:  time.Now(),
	}

	if err := s.linkRepo.Create(link); err != nil {
		return nil, err
	}
	return link, nil
}

func (s *inviteLinkService) GetInviteLinks(chatID, userID string) ([]*chatdomain.InviteLink, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}
	return s.linkRepo.FindByChatID(chatID)
}

func (s *inviteLinkService) DeleteInviteLink(linkID, userID string) error {
	return s.linkRepo.Delete(linkID)
}

func (s *inviteLinkService) JoinByInviteLink(code, userID string) error {
	link, err := s.linkRepo.FindByCode(code)
	if err != nil {
		return errors.New("invalid invite link")
	}
	if !link.Active {
		return errors.New("invite link is deactivated")
	}
	if link.ExpiresAt != nil && time.Now().After(*link.ExpiresAt) {
		return errors.New("invite link has expired")
	}
	if link.UsageLimit > 0 && link.UsageCount >= link.UsageLimit {
		return errors.New("invite link usage limit reached")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(link.ChatID, userID)
	if isParticipant {
		return errors.New("you are already a member of this chat")
	}

	if err := s.chatRepo.AddParticipant(link.ChatID, userID, "member"); err != nil {
		return err
	}

	return s.linkRepo.IncrementUsage(link.ID)
}
