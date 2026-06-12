package service

import (
	authdomain "ChatServerGolang/internal/domain/auth"
	userdomain "ChatServerGolang/internal/domain/user"
	chatdomain "ChatServerGolang/internal/domain/chat"
	messagedomain "ChatServerGolang/internal/domain/message"
	calldomain "ChatServerGolang/internal/domain/call"
	polldomain "ChatServerGolang/internal/domain/poll"
	stickerdomain "ChatServerGolang/internal/domain/sticker"
	contactdomain "ChatServerGolang/internal/domain/contact"
	botdomain "ChatServerGolang/internal/domain/bot"
	draftdomain "ChatServerGolang/internal/domain/draft"
	sessiondomain "ChatServerGolang/internal/domain/session"
)

type AuthService interface {
	Register(req *authdomain.RegisterRequest) (*authdomain.AuthResponse, error)
	RegisterAdmin(req *authdomain.AdminRegisterRequest) (*authdomain.AuthResponse, error)
	Login(req *authdomain.LoginRequest) (*authdomain.AuthResponse, error)
	RefreshToken(userID string) (string, error)
	ChangePassword(userID, oldPassword, newPassword string) error
	ValidateToken(tokenString string) (string, error)
}

type UserService interface {
	GetProfile(userID string) (*userdomain.UserResponse, error)
	UpdateProfile(userID string, req *userdomain.UpdateProfileRequest) (*userdomain.UserResponse, error)
	SearchUsers(query string, limit, offset int) ([]*userdomain.UserResponse, int, error)
	UpdatePushToken(userID, token, provider string) error
	GetUserByID(id string) (*userdomain.User, error)
	GetUsersByIDs(ids []string) (map[string]*userdomain.UserResponse, error)
	UpdateStatus(userID, status string) (*userdomain.UserResponse, error)
	GetByUsername(username string) (*userdomain.UserResponse, error)
	DeleteAccount(userID string) error
	BlockUser(userID, blockedID string) error
	UnblockUser(userID, blockedID string) error
	GetBlockedUsers(userID string) ([]*userdomain.UserResponse, error)
	IsBlocked(userID, blockedID string) (bool, error)
	SetNotificationMuted(userID, chatID string, muted bool) error
	IsNotificationMuted(userID, chatID string) (bool, error)
	GetAccountSetting(userID string) (*userdomain.AccountSetting, error)
	UpdateAccountSetting(userID string, req *userdomain.UpdateAccountSettingRequest) (*userdomain.AccountSetting, error)
}

type ChatService interface {
	CreateChat(userID string, req *chatdomain.CreateChatRequest) (*chatdomain.ChatResponse, error)
	GetChat(chatID, userID string) (*chatdomain.ChatResponse, error)
	ListChats(userID string) ([]*chatdomain.ChatResponse, error)
	SearchChats(userID, query string) ([]*chatdomain.ChatResponse, error)
	ListArchivedChats(userID string) ([]*chatdomain.ChatResponse, error)
	DeleteChat(chatID, userID string) error
	AddParticipant(chatID, userID, requesterID string) error
	RemoveParticipant(chatID, userID, requesterID string) error
	MarkAsRead(chatID, userID string) error
	GetUnreadCount(chatID, userID string) (int, error)
	SetRole(chatID, targetUserID, requesterID, role string) error
	LeaveGroup(chatID, userID string) error
	UpdateGroup(chatID, userID string, req *chatdomain.UpdateGroupRequest) error
	GetParticipants(chatID string) ([]*chatdomain.ChatParticipant, error)
	HideChat(chatID, userID string) error
	PinChat(chatID, userID string) error
	UnpinChat(chatID, userID string) error
	ArchiveChat(chatID, userID string) error
	UnarchiveChat(chatID, userID string) error
	TransferOwnership(chatID, fromUserID, toUserID string) error
	SetSlowMode(chatID, userID string, seconds int) error
}

type MessageService interface {
	SendMessage(chatID, senderID string, req *messagedomain.SendMessageRequest) (*messagedomain.MessageResponse, error)
	SendFileMessage(chatID, senderID, fileName, filePath string, fileSize int64, replyToID *string) (*messagedomain.MessageResponse, error)
	GetMessages(chatID, userID string, limit, offset int) ([]*messagedomain.MessageResponse, int, error)
	SearchMessages(chatID, userID, query string, limit, offset int) ([]*messagedomain.MessageResponse, int, error)
	SearchAllMessages(userID, query string, limit, offset int) ([]*messagedomain.MessageResponse, int, error)
	ForwardMessage(msgID, fromChatID, toChatID, userID string) (*messagedomain.MessageResponse, error)
	ResendMessage(chatID, userID, msgID string) (*messagedomain.MessageResponse, error)
	EditMessage(msgID, userID string, req *messagedomain.EditMessageRequest) (*messagedomain.MessageResponse, error)
	GetMessageByID(msgID, userID string) (*messagedomain.MessageResponse, error)
	DeleteMessage(msgID, userID string) error
	DeleteMessageForMe(msgID, userID string) error
	AddReaction(msgID, userID, emoji string) (*messagedomain.MessageResponse, error)
	RemoveReaction(msgID, userID, emoji string) (*messagedomain.MessageResponse, error)
	TogglePin(msgID, userID string, pin bool) (*messagedomain.MessageResponse, error)
	GetPinnedMessages(chatID, userID string) ([]*messagedomain.MessageResponse, error)
	MarkMessageRead(msgID, userID string) error
	StarMessage(msgID, userID string) (*messagedomain.MessageResponse, error)
	UnstarMessage(msgID, userID string) error
	GetStarredMessages(userID string) ([]*chatdomain.StarredMessageResponse, error)
	ExportChat(chatID, userID string) ([]*messagedomain.MessageResponse, error)
	GetChatMedia(chatID, userID, mediaType string, limit, offset int) ([]*messagedomain.MessageResponse, int, error)
}

