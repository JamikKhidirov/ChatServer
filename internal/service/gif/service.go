package gifservice

import (
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"
)

type savedGifService struct {
	gifRepo repository.SavedGifRepository
}

func NewSavedGifService(gifRepo repository.SavedGifRepository) service.SavedGifService {
	return &savedGifService{gifRepo: gifRepo}
}

func (s *savedGifService) SaveGif(userID, gifURL string) error {
	return s.gifRepo.Save(userID, gifURL)
}

func (s *savedGifService) GetSavedGifs(userID string) ([]string, error) {
	return s.gifRepo.FindByUserID(userID)
}

func (s *savedGifService) DeleteGif(userID, gifURL string) error {
	return s.gifRepo.Delete(userID, gifURL)
}
