package voicechatservice

import (
	"errors"
	"time"

	chatdomain "ChatServerGolang/internal/domain/chat"
	voicechatdomain "ChatServerGolang/internal/domain/voicechat"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"

	"github.com/google/uuid"
)

type voiceChatService struct {
	vcRepo  repository.VoiceChatRepository
	chatRepo repository.ChatRepository
	userRepo repository.UserRepository
}

func NewVoiceChatService(
	vcRepo repository.VoiceChatRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
) service.VoiceChatService {
	return &voiceChatService{
		vcRepo:   vcRepo,
		chatRepo: chatRepo,
		userRepo: userRepo,
	}
}

func (s *voiceChatService) CreateVoiceChat(chatID, userID string, req *voicechatdomain.CreateVoiceChatRequest) (*voicechatdomain.VoiceChatResponse, error) {
	chat, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}
	if chat.Type != chatdomain.ChatGroup {
		return nil, errors.New("voice chats are only available in groups")
	}

	now := time.Now()
	vc := &voicechatdomain.VoiceChat{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		StartedBy: userID,
		Title:     req.Title,
		Status:    voicechatdomain.VoiceChatActive,
		CreatedAt: now,
		StartedAt: &now,
	}

	if req.ScheduledInMins > 0 {
		scheduled := now.Add(time.Duration(req.ScheduledInMins) * time.Minute)
		vc.ScheduledAt = &scheduled
		vc.Status = voicechatdomain.VoiceChatScheduled
		vc.StartedAt = nil
	}

	if err := s.vcRepo.Create(vc); err != nil {
		return nil, err
	}

	if err := s.vcRepo.AddParticipant(vc.ID, userID); err != nil {
		return nil, err
	}

	return s.buildResponse(vc), nil
}

func (s *voiceChatService) GetVoiceChat(id string) (*voicechatdomain.VoiceChatResponse, error) {
	vc, err := s.vcRepo.FindByID(id)
	if err != nil {
		return nil, err
	}
	return s.buildResponse(vc), nil
}

func (s *voiceChatService) GetActiveVoiceChats(chatID string) ([]*voicechatdomain.VoiceChatResponse, error) {
	vcs, err := s.vcRepo.FindActiveByChatID(chatID)
	if err != nil {
		return nil, err
	}
	return s.buildResponses(vcs), nil
}

func (s *voiceChatService) GetVoiceChatHistory(chatID string) ([]*voicechatdomain.VoiceChatResponse, error) {
	vcs, err := s.vcRepo.FindByChatID(chatID)
	if err != nil {
		return nil, err
	}
	return s.buildResponses(vcs), nil
}

func (s *voiceChatService) JoinVoiceChat(vcID, userID string) error {
	return s.vcRepo.AddParticipant(vcID, userID)
}

func (s *voiceChatService) LeaveVoiceChat(vcID, userID string) error {
	return s.vcRepo.RemoveParticipant(vcID, userID)
}

func (s *voiceChatService) EndVoiceChat(vcID, userID string) error {
	return s.vcRepo.UpdateStatus(vcID, voicechatdomain.VoiceChatEnded)
}

func (s *voiceChatService) MuteParticipant(vcID, userID string, muted bool) error {
	return s.vcRepo.SetParticipantMuted(vcID, userID, muted)
}

func (s *voiceChatService) buildResponse(vc *voicechatdomain.VoiceChat) *voicechatdomain.VoiceChatResponse {
	return &voicechatdomain.VoiceChatResponse{
		ID:               vc.ID,
		ChatID:           vc.ChatID,
		StartedBy:        vc.StartedBy,
		Title:            vc.Title,
		Status:           vc.Status,
		ParticipantCount: vc.ParticipantCount,
		ScheduledAt:      vc.ScheduledAt,
		StartedAt:        vc.StartedAt,
		EndedAt:          vc.EndedAt,
		CreatedAt:        vc.CreatedAt,
	}
}

func (s *voiceChatService) buildResponses(vcs []*voicechatdomain.VoiceChat) []*voicechatdomain.VoiceChatResponse {
	res := make([]*voicechatdomain.VoiceChatResponse, 0, len(vcs))
	for _, vc := range vcs {
		res = append(res, s.buildResponse(vc))
	}
	return res
}
