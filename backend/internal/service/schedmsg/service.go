package schedmsgservice

import (
	"errors"
	"log"
	"time"

	draftdomain "ChatServerGolang/backend/internal/domain/draft"
	messagedomain "ChatServerGolang/backend/internal/domain/message"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type scheduledMessageService struct {
	schedRepo  repository.ScheduledMessageRepository
	messageRepo repository.MessageRepository
	chatRepo   repository.ChatRepository
}

func NewScheduledMessageService(schedRepo repository.ScheduledMessageRepository, messageRepo repository.MessageRepository, chatRepo repository.ChatRepository) service.ScheduledMessageService {
	return &scheduledMessageService{schedRepo: schedRepo, messageRepo: messageRepo, chatRepo: chatRepo}
}

func (s *scheduledMessageService) Schedule(userID string, req *draftdomain.ScheduleMessageRequest) (*draftdomain.ScheduledMessage, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(req.ChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	t, err := time.Parse(time.RFC3339, req.ScheduledAt)
	if err != nil {
		return nil, errors.New("invalid scheduled_at format, use RFC3339")
	}
	if t.Before(time.Now()) {
		return nil, errors.New("scheduled time must be in the future")
	}

	msg := &draftdomain.ScheduledMessage{
		ID:          uuid.New().String(),
		ChatID:      req.ChatID,
		SenderID:    userID,
		Content:     req.Content,
		Type:        req.Type,
		ReplyToID:   req.ReplyToID,
		ScheduledAt: req.ScheduledAt,
		CreatedAt:   time.Now(),
	}
	if err := s.schedRepo.Create(msg); err != nil {
		return nil, err
	}
	return msg, nil
}

func (s *scheduledMessageService) GetScheduled(userID string) ([]*draftdomain.ScheduledMessage, error) {
	return s.schedRepo.FindByUserID(userID)
}

func (s *scheduledMessageService) CancelScheduled(id, userID string) error {
	return s.schedRepo.Delete(id)
}

// SchedulerProcess runs as a goroutine to send pending scheduled messages
func (s *scheduledMessageService) SchedulerProcess() {
	messages, err := s.schedRepo.FindPending()
	if err != nil {
		log.Printf("Scheduler error: %v", err)
		return
	}

	for _, sm := range messages {
		isParticipant, _ := s.chatRepo.IsParticipant(sm.ChatID, sm.SenderID)
		if !isParticipant {
			s.schedRepo.MarkAsSent(sm.ID)
			continue
		}

		now := time.Now()
		msg := &messagedomain.Message{
			ID:        uuid.New().String(),
			ChatID:    sm.ChatID,
			SenderID:  sm.SenderID,
			Content:   sm.Content,
			Type:      sm.Type,
			ReplyToID: sm.ReplyToID,
			CreatedAt: now,
			UpdatedAt: now,
		}

		if err := s.messageRepo.Create(msg); err != nil {
			log.Printf("Scheduler failed to send message %s: %v", sm.ID, err)
			continue
		}

		s.schedRepo.MarkAsSent(sm.ID)
	}
}


