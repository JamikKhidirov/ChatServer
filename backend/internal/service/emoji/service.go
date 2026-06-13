package emojiservice

import (
	"time"

	emojidomain "ChatServerGolang/backend/internal/domain/emoji"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type customEmojiService struct {
	emojiRepo repository.CustomEmojiRepository
}

func NewCustomEmojiService(emojiRepo repository.CustomEmojiRepository) service.CustomEmojiService {
	return &customEmojiService{emojiRepo: emojiRepo}
}

func (s *customEmojiService) CreateEmoji(userID, shortcode string, filePath, fileURL string) (*emojidomain.CustomEmojiResponse, error) {
	emoji := &emojidomain.CustomEmoji{
		ID:        uuid.New().String(),
		UserID:    userID,
		Shortcode: shortcode,
		FileURL:   fileURL,
		FilePath:  filePath,
		Animated:  false,
		CreatedAt: time.Now(),
	}
	if err := s.emojiRepo.Create(emoji); err != nil {
		return nil, err
	}
	return s.toResponse(emoji), nil
}

func (s *customEmojiService) GetMyEmojis(userID string) ([]*emojidomain.CustomEmojiResponse, error) {
	emojis, err := s.emojiRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.toResponses(emojis), nil
}

func (s *customEmojiService) GetAllEmojis() ([]*emojidomain.CustomEmojiResponse, error) {
	emojis, err := s.emojiRepo.FindAll()
	if err != nil {
		return nil, err
	}
	return s.toResponses(emojis), nil
}

func (s *customEmojiService) DeleteEmoji(id, userID string) error {
	return s.emojiRepo.Delete(id, userID)
}

func (s *customEmojiService) toResponse(e *emojidomain.CustomEmoji) *emojidomain.CustomEmojiResponse {
	return &emojidomain.CustomEmojiResponse{
		ID:        e.ID,
		Shortcode: e.Shortcode,
		FileURL:   e.FileURL,
		Animated:  e.Animated,
		CreatedAt: e.CreatedAt,
	}
}

func (s *customEmojiService) toResponses(emojis []*emojidomain.CustomEmoji) []*emojidomain.CustomEmojiResponse {
	res := make([]*emojidomain.CustomEmojiResponse, 0, len(emojis))
	for _, e := range emojis {
		res = append(res, s.toResponse(e))
	}
	return res
}
