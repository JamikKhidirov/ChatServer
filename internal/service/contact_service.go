package service

import (
	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"
)

type contactService struct {
	contactRepo repository.ContactRepository
}

func NewContactService(contactRepo repository.ContactRepository) ContactService {
	return &contactService{contactRepo: contactRepo}
}

func (s *contactService) SyncContacts(userID string, req *domain.SyncContactsRequest) error {
	return s.contactRepo.SyncContacts(userID, req.Contacts)
}

func (s *contactService) GetContacts(userID string) ([]*domain.ContactResponse, error) {
	return s.contactRepo.GetContacts(userID)
}

func (s *contactService) SearchByPhone(userID, query string) ([]*domain.ContactResponse, error) {
	return s.contactRepo.SearchByPhone(userID, query)
}

func (s *contactService) FindRegisteredByPhone(userID string) ([]*domain.UserResponse, error) {
	contacts, err := s.contactRepo.GetContacts(userID)
	if err != nil {
		return nil, err
	}
	if len(contacts) == 0 {
		return nil, nil
	}
	phones := make([]string, len(contacts))
	for i, c := range contacts {
		phones[i] = c.Phone
	}
	return s.contactRepo.FindRegisteredByPhone(phones)
}
