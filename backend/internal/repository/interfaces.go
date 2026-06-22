package repository

import (
	"time"

	authdomain "ChatServerGolang/backend/internal/domain/auth"
	userdomain "ChatServerGolang/backend/internal/domain/user"
	chatdomain "ChatServerGolang/backend/internal/domain/chat"
	channeldomain "ChatServerGolang/backend/internal/domain/channel"
	messagedomain "ChatServerGolang/backend/internal/domain/message"
	polldomain "ChatServerGolang/backend/internal/domain/poll"
	stickerdomain "ChatServerGolang/backend/internal/domain/sticker"
	contactdomain "ChatServerGolang/backend/internal/domain/contact"
	botdomain "ChatServerGolang/backend/internal/domain/bot"
	calldomain "ChatServerGolang/backend/internal/domain/call"
	storydomain "ChatServerGolang/backend/internal/domain/story"
	draftdomain "ChatServerGolang/backend/internal/domain/draft"
	sessiondomain "ChatServerGolang/backend/internal/domain/session"
	verificationdomain "ChatServerGolang/backend/internal/domain/verification"
	emojidomain "ChatServerGolang/backend/internal/domain/emoji"
	voicechatdomain "ChatServerGolang/backend/internal/domain/voicechat"
)

type InviteLinkRepository interface {
	Create(link *chatdomain.InviteLink) error
	FindByCode(code string) (*chatdomain.InviteLink, error)
	FindByChatID(chatID string) ([]*chatdomain.InviteLink, error)
	IncrementUsage(id string) error
	Deactivate(id string) error
	Delete(id string) error
}

type ChatFolderRepository interface {
	Create(folder *chatdomain.ChatFolder) error
	FindByUserID(userID string) ([]*chatdomain.ChatFolder, error)
	FindByID(id string) (*chatdomain.ChatFolder, error)
	Update(folder *chatdomain.ChatFolder) error
	Delete(id string) error
	AddChatToFolder(folderID, chatID string) error
	RemoveChatFromFolder(folderID, chatID string) error
	GetChatIDsByFolder(folderID string) ([]string, error)
	SetChatsForFolder(folderID string, chatIDs []string) error
}

type UserRepository interface {
	Create(user *userdomain.User) error
	FindByID(id string) (*userdomain.User, error)
	FindByIDs(ids []string) (map[string]*userdomain.User, error)
	FindByEmail(email string) (*userdomain.User, error)
	FindByUsername(username string) (*userdomain.User, error)
	Search(query string, limit, offset int) ([]*userdomain.User, error)
	SearchTotalCount(query string) (int, error)
	Update(user *userdomain.User) error
	UpdatePushToken(userID, token, provider string) error
	SetOnline(userID string, online bool) error
	SoftDelete(userID string) error
	UpdatePassword(userID, hash string) error
	GetParticipantsInChat(chatID string) ([]*userdomain.User, error)
	BlockUser(userID, blockedID string) error
	UnblockUser(userID, blockedID string) error
	GetBlockedUsers(userID string) ([]*userdomain.User, error)
	IsBlocked(userID, blockedID string) (bool, error)
	FindByIDIncludeDeleted(id string) (*userdomain.User, error)
	FindByPhone(phone string) (*userdomain.User, error)
}

type ChatRepository interface {
	Create(chat *chatdomain.Chat) error
	FindByID(id string) (*chatdomain.Chat, error)
	FindByUserID(userID string) ([]*chatdomain.Chat, error)
	Update(chat *chatdomain.Chat) error
	Delete(id string) error
	AddParticipant(chatID, userID, role string) error
	RemoveParticipant(chatID, userID string) error
	GetParticipants(chatID string) ([]*chatdomain.ChatParticipant, error)
	GetParticipantsByChatIDs(chatIDs []string) (map[string][]*chatdomain.ChatParticipant, error)
	IsParticipant(chatID, userID string) (bool, error)
	GetPrivateChat(user1ID, user2ID string) (*chatdomain.Chat, error)
	SetRole(chatID, userID, role string) error
	UpdateLastRead(chatID, userID string) error
	GetUnreadCount(chatID, userID string) (int, error)
	GetUnreadCounts(userID string, chatIDs []string) (map[string]int, error)
	SetNotificationMuted(userID, chatID string, muted bool) error
	IsNotificationMuted(userID, chatID string) (bool, error)
	SetSlowMode(chatID string, seconds int) error
	HideChat(userID, chatID string) error
	IsHidden(userID, chatID string) (bool, error)
	FindByUserIDExcludeHidden(userID string) ([]*chatdomain.Chat, error)
	SearchByName(userID, query string) ([]*chatdomain.Chat, error)
	PinChat(userID, chatID string) error
	UnpinChat(userID, chatID string) error
	GetPinnedChatIDs(userID string) ([]string, error)
	ArchiveChat(userID, chatID string) error
	UnarchiveChat(userID, chatID string) error
	IsArchived(userID, chatID string) (bool, error)
	FindByUserIDArchived(userID string) ([]*chatdomain.Chat, error)
}

