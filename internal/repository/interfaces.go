package repository

import (
	"time"

	"ChatServerGolang/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
	FindByIDs(ids []string) (map[string]*domain.User, error)
	FindByEmail(email string) (*domain.User, error)
	FindByUsername(username string) (*domain.User, error)
	Search(query string, limit, offset int) ([]*domain.User, error)
	Update(user *domain.User) error
	UpdatePushToken(userID, token, provider string) error
	SetOnline(userID string, online bool) error
	SoftDelete(userID string) error
	UpdatePassword(userID, hash string) error
	GetParticipantsInChat(chatID string) ([]*domain.User, error)
	BlockUser(userID, blockedID string) error
	UnblockUser(userID, blockedID string) error
	GetBlockedUsers(userID string) ([]*domain.User, error)
	IsBlocked(userID, blockedID string) (bool, error)
	FindByIDIncludeDeleted(id string) (*domain.User, error)
	FindByPhone(phone string) (*domain.User, error)
}

type ChatRepository interface {
	Create(chat *domain.Chat) error
	FindByID(id string) (*domain.Chat, error)
	FindByUserID(userID string) ([]*domain.Chat, error)
	Update(chat *domain.Chat) error
	Delete(id string) error
	AddParticipant(chatID, userID, role string) error
	RemoveParticipant(chatID, userID string) error
	GetParticipants(chatID string) ([]*domain.ChatParticipant, error)
	GetParticipantsByChatIDs(chatIDs []string) (map[string][]*domain.ChatParticipant, error)
	IsParticipant(chatID, userID string) (bool, error)
	GetPrivateChat(user1ID, user2ID string) (*domain.Chat, error)
	SetRole(chatID, userID, role string) error
	UpdateLastRead(chatID, userID string) error
	GetUnreadCount(chatID, userID string) (int, error)
	GetUnreadCounts(userID string, chatIDs []string) (map[string]int, error)
	SetNotificationMuted(userID, chatID string, muted bool) error
	IsNotificationMuted(userID, chatID string) (bool, error)
	HideChat(userID, chatID string) error
	IsHidden(userID, chatID string) (bool, error)
	FindByUserIDExcludeHidden(userID string) ([]*domain.Chat, error)
	SearchByName(userID, query string) ([]*domain.Chat, error)
	PinChat(userID, chatID string) error
	UnpinChat(userID, chatID string) error
	GetPinnedChatIDs(userID string) ([]string, error)
	ArchiveChat(userID, chatID string) error
	UnarchiveChat(userID, chatID string) error
	IsArchived(userID, chatID string) (bool, error)
	FindByUserIDArchived(userID string) ([]*domain.Chat, error)
}

type MessageRepository interface {
	Create(msg *domain.Message) error
	FindByID(id string) (*domain.Message, error)
	FindByIDs(ids []string) (map[string]*domain.Message, error)
	FindByChatID(chatID string, limit, offset int) ([]*domain.Message, error)
	Search(chatID, query string, limit, offset int) ([]*domain.Message, error)
	SearchByUser(userID, query string, limit, offset int) ([]*domain.Message, error)
	Update(msg *domain.Message) error
	SoftDelete(id string) error
	GetLastMessagesByChatIDs(chatIDs []string) (map[string]*domain.Message, error)
	GetLastMessage(chatID string) (*domain.Message, error)
	TogglePin(msgID string, pinned bool) error
	GetPinned(chatID string) ([]*domain.Message, error)
	AddReaction(msgID, userID, emoji string) error
	RemoveReaction(msgID, userID, emoji string) error
	GetReactions(msgID string) ([]*domain.Reaction, error)
	GetReactionsByMessageIDs(ids []string) (map[string][]*domain.Reaction, error)
	AddReadReceipt(msgID, userID string) error
	GetReadReceipts(msgID string) ([]*domain.ReadReceipt, error)
	GetReadReceiptsByMessageIDs(ids []string) (map[string][]*domain.ReadReceipt, error)
	StarMessage(userID, messageID, chatID string) error
	UnstarMessage(userID, messageID string) error
	GetStarredMessages(userID string) ([]*domain.StarredMessage, error)
	DeleteMessageForMe(userID, messageID string) error
	FindDeletedForMe(userID string, messageIDs []string) (map[string]bool, error)
	SaveMention(messageID, userID, username string) error
	GetMentionsByMessageID(messageID string) ([]*domain.Mention, error)
	FindMediaByChatID(chatID string, mediaType string, limit, offset int) ([]*domain.Message, error)
}

