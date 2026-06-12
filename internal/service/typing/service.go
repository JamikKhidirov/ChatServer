package typingservice

import (
	"sync"
	"time"

	"ChatServerGolang/internal/ws"
)

type TypingService interface {
	UserTyping(userID, chatID string)
	UserStopTyping(userID, chatID string)
}

type typingService struct {
	hub        *ws.Hub
	mu         sync.Mutex
	typingTimers map[string]*time.Timer
}

func NewTypingService(hub *ws.Hub) TypingService {
	return &typingService{
		hub:          hub,
		typingTimers: make(map[string]*time.Timer),
	}
}

func (s *typingService) key(userID, chatID string) string {
	return userID + ":" + chatID
}

func (s *typingService) UserTyping(userID, chatID string) {
	s.mu.Lock()
	key := s.key(userID, chatID)
	if t, ok := s.typingTimers[key]; ok {
		t.Stop()
	}
	s.typingTimers[key] = time.AfterFunc(4*time.Second, func() {
		s.mu.Lock()
		delete(s.typingTimers, key)
		s.mu.Unlock()
		s.hub.BroadcastToChatExcept([]string{userID}, ws.WSOutgoingMessage{
			Type:    "user:stop_typing",
			Payload: map[string]string{"chatId": chatID, "userId": userID},
		})
	})
	s.mu.Unlock()

	s.hub.BroadcastToChatExcept([]string{userID}, ws.WSOutgoingMessage{
		Type:    "user:typing",
		Payload: map[string]string{"chatId": chatID, "userId": userID},
	})
}

func (s *typingService) UserStopTyping(userID, chatID string) {
	s.mu.Lock()
	key := s.key(userID, chatID)
	if t, ok := s.typingTimers[key]; ok {
		t.Stop()
		delete(s.typingTimers, key)
	}
	s.mu.Unlock()

	s.hub.BroadcastToChatExcept([]string{userID}, ws.WSOutgoingMessage{
		Type:    "user:stop_typing",
		Payload: map[string]string{"chatId": chatID, "userId": userID},
	})
}


