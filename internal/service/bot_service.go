package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
)

type botService struct {
	botRepo repository.BotRepository
}

func NewBotService(botRepo repository.BotRepository) BotService {
	return &botService{botRepo: botRepo}
}

func (s *botService) CreateBot(userID string, req *domain.CreateBotRequest) (*domain.Bot, error) {
	token, _ := generateBotToken()
	bot := &domain.Bot{
		ID:         uuid.New().String(),
		Token:      token,
		OwnerID:    userID,
		Name:       req.Name,
		WebhookURL: req.WebhookURL,
		CreatedAt:  time.Now(),
		Active:     true,
	}
	if err := s.botRepo.Create(bot); err != nil {
		return nil, err
	}
	return bot, nil
}

func (s *botService) GetMyBots(userID string) ([]*domain.Bot, error) {
	bots, err := s.botRepo.FindByOwnerID(userID)
	if err != nil {
		return nil, err
	}
	if bots == nil {
		return []*domain.Bot{}, nil
	}
	return bots, nil
}

func (s *botService) UpdateBot(botID, userID string, req *domain.UpdateBotRequest) error {
	bot, err := s.botRepo.FindByID(botID)
	if err != nil {
		return errors.New("bot not found")
	}
	if bot.OwnerID != userID {
		return errors.New("access denied")
	}
	if req.Name != "" {
		bot.Name = req.Name
	}
	if req.AvatarURL != "" {
		bot.AvatarURL = req.AvatarURL
	}
	if req.WebhookURL != "" {
		bot.WebhookURL = req.WebhookURL
	}
	return s.botRepo.Update(bot)
}

func (s *botService) DeleteBot(botID, userID string) error {
	bot, err := s.botRepo.FindByID(botID)
	if err != nil {
		return errors.New("bot not found")
	}
	if bot.OwnerID != userID {
		return errors.New("access denied")
	}
	return s.botRepo.Delete(botID)
}

func (s *botService) RegenerateToken(botID, userID string) error {
	bot, err := s.botRepo.FindByID(botID)
	if err != nil {
		return errors.New("bot not found")
	}
	if bot.OwnerID != userID {
		return errors.New("access denied")
	}
	token, _ := generateBotToken()
	return s.botRepo.RegenerateToken(botID, token)
}

func (s *botService) ValidateBotToken(token string) (string, error) {
	bot, err := s.botRepo.FindByToken(token)
	if err != nil {
		return "", errors.New("invalid bot token")
	}
	if !bot.Active {
		return "", errors.New("bot is deactivated")
	}
	return bot.ID, nil
}

func generateBotToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "bot_" + hex.EncodeToString(bytes), nil
}
