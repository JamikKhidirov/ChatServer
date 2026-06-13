package channelservice

import (
	"errors"

	channeldomain "ChatServerGolang/backend/internal/domain/channel"
	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"
)

type channelService struct {
	channelRepo repository.ChannelSubscriberRepository
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
}

func NewChannelService(channelRepo repository.ChannelSubscriberRepository, chatRepo repository.ChatRepository, userRepo repository.UserRepository) service.ChannelService {
	return &channelService{
		channelRepo: channelRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
	}
}

func (s *channelService) Subscribe(channelID, userID string) error {
	chat, err := s.chatRepo.FindByID(channelID)
	if err != nil {
		return errors.New("channel not found")
	}
	if chat.Type != chatdomain.ChatChannel {
		return errors.New("not a channel")
	}
	existingRole, _ := s.channelRepo.GetRole(channelID, userID)
	if existingRole != "" {
		return errors.New("already subscribed")
	}
	return s.channelRepo.Subscribe(channelID, userID, string(chatdomain.RoleSubscriber))
}

func (s *channelService) Unsubscribe(channelID, userID string) error {
	chat, err := s.chatRepo.FindByID(channelID)
	if err != nil {
		return errors.New("channel not found")
	}
	if chat.Type != chatdomain.ChatChannel {
		return errors.New("not a channel")
	}
	role, _ := s.channelRepo.GetRole(channelID, userID)
	if role == "" {
		return errors.New("not subscribed")
	}
	if role == string(chatdomain.RoleOwner) || role == string(chatdomain.RoleAdmin) {
		return errors.New("admins cannot unsubscribe")
	}
	return s.channelRepo.Unsubscribe(channelID, userID)
}

func (s *channelService) GetSubscribers(channelID, userID string) ([]*channeldomain.ChannelSubscriber, error) {
	role, _ := s.channelRepo.GetRole(channelID, userID)
	if !chatdomain.IsAdminRole(role) {
		return nil, errors.New("access denied")
	}
	return s.channelRepo.GetSubscribers(channelID)
}

func (s *channelService) GetSubscribedChannels(userID string) ([]*chatdomain.ChatResponse, error) {
	channelIDs, err := s.channelRepo.GetSubscribedChannels(userID)
	if err != nil {
		return nil, err
	}
	chats := make([]*chatdomain.ChatResponse, 0, len(channelIDs))
	for _, cid := range channelIDs {
		chat, err := s.chatRepo.FindByID(cid)
		if err != nil {
			continue
		}
		chats = append(chats, &chatdomain.ChatResponse{
			ID:        chat.ID,
			Name:      chat.Name,
			Type:      chat.Type,
			CreatedBy: chat.CreatedBy,
		})
	}
	return chats, nil
}

func (s *channelService) SetSubscriberRole(channelID, targetUserID, requesterID, role string) error {
	reqRole, _ := s.channelRepo.GetRole(channelID, requesterID)
	if !chatdomain.IsAdminRole(reqRole) {
		return errors.New("access denied")
	}
	return s.channelRepo.SetRole(channelID, targetUserID, role)
}

func (s *channelService) IsSubscribed(channelID, userID string) (bool, error) {
	return s.channelRepo.IsSubscribed(channelID, userID)
}
