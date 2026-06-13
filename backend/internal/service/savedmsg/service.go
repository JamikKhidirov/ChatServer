package savedmsgservice

import (
	"errors"
	"time"

	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	messagedomain "ChatServerGolang/backend/internal/domain/message"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"
	userdomain "ChatServerGolang/backend/internal/domain/user"

	"github.com/google/uuid"
)

type savedMessageService struct {
	savedMsgRepo repository.SavedMessageRepository
	msgRepo      repository.MessageRepository
	chatRepo     repository.ChatRepository
	userRepo     repository.UserRepository
}

func NewSavedMessageService(
	savedMsgRepo repository.SavedMessageRepository,
	msgRepo repository.MessageRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
) service.SavedMessageService {
	return &savedMessageService{
		savedMsgRepo: savedMsgRepo,
		msgRepo:      msgRepo,
		chatRepo:     chatRepo,
		userRepo:     userRepo,
	}
}

func (s *savedMessageService) SaveMessage(userID, messageID, chatID string) (*chatdomain.SavedMessageResponse, error) {
	exists, err := s.savedMsgRepo.Exists(userID, messageID)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("message already saved")
	}

	msg := &chatdomain.SavedMessage{
		ID:        uuid.New().String(),
		UserID:    userID,
		MessageID: messageID,
		ChatID:    chatID,
		CreatedAt: time.Now(),
	}
	if err := s.savedMsgRepo.Save(msg); err != nil {
		return nil, err
	}
	return s.buildResponse(msg)
}

func (s *savedMessageService) GetSavedMessages(userID string, limit, offset int) ([]*chatdomain.SavedMessageResponse, int, error) {
	saved, err := s.savedMsgRepo.FindByUserID(userID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	total, err := s.savedMsgRepo.CountByUserID(userID)
	if err != nil {
		return nil, 0, err
	}

	responses := make([]*chatdomain.SavedMessageResponse, 0, len(saved))
	for _, sm := range saved {
		r, err := s.buildResponse(sm)
		if err != nil {
			continue
		}
		responses = append(responses, r)
	}
	return responses, total, nil
}

func (s *savedMessageService) DeleteSavedMessage(id, userID string) error {
	return s.savedMsgRepo.Delete(id, userID)
}

func (s *savedMessageService) buildResponse(sm *chatdomain.SavedMessage) (*chatdomain.SavedMessageResponse, error) {
	msg, err := s.msgRepo.FindByID(sm.MessageID)
	if err != nil {
		return nil, err
	}
	chat, err := s.chatRepo.FindByID(sm.ChatID)
	if err != nil {
		return nil, err
	}

	sender, err := s.userRepo.FindByID(msg.SenderID)
	if err != nil {
		return nil, err
	}

	msgResp := &messagedomain.MessageResponse{
		ID:      msg.ID,
		ChatID:  msg.ChatID,
		Content: msg.Content,
		Type:    msg.Type,
		Sender:  &userdomain.UserResponse{ID: sender.ID, Username: sender.Username, AvatarURL: sender.AvatarURL},
		CreatedAt: msg.CreatedAt,
	}

	chatResp := &chatdomain.ChatResponse{
		ID:   chat.ID,
		Name: chat.Name,
		Type: chat.Type,
	}

	return &chatdomain.SavedMessageResponse{
		ID:        sm.ID,
		Message:   msgResp,
		Chat:      chatResp,
		CreatedAt: sm.CreatedAt,
	}, nil
}
