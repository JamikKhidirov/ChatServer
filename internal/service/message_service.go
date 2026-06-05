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

	// Check if blocked
	blocked, _ := s.userService.IsBlocked(senderID, chatID)
	if blocked {
		return nil, errors.New("you are blocked from sending messages")
	}

	if req.ReplyToID != nil && *req.ReplyToID != "" {
		replyMsg, err := s.messageRepo.FindByID(*req.ReplyToID)
		if err != nil || replyMsg.ChatID != chatID {
			return nil, errors.New("invalid reply message")
		}
	}

	var forwardFrom *string
	if req.ForwardMsgID != nil && *req.ForwardMsgID != "" {
		forwardMsg, err := s.messageRepo.FindByID(*req.ForwardMsgID)
		if err != nil {
			return nil, errors.New("invalid forwarded message")
		}
		if forwardMsg.ChatID != chatID {
			return nil, errors.New("forwarded message is not in this chat")
		}
		forwardFrom = &forwardMsg.SenderID
	}

	now := time.Now()
	msg := &domain.Message{
		ID:          uuid.New().String(),
		ChatID:      chatID,
		SenderID:    senderID,
		Content:     req.Content,
		Type:        req.Type,
		ReplyToID:   req.ReplyToID,
		ForwardFrom: forwardFrom,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	s.chatRepo.Update(&domain.Chat{ID: chatID, UpdatedAt: now})

	return s.getMessageResponse(msg)
}

func (s *MessageService) SendFileMessage(chatID, senderID, fileName, filePath string, fileSize int64, replyToID *string) (*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, senderID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	now := time.Now()
	fileType := domain.MessageFile
	msg := &domain.Message{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   fileName,
		Type:      fileType,
		ReplyToID: replyToID,
		FileName:  fileName,
		FileSize:  fileSize,
		FilePath:  filePath,
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

func (s *MessageService) SearchMessages(chatID, userID, query string, limit, offset int) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.Search(chatID, query, limit, offset)
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

func (s *MessageService) ResendMessage(chatID, userID, msgID string) (*domain.MessageResponse, error) {
	original, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	now := time.Now()
	msg := &domain.Message{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		SenderID:  userID,
		Content:   original.Content,
		Type:      original.Type,
		CreatedAt: now,
		UpdatedAt: now,
	}

	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	s.chatRepo.Update(&domain.Chat{ID: chatID, UpdatedAt: now})

	return s.getMessageResponse(msg)
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

// Reactions
func (s *MessageService) AddReaction(msgID, userID, emoji string) (*domain.MessageResponse, error) {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	if err := s.messageRepo.AddReaction(msgID, userID, emoji); err != nil {
		return nil, err
	}

	updated, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, err
	}

	return s.getMessageResponse(updated)
}

func (s *MessageService) RemoveReaction(msgID, userID, emoji string) (*domain.MessageResponse, error) {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	if err := s.messageRepo.RemoveReaction(msgID, userID, emoji); err != nil {
		return nil, err
	}

	updated, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, err
	}

	return s.getMessageResponse(updated)
}

// Pinned messages
func (s *MessageService) TogglePin(msgID, userID string, pin bool) (*domain.MessageResponse, error) {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	chat, err := s.chatRepo.FindByID(msg.ChatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	if chat.CreatedBy != userID {
		participants, _ := s.chatRepo.GetParticipants(msg.ChatID)
		isAdmin := false
		for _, p := range participants {
			if p.UserID == userID && p.Role == "admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			return nil, errors.New("only admins can pin messages")
		}
	}

	if err := s.messageRepo.TogglePin(msgID, pin); err != nil {
		return nil, err
	}

	updated, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, err
	}

	return s.getMessageResponse(updated)
}

func (s *MessageService) GetPinnedMessages(chatID, userID string) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.GetPinned(chatID)
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

// Read receipts
func (s *MessageService) MarkMessageRead(msgID, userID string) error {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}

	return s.messageRepo.AddReadReceipt(msgID, userID)
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
		FileName:  msg.FileName,
		FileSize:  msg.FileSize,
		Pinned:    msg.Pinned,
		CreatedAt: msg.CreatedAt,
		UpdatedAt: msg.UpdatedAt,
		Edited:    msg.UpdatedAt.After(msg.CreatedAt),
		Deleted:   msg.DeletedAt != nil,
	}

	if msg.FilePath != "" {
		resp.FileURL = "/uploads/" + msg.FilePath
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

	if msg.ForwardFrom != nil && *msg.ForwardFrom != "" {
		fwdUser, err := s.userRepo.FindByID(*msg.ForwardFrom)
		if err == nil {
			resp.ForwardFrom = fwdUser.ToResponse()
		}
	}

	// Get reactions
	reactions, _ := s.messageRepo.GetReactions(msg.ID)
	for _, r := range reactions {
		u, err := s.userRepo.FindByID(r.UserID)
		if err == nil {
			r.User = u.ToResponse()
		}
	}
	resp.Reactions = reactions

	// Get read receipts
	receipts, _ := s.messageRepo.GetReadReceipts(msg.ID)
	for _, r := range receipts {
		u, err := s.userRepo.FindByID(r.UserID)
		if err == nil {
			resp.ReadBy = append(resp.ReadBy, u.ToResponse())
		}
	}

	if resp.Deleted {
		resp.Content = ""
	}

	return resp, nil
}
