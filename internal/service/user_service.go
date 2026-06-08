package service

import (
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"
)

type userService struct {
	userRepo repository.UserRepository
	chatRepo repository.ChatRepository
	accRepo  repository.AccountSettingRepository
}

func NewUserService(userRepo repository.UserRepository, chatRepo repository.ChatRepository, accRepo repository.AccountSettingRepository) UserService {
	return &userService{userRepo: userRepo, chatRepo: chatRepo, accRepo: accRepo}
}

func (s *userService) GetProfile(userID string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

func (s *userService) UpdateProfile(userID string, req *domain.UpdateProfileRequest) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	if req.DisplayName != "" {
		user.DisplayName = req.DisplayName
	}
	if req.Bio != "" {
		user.Bio = req.Bio
	}
	if req.Phone != "" {
		user.Phone = req.Phone
	}
	if req.Gender != "" {
		user.Gender = req.Gender
	}
	if req.DateOfBirth != "" {
		user.DateOfBirth = req.DateOfBirth
	}
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *userService) SearchUsers(query string, limit, offset int) ([]*domain.UserResponse, error) {
	if query == "" {
		return nil, errors.New("search query is required")
	}
	users, err := s.userRepo.Search(query, limit, offset)
	if err != nil {
		return nil, err
	}
	responses := make([]*domain.UserResponse, 0)
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}
	return responses, nil
}

func (s *userService) UpdatePushToken(userID, token, provider string) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.UpdatePushToken(userID, token, provider)
}

func (s *userService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *userService) GetUsersByIDs(ids []string) (map[string]*domain.UserResponse, error) {
	result := make(map[string]*domain.UserResponse)
	for _, id := range ids {
		u, err := s.userRepo.FindByID(id)
		if err != nil {
			continue
		}
		result[id] = u.ToResponse()
	}
	return result, nil
}

func (s *userService) UpdateStatus(userID, status string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	user.Status = status
	user.UpdatedAt = time.Now()
	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}
	return user.ToResponse(), nil
}

func (s *userService) GetByUsername(username string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

func (s *userService) DeleteAccount(userID string) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.SoftDelete(userID)
}

func (s *userService) BlockUser(userID, blockedID string) error {
	if userID == blockedID {
		return errors.New("cannot block yourself")
	}
	_, err := s.userRepo.FindByID(blockedID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.BlockUser(userID, blockedID)
}

func (s *userService) UnblockUser(userID, blockedID string) error {
	return s.userRepo.UnblockUser(userID, blockedID)
}

func (s *userService) GetBlockedUsers(userID string) ([]*domain.UserResponse, error) {
	users, err := s.userRepo.GetBlockedUsers(userID)
	if err != nil {
		return nil, err
	}
	responses := make([]*domain.UserResponse, 0)
	for _, u := range users {
		responses = append(responses, u.ToResponse())
	}
	return responses, nil
}

func (s *userService) IsBlocked(userID, blockedID string) (bool, error) {
	blocked, err := s.userRepo.IsBlocked(blockedID, userID)
	if err != nil {
		return false, err
	}
	if blocked {
		return true, nil
	}
	return s.userRepo.IsBlocked(userID, blockedID)
}

func (s *userService) SetNotificationMuted(userID, chatID string, muted bool) error {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	return s.chatRepo.SetNotificationMuted(userID, chatID, muted)
}

func (s *userService) IsNotificationMuted(userID, chatID string) (bool, error) {
	return s.chatRepo.IsNotificationMuted(userID, chatID)
}

func (s *userService) GetAccountSetting(userID string) (*domain.AccountSetting, error) {
	setting, err := s.accRepo.GetByUserID(userID)
	if err != nil {
		return &domain.AccountSetting{
			UserID:        userID,
			Language:      "en",
			Theme:         "light",
			Notifications: true,
			SoundEnabled:  true,
			LastSeenMode:  "everyone",
		}, nil
	}
	return setting, nil
}

func (s *userService) UpdateAccountSetting(userID string, req *domain.UpdateAccountSettingRequest) (*domain.AccountSetting, error) {
	setting, err := s.accRepo.GetByUserID(userID)
	if err != nil {
		setting = &domain.AccountSetting{
			UserID:        userID,
			Language:      "en",
			Theme:         "light",
			Notifications: true,
			SoundEnabled:  true,
			LastSeenMode:  "everyone",
		}
	}

	if req.Language != nil {
		setting.Language = *req.Language
	}
	if req.Theme != nil {
		setting.Theme = *req.Theme
	}
	if req.Notifications != nil {
		setting.Notifications = *req.Notifications
	}
	if req.SoundEnabled != nil {
		setting.SoundEnabled = *req.SoundEnabled
	}
	if req.LastSeenMode != nil {
		setting.LastSeenMode = *req.LastSeenMode
	}

	if err := s.accRepo.Upsert(setting); err != nil {
		return nil, err
	}

	return setting, nil
}
