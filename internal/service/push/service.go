package pushservice

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"

	"ChatServerGolang/config"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"
)

type pushService struct {
	userRepo repository.UserRepository
	cfg      *config.Config
	client   *http.Client
}

func NewPushService(userRepo repository.UserRepository, cfg *config.Config) service.PushService {
	return &pushService{
		userRepo: userRepo,
		cfg:      cfg,
		client:   &http.Client{},
	}
}

type FCMRequest struct {
	To           string           `json:"to"`
	Priority     string           `json:"priority"`
	Data         *FCMData         `json:"data,omitempty"`
	Notification *FCMNotification `json:"notification,omitempty"`
}

type FCMData struct {
	Type    string `json:"type"`
	ChatID  string `json:"chatId,omitempty"`
	Message string `json:"message,omitempty"`
	CallID  string `json:"callId,omitempty"`
	Title   string `json:"title,omitempty"`
	Body    string `json:"body,omitempty"`
}

type FCMNotification struct {
	Title string `json:"title"`
	Body  string `json:"body"`
}

func (s *pushService) SendMessageNotification(senderID, chatID, msgID, msgContent, msgType string) {
	if !s.cfg.PushEnabled {
		return
	}

	sender, err := s.userRepo.FindByID(senderID)
	if err != nil {
		return
	}

	participants, err := s.userRepo.GetParticipantsInChat(chatID)
	if err != nil {
		log.Printf("Failed to get participants for push: %v", err)
		return
	}

	for _, user := range participants {
		if user.ID == senderID || user.PushToken == "" {
			continue
		}

		title := sender.DisplayName
		var body string
		switch msgType {
		case "text":
			body = msgContent
		case "image":
			body = sender.DisplayName + " sent a photo"
		case "file":
			body = sender.DisplayName + " sent a file"
		default:
			body = msgContent
		}
		if len(body) > 200 {
			body = body[:200] + "..."
		}

		if user.PushProvider == "fcm" {
			s.sendFCM(user.PushToken, title, body, &FCMData{
				Type:    "message",
				ChatID:  chatID,
				Title:   title,
				Body:    body,
				Message: msgID,
			})
		}
	}
}

func (s *pushService) SendCallNotification(callerID, chatID, callID, callType string) {
	if !s.cfg.PushEnabled {
		return
	}

	caller, err := s.userRepo.FindByID(callerID)
	if err != nil {
		return
	}

	participants, err := s.userRepo.GetParticipantsInChat(chatID)
	if err != nil {
		return
	}

	title := caller.DisplayName
	body := fmt.Sprintf("Incoming %s call...", callType)

	for _, user := range participants {
		if user.ID == callerID || user.PushToken == "" {
			continue
		}
		if user.PushProvider == "fcm" {
			s.sendFCM(user.PushToken, title, body, &FCMData{
				Type:   "call",
				CallID: callID,
				ChatID: chatID,
				Title:  title,
				Body:   body,
			})
		}
	}
}

func (s *pushService) SendTestPush(userID, title, body string) {
	sender, err := s.userRepo.FindByID(userID)
	if err != nil || sender.PushToken == "" {
		log.Printf("[Push] No push token for user %s", userID)
		return
	}
	s.sendFCM(sender.PushToken, title, body, &FCMData{
		Type:  "test",
		Title: title,
		Body:  body,
	})
}

func (s *pushService) sendFCM(token, title, body string, data *FCMData) {
	if s.cfg.FirebaseCredentials == "" {
		log.Printf("[FCM] Would send push (no credentials configured): title=%s, body=%s", title, body)
		return
	}

	req := &FCMRequest{
		To:       token,
		Priority: "high",
		Data:     data,
		Notification: &FCMNotification{
			Title: title,
			Body:  body,
		},
	}

	bodyBytes, _ := json.Marshal(req)
	httpReq, err := http.NewRequest("POST", "https://fcm.googleapis.com/fcm/send", strings.NewReader(string(bodyBytes)))
	if err != nil {
		log.Printf("Failed to create FCM request: %v", err)
		return
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "key="+s.cfg.FirebaseCredentials)

	resp, err := s.client.Do(httpReq)
	if err != nil {
		log.Printf("Failed to send FCM push: %v", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("FCM push sent, status: %d", resp.StatusCode)
}


