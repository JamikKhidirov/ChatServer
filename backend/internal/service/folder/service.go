package folderservice

import (
	"errors"
	"time"

	"ChatServerGolang/backend/internal/domain/chat"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type chatFolderService struct {
	folderRepo repository.ChatFolderRepository
	chatRepo   repository.ChatRepository
}

func NewChatFolderService(folderRepo repository.ChatFolderRepository, chatRepo repository.ChatRepository) service.ChatFolderService {
	return &chatFolderService{folderRepo: folderRepo, chatRepo: chatRepo}
}

func (s *chatFolderService) Create(userID string, req *chatdomain.CreateChatFolderRequest) (*chatdomain.ChatFolderWithChats, error) {
	folders, err := s.folderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	folder := &chatdomain.ChatFolder{
		ID:        uuid.New().String(),
		UserID:    userID,
		Name:      req.Name,
		Emoji:     req.Emoji,
		Order:     len(folders),
		CreatedAt: time.Now(),
	}

	if err := s.folderRepo.Create(folder); err != nil {
		return nil, err
	}

	if len(req.ChatIDs) > 0 {
		if err := s.folderRepo.SetChatsForFolder(folder.ID, req.ChatIDs); err != nil {
			return nil, err
		}
	}

	return s.GetWithChats(folder.ID, userID)
}

func (s *chatFolderService) List(userID string) ([]*chatdomain.ChatFolderWithChats, error) {
	folders, err := s.folderRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}

	result := make([]*chatdomain.ChatFolderWithChats, 0, len(folders))
	for _, f := range folders {
		withChats, err := s.buildWithChats(f)
		if err != nil {
			continue
		}
		result = append(result, withChats)
	}
	return result, nil
}

func (s *chatFolderService) Update(folderID, userID string, req *chatdomain.UpdateChatFolderRequest) (*chatdomain.ChatFolderWithChats, error) {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}
	if folder.UserID != userID {
		return nil, errors.New("access denied")
	}

	if req.Name != "" {
		folder.Name = req.Name
	}
	if req.Emoji != "" {
		folder.Emoji = req.Emoji
	}
	if req.Order != 0 {
		folder.Order = req.Order
	}

	if err := s.folderRepo.Update(folder); err != nil {
		return nil, err
	}

	if req.ChatIDs != nil {
		if err := s.folderRepo.SetChatsForFolder(folderID, req.ChatIDs); err != nil {
			return nil, err
		}
	}

	return s.GetWithChats(folderID, userID)
}

func (s *chatFolderService) Delete(folderID, userID string) error {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return errors.New("folder not found")
	}
	if folder.UserID != userID {
		return errors.New("access denied")
	}
	return s.folderRepo.Delete(folderID)
}

func (s *chatFolderService) GetWithChats(folderID, userID string) (*chatdomain.ChatFolderWithChats, error) {
	folder, err := s.folderRepo.FindByID(folderID)
	if err != nil {
		return nil, errors.New("folder not found")
	}
	if folder.UserID != userID {
		return nil, errors.New("access denied")
	}
	return s.buildWithChats(folder)
}

func (s *chatFolderService) buildWithChats(folder *chatdomain.ChatFolder) (*chatdomain.ChatFolderWithChats, error) {
	chatIDs, err := s.folderRepo.GetChatIDsByFolder(folder.ID)
	if err != nil {
		return nil, err
	}

	return &chatdomain.ChatFolderWithChats{
		ChatFolder: *folder,
		ChatIDs:    chatIDs,
	}, nil
}
