package service

import (
	"errors"
	"sort"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
)

type ChatService struct {
	chatRepo    *repository.ChatRepository
	userRepo    *repository.UserRepository
	messageRepo *repository.MessageRepository
	userService *UserService
}

func NewChatService(
	chatRepo *repository.ChatRepository,
	userRepo *repository.UserRepository,
	messageRepo *repository.MessageRepository,
	userService *UserService,
) *ChatService {
	return &ChatService{
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		userService: userService,
	}
}

func (s *ChatService) CreateChat(userID string, req *domain.CreateChatRequest) (*domain.ChatResponse, error) {
	if req.Type == domain.ChatPrivate && len(req.ParticipantIDs) != 1 {
		return nil, errors.New("private chat must have exactly 2 participants")
	}

	if req.Type == domain.ChatPrivate {
		existing, _ := s.chatRepo.GetPrivateChat(userID, req.ParticipantIDs[0])
		if existing != nil {
			return s.GetChat(existing.ID, userID)
		}
	}

	now := time.Now()
	chat := &domain.Chat{
		ID:          uuid.New().String(),
		Name:        req.Name,
		Description: req.Description,
		Type:        req.Type,
		CreatedBy:   userID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.chatRepo.Create(chat); err != nil {
		return nil, err
	}

	participantIDs := append([]string{userID}, req.ParticipantIDs...)
	for _, pid := range participantIDs {
		role := "member"
		if pid == userID {
			role = "owner"
		}
		if err := s.chatRepo.AddParticipant(chat.ID, pid, role); err != nil {
			return nil, err
		}
	}

	return s.GetChat(chat.ID, userID)
}

func (s *ChatService) GetChat(chatID, userID string) (*domain.ChatResponse, error) {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	participants, err := s.chatRepo.GetParticipants(chatID)
	if err != nil {
		return nil, err
	}

	var userResponses []*domain.UserResponse
	for _, p := range participants {
		u, err := s.userRepo.FindByID(p.UserID)
		if err != nil {
			continue
		}
		userResponses = append(userResponses, u.ToResponse())
	}

	lastMsg, _ := s.messageRepo.GetLastMessage(chatID)
	var lastMsgResponse *domain.MessageResponse
	if lastMsg != nil {
		sender, _ := s.userRepo.FindByID(lastMsg.SenderID)
		senderResp := sender.ToResponse()
		lastMsgResponse = &domain.MessageResponse{
			ID:        lastMsg.ID,
			ChatID:    lastMsg.ChatID,
			Sender:    senderResp,
			Content:   lastMsg.Content,
			Type:      lastMsg.Type,
			CreatedAt: lastMsg.CreatedAt,
			UpdatedAt: lastMsg.UpdatedAt,
			Deleted:   lastMsg.DeletedAt != nil,
		}
	}

	unreadCount, _ := s.chatRepo.GetUnreadCount(chatID, userID)

	return &domain.ChatResponse{
		ID:           chat.ID,
		Name:         chat.Name,
		Description:  chat.Description,
		AvatarURL:    chat.AvatarURL,
		Type:         chat.Type,
		CreatedBy:    chat.CreatedBy,
		Participants: userResponses,
		LastMessage:  lastMsgResponse,
		UnreadCount:  unreadCount,
		CreatedAt:    chat.CreatedAt,
	}, nil
}

func (s *ChatService) ListChats(userID string) ([]*domain.ChatResponse, error) {
	chats, err := s.chatRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.ChatResponse, 0)
	for _, chat := range chats {
		resp, err := s.GetChat(chat.ID, userID)
		if err != nil {
			continue
		}
		responses = append(responses, resp)
	}

	if len(responses) > 0 {
			sort.Slice(responses, func(i, j int) bool {
			if responses[i].LastMessage != nil && responses[j].LastMessage != nil {
				return responses[i].LastMessage.CreatedAt.After(responses[j].LastMessage.CreatedAt)
			}
			if responses[i].LastMessage != nil {
				return true
			}
			return false
		})
	}

	return responses, nil
}

func (s *ChatService) DeleteChat(chatID, userID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.CreatedBy != userID {
		return errors.New("only the creator can delete the chat")
	}

	participants, _ := s.chatRepo.GetParticipants(chatID)
	for _, p := range participants {
		s.chatRepo.RemoveParticipant(chatID, p.UserID)
	}

	return s.chatRepo.Delete(chatID)
}

func (s *ChatService) AddParticipant(chatID, userID, requesterID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != domain.ChatGroup {
		return errors.New("can only add participants to group chats")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(chatID, requesterID)
	if !isParticipant {
		return errors.New("access denied")
	}

	_, err = s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}

	return s.chatRepo.AddParticipant(chatID, userID, "member")
}

func (s *ChatService) RemoveParticipant(chatID, userID, requesterID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.CreatedBy != requesterID && requesterID != userID {
		return errors.New("access denied")
	}

	if userID == chat.CreatedBy {
		return errors.New("cannot remove the creator")
	}

	return s.chatRepo.RemoveParticipant(chatID, userID)
}

func (s *ChatService) MarkAsRead(chatID, userID string) error {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	return s.chatRepo.UpdateLastRead(chatID, userID)
}

func (s *ChatService) GetUnreadCount(chatID, userID string) (int, error) {
	return s.chatRepo.GetUnreadCount(chatID, userID)
}

func (s *ChatService) SetRole(chatID, targetUserID, requesterID, role string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != domain.ChatGroup {
		return errors.New("only group chats have roles")
	}

	if chat.CreatedBy != requesterID {
		requesterRole, _ := s.chatRepo.GetParticipants(chatID)
		isAdmin := false
		for _, p := range requesterRole {
			if p.UserID == requesterID && p.Role == "admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			return errors.New("only admins can change roles")
		}
	}

	if targetUserID == chat.CreatedBy {
		return errors.New("cannot change creator's role")
	}

	return s.chatRepo.SetRole(chatID, targetUserID, role)
}

func (s *ChatService) LeaveGroup(chatID, userID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != domain.ChatGroup {
		return errors.New("can only leave group chats")
	}

	if chat.CreatedBy == userID {
		return errors.New("creator cannot leave; transfer ownership or delete the chat")
	}

	return s.chatRepo.RemoveParticipant(chatID, userID)
}

func (s *ChatService) UpdateGroup(chatID, userID string, req *domain.UpdateGroupRequest) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != domain.ChatGroup {
		return errors.New("only group chats can be updated")
	}

	if chat.CreatedBy != userID {
		return errors.New("only the creator can update the group")
	}

	if req.Name != "" {
		chat.Name = req.Name
	}
	if req.Description != "" {
		chat.Description = req.Description
	}
	if req.AvatarURL != "" {
		chat.AvatarURL = req.AvatarURL
	}

	return s.chatRepo.Update(chat)
}
