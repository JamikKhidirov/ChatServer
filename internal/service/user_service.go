package service

import (
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	userRepo *repository.UserRepository
	chatRepo *repository.ChatRepository
}

func NewUserService(userRepo *repository.UserRepository, chatRepo *repository.ChatRepository) *UserService {
	return &UserService{userRepo: userRepo, chatRepo: chatRepo}
}

func (s *UserService) GetProfile(userID string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

func (s *UserService) UpdateProfile(userID string, req *domain.UpdateProfileRequest) (*domain.UserResponse, error) {
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
	if req.AvatarURL != "" {
		user.AvatarURL = req.AvatarURL
	}
	user.UpdatedAt = time.Now()

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user.ToResponse(), nil
}

func (s *UserService) SearchUsers(query string, limit, offset int) ([]*domain.UserResponse, error) {
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

func (s *UserService) UpdatePushToken(userID, token, provider string) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.UpdatePushToken(userID, token, provider)
}

func (s *UserService) GetUserByID(id string) (*domain.User, error) {
	return s.userRepo.FindByID(id)
}

func (s *UserService) GetUsersByIDs(ids []string) (map[string]*domain.UserResponse, error) {
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

func (s *UserService) UpdateStatus(userID, status string) (*domain.UserResponse, error) {
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

func (s *UserService) GetByUsername(username string) (*domain.UserResponse, error) {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return nil, errors.New("user not found")
	}
	return user.ToResponse(), nil
}

func (s *UserService) GetUserByUsername(username string) (*domain.User, error) {
	return s.userRepo.FindByUsername(username)
}

func (s *UserService) CreateUser(username, email, password, displayName string) (*domain.User, error) {
	existing, _ := s.userRepo.FindByEmail(email)
	if existing != nil {
		return nil, errors.New("email already registered")
	}

	existing, _ = s.userRepo.FindByUsername(username)
	if existing != nil {
		return nil, errors.New("username already taken")
	}

	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	now := time.Now()
	user := &domain.User{
		ID:           uuid.New().String(),
		Username:     username,
		Email:        email,
		PasswordHash: string(hash),
		DisplayName:  displayName,
		Status:       "Available",
		LastSeen:     now,
		CreatedAt:    now,
		UpdatedAt:    now,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}
	return user, nil
}

// Delete account
func (s *UserService) DeleteAccount(userID string) error {
	_, err := s.userRepo.FindByID(userID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.SoftDelete(userID)
}

// Block/unblock
func (s *UserService) BlockUser(userID, blockedID string) error {
	if userID == blockedID {
		return errors.New("cannot block yourself")
	}
	_, err := s.userRepo.FindByID(blockedID)
	if err != nil {
		return errors.New("user not found")
	}
	return s.userRepo.BlockUser(userID, blockedID)
}

func (s *UserService) UnblockUser(userID, blockedID string) error {
	return s.userRepo.UnblockUser(userID, blockedID)
}

func (s *UserService) GetBlockedUsers(userID string) ([]*domain.UserResponse, error) {
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

func (s *UserService) IsBlocked(userID, blockedID string) (bool, error) {
	blocked, err := s.userRepo.IsBlocked(blockedID, userID)
	if err != nil {
		return false, err
	}
	if blocked {
		return true, nil
	}
	return s.userRepo.IsBlocked(userID, blockedID)
}

// Notification settings
func (s *UserService) SetNotificationMuted(userID, chatID string, muted bool) error {
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	return s.chatRepo.SetNotificationMuted(userID, chatID, muted)
}

func (s *UserService) IsNotificationMuted(userID, chatID string) (bool, error) {
	return s.chatRepo.IsNotificationMuted(userID, chatID)
}
