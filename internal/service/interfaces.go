package service

import "ChatServerGolang/internal/domain"

type AuthService interface {
	Register(req *domain.RegisterRequest) (*domain.AuthResponse, error)
	Login(req *domain.LoginRequest) (*domain.AuthResponse, error)
	RefreshToken(userID string) (string, error)
	ChangePassword(userID, oldPassword, newPassword string) error
	ValidateToken(tokenString string) (string, error)
}

type UserService interface {
	GetProfile(userID string) (*domain.UserResponse, error)
	UpdateProfile(userID string, req *domain.UpdateProfileRequest) (*domain.UserResponse, error)
	SearchUsers(query string, limit, offset int) ([]*domain.UserResponse, error)
	UpdatePushToken(userID, token, provider string) error
	GetUserByID(id string) (*domain.User, error)
	GetUsersByIDs(ids []string) (map[string]*domain.UserResponse, error)
	UpdateStatus(userID, status string) (*domain.UserResponse, error)
	GetByUsername(username string) (*domain.UserResponse, error)
	DeleteAccount(userID string) error
	BlockUser(userID, blockedID string) error
	UnblockUser(userID, blockedID string) error
	GetBlockedUsers(userID string) ([]*domain.UserResponse, error)
	IsBlocked(userID, blockedID string) (bool, error)
	SetNotificationMuted(userID, chatID string, muted bool) error
	IsNotificationMuted(userID, chatID string) (bool, error)
	GetAccountSetting(userID string) (*domain.AccountSetting, error)
	UpdateAccountSetting(userID string, req *domain.UpdateAccountSettingRequest) (*domain.AccountSetting, error)
}

type ChatService interface {
	CreateChat(userID string, req *domain.CreateChatRequest) (*domain.ChatResponse, error)
	GetChat(chatID, userID string) (*domain.ChatResponse, error)
	ListChats(userID string) ([]*domain.ChatResponse, error)
	DeleteChat(chatID, userID string) error
	AddParticipant(chatID, userID, requesterID string) error
	RemoveParticipant(chatID, userID, requesterID string) error
	MarkAsRead(chatID, userID string) error
	GetUnreadCount(chatID, userID string) (int, error)
	SetRole(chatID, targetUserID, requesterID, role string) error
	LeaveGroup(chatID, userID string) error
	UpdateGroup(chatID, userID string, req *domain.UpdateGroupRequest) error
	GetParticipants(chatID string) ([]*domain.ChatParticipant, error)
	HideChat(chatID, userID string) error
}

type MessageService interface {
	SendMessage(chatID, senderID string, req *domain.SendMessageRequest) (*domain.MessageResponse, error)
	SendFileMessage(chatID, senderID, fileName, filePath string, fileSize int64, replyToID *string) (*domain.MessageResponse, error)
	GetMessages(chatID, userID string, limit, offset int) ([]*domain.MessageResponse, error)
	SearchMessages(chatID, userID, query string, limit, offset int) ([]*domain.MessageResponse, error)
	ResendMessage(chatID, userID, msgID string) (*domain.MessageResponse, error)
	EditMessage(msgID, userID string, req *domain.EditMessageRequest) (*domain.MessageResponse, error)
	GetMessageByID(msgID, userID string) (*domain.MessageResponse, error)
	DeleteMessage(msgID, userID string) error
	AddReaction(msgID, userID, emoji string) (*domain.MessageResponse, error)
	RemoveReaction(msgID, userID, emoji string) (*domain.MessageResponse, error)
	TogglePin(msgID, userID string, pin bool) (*domain.MessageResponse, error)
	GetPinnedMessages(chatID, userID string) ([]*domain.MessageResponse, error)
	MarkMessageRead(msgID, userID string) error
}

type CallService interface {
	InitiateCall(chatID, callerID string, callType domain.CallType) (*domain.Call, error)
	AcceptCall(callID, userID string) error
	EndCall(callID, userID string) error
	MissCall(callID string) error
	RejectCall(callID, userID string) error
	GetCallByID(callID string) (*domain.Call, error)
	GetCallHistory(chatID, userID string) ([]*domain.CallResponse, error)
}

type PushService interface {
	SendMessageNotification(senderID, chatID, msgID, msgContent, msgType string)
	SendCallNotification(callerID, chatID, callID, callType string)
	SendTestPush(userID, title, body string)
}

type ContactService interface {
	SyncContacts(userID string, req *domain.SyncContactsRequest) error
	GetContacts(userID string) ([]*domain.ContactResponse, error)
	SearchByPhone(userID, query string) ([]*domain.ContactResponse, error)
	FindRegisteredByPhone(userID string) ([]*domain.UserResponse, error)
}