type MessageRepository interface {
	Create(msg *messagedomain.Message) error
	FindByID(id string) (*messagedomain.Message, error)
	FindByIDs(ids []string) (map[string]*messagedomain.Message, error)
	FindByChatID(chatID string, limit, offset int) ([]*messagedomain.Message, error)
	CountByChatID(chatID string) (int, error)
	CountChatMedia(chatID, mediaType string) (int, error)
	Search(chatID, query string, limit, offset int) ([]*messagedomain.Message, error)
	SearchByUser(userID, query string, limit, offset int) ([]*messagedomain.Message, error)
	Update(msg *messagedomain.Message) error
	SoftDelete(id string) error
	GetLastMessagesByChatIDs(chatIDs []string) (map[string]*messagedomain.Message, error)
	GetLastMessage(chatID string) (*messagedomain.Message, error)
	TogglePin(msgID string, pinned bool) error
	GetPinned(chatID string) ([]*messagedomain.Message, error)
	AddReaction(msgID, userID, emoji string) error
	RemoveReaction(msgID, userID, emoji string) error
	GetReactions(msgID string) ([]*messagedomain.Reaction, error)
	GetReactionsByMessageIDs(ids []string) (map[string][]*messagedomain.Reaction, error)
	AddReadReceipt(msgID, userID string) error
	GetReadReceipts(msgID string) ([]*messagedomain.ReadReceipt, error)
	GetReadReceiptsByMessageIDs(ids []string) (map[string][]*messagedomain.ReadReceipt, error)
	StarMessage(userID, messageID, chatID string) error
	UnstarMessage(userID, messageID string) error
	GetStarredMessages(userID string) ([]*chatdomain.StarredMessage, error)
	DeleteMessageForMe(userID, messageID string) error
	FindDeletedForMe(userID string, messageIDs []string) (map[string]bool, error)
	SaveMention(messageID, userID, username string) error
	GetMentionsByMessageID(messageID string) ([]*messagedomain.Mention, error)
	FindMediaByChatID(chatID string, mediaType string, limit, offset int) ([]*messagedomain.Message, error)
	SetSelfDestruct(msgID, chatID string, deleteAt time.Time) error
	GetExpiredSelfDestruct() ([]messagedomain.MessageSelfDestruct, error)
	DeleteSelfDestructByMessageID(messageID string) error
}

type PollRepository interface {
	Create(poll *polldomain.Poll) error
	FindByID(id string) (*polldomain.Poll, error)
	FindByChatID(chatID string) ([]*polldomain.Poll, error)
	Update(poll *polldomain.Poll) error
	AddVote(vote *polldomain.PollVote) error
	HasVoted(pollID, userID string) (bool, error)
	GetVoteCount(pollID string, optionIndex int) (int, error)
	GetUserVote(pollID, userID string) (*polldomain.PollVote, error)
	GetTotalVotes(pollID string) (int, error)
	GetAllVotes(pollID string) ([]*polldomain.PollVote, error)
}

type StickerRepository interface {
	CreatePack(pack *stickerdomain.StickerPack) error
	GetPackByID(id string) (*stickerdomain.StickerPack, error)
	GetPacksByUserID(userID string) ([]*stickerdomain.StickerPack, error)
	ListPacks() ([]*stickerdomain.StickerPack, error)
	AddSticker(sticker *stickerdomain.Sticker) error
	GetStickersByPackID(packID string) ([]*stickerdomain.Sticker, error)
	DeletePack(id string) error
	DeleteSticker(id string) error
	AddToUserLibrary(userID, stickerID string) error
	GetUserLibrary(userID string) ([]*stickerdomain.Sticker, error)
}

type DraftRepository interface {
	Save(draft *draftdomain.Draft) error
	FindByUserAndChat(userID, chatID string) (*draftdomain.Draft, error)
	Delete(id string) error
	DeleteByUserAndChat(userID, chatID string) error
}

type ScheduledMessageRepository interface {
	Create(msg *draftdomain.ScheduledMessage) error
	FindPending() ([]*draftdomain.ScheduledMessage, error)
	FindByUserID(userID string) ([]*draftdomain.ScheduledMessage, error)
	MarkAsSent(id string) error
	Delete(id string) error
}

type SessionRepository interface {
	Create(session *sessiondomain.Session) error
	FindByID(id string) (*sessiondomain.Session, error)
	FindByUserID(userID string) ([]*sessiondomain.Session, error)
	UpdateLastActive(id string) error
	Delete(id string) error
	DeleteByUserID(userID string) error
}

type BotRepository interface {
	Create(bot *botdomain.Bot) error
	FindByID(id string) (*botdomain.Bot, error)
	FindByOwnerID(ownerID string) ([]*botdomain.Bot, error)
	Update(bot *botdomain.Bot) error
	RegenerateToken(id, token string) error
	Delete(id string) error
	FindByToken(token string) (*botdomain.Bot, error)
}

type SavedGifRepository interface {
	Save(userID, gifURL string) error
	FindByUserID(userID string) ([]string, error)
	Delete(userID, gifURL string) error
}

