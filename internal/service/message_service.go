package service

import (
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
)

type MessageService struct {
	messageRepo *repository.MessageRepository
	chatRepo    *repository.ChatRepository
	userRepo    *repository.UserRepository
	userService *UserService
}

func NewMessageService(
	messageRepo *repository.MessageRepository,
	chatRepo *repository.ChatRepository,
	userRepo *repository.UserRepository,
	userService *UserService,
) *MessageService {
	return &MessageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		userService: userService,
	}
}

func (s *MessageService) SendMessage(chatID, senderID string, req *domain.SendMessageRequest) (*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, senderID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	if req.ReplyToID != nil && *req.ReplyToID != "" {
		replyMsg, err := s.messageRepo.FindByID(*req.ReplyToID)
		if err != nil || replyMsg.ChatID != chatID {
			return nil, errors.New("invalid reply message")
		}
	}

	now := time.Now()
	msg := &domain.Message{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   req.Content,
		Type:      req.Type,
		ReplyToID: req.ReplyToID,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	s.chatRepo.Update(&domain.Chat{ID: chatID, UpdatedAt: now})

	return s.getMessageResponse(msg)
}

func (s *MessageService) GetMessages(chatID, userID string, limit, offset int) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.FindByChatID(chatID, limit, offset)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.MessageResponse, 0)
	for _, msg := range messages {
		resp, err := s.getMessageResponse(msg)
		if err != nil {
			continue
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *MessageService) EditMessage(msgID, userID string, req *domain.EditMessageRequest) (*domain.MessageResponse, error) {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	if msg.SenderID != userID {
		return nil, errors.New("cannot edit another user's message")
	}

	if msg.DeletedAt != nil {
		return nil, errors.New("cannot edit deleted message")
	}

	msg.Content = req.Content
	msg.UpdatedAt = time.Now()

	if err := s.messageRepo.Update(msg); err != nil {
		return nil, err
	}

	return s.getMessageResponse(msg)
}

func (s *MessageService) GetMessageByID(msgID, userID string) (*domain.MessageResponse, error) {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	return s.getMessageResponse(msg)
}

func (s *MessageService) DeleteMessage(msgID, userID string) error {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return errors.New("message not found")
	}

	if msg.SenderID != userID {
		chat, err := s.chatRepo.FindByID(msg.ChatID)
		if err != nil || chat.CreatedBy != userID {
			return errors.New("access denied")
		}
	}

	return s.messageRepo.SoftDelete(msgID)
}

func (s *MessageService) getMessageResponse(msg *domain.Message) (*domain.MessageResponse, error) {
	sender, err := s.userRepo.FindByID(msg.SenderID)
	if err != nil {
		return nil, err
	}

	resp := &domain.MessageResponse{
		ID:        msg.ID,
		ChatID:    msg.ChatID,
		Sender:    sender.ToResponse(),
		Content:   msg.Content,
		Type:      msg.Type,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		Edited:    msg.UpdatedAt.After(msg.CreatedAt),
		Deleted:   msg.DeletedAt != nil,
	}

	if msg.ReplyToID != nil && *msg.ReplyToID != "" {
		replyMsg, err := s.messageRepo.FindByID(*msg.ReplyToID)
		if err == nil {
			replySender, _ := s.userRepo.FindByID(replyMsg.SenderID)
			resp.ReplyTo = &domain.MessageResponse{
				ID:      replyMsg.ID,
				Content: replyMsg.Content,
				Type:    replyMsg.Type,
				Sender:  replySender.ToResponse(),
			}
		}
	}

	if resp.Deleted {
		resp.Content = ""
	}

	return resp, nil
}
