package callservice

import (
	"errors"
	"time"

	"ChatServerGolang/backend/internal/domain/call"
	"ChatServerGolang/backend/internal/repository"
	"ChatServerGolang/backend/internal/service"

	"github.com/google/uuid"
)

type callService struct {
	callRepo    repository.CallRepository
	chatRepo    repository.ChatRepository
	userRepo    repository.UserRepository
	userService service.UserService
}

func NewCallService(
	callRepo repository.CallRepository,
	chatRepo repository.ChatRepository,
	userRepo repository.UserRepository,
	userService service.UserService,
) service.CallService {
	return &callService{
		callRepo:    callRepo,
		chatRepo:    chatRepo,
		userRepo:    userRepo,
		userService: userService,
	}
}

func (s *callService) InitiateCall(chatID, callerID string, callType calldomain.CallType) (*calldomain.Call, error) {
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

	if callType == "" {
		callType = calldomain.CallAudio
	}

	now := time.Now()
	call := &calldomain.Call{
		ID:        uuid.New().String(),
		ChatID:    chatID,
		CallerID:  callerID,
		CalleeID:  calleeID,
		Type:      callType,
		Status:    calldomain.CallInitiated,
		StartedAt: now,
	}

	if err := s.callRepo.Create(call); err != nil {
		return nil, err
	}

	// Schedule auto-miss after 30 seconds
	go func() {
		time.Sleep(30 * time.Second)
		s.MissCall(call.ID)
	}()

	return call, nil
}

func (s *callService) AcceptCall(callID, userID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.CalleeID != userID {
		return errors.New("only the callee can accept the call")
	}

	if call.Status != calldomain.CallInitiated {
		return errors.New("call is not in initiated state")
	}

	return s.callRepo.UpdateStatus(callID, calldomain.CallOngoing)
}

func (s *callService) EndCall(callID, userID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.CallerID != userID && call.CalleeID != userID {
		return errors.New("access denied")
	}

	if call.Status == calldomain.CallEnded || call.Status == calldomain.CallMissed {
		return nil
	}

	return s.callRepo.UpdateStatus(callID, calldomain.CallEnded)
}

func (s *callService) MissCall(callID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return nil
	}

	if call.Status != calldomain.CallInitiated {
		return nil
	}

	return s.callRepo.UpdateStatus(callID, calldomain.CallMissed)
}

func (s *callService) RejectCall(callID, userID string) error {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return errors.New("call not found")
	}

	if call.CalleeID != userID {
		return errors.New("only the callee can reject the call")
	}

	return s.callRepo.UpdateStatus(callID, calldomain.CallRejected)
}

func (s *callService) GetCallByID(callID string) (*calldomain.Call, error) {
	call, err := s.callRepo.FindByID(callID)
	if err != nil {
		return nil, errors.New("call not found")
	}
	return call, nil
}

func (s *callService) GetCallHistory(chatID, userID string) ([]*calldomain.CallResponse, error) {
	calls, err := s.callRepo.FindByChatAndUser(chatID, userID)
	if err != nil {
		return nil, err
	}

	responses := make([]*calldomain.CallResponse, 0, len(calls))
	for _, call := range calls {
		resp, err := s.toCallResponse(call)
		if err != nil {
			continue
		}
		responses = append(responses, resp)
	}
	return responses, nil
}

func (s *callService) toCallResponse(call *calldomain.Call) (*calldomain.CallResponse, error) {
	caller, err := s.userRepo.FindByID(call.CallerID)
	if err != nil {
		return nil, err
	}
	callee, err := s.userRepo.FindByID(call.CalleeID)
	if err != nil {
		return nil, err
	}

	var duration int
	if call.EndedAt != nil {
		duration = int(call.EndedAt.Sub(call.StartedAt).Seconds())
	}

	return &calldomain.CallResponse{
		ID:        call.ID,
		ChatID:    call.ChatID,
		Caller:    caller.ToResponse(),
		Callee:    callee.ToResponse(),
		Type:      call.Type,
		Status:    call.Status,
		StartedAt: call.StartedAt,
		EndedAt:   call.EndedAt,
		Duration:  duration,
	}, nil
}