type PollRepository interface {
	Create(poll *domain.Poll) error
	FindByID(id string) (*domain.Poll, error)
	FindByChatID(chatID string) ([]*domain.Poll, error)
	Update(poll *domain.Poll) error
	AddVote(vote *domain.PollVote) error
	HasVoted(pollID, userID string) (bool, error)
	GetVoteCount(pollID string, optionIndex int) (int, error)
	GetUserVote(pollID, userID string) (*domain.PollVote, error)
	GetTotalVotes(pollID string) (int, error)
	GetAllVotes(pollID string) ([]*domain.PollVote, error)
}

type StickerRepository interface {
	CreatePack(pack *domain.StickerPack) error
	GetPackByID(id string) (*domain.StickerPack, error)
	GetPacksByUserID(userID string) ([]*domain.StickerPack, error)
	ListPacks() ([]*domain.StickerPack, error)
	AddSticker(sticker *domain.Sticker) error
	GetStickersByPackID(packID string) ([]*domain.Sticker, error)
	DeletePack(id string) error
	DeleteSticker(id string) error
	AddToUserLibrary(userID, stickerID string) error
	GetUserLibrary(userID string) ([]*domain.Sticker, error)
}

type DraftRepository interface {
	Save(draft *domain.Draft) error
	FindByUserAndChat(userID, chatID string) (*domain.Draft, error)
	Delete(id string) error
	DeleteByUserAndChat(userID, chatID string) error
}

type ScheduledMessageRepository interface {
	Create(msg *domain.ScheduledMessage) error
	FindPending() ([]*domain.ScheduledMessage, error)
	FindByUserID(userID string) ([]*domain.ScheduledMessage, error)
	MarkAsSent(id string) error
	Delete(id string) error
}

type SessionRepository interface {
	Create(session *domain.Session) error
	FindByID(id string) (*domain.Session, error)
	FindByUserID(userID string) ([]*domain.Session, error)
	UpdateLastActive(id string) error
	Delete(id string) error
	DeleteByUserID(userID string) error
}

type BotRepository interface {
	Create(bot *domain.Bot) error
	FindByID(id string) (*domain.Bot, error)
	FindByOwnerID(ownerID string) ([]*domain.Bot, error)
	Update(bot *domain.Bot) error
	RegenerateToken(id, token string) error
	Delete(id string) error
	FindByToken(token string) (*domain.Bot, error)
}

type SavedGifRepository interface {
	Save(userID, gifURL string) error
	FindByUserID(userID string) ([]string, error)
	Delete(userID, gifURL string) error
}

type CallRepository interface {
	Create(call *domain.Call) error
	FindByID(id string) (*domain.Call, error)
	FindActiveByUser(userID string) (*domain.Call, error)
	FindByChatAndUser(chatID, userID string) ([]*domain.Call, error)
	UpdateStatus(id string, status domain.CallStatus) error
}

type AccountSettingRepository interface {
	GetByUserID(userID string) (*domain.AccountSetting, error)
	Upsert(setting *domain.AccountSetting) error
}

type ContactRepository interface {
	SyncContacts(userID string, contacts []domain.ContactInput) error
	GetContacts(userID string) ([]*domain.ContactResponse, error)
	SearchByPhone(userID, phoneQuery string) ([]*domain.ContactResponse, error)
	FindRegisteredByPhone(phones []string) ([]*domain.UserResponse, error)
}

type VerificationRepository interface {
	CreateEmail(ver *domain.EmailVerification) error
	FindEmailByUserID(userID string) (*domain.EmailVerification, error)
	VerifyEmail(id string) error
	CreatePhone(ver *domain.PhoneVerification) error
	FindPhoneByUserID(userID string) (*domain.PhoneVerification, error)
	VerifyPhone(id string) error
	CreateEmailLoginCode(code *domain.EmailLoginCode) error
	FindEmailLoginCode(email string) (*domain.EmailLoginCode, error)
	VerifyEmailLoginCode(id string) error
	CreatePhoneLoginCode(code *domain.PhoneLoginCode) error
	FindPhoneLoginCode(phone string) (*domain.PhoneLoginCode, error)
	VerifyPhoneLoginCode(id string) error
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func parseTime(s string) time.Time {
	t, _ := time.Parse(time.RFC3339, s)
	return t
}
