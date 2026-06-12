package groupcallservice

import (
	"errors"
	"time"

	calldomain "ChatServerGolang/internal/domain/call"
	"ChatServerGolang/internal/repository"
	"ChatServerGolang/internal/service"

	"github.com/google/uuid"
)

type groupCallService struct {
	groupCallRepo repository.GroupCallRepository
	chatRepo      repository.ChatRepository
	userRepo      repository.UserRepository
}

func NewGroupCallService(groupCallRepo repository.GroupCallRepository, chatRepo repository.ChatRepository, userRepo repository.UserRepository) service.GroupCallService {
	return &groupCallService{
		groupCallRepo: groupCallRepo,
		chatRepo:      chatRepo,
		userRepo:      userRepo,
	}
}

func (s *groupCallService) InitiateGroupCall(chatID, callerID string, callType calldomain.CallType) (*calldomain.GroupCallResponse, error) {
	_, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}
	isParticipant, _ := s.chatRepo.IsParticipant(chatID, callerID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}
	active, _ := s.groupCallRepo.FindActiveByChatID(chatID)
	if len(active) > 0 {
		return nil, errors.New("active call already exists")
	}
	now := time.Now()
	call := &calldomain.GroupCall{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		CallerID:  callerID,
		Type:      callType,
		Status:    calldomain.CallInitiated,
		StartedAt: now,
	}
	if err := s.groupCallRepo.Create(call); err != nil {
		return nil, err
	}
	s.groupCallRepo.AddParticipant(call.ID, callerID)
	return s.getResponse(call.ID)
}

func (s *groupCallService) JoinGroupCall(callID, userID string) error {
	call, err := s.groupCallRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}
	isParticipant, _ := s.chatRepo.IsParticipant(call.ChatID, userID)
	if !isParticipant {
		return errors.New("access denied")
	}
	if call.Status != calldomain.CallInitiated && call.Status != calldomain.CallOngoing {
		return errors.New("call is not active")
	}
	if call.Status == calldomain.CallInitiated {
		s.groupCallRepo.UpdateStatus(callID, calldomain.CallOngoing)
	}
	return s.groupCallRepo.AddParticipant(callID, userID)
}

func (s *groupCallService) LeaveGroupCall(callID, userID string) error {
	if err := s.groupCallRepo.RemoveParticipant(callID, userID); err != nil {
		return err
	}
	participants, _ := s.groupCallRepo.GetParticipants(callID)
	activeCount := 0
	for _, p := range participants {
		if p.LeftAt == nil {
			activeCount++
		}
	}
	if activeCount == 0 {
		s.groupCallRepo.UpdateStatus(callID, calldomain.CallEnded)
	}
	return nil
}

func (s *groupCallService) EndGroupCall(callID, userID string) error {
	call, err := s.groupCallRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}
	if call.CallerID != userID {
		return errors.New("only the caller can end the call")
	}
	return s.groupCallRepo.UpdateStatus(callID, calldomain.CallEnded)
}

func (s *groupCallService) MuteParticipant(callID, userID string, audioMuted, videoMuted bool) error {
	return s.groupCallRepo.UpdateParticipantMute(callID, userID, audioMuted, videoMuted)
}

func (s *groupCallService) GetGroupCallByID(callID string) (*calldomain.GroupCallResponse, error) {
	return s.getResponse(callID)
}

func (s *groupCallService) GetActiveGroupCalls(chatID, userID string) ([]*calldomain.GroupCallResponse, error) {
	calls, err := s.groupCallRepo.FindActiveByChatID(chatID)
	if err != nil {
		return nil, err
	}
	responses := make([]*calldomain.GroupCallResponse, 0, len(calls))
	for _, call := range calls {
		resp, _ := s.getResponse(call.ID)
		if resp != nil {
			responses = append(responses, resp)
		}
	}
	return responses, nil
}

func (s *groupCallService) getResponse(callID string) (*calldomain.GroupCallResponse, error) {
	call, err := s.groupCallRepo.FindByID(callID)
	if err != nil {
		return nil, errors.New("call not found")
	}
	participants, err := s.groupCallRepo.GetParticipants(callID)
	if err != nil {
		return nil, err
	}
	var duration int
	if call.EndedAt != nil {
		duration = int(call.EndedAt.Sub(call.StartedAt).Seconds())
	}
	resp := &calldomain.GroupCallResponse{
		ID:           call.ID,
		ChatID:       call.ChatID,
		CallerID:     call.CallerID,
		Type:         call.Type,
		Status:       call.Status,
		Participants: participants,
		StartedAt:    call.StartedAt,
		EndedAt:      call.EndedAt,
	}
	_ = duration
	return resp, nil
}
