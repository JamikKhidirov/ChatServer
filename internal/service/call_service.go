package service

import (
	"errors"
	"time"

	"ChatServerGolang/internal/domain"
	"ChatServerGolang/internal/repository"

	"github.com/google/uuid"
)

type CallService struct {
	callRepo    *repository.CallRepository
	chatRepo    *repository.ChatRepository
	userRepo    *repository.UserRepository
	userService *UserService
}

func NewCallService(
	callRepo *repository.CallRepository,
	chatRepo *repository.ChatRepository,
	userRepo *repository.UserRepository,
	userService *UserService,
) *CallService {
	return &CallService{
		callRepo:    callRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		userService: userService,
	}
}

func (s *CallService) InitiateCall(chatID, callerID string) (*domain.Call, error) {
	_, err := s.chatRepo.FindByID(chatID)
	if err != nil {
		return nil, errors.New("chat not found")
	}

	isParticipant, _ := s.chatRepo.IsParticipant(chatID, callerID)
	if !isParticipant {
		return nil, errors.New("access denied")
	}

	participants, _ := s.chatRepo.GetParticipants(chatID)
	var calleeID string
	for _, p := range participants {
		if p.UserID != callerID {
			calleeID = p.UserID
			break
		}
	}

	if calleeID == "" {
		return nil, errors.New("no other participants to call")
	}

	activeCall, _ := s.callRepo.FindActiveByUser(calleeID)
	if activeCall != nil {
		return nil, errors.New("callee is already in a call")
	}

	now := time.Now()
	call := &domain.Call{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		CallerID:  callerID,
		CalleeID:  calleeID,
		Status:    domain.CallInitiated,
		StartedAt: now,
	}

	if err := s.callRepo.Create(call); err != nil {
		return nil, err
	}

	return call, nil
}

func (s *CallService) AcceptCall(callID, userID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.CalleeID != userID {
		return errors.New("only the callee can accept the call")
	}

	if call.Status != domain.CallInitiated {
		return errors.New("call is not in initiated state")
	}

	return s.callRepo.UpdateStatus(callID, domain.CallOngoing)
}

func (s *CallService) EndCall(callID, userID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.CallerID != userID && call.CalleeID != userID {
		return errors.New("access denied")
	}

	if call.Status == domain.CallEnded || call.Status == domain.CallMissed {
		return nil
	}

	return s.callRepo.UpdateStatus(callID, domain.CallEnded)
}

func (s *CallService) MissCall(callID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.Status != domain.CallInitiated {
		return nil
	}

	return s.callRepo.UpdateStatus(callID, domain.CallMissed)
}

func (s *CallService) RejectCall(callID, userID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.CalleeID != userID {
		return errors.New("only the callee can reject the call")
	}

	return s.callRepo.UpdateStatus(callID, domain.CallRejected)
}

func (s *CallService) GetCallByID(callID string) (*domain.Call, error) {
	return s.callRepo.FindByID(callID)
}

func (s *CallService) GetCallHistory(chatID, userID string) ([]*domain.CallResponse, error) {
	calls, err := s.callRepo.FindByChatAndUser(chatID, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*domain.CallResponse, 0)
	for _, call := range calls {
		resp, err := s.toCallResponse(call)
		if err != nil {
			continue
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *CallService) toCallResponse(call *domain.Call) (*domain.CallResponse, error) {
	caller, err := s.userRepo.FindByID(call.CallerID)
	if err != nil {
		return nil, err
	}
	callee, err := s.userRepo.FindByID(call.CalleeID)
	if err != nil {
		return nil, err
	}

	return &domain.CallResponse{
		ID:        call.ID,
		ChatID:    call.ChatID,
		Caller:    caller.ToResponse(),
		Callee:    callee.ToResponse(),
		Status:    call.Status,
		StartedAt: call.StartedAt,
		EndedAt:   call.EndedAt,
	}, nil
}
