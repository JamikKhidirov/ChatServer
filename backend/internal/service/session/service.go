package sessionservice

import (
	"errors"
	"time"

	"ChatServerGolang/backend/internal/domain/session"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type sessionService struct {
	sessionRepo repository.SessionRepository
}

func NewSessionService(sessionRepo repository.SessionRepository) service.SessionService {
	return &sessionService{sessionRepo: sessionRepo}
}

func (s *sessionService) CreateSession(userID, deviceName, ipAddress string) (*sessiondomain.Session, error) {
	now := time.Now()
	session := &sessiondomain.Session{
		ID:         uuid.New().String(),
		UserID:     userID,
		DeviceName: deviceName,
		IPAddress:  ipAddress,
		LastActive: now,
		CreatedAt:  now,
	}
	if err := s.sessionRepo.Create(session); err != nil {
		return nil, err
	}
	return session, nil
}

func (s *sessionService) GetSessions(userID string) ([]*sessiondomain.Session, error) {
	sessions, err := s.sessionRepo.FindByUserID(userID)
	if err != nil {
		return nil, err
	}
	if sessions == nil {
		return []*sessiondomain.Session{}, nil
	}
	return sessions, nil
}

func (s *sessionService) DeleteSession(sessionID, userID string) error {
	session, err := s.sessionRepo.FindByID(sessionID)
	if err != nil {
		return errors.New("session not found")
	}
	if session.UserID != userID {
		return errors.New("access denied")
	}
	return s.sessionRepo.Delete(sessionID)
}

func (s *sessionService) DeleteAllSessions(userID string) error {
	return s.sessionRepo.DeleteByUserID(userID)
}

func (s *sessionService) UpdateLastActive(sessionID string) error {
	return s.sessionRepo.UpdateLastActive(sessionID)
}


