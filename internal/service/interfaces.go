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
	SearchChats(userID, query string) ([]*domain.ChatResponse, error)
	ListArchivedChats(userID string) ([]*domain.ChatResponse, error)
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
	PinChat(chatID, userID string) error
	UnpinChat(chatID, userID string) error
	ArchiveChat(chatID, userID string) error
	UnarchiveChat(chatID, userID string) error
	TransferOwnership(chatID, fromUserID, toUserID string) error
}

type MessageService interface {
	SendMessage(chatID, senderID string, req *domain.SendMessageRequest) (*domain.MessageResponse, error)
	SendFileMessage(chatID, senderID, fileName, filePath string, fileSize int64, replyToID *string) (*domain.MessageResponse, error)
	GetMessages(chatID, userID string, limit, offset int) ([]*domain.MessageResponse, error)
	SearchMessages(chatID, userID, query string, limit, offset int) ([]*domain.MessageResponse, error)
	SearchAllMessages(userID, query string, limit, offset int) ([]*domain.MessageResponse, error)
	ForwardMessage(msgID, fromChatID, toChatID, userID string) (*domain.MessageResponse, error)
	ResendMessage(chatID, userID, msgID string) (*domain.MessageResponse, error)
	EditMessage(msgID, userID string, req *domain.EditMessageRequest) (*domain.MessageResponse, error)
	GetMessageByID(msgID, userID string) (*domain.MessageResponse, error)
	DeleteMessage(msgID, userID string) error
	DeleteMessageForMe(msgID, userID string) error
	AddReaction(msgID, userID, emoji string) (*domain.MessageResponse, error)
	RemoveReaction(msgID, userID, emoji string) (*domain.MessageResponse, error)
	TogglePin(msgID, userID string, pin bool) (*domain.MessageResponse, error)
	GetPinnedMessages(chatID, userID string) ([]*domain.MessageResponse, error)
	MarkMessageRead(msgID, userID string) error
	StarMessage(msgID, userID string) (*domain.MessageResponse, error)
	UnstarMessage(msgID, userID string) error
	GetStarredMessages(userID string) ([]*domain.StarredMessageResponse, error)
	ExportChat(chatID, userID string) ([]*domain.MessageResponse, error)
	GetChatMedia(chatID, userID, mediaType string, limit, offset int) ([]*domain.MessageResponse, error)
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

type PollService interface {
	CreatePoll(userID string, req *domain.CreatePollRequest) (*domain.PollWithResults, error)
	GetPollsByChatID(chatID, userID string) ([]*domain.PollWithResults, error)
	Vote(pollID, userID string, optionIndex int) error
	ClosePoll(pollID, userID string) error
}

type StickerService interface {
	CreatePack(userID string, req *domain.CreateStickerPackRequest) (*domain.StickerPackWithStickers, error)
	GetPacks() ([]*domain.StickerPackWithStickers, error)
	GetMyPacks(userID string) ([]*domain.StickerPackWithStickers, error)
	GetPackByID(id string) (*domain.StickerPackWithStickers, error)
	AddSticker(packID, userID string, req *domain.AddStickerRequest) (*domain.Sticker, error)
	DeletePack(id, userID string) error
	DeleteSticker(id, userID string) error
	AddToLibrary(userID, stickerID string) error
	GetLibrary(userID string) ([]*domain.Sticker, error)
}

type DraftService interface {
	SaveDraft(userID string, req *domain.SaveDraftRequest) (*domain.Draft, error)
	GetDraft(userID, chatID string) (*domain.Draft, error)
	DeleteDraft(userID, draftID string) error
}

type ScheduledMessageService interface {
	Schedule(userID string, req *domain.ScheduleMessageRequest) (*domain.ScheduledMessage, error)
	GetScheduled(userID string) ([]*domain.ScheduledMessage, error)
	CancelScheduled(id, userID string) error
	SchedulerProcess()
}

type SessionService interface {
	CreateSession(userID, deviceName, ipAddress string) (*domain.Session, error)
	GetSessions(userID string) ([]*domain.Session, error)
	DeleteSession(sessionID, userID string) error
	DeleteAllSessions(userID string) error
	UpdateLastActive(sessionID string) error
}

type BotService interface {
	CreateBot(userID string, req *domain.CreateBotRequest) (*domain.Bot, error)
	GetMyBots(userID string) ([]*domain.Bot, error)
	UpdateBot(botID, userID string, req *domain.UpdateBotRequest) error
	DeleteBot(botID, userID string) error
	RegenerateToken(botID, userID string) error
	ValidateBotToken(token string) (string, error)
}

type SavedGifService interface {
	SaveGif(userID, gifURL string) error
	GetSavedGifs(userID string) ([]string, error)
	DeleteGif(userID, gifURL string) error
}

type ThemeService interface {
	GetTheme(userID string) (string, error)
	SetTheme(userID, theme string) error
}
