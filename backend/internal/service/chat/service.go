package chatservice

import (
	"errors"
	"sort"
	"time"

	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	messagedomain "ChatServerGolang/backend/internal/domain/message"
	userdomain "ChatServerGolang/backend/internal/domain/user"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type chatService struct {
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	messageRepo repository.MessageRepository
	userService service.UserService
}

func NewChatService(
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	messageRepo repository.MessageRepository,
	userService service.UserService,
) service.ChatService {
	return &chatService{
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		messageRepo: messageRepo,
		userService: userService,
	}
}

func (s *chatService) CreateChat(userID string, req *chatdomain.CreateChatRequest) (*chatdomain.ChatResponse, error) {
	if req.Type == chatdomain.ChatPrivate && len(req.ParticipantIDs) != 1 {
		return nil, errors.New("private chat must have exactly 2 participants")
	}

	if req.Type == chatdomain.ChatPrivate {
		existing, _ := s.chatRepo.GetPrivateChat(userID, req.ParticipantIDs[0])
		if existing != nil {
			return s.GetChat(existing.ID, userID)
		}
	}

	now := time.Now()
	chat := &chatdomain.Chat{
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

func (s *chatService) GetChat(chatID, userID string) (*chatdomain.ChatResponse, error) {
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

	userResponses := s.fetchUserResponses(participants)

	lastMsg, _ := s.messageRepo.GetLastMessage(chatID)
	var lastMsgResponse *messagedomain.MessageResponse
	if lastMsg != nil {
		sender, _ := s.userRepo.FindByID(lastMsg.SenderID)
		if sender != nil {
			lastMsgResponse = &messagedomain.MessageResponse{
				ID:        lastMsg.ID,
				ChatID:    lastMsg.ChatID,
				Sender:    sender.ToResponse(),
				Content:   lastMsg.Content,
				Type:      lastMsg.Type,
				CreatedAt: lastMsg.CreatedAt,
				UpdatedAt: lastMsg.UpdatedAt,
				Deleted:   lastMsg.DeletedAt != nil,
			}
		}
	}

	unreadCount, _ := s.chatRepo.GetUnreadCount(chatID, userID)

	return &chatdomain.ChatResponse{
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

func (s *chatService) ListChats(userID string) ([]*chatdomain.ChatResponse, error) {
	chats, err := s.chatRepo.FindByUserIDExcludeHidden(userID)
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return []*chatdomain.ChatResponse{}, nil
	}

	chatIDs := make([]string, len(chats))
	chatMap := make(map[string]*chatdomain.Chat, len(chats))
	for i, chat := range chats {
		chatIDs[i] = chat.ID
		chatMap[chat.ID] = chat
	}

	// Batch load all participants
	participantsMap, _ := s.chatRepo.GetParticipantsByChatIDs(chatIDs)

	// Batch load all users that are participants
	allUserIDs := make([]string, 0)
	for _, participants := range participantsMap {
		for _, p := range participants {
			allUserIDs = append(allUserIDs, p.UserID)
		}
	}
	allUserIDs = uniqueStrings(allUserIDs)
	userResponses := make(map[string]*userdomain.UserResponse)
	if users, err := s.userRepo.FindByIDs(allUserIDs); err == nil {
		for id, u := range users {
			userResponses[id] = u.ToResponse()
		}
	}

	// Batch load last messages
	lastMsgs, _ := s.messageRepo.GetLastMessagesByChatIDs(chatIDs)

	// Build sender map for last messages
	lastMsgSenderIDs := make([]string, 0, len(lastMsgs))
	for _, msg := range lastMsgs {
		lastMsgSenderIDs = append(lastMsgSenderIDs, msg.SenderID)
	}
	lastMsgSenderIDs = uniqueStrings(lastMsgSenderIDs)
	lastMsgSenders, _ := s.userRepo.FindByIDs(lastMsgSenderIDs)

	// Batch load unread counts
	unreadCounts, _ := s.chatRepo.GetUnreadCounts(userID, chatIDs)

	responses := make([]*chatdomain.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		participants := participantsMap[chat.ID]
		partResponses := make([]*userdomain.UserResponse, 0, len(participants))
		for _, p := range participants {
			if ur, ok := userResponses[p.UserID]; ok {
				partResponses = append(partResponses, ur)
			}
		}

		var lastMsgResponse *messagedomain.MessageResponse
		if lastMsg, ok := lastMsgs[chat.ID]; ok {
			if sender, ok := lastMsgSenders[lastMsg.SenderID]; ok {
				lastMsgResponse = &messagedomain.MessageResponse{
					ID:        lastMsg.ID,
					ChatID:    lastMsg.ChatID,
					Sender:    sender.ToResponse(),
					Content:   lastMsg.Content,
					Type:      lastMsg.Type,
					CreatedAt: lastMsg.CreatedAt,
					UpdatedAt: lastMsg.UpdatedAt,
					Deleted:   lastMsg.DeletedAt != nil,
				}
			}
		}

		unreadCount := unreadCounts[chat.ID]

		responses = append(responses, &chatdomain.ChatResponse{
			ID:           chat.ID,
			Name:         chat.Name,
			Description:  chat.Description,
			AvatarURL:    chat.AvatarURL,
			Type:         chat.Type,
			CreatedBy:    chat.CreatedBy,
			Participants: partResponses,
			LastMessage:  lastMsgResponse,
			UnreadCount:  unreadCount,
			CreatedAt:    chat.CreatedAt,
		})
	}

	sort.Slice(responses, func(i, j int) bool {
		ti := responses[i].LastMessage
		tj := responses[j].LastMessage
		if ti != nil && tj != nil {
			return ti.CreatedAt.After(tj.CreatedAt)
		}
		if ti != nil {
			return true
		}
		return false
	})

	return responses, nil
}

func (s *chatService) buildChatResponse(chat *chatdomain.Chat, userID string) (*chatdomain.ChatResponse, error) {
	participants, err := s.chatRepo.GetParticipants(chat.ID)
	if err != nil {
		return nil, err
	}

	userResponses := s.fetchUserResponses(participants)

	lastMsg, _ := s.messageRepo.GetLastMessage(chat.ID)
	var lastMsgResponse *messagedomain.MessageResponse
	if lastMsg != nil {
		sender, _ := s.userRepo.FindByID(lastMsg.SenderID)
		if sender != nil {
			lastMsgResponse = &messagedomain.MessageResponse{
				ID:        lastMsg.ID,
				ChatID:    lastMsg.ChatID,
				Sender:    sender.ToResponse(),
				Content:   lastMsg.Content,
				Type:      lastMsg.Type,
				CreatedAt: lastMsg.CreatedAt,
				UpdatedAt: lastMsg.UpdatedAt,
				Deleted:   lastMsg.DeletedAt != nil,
			}
		}
	}

	unreadCount, _ := s.chatRepo.GetUnreadCount(chat.ID, userID)

	return &chatdomain.ChatResponse{
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

func (s *chatService) fetchUserResponses(participants []*chatdomain.ChatParticipant) []*userdomain.UserResponse {
	responses := make([]*userdomain.UserResponse, 0, len(participants))
	for _, p := range participants {
		u, err := s.userRepo.FindByID(p.UserID)
		if err != nil {
			continue
		}
		responses = append(responses, u.ToResponse())
	}
	return responses
}

func (s *chatService) DeleteChat(chatID, userID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.CreatedBy != userID {
		return errors.New("only the creator can delete the chat")
	}

	return s.chatRepo.Delete(chatID)
}

func (s *chatService) AddParticipant(chatID, userID, requesterID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != chatdomain.ChatGroup {
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

func (s *chatService) RemoveParticipant(chatID, userID, requesterID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.CreatedBy == requesterID {
		// Creator can remove anyone except themselves
	} else if requesterID == userID {
		// User can remove themselves (leave)
	} else {
		// Check if requester is admin
		participants, _ := s.chatRepo.GetParticipants(chatID)
		isAdmin := false
		for _, p := range participants {
			if p.UserID == requesterID && p.Role == "admin" {
				isAdmin = true
				break
			}
		}
		if !isAdmin {
			return errors.New("access denied")
		}
	}

	if userID == chat.CreatedBy {
		return errors.New("cannot remove the creator")
	}

	return s.chatRepo.RemoveParticipant(chatID, userID)
}

func (s *chatService) MarkAsRead(chatID, userID string) error {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	return s.chatRepo.UpdateLastRead(chatID, userID)
}

func (s *chatService) GetUnreadCount(chatID, userID string) (int, error) {
	return s.chatRepo.GetUnreadCount(chatID, userID)
}

func (s *chatService) SetRole(chatID, targetUserID, requesterID, role string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != chatdomain.ChatGroup {
		return errors.New("only group chats have roles")
	}

	if chat.CreatedBy != requesterID {
		participants, _ := s.chatRepo.GetParticipants(chatID)
		isAdmin := false
		for _, p := range participants {
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

func (s *chatService) LeaveGroup(chatID, userID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != chatdomain.ChatGroup {
		return errors.New("can only leave group chats")
	}

	if chat.CreatedBy == userID {
		return errors.New("creator cannot leave; transfer ownership or delete the chat")
	}

	return s.chatRepo.RemoveParticipant(chatID, userID)
}

func (s *chatService) UpdateGroup(chatID, userID string, req *chatdomain.UpdateGroupRequest) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}

	if chat.Type != chatdomain.ChatGroup {
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

func (s *chatService) HideChat(chatID, userID string) error {
	return s.chatRepo.HideChat(userID, chatID)
}

func (s *chatService) SearchChats(userID, query string) ([]*chatdomain.ChatResponse, error) {
	if query == "" {
		return s.ListChats(userID)
	}
	chats, err := s.chatRepo.SearchByName(userID, query)
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return []*chatdomain.ChatResponse{}, nil
	}

	chatIDs := make([]string, len(chats))
	chatMap := make(map[string]*chatdomain.Chat, len(chats))
	for i, chat := range chats {
		chatIDs[i] = chat.ID
		chatMap[chat.ID] = chat
	}

	participantsMap, _ := s.chatRepo.GetParticipantsByChatIDs(chatIDs)
	allUserIDs := make([]string, 0)
	for _, participants := range participantsMap {
		for _, p := range participants {
			allUserIDs = append(allUserIDs, p.UserID)
		}
	}
	allUserIDs = uniqueStrings(allUserIDs)
	userResponses := make(map[string]*userdomain.UserResponse)
	if users, err := s.userRepo.FindByIDs(allUserIDs); err == nil {
		for id, u := range users {
			userResponses[id] = u.ToResponse()
		}
	}

	lastMsgs, _ := s.messageRepo.GetLastMessagesByChatIDs(chatIDs)
	lastMsgSenderIDs := make([]string, 0, len(lastMsgs))
	for _, msg := range lastMsgs {
		lastMsgSenderIDs = append(lastMsgSenderIDs, msg.SenderID)
	}
	lastMsgSenderIDs = uniqueStrings(lastMsgSenderIDs)
	lastMsgSenders, _ := s.userRepo.FindByIDs(lastMsgSenderIDs)

	unreadCounts, _ := s.chatRepo.GetUnreadCounts(userID, chatIDs)

	responses := make([]*chatdomain.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		participants := participantsMap[chat.ID]
		partResponses := make([]*userdomain.UserResponse, 0, len(participants))
		for _, p := range participants {
			if ur, ok := userResponses[p.UserID]; ok {
				partResponses = append(partResponses, ur)
			}
		}

		var lastMsgResponse *messagedomain.MessageResponse
		if lastMsg, ok := lastMsgs[chat.ID]; ok {
			if sender, ok := lastMsgSenders[lastMsg.SenderID]; ok {
				lastMsgResponse = &messagedomain.MessageResponse{
					ID:        lastMsg.ID,
					ChatID:    lastMsg.ChatID,
					Sender:    sender.ToResponse(),
					Content:   lastMsg.Content,
					Type:      lastMsg.Type,
					CreatedAt: lastMsg.CreatedAt,
					UpdatedAt: lastMsg.UpdatedAt,
					Deleted:   lastMsg.DeletedAt != nil,
				}
			}
		}

		responses = append(responses, &chatdomain.ChatResponse{
			ID:           chat.ID,
			Name:         chat.Name,
			Description:  chat.Description,
			AvatarURL:    chat.AvatarURL,
			Type:         chat.Type,
			CreatedBy:    chat.CreatedBy,
			Participants: partResponses,
			LastMessage:  lastMsgResponse,
			UnreadCount:  unreadCounts[chat.ID],
			CreatedAt:    chat.CreatedAt,
		})
	}

	return responses, nil
}

func (s *chatService) PinChat(chatID, userID string) error {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	return s.chatRepo.PinChat(userID, chatID)
}

func (s *chatService) UnpinChat(chatID, userID string) error {
	return s.chatRepo.UnpinChat(userID, chatID)
}

func (s *chatService) ArchiveChat(chatID, userID string) error {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	return s.chatRepo.ArchiveChat(userID, chatID)
}

func (s *chatService) UnarchiveChat(chatID, userID string) error {
	return s.chatRepo.UnarchiveChat(userID, chatID)
}

func (s *chatService) ListArchivedChats(userID string) ([]*chatdomain.ChatResponse, error) {
	chats, err := s.chatRepo.FindByUserIDArchived(userID)
	if err != nil {
		return nil, err
	}
	if len(chats) == 0 {
		return []*chatdomain.ChatResponse{}, nil
	}

	chatIDs := make([]string, len(chats))
	for i, chat := range chats {
		chatIDs[i] = chat.ID
	}

	participantsMap, _ := s.chatRepo.GetParticipantsByChatIDs(chatIDs)
	allUserIDs := make([]string, 0)
	for _, participants := range participantsMap {
		for _, p := range participants {
			allUserIDs = append(allUserIDs, p.UserID)
		}
	}
	allUserIDs = uniqueStrings(allUserIDs)
	userResponses := make(map[string]*userdomain.UserResponse)
	if users, err := s.userRepo.FindByIDs(allUserIDs); err == nil {
		for id, u := range users {
			userResponses[id] = u.ToResponse()
		}
	}

	lastMsgs, _ := s.messageRepo.GetLastMessagesByChatIDs(chatIDs)
	lastMsgSenderIDs := make([]string, 0, len(lastMsgs))
	for _, msg := range lastMsgs {
		lastMsgSenderIDs = append(lastMsgSenderIDs, msg.SenderID)
	}
	lastMsgSenderIDs = uniqueStrings(lastMsgSenderIDs)
	lastMsgSenders, _ := s.userRepo.FindByIDs(lastMsgSenderIDs)

	unreadCounts, _ := s.chatRepo.GetUnreadCounts(userID, chatIDs)

	responses := make([]*chatdomain.ChatResponse, 0, len(chats))
	for _, chat := range chats {
		participants := participantsMap[chat.ID]
		partResponses := make([]*userdomain.UserResponse, 0, len(participants))
		for _, p := range participants {
			if ur, ok := userResponses[p.UserID]; ok {
				partResponses = append(partResponses, ur)
			}
		}

		var lastMsgResponse *messagedomain.MessageResponse
		if lastMsg, ok := lastMsgs[chat.ID]; ok {
			if sender, ok := lastMsgSenders[lastMsg.SenderID]; ok {
				lastMsgResponse = &messagedomain.MessageResponse{
					ID:        lastMsg.ID,
					ChatID:    lastMsg.ChatID,
					Sender:    sender.ToResponse(),
					Content:   lastMsg.Content,
					Type:      lastMsg.Type,
					CreatedAt: lastMsg.CreatedAt,
					UpdatedAt: lastMsg.UpdatedAt,
					Deleted:   lastMsg.DeletedAt != nil,
				}
			}
		}

		responses = append(responses, &chatdomain.ChatResponse{
			ID:           chat.ID,
			Name:         chat.Name,
			Description:  chat.Description,
			AvatarURL:    chat.AvatarURL,
			Type:         chat.Type,
			CreatedBy:    chat.CreatedBy,
			Participants: partResponses,
			LastMessage:  lastMsgResponse,
			UnreadCount:  unreadCounts[chat.ID],
			CreatedAt:    chat.CreatedAt,
		})
	}

	sort.Slice(responses, func(i, j int) bool {
		ti := responses[i].LastMessage
		tj := responses[j].LastMessage
		if ti != nil && tj != nil {
			return ti.CreatedAt.After(tj.CreatedAt)
		}
		if ti != nil {
			return true
		}
		return false
	})

	return responses, nil
}

func (s *chatService) TransferOwnership(chatID, fromUserID, toUserID string) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}
	if chat.CreatedBy != fromUserID {
		return errors.New("only the creator can transfer ownership")
	}
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, toUserID)
	if !isParticipant {
		return errors.New("target user is not a participant")
	}
	if err := s.chatRepo.SetRole(chatID, fromUserID, "admin"); err != nil {
		return err
	}
	if err := s.chatRepo.SetRole(chatID, toUserID, "owner"); err != nil {
		return err
	}
	chat.CreatedBy = toUserID
	return s.chatRepo.Update(chat)
}

func (s *chatService) SetSlowMode(chatID, userID string, seconds int) error {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return errors.New("chat not found")
	}
	if chat.CreatedBy != userID {
		return errors.New("only the creator can set slow mode")
	}
	return s.chatRepo.SetSlowMode(chatID, seconds)
}

func (s *chatService) GetParticipants(chatID string) ([]*chatdomain.ChatParticipant, error) {
	return s.chatRepo.GetParticipants(chatID)
}

func uniqueStrings(slice []string) []string {
	seen := make(map[string]struct{}, len(slice))
	res := make([]string, 0, len(slice))
	for _, s := range slice {
		if _, ok := seen[s]; !ok {
			seen[s] = struct{}{}
			res = append(res, s)
		}
	}
	return res
}

