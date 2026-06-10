package service

import "ChatServerGolang/internal/repository"

type savedGifService struct {
	gifRepo repository.SavedGifRepository
}

func NewSavedGifService(gifRepo repository.SavedGifRepository) SavedGifService {
	return &savedGifService{gifRepo: gifRepo}
}

func (s *savedGifService) SaveGif(userID, gifURL string) error {
	return s.gifRepo.Save(userID, gifURL)
}

func (s *savedGifService) GetSavedGifs(userID string) ([]string, error) {
	gifs, err := s.gifRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if gifs == nil {
		return []string{}, nil
	}
	return gifs, nil
}

func (s *savedGifService) DeleteGif(userID, gifURL string) error {
	return s.gifRepo.Delete(userID, gifURL)
}