type CallRepository interface {
	Create(call *calldomain.Call) error
	FindByID(id string) (*calldomain.Call, error)
	FindActiveByUser(userID string) (*calldomain.Call, error)
	FindByChatAndUser(chatID, userID string) ([]*calldomain.Call, error)
	UpdateStatus(id string, status calldomain.CallStatus) error
}

type AccountSettingRepository interface {
	GetByUserID(userID string) (*userdomain.AccountSetting, error)
	Upsert(setting *userdomain.AccountSetting) error
}

type ContactRepository interface {
	SyncContacts(userID string, contacts []contactdomain.ContactInput) error
	GetContacts(userID string) ([]*contactdomain.ContactResponse, error)
	SearchByPhone(userID, phoneQuery string) ([]*contactdomain.ContactResponse, error)
	FindRegisteredByPhone(phones []string) ([]*userdomain.UserResponse, error)
	UpdateContactPhoto(userID, phone, photoURL string) error
}

type VerificationRepository interface {
	CreateEmail(ver *verificationdomain.EmailVerification) error
	FindEmailByUserID(userID string) (*verificationdomain.EmailVerification, error)
	VerifyEmail(id string) error
	CreatePhone(ver *verificationdomain.PhoneVerification) error
	FindPhoneByUserID(userID string) (*verificationdomain.PhoneVerification, error)
	VerifyPhone(id string) error
	CreateEmailLoginCode(code *authdomain.EmailLoginCode) error
	FindEmailLoginCode(email string) (*authdomain.EmailLoginCode, error)
	VerifyEmailLoginCode(id string) error
	CreatePhoneLoginCode(code *authdomain.PhoneLoginCode) error
	FindPhoneLoginCode(phone string) (*authdomain.PhoneLoginCode, error)
	VerifyPhoneLoginCode(id string) error
}

type StoryRepository interface {
	Create(story *storydomain.Story) error
	FindByID(id string) (*storydomain.Story, error)
	FindActiveByUserID(userID string) ([]*storydomain.Story, error)
	FindActiveByFollowing(userIDs []string) ([]*storydomain.Story, error)
	MarkExpired() error
	Delete(id string) error
	AddView(storyID, userID string) error
	GetViewCount(storyID string) (int64, error)
	GetViews(storyID string) ([]*storydomain.StoryView, error)
	HasViewed(storyID, userID string) (bool, error)
}

type GroupCallRepository interface {
	Create(call *calldomain.GroupCall) error
	FindByID(id string) (*calldomain.GroupCall, error)
	FindActiveByChatID(chatID string) ([]*calldomain.GroupCall, error)
	FindActiveByUserID(userID string) (*calldomain.GroupCall, error)
	UpdateStatus(id string, status calldomain.CallStatus) error
	AddParticipant(callID, userID string) error
	RemoveParticipant(callID, userID string) error
	UpdateParticipantMute(callID, userID string, audioMuted, videoMuted bool) error
	GetParticipants(callID string) ([]*calldomain.GroupCallParticipant, error)
	FindByChatAndUser(chatID, userID string) ([]*calldomain.GroupCall, error)
}

type ChannelSubscriberRepository interface {
	Subscribe(channelID, userID string, role string) error
	Unsubscribe(channelID, userID string) error
	IsSubscribed(channelID, userID string) (bool, error)
	GetSubscribers(channelID string) ([]*channeldomain.ChannelSubscriber, error)
	GetSubscribedChannels(userID string) ([]string, error)
	SetRole(channelID, userID, role string) error
	GetRole(channelID, userID string) (string, error)
}

type SavedMessageRepository interface {
	Save(msg *chatdomain.SavedMessage) error
	FindByUserID(userID string, limit, offset int) ([]*chatdomain.SavedMessage, error)
	CountByUserID(userID string) (int, error)
	Delete(id, userID string) error
	Exists(userID, messageID string) (bool, error)
}

type CustomEmojiRepository interface {
	Create(emoji *emojidomain.CustomEmoji) error
	FindByID(id string) (*emojidomain.CustomEmoji, error)
	FindByUserID(userID string) ([]*emojidomain.CustomEmoji, error)
	FindAll() ([]*emojidomain.CustomEmoji, error)
	Delete(id, userID string) error
}

type VoiceChatRepository interface {
	Create(vc *voicechatdomain.VoiceChat) error
	FindByID(id string) (*voicechatdomain.VoiceChat, error)
	FindActiveByChatID(chatID string) ([]*voicechatdomain.VoiceChat, error)
	FindByChatID(chatID string) ([]*voicechatdomain.VoiceChat, error)
	UpdateStatus(id string, status voicechatdomain.VoiceChatStatus) error
	AddParticipant(vcID, userID string) error
	RemoveParticipant(vcID, userID string) error
	IsParticipant(vcID, userID string) (bool, error)
	GetParticipants(vcID string) ([]*voicechatdomain.VoiceChatParticipant, error)
	GetParticipantCount(vcID string) (int, error)
	SetParticipantMuted(vcID, userID string, muted bool) error
}

func BoolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func ParseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}

type Scanner interface {
	Scan(dest ...interface{}) error
}

type MessageScanner interface {
	Scan(dest ...interface{}) error
}
