package storyservice

import (
	"errors"
	"os"
	"path/filepath"
	"time"

	storydomain "ChatServerGolang/internal/domain/story"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"

	"github.com/google/uuid"
)

type storyService struct {
	storyRepo repository.StoryRepository
	userRepo  repository.UserRepository
	chatRepo  repository.ChatRepository
}

func NewStoryService(storyRepo repository.StoryRepository, userRepo repository.UserRepository, chatRepo repository.ChatRepository) service.StoryService {
	return &storyService{
		storyRepo: storyRepo,
		userRepo:  userRepo,
		chatRepo:  chatRepo,
	}
}

func (s *storyService) CreateStory(userID string, req *storydomain.CreateStoryRequest, filePath, fileURL string) (*storydomain.StoryResponse, error) {
	now := time.Now()
	story := &storydomain.Story{
		ID:        uuid.New().String(),
		UserID:    userID,
		FilePath:  filePath,
		FileURL:   fileURL,
		Type:      req.Type,
		Caption:   req.Caption,
		ExpiresAt: now.Add(24 * time.Hour),
		CreatedAt: now,
	}
	if err := s.storyRepo.Create(story); err != nil {
		return nil, err
	}
	return s.toResponse(story), nil
}

func (s *storyService) GetMyStories(userID string) ([]*storydomain.StoryResponse, error) {
	stories, err := s.storyRepo.FindActiveByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(stories, userID), nil
}

func (s *storyService) GetFollowingStories(userID string) ([]*storydomain.StoryResponse, error) {
	chats, err := s.chatRepo.FindByUserID(userID)
	if err != nil {
		return []*storydomain.StoryResponse{}, nil
	}
	contactIDs := make([]string, 0)
	seen := make(map[string]bool)
	for _, chat := range chats {
		participants, _ := s.chatRepo.GetParticipants(chat.ID)
		for _, p := range participants {
			if p.UserID != userID && !seen[p.UserID] {
				seen[p.UserID] = true
				contactIDs = append(contactIDs, p.UserID)
			}
		}
	}
	if len(contactIDs) == 0 {
		return []*storydomain.StoryResponse{}, nil
	}
	stories, err := s.storyRepo.FindActiveByFollowing(contactIDs)
	if err != nil {
		return []*storydomain.StoryResponse{}, nil
	}
	return s.toResponses(stories, userID), nil
}

func (s *storyService) GetStoryByID(storyID, userID string) (*storydomain.StoryResponse, error) {
	story, err := s.storyRepo.FindByID(storyID)
	if err != nil {
		return nil, errors.New("story not found")
	}
	s.storyRepo.AddView(storyID, userID)
	return s.toResponse(story), nil
}

func (s *storyService) DeleteStory(storyID, userID string) error {
	story, err := s.storyRepo.FindByID(storyID)
	if err != nil {
		return errors.New("story not found")
	}
	if story.UserID != userID {
		return errors.New("access denied")
	}
	if story.FilePath != "" {
		os.Remove(filepath.Join("uploads/stories", story.FilePath))
	}
	return s.storyRepo.Delete(storyID)
}

func (s *storyService) GetStoryViews(storyID, userID string) ([]*storydomain.StoryView, error) {
	story, err := s.storyRepo.FindByID(storyID)
	if err != nil {
		return nil, errors.New("story not found")
	}
	if story.UserID != userID {
		return nil, errors.New("access denied")
	}
	return s.storyRepo.GetViews(storyID)
}

func (s *storyService) toResponse(story *storydomain.Story) *storydomain.StoryResponse {
	return &storydomain.StoryResponse{
		ID:        story.ID,
		UserID:    story.UserID,
		FileURL:   story.FileURL,
		Type:      story.Type,
		Caption:   story.Caption,
		ExpiresAt: story.ExpiresAt,
		CreatedAt: story.CreatedAt,
	}
}

func (s *storyService) toResponses(stories []*storydomain.Story, userID string) []*storydomain.StoryResponse {
	responses := make([]*storydomain.StoryResponse, 0, len(stories))
	for _, story := range stories {
		resp := s.toResponse(story)
		count, _ := s.storyRepo.GetViewCount(story.ID)
		resp.Views = count
		viewed, _ := s.storyRepo.HasViewed(story.ID, userID)
		resp.Viewed = viewed
		responses = append(responses, resp)
	}
	return responses
}
