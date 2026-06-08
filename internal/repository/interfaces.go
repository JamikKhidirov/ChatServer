package repository

import (
	"time"

	"ChatServerGolang/internal/domain"
)

type UserRepository interface {
	Create(user *domain.User) error
	FindByID(id string) (*domain.User, error)
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
	IsParticipant(chatID, userID string) (bool, error)
	GetPrivateChat(user1ID, user2ID string) (*domain.Chat, error)
	SetRole(chatID, userID, role string) error
	UpdateLastRead(chatID, userID string) error
	GetUnreadCount(chatID, userID string) (int, error)
	SetNotificationMuted(userID, chatID string, muted bool) error
	IsNotificationMuted(userID, chatID string) (bool, error)
	HideChat(userID, chatID string) error
	IsHidden(userID, chatID string) (bool, error)
	FindByUserIDExcludeHidden(userID string) ([]*domain.Chat, error)
}

type MessageRepository interface {
	Create(msg *domain.Message) error
	FindByID(id string) (*domain.Message, error)
	FindByChatID(chatID string, limit, offset int) ([]*domain.Message, error)
	Search(chatID, query string, limit, offset int) ([]*domain.Message, error)
	Update(msg *domain.Message) error
	SoftDelete(id string) error
	GetLastMessage(chatID string) (*domain.Message, error)
	TogglePin(msgID string, pinned bool) error
	GetPinned(chatID string) ([]*domain.Message, error)
	AddReaction(msgID, userID, emoji string) error
	RemoveReaction(msgID, userID, emoji string) error
	GetReactions(msgID string) ([]*domain.Reaction, error)
	AddReadReceipt(msgID, userID string) error
	GetReadReceipts(msgID string) ([]*domain.ReadReceipt, error)
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
