package systemmsgservice

import (
	"encoding/json"
	"time"

	"ChatServerGolang/internal/domain/message"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/ws"

	"github.com/google/uuid"
)

type SystemMessageService interface {
	NotifyUserJoined(chatID, userID, userName string)
	NotifyUserLeft(chatID, userID, userName string)
	NotifyUserRemoved(chatID, userID, userName, byUserID, byUserName string)
	NotifyUserAdded(chatID, userID, userName, byUserID, byUserName string)
	NotifyRoleChanged(chatID, userID, userName, newRole string)
	NotifyChatRenamed(chatID, userID, userName, newName string)
	NotifyChatPhotoChanged(chatID, userID, userName string)
	NotifyMessagePinned(chatID, userID, userName string)
	NotifyMessageUnpinned(chatID, userID, userName string)
	NotifyChatCreated(chatID, userID, userName string)
}

type systemMessageService struct {
	messageRepo repository.MessageRepository
	chatRepo    repository.ChatRepository
	hub         *ws.Hub
}

func NewSystemMessageService(messageRepo repository.MessageRepository, chatRepo repository.ChatRepository, hub *ws.Hub) SystemMessageService {
	return &systemMessageService{
		messageRepo: messageRepo,
		chatRepo:    chatRepo,
		hub:         hub,
	}
}

func (s *systemMessageService) createSystemMessage(chatID, senderID string, content messagedomain.SystemMessageContent) {
	now := time.Now()
	data, _ := json.Marshal(content)
	msg := &messagedomain.Message{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		SenderID:  senderID,
		Content:   string(data),
		Type:      messagedomain.MessageSystem,
		CreatedAt: now,
		UpdatedAt: now,
	}
	if err := s.messageRepo.Create(msg); err != nil {
		return
	}

	participants, _ := s.chatRepo.GetParticipants(chatID)
	userIDs := make([]string, 0, len(participants))
	for _, p := range participants {
		userIDs = append(userIDs, p.UserID)
	}

	s.hub.BroadcastToChat(userIDs, ws.WSOutgoingMessage{
		Type:    "message:new",
		Payload: msg,
	})
}

func (s *systemMessageService) NotifyUserJoined(chatID, userID, userName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysUserJoined,
		UserID:   userID,
		UserName: userName,
	})
}

func (s *systemMessageService) NotifyUserLeft(chatID, userID, userName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysUserLeft,
		UserID:   userID,
		UserName: userName,
	})
}

func (s *systemMessageService) NotifyUserRemoved(chatID, userID, userName, byUserID, byUserName string) {
	s.createSystemMessage(chatID, byUserID, messagedomain.SystemMessageContent{
		Action:     messagedomain.SysUserRemoved,
		UserID:     byUserID,
		UserName:   byUserName,
		TargetID:   userID,
		TargetName: userName,
	})
}

func (s *systemMessageService) NotifyUserAdded(chatID, userID, userName, byUserID, byUserName string) {
	s.createSystemMessage(chatID, byUserID, messagedomain.SystemMessageContent{
		Action:     messagedomain.SysUserAdded,
		UserID:     byUserID,
		UserName:   byUserName,
		TargetID:   userID,
		TargetName: userName,
	})
}

func (s *systemMessageService) NotifyRoleChanged(chatID, userID, userName, newRole string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysRoleChanged,
		UserID:   userID,
		UserName: userName,
		Extra:    newRole,
	})
}

func (s *systemMessageService) NotifyChatRenamed(chatID, userID, userName, newName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysChatRenamed,
		UserID:   userID,
		UserName: userName,
		Extra:    newName,
	})
}

func (s *systemMessageService) NotifyChatPhotoChanged(chatID, userID, userName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysChatPhotoChanged,
		UserID:   userID,
		UserName: userName,
	})
}

func (s *systemMessageService) NotifyMessagePinned(chatID, userID, userName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysMessagePinned,
		UserID:   userID,
		UserName: userName,
	})
}

func (s *systemMessageService) NotifyMessageUnpinned(chatID, userID, userName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysMessageUnpinned,
		UserID:   userID,
		UserName: userName,
	})
}

func (s *systemMessageService) NotifyChatCreated(chatID, userID, userName string) {
	s.createSystemMessage(chatID, userID, messagedomain.SystemMessageContent{
		Action:   messagedomain.SysChatCreated,
		UserID:   userID,
		UserName: userName,
	})
}




