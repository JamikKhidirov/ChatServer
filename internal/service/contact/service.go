package contactservice

import (
	contactdomain "ChatServerGolang/internal/domain/contact"
	userdomain "ChatServerGolang/internal/domain/user"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"
)

type contactService struct {
	contactRepo repository.ContactRepository
}

func NewContactService(contactRepo repository.ContactRepository) service.ContactService {
	return &contactService{contactRepo: contactRepo}
}

func (s *contactService) SyncContacts(userID string, req *contactdomain.SyncContactsRequest) error {
	return s.contactRepo.SyncContacts(userID, req.Contacts)
}

func (s *contactService) GetContacts(userID string) ([]*contactdomain.ContactResponse, error) {
	return s.contactRepo.GetContacts(userID)
}

func (s *contactService) SearchByPhone(userID, query string) ([]*contactdomain.ContactResponse, error) {
	return s.contactRepo.SearchByPhone(userID, query)
}

func (s *contactService) UpdateContactPhoto(userID, phone, photoURL string) error {
	return s.contactRepo.UpdateContactPhoto(userID, phone, photoURL)
}

func (s *contactService) FindRegisteredByPhone(userID string) ([]*userdomain.UserResponse, error) {
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


