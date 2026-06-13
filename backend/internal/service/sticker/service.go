package stickerservice

import (
	"errors"
	"time"

	"ChatServerGolang/backend/internal/domain/sticker"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type stickerService struct {
	stickerRepo repository.StickerRepository
}

func NewStickerService(stickerRepo repository.StickerRepository) service.StickerService {
	return &stickerService{stickerRepo: stickerRepo}
}

func (s *stickerService) CreatePack(userID string, req *stickerdomain.CreateStickerPackRequest) (*stickerdomain.StickerPackWithStickers, error) {
	pack := &stickerdomain.StickerPack{
		ID:        uuid.New().String(),
		Name:      req.Name,
		CreatorID: userID,
		Animated:  req.Animated,
		CreatedAt: time.Now(),
	}
	if err := s.stickerRepo.CreatePack(pack); err != nil {
		return nil, err
	}
	return &stickerdomain.StickerPackWithStickers{
		StickerPack: *pack,
		Stickers:    []*stickerdomain.Sticker{},
	}, nil
}

func (s *stickerService) GetPacks() ([]*stickerdomain.StickerPackWithStickers, error) {
	packs, err := s.stickerRepo.ListPacks()
	if err != nil {
		return nil, err
	}
	return s.enrichPacks(packs)
}

func (s *stickerService) GetMyPacks(userID string) ([]*stickerdomain.StickerPackWithStickers, error) {
	packs, err := s.stickerRepo.GetPacksByUserID(userID)
	if err != nil {
		return nil, err
	}
	return s.enrichPacks(packs)
}

func (s *stickerService) GetPackByID(id string) (*stickerdomain.StickerPackWithStickers, error) {
	pack, err := s.stickerRepo.GetPackByID(id)
	if err != nil {
		return nil, errors.New("pack not found")
	}
	stickers, err := s.stickerRepo.GetStickersByPackID(id)
	if err != nil {
		return nil, err
	}
	if stickers == nil {
		stickers = []*stickerdomain.Sticker{}
	}
	return &stickerdomain.StickerPackWithStickers{
		StickerPack: *pack,
		Stickers:    stickers,
	}, nil
}

func (s *stickerService) AddSticker(packID, userID string, req *stickerdomain.AddStickerRequest) (*stickerdomain.Sticker, error) {
	pack, err := s.stickerRepo.GetPackByID(packID)
	if err != nil {
		return nil, errors.New("pack not found")
	}
	if pack.CreatorID != userID {
		return nil, errors.New("only the pack owner can add stickers")
	}

	sticker := &stickerdomain.Sticker{
		ID:       uuid.New().String(),
		PackID:   packID,
		Emoji:    req.Emoji,
		ImageURL: req.ImageURL,
	}
	if err := s.stickerRepo.AddSticker(sticker); err != nil {
		return nil, err
	}
	return sticker, nil
}

func (s *stickerService) DeletePack(id, userID string) error {
	pack, err := s.stickerRepo.GetPackByID(id)
	if err != nil {
		return errors.New("pack not found")
	}
	if pack.CreatorID != userID {
		return errors.New("only the pack owner can delete")
	}
	return s.stickerRepo.DeletePack(id)
}

func (s *stickerService) DeleteSticker(id, userID string) error {
	return s.stickerRepo.DeleteSticker(id)
}

func (s *stickerService) AddToLibrary(userID, stickerID string) error {
	return s.stickerRepo.AddToUserLibrary(userID, stickerID)
}

func (s *stickerService) GetLibrary(userID string) ([]*stickerdomain.Sticker, error) {
	stickers, err := s.stickerRepo.GetUserLibrary(userID)
	if err != nil {
		return nil, err
	}
	if stickers == nil {
		return []*stickerdomain.Sticker{}, nil
	}
	return stickers, nil
}

func (s *stickerService) enrichPacks(packs []*stickerdomain.StickerPack) ([]*stickerdomain.StickerPackWithStickers, error) {
	result := make([]*stickerdomain.StickerPackWithStickers, 0, len(packs))
	for _, p := range packs {
		stickers, err := s.stickerRepo.GetStickersByPackID(p.ID)
		if err != nil {
			return nil, err
		}
			if stickers == nil {
			stickers = []*stickerdomain.Sticker{}
		}
		result = append(result, &stickerdomain.StickerPackWithStickers{
			StickerPack: *p,
			Stickers:    stickers,
		})
	}
	return result, nil
}


