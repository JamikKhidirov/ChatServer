package service

import (
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
)

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

type messageService struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	userService UserService
}

func NewMessageService(
	messageRepo repository.MessageRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	userService UserService,
) MessageService {
	return &messageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		userService: userService,
	}
}

func (s *messageService) getMessageResponse(msg *domain.Message) (*domain.MessageResponse, error) {
	responses, err := s.buildMessageResponses([]*domain.Message{msg})
	if err != nil {
		return nil, err
	}
	if len(responses) == 0 {
		return nil, errors.New("failed to build message response")
	}
	return responses[0], nil
}

func (s *messageService) SendMessage(chatID, senderID string, req *domain.SendMessageRequest) (*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, senderID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	// Check block: get other participants and check each direction
	participants, err := s.chatRepo.GetParticipants(chatID)
	if err == nil {
		for _, p := range participants {
			if p.UserID != senderID {
				blocked, _ := s.userService.IsBlocked(senderID, p.UserID)
				if blocked {
					return nil, errors.New("you are blocked from sending messages")
				}
			}
		}
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

func (s *messageService) SendFileMessage(chatID, senderID, fileName, filePath string, fileSize int64, replyToID *string) (*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, senderID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	// Check block
	participants, err := s.chatRepo.GetParticipants(chatID)
	if err == nil {
		for _, p := range participants {
			if p.UserID != senderID {
				blocked, _ := s.userService.IsBlocked(senderID, p.UserID)
				if blocked {
					return nil, errors.New("you are blocked from sending messages")
				}
			}
		}
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

func (s *messageService) GetMessages(chatID, userID string, limit, offset int) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.FindByChatID(chatID, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

func (s *messageService) SearchMessages(chatID, userID, query string, limit, offset int) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.Search(chatID, query, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

func (s *messageService) ResendMessage(chatID, userID, msgID string) (*domain.MessageResponse, error) {
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

func (s *messageService) EditMessage(msgID, userID string, req *domain.EditMessageRequest) (*domain.MessageResponse, error) {
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

func (s *messageService) GetMessageByID(msgID, userID string) (*domain.MessageResponse, error) {
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

func (s *messageService) DeleteMessage(msgID, userID string) error {
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

func (s *messageService) AddReaction(msgID, userID, emoji string) (*domain.MessageResponse, error) {
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

func (s *messageService) RemoveReaction(msgID, userID, emoji string) (*domain.MessageResponse, error) {
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

func (s *messageService) TogglePin(msgID, userID string, pin bool) (*domain.MessageResponse, error) {
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

func (s *messageService) GetPinnedMessages(chatID, userID string) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.GetPinned(chatID)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

func (s *messageService) MarkMessageRead(msgID, userID string) error {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}

	if err := s.messageRepo.AddReadReceipt(msgID, userID); err != nil {
		return err
	}

	return s.chatRepo.UpdateLastRead(msg.ChatID, userID)
}

func (s *messageService) SearchAllMessages(userID, query string, limit, offset int) ([]*domain.MessageResponse, error) {
	messages, err := s.messageRepo.SearchByUser(userID, query, limit, offset)
	if err != nil {
		return nil, err
	}
	return s.buildMessageResponses(messages)
}

func (s *messageService) ForwardMessage(msgID, fromChatID, toChatID, userID string) (*domain.MessageResponse, error) {
	original, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(fromChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied to source chat")
	}

	isParticipant, _ = s.chatRepo.IsParticipant(toChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied to target chat")
	}

	now := time.Now()
	msg := &domain.Message{
		ID:          uuid.New().String(),
		ChatID:      toChatID,
		SenderID:    userID,
		Content:     original.Content,
		Type:        original.Type,
		FileName:    original.FileName,
		FileSize:    original.FileSize,
		FilePath:    original.FilePath,
		ForwardFrom: &original.SenderID,
		CreatedAt:   now,
		UpdatedAt:   now,
	}

	if err := s.messageRepo.Create(msg); err != nil {
		return nil, err
	}

	s.chatRepo.Update(&domain.Chat{ID: toChatID, UpdatedAt: now})

	return s.getMessageResponse(msg)
}

func (s *messageService) DeleteMessageForMe(msgID, userID string) error {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}

	return s.messageRepo.DeleteMessageForMe(userID, msgID)
}

func (s *messageService) StarMessage(msgID, userID string) (*domain.MessageResponse, error) {
	msg, err := s.messageRepo.FindByID(msgID)
	if err != nil {
		return nil, errors.New("message not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(msg.ChatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	if err := s.messageRepo.StarMessage(userID, msgID, msg.ChatID); err != nil {
		return nil, err
	}

	return s.getMessageResponse(msg)
}

func (s *messageService) UnstarMessage(msgID, userID string) error {
	return s.messageRepo.UnstarMessage(userID, msgID)
}

func (s *messageService) GetStarredMessages(userID string) ([]*domain.StarredMessageResponse, error) {
	starred, err := s.messageRepo.GetStarredMessages(userID)
	if err != nil {
		return nil, err
	}
	if len(starred) == 0 {
		return []*domain.StarredMessageResponse{}, nil
	}

	msgIDs := make([]string, len(starred))
	chatIDs := make([]string, 0, len(starred))
	for _, sm := range starred {
		msgIDs = append(msgIDs, sm.MessageID)
		chatIDs = append(chatIDs, sm.ChatID)
	}
	chatIDs = uniqueStrings(chatIDs)

	msgMap, _ := s.messageRepo.FindByIDs(msgIDs)
	chatMap := make(map[string]*domain.Chat, len(chatIDs))
	for _, cid := range chatIDs {
		if chat, err := s.chatRepo.FindByID(cid); err == nil {
			chatMap[cid] = chat
		}
	}

	allUserIDs := make([]string, 0)
	for _, m := range msgMap {
		allUserIDs = append(allUserIDs, m.SenderID)
	}
	allUserIDs = uniqueStrings(allUserIDs)
	userMap, _ := s.userRepo.FindByIDs(allUserIDs)

	responses := make([]*domain.StarredMessageResponse, 0, len(starred))
	for _, sm := range starred {
		msg, ok := msgMap[sm.MessageID]
		if !ok {
			continue
		}

		msgResp, _ := s.getMessageResponse(msg)

		var chatResp *domain.ChatResponse
		if chat, ok := chatMap[sm.ChatID]; ok {
			participants, _ := s.chatRepo.GetParticipants(chat.ID)
			userResponses := make([]*domain.UserResponse, 0, len(participants))
			for _, p := range participants {
				if u, ok := userMap[p.UserID]; ok {
					userResponses = append(userResponses, u.ToResponse())
				}
			}
			chatResp = &domain.ChatResponse{
				ID:           chat.ID,
				Name:         chat.Name,
				Type:         chat.Type,
				Participants: userResponses,
			}
		}

		responses = append(responses, &domain.StarredMessageResponse{
			Message:   msgResp,
			Chat:      chatResp,
			CreatedAt: sm.CreatedAt,
		})
	}

	return responses, nil
}

func (s *messageService) GetChatMedia(chatID, userID, mediaType string, limit, offset int) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.FindMediaByChatID(chatID, mediaType, limit, offset)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

func (s *messageService) ExportChat(chatID, userID string) ([]*domain.MessageResponse, error) {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	messages, err := s.messageRepo.FindByChatID(chatID, 10000, 0)
	if err != nil {
		return nil, err
	}

	return s.buildMessageResponses(messages)
}

func (s *messageService) buildMessageResponses(messages []*domain.Message) ([]*domain.MessageResponse, error) {
	if len(messages) == 0 {
		return []*domain.MessageResponse{}, nil
	}

	// Collect all needed IDs upfront
	senderIDs := make([]string, 0, len(messages))
	replyToIDs := make([]string, 0)
	forwardFromIDs := make([]string, 0)
	msgIDs := make([]string, len(messages))
	for i, msg := range messages {
		senderIDs = append(senderIDs, msg.SenderID)
		msgIDs[i] = msg.ID
		if msg.ReplyToID != nil && *msg.ReplyToID != "" {
			replyToIDs = append(replyToIDs, *msg.ReplyToID)
		}
		if msg.ForwardFrom != nil && *msg.ForwardFrom != "" {
			forwardFromIDs = append(forwardFromIDs, *msg.ForwardFrom)
		}
	}

	// Batch load users for senders, forward-froms
	allUserIDs := append(senderIDs, forwardFromIDs...)
	allUserIDs = uniqueStrings(allUserIDs)
	userMap, _ := s.userRepo.FindByIDs(allUserIDs)

	// Batch load reply-to messages
	replyMsgIDs := uniqueStrings(replyToIDs)
	replyMsgs := make(map[string]*domain.Message)
	var replySenders map[string]*domain.User
	if len(replyMsgIDs) > 0 {
		replyMsgs, _ = s.messageRepo.FindByIDs(replyMsgIDs)
		replySenderIDs := make([]string, 0, len(replyMsgs))
		for _, m := range replyMsgs {
			replySenderIDs = append(replySenderIDs, m.SenderID)
		}
		replySenderIDs = uniqueStrings(replySenderIDs)
		replySenders, _ = s.userRepo.FindByIDs(replySenderIDs)
	}

	// Batch load reactions
	reactionsMap, _ := s.messageRepo.GetReactionsByMessageIDs(msgIDs)
	reactionUserIDs := make([]string, 0)
	for _, reactions := range reactionsMap {
		for _, r := range reactions {
			reactionUserIDs = append(reactionUserIDs, r.UserID)
		}
	}
	reactionUserIDs = uniqueStrings(reactionUserIDs)
	reactionUsers, _ := s.userRepo.FindByIDs(reactionUserIDs)

	// Batch load read receipts
	receiptsMap, _ := s.messageRepo.GetReadReceiptsByMessageIDs(msgIDs)
	receiptUserIDs := make([]string, 0)
	for _, receipts := range receiptsMap {
		for _, r := range receipts {
			receiptUserIDs = append(receiptUserIDs, r.UserID)
		}
	}
	receiptUserIDs = uniqueStrings(receiptUserIDs)
	receiptUsers, _ := s.userRepo.FindByIDs(receiptUserIDs)

	// Build responses
	responses := make([]*domain.MessageResponse, 0, len(messages))
	for _, msg := range messages {
		resp := s.buildSingleResponse(msg, userMap, replyMsgs, replySenders, reactionsMap, reactionUsers, receiptsMap, receiptUsers)
		if resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses, nil
}

func (s *messageService) buildSingleResponse(
	msg *domain.Message,
	userMap map[string]*domain.User,
	replyMsgs map[string]*domain.Message,
	replySenders map[string]*domain.User,
	reactionsMap map[string][]*domain.Reaction,
	reactionUsers map[string]*domain.User,
	receiptsMap map[string][]*domain.ReadReceipt,
	receiptUsers map[string]*domain.User,
) *domain.MessageResponse {
	sender, ok := userMap[msg.SenderID]
	if !ok {
		return nil
	}

	edited := msg.UpdatedAt.Sub(msg.CreatedAt) > time.Second

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
		Edited:    edited,
		Deleted:   msg.DeletedAt != nil,
	}

	if msg.FilePath != "" {
		resp.FileURL = "/uploads/" + msg.FilePath
	}

	if msg.ReplyToID != nil && *msg.ReplyToID != "" {
		if replyMsg, ok := replyMsgs[*msg.ReplyToID]; ok {
			if replySender, ok := replySenders[replyMsg.SenderID]; ok {
				resp.ReplyTo = &domain.MessageResponse{
					ID:      replyMsg.ID,
					Content: replyMsg.Content,
					Type:    replyMsg.Type,
					Sender:  replySender.ToResponse(),
				}
			}
		}
	}

	if msg.ForwardFrom != nil && *msg.ForwardFrom != "" {
		if fwdUser, ok := userMap[*msg.ForwardFrom]; ok {
			resp.ForwardFrom = fwdUser.ToResponse()
		}
	}

	if reactions, ok := reactionsMap[msg.ID]; ok {
		for _, r := range reactions {
			if u, ok := reactionUsers[r.UserID]; ok {
				r.User = u.ToResponse()
			}
		}
		resp.Reactions = reactions
	}

	if receipts, ok := receiptsMap[msg.ID]; ok {
		for _, r := range receipts {
			if u, ok := receiptUsers[r.UserID]; ok {
				resp.ReadBy = append(resp.ReadBy, u.ToResponse())
			}
		}
	}

	if resp.Deleted {
		resp.Content = ""
	}

	return resp
}