type CallService interface {
	InitiateCall(chatID, callerID string, callType calldomain.CallType) (*calldomain.Call, error)
	AcceptCall(callID, userID string) error
	EndCall(callID, userID string) error
	MissCall(callID string) error
	RejectCall(callID, userID string) error
	GetCallByID(callID string) (*calldomain.Call, error)
	GetCallHistory(chatID, userID string) ([]*calldomain.CallResponse, error)
}

type PushService interface {
	SendMessageNotification(senderID, chatID, msgID, msgContent, msgType string)
	SendCallNotification(callerID, chatID, callID, callType string)
	SendTestPush(userID, title, body string)
}

type ContactService interface {
	SyncContacts(userID string, req *contactdomain.SyncContactsRequest) error
	GetContacts(userID string) ([]*contactdomain.ContactResponse, error)
	SearchByPhone(userID, query string) ([]*contactdomain.ContactResponse, error)
	FindRegisteredByPhone(userID string) ([]*userdomain.UserResponse, error)
	UpdateContactPhoto(userID, phone, photoURL string) error
}

type PollService interface {
	CreatePoll(userID string, req *polldomain.CreatePollRequest) (*polldomain.PollWithResults, error)
	GetPollsByChatID(chatID, userID string) ([]*polldomain.PollWithResults, error)
	Vote(pollID, userID string, optionIndex int) error
	ClosePoll(pollID, userID string) error
}

type StickerService interface {
	CreatePack(userID string, req *stickerdomain.CreateStickerPackRequest) (*stickerdomain.StickerPackWithStickers, error)
	GetPacks() ([]*stickerdomain.StickerPackWithStickers, error)
	GetMyPacks(userID string) ([]*stickerdomain.StickerPackWithStickers, error)
	GetPackByID(id string) (*stickerdomain.StickerPackWithStickers, error)
	AddSticker(packID, userID string, req *stickerdomain.AddStickerRequest) (*stickerdomain.Sticker, error)
	DeletePack(id, userID string) error
	DeleteSticker(id, userID string) error
	AddToLibrary(userID, stickerID string) error
	GetLibrary(userID string) ([]*stickerdomain.Sticker, error)
}

type DraftService interface {
	SaveDraft(userID string, req *draftdomain.SaveDraftRequest) (*draftdomain.Draft, error)
	GetDraft(userID, chatID string) (*draftdomain.Draft, error)
	DeleteDraft(userID, draftID string) error
}

type ScheduledMessageService interface {
	Schedule(userID string, req *draftdomain.ScheduleMessageRequest) (*draftdomain.ScheduledMessage, error)
	GetScheduled(userID string) ([]*draftdomain.ScheduledMessage, error)
	CancelScheduled(id, userID string) error
	SchedulerProcess()
}

type SessionService interface {
	CreateSession(userID, deviceName, ipAddress string) (*sessiondomain.Session, error)
	GetSessions(userID string) ([]*sessiondomain.Session, error)
	DeleteSession(sessionID, userID string) error
	DeleteAllSessions(userID string) error
	UpdateLastActive(sessionID string) error
}

type BotService interface {
	CreateBot(userID string, req *botdomain.CreateBotRequest) (*botdomain.Bot, error)
	GetMyBots(userID string) ([]*botdomain.Bot, error)
	UpdateBot(botID, userID string, req *botdomain.UpdateBotRequest) error
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

type InviteLinkService interface {
	CreateInviteLink(chatID, userID string, req *chatdomain.CreateInviteLinkRequest) (*chatdomain.InviteLink, error)
	GetInviteLinks(chatID, userID string) ([]*chatdomain.InviteLink, error)
	DeleteInviteLink(linkID, userID string) error
	JoinByInviteLink(code, userID string) error
}

type ChatFolderService interface {
	Create(userID string, req *chatdomain.CreateChatFolderRequest) (*chatdomain.ChatFolderWithChats, error)
	List(userID string) ([]*chatdomain.ChatFolderWithChats, error)
	Update(folderID, userID string, req *chatdomain.UpdateChatFolderRequest) (*chatdomain.ChatFolderWithChats, error)
	Delete(folderID, userID string) error
	GetWithChats(folderID, userID string) (*chatdomain.ChatFolderWithChats, error)
}

type VerificationService interface {
	SendEmailVerification(userID, email string) error
	VerifyEmail(userID, code string) error
	SendPhoneVerification(userID, phone string) error
	VerifyPhone(userID, code string) error
	IsEmailVerified(userID string) (bool, error)
	IsPhoneVerified(userID string) (bool, error)
	LoginSendEmailCode(email string) (string, error)
	LoginVerifyEmailCode(email, code string) (string, error)
	LoginSendPhoneCode(phone string) (string, error)
	LoginVerifyPhoneCode(phone, code string) (string, error)
}
