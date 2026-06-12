package groupcallrepo

import (
	"database/sql"
	"time"

	calldomain "ChatServerGolang/internal/domain/call"
	"ChatServerGolang/internal/repository"
)

type groupCallRepository struct {
	db *sql.DB
}

func NewGroupCallRepository(db *sql.DB) repository.GroupCallRepository {
	return &groupCallRepository{db: db}
}

func (r *groupCallRepository) Create(call *calldomain.GroupCall) error {
	_, err := r.db.Exec(
		`INSERT INTO group_calls (id, chat_id, caller_id, type, status, started_at, ended_at)
		VALUES (?, ?, ?, ?, ?, ?, ?)`,
		call.ID, call.ChatID, call.CallerID, call.Type, call.Status,
		call.StartedAt.Format(time.RFC3339), nil,
	)
	return err
}

func (r *groupCallRepository) FindByID(id string) (*calldomain.GroupCall, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, caller_id, type, status, started_at, ended_at
		FROM group_calls WHERE id = ?`, id,
	)
	return scanGroupCall(row)
}

func (r *groupCallRepository) FindActiveByChatID(chatID string) ([]*calldomain.GroupCall, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, caller_id, type, status, started_at, ended_at
		FROM group_calls WHERE chat_id = ? AND status IN ('initiated','ongoing')
		ORDER BY started_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []*calldomain.GroupCall
	for rows.Next() {
		c, err := scanGroupCall(rows)
		if err != nil {
			return nil, err
		}
		calls = append(calls, c)
	}
	return calls, nil
}

func (r *groupCallRepository) FindActiveByUserID(userID string) (*calldomain.GroupCall, error) {
	row := r.db.QueryRow(
		`SELECT gc.id, gc.chat_id, gc.caller_id, gc.type, gc.status, gc.started_at, gc.ended_at
		FROM group_calls gc
		INNER JOIN group_call_participants gcp ON gcp.call_id = gc.id
		WHERE gcp.user_id = ? AND gc.status IN ('initiated','ongoing')
		LIMIT 1`, userID,
	)
	return scanGroupCall(row)
}

func (r *groupCallRepository) UpdateStatus(id string, status calldomain.CallStatus) error {
	now := time.Now().Format(time.RFC3339)
	var err error
	if status == calldomain.CallEnded || status == calldomain.CallMissed || status == calldomain.CallRejected {
		_, err = r.db.Exec(`UPDATE group_calls SET status=?, ended_at=? WHERE id=?`, status, now, id)
	} else {
		_, err = r.db.Exec(`UPDATE group_calls SET status=? WHERE id=?`, status, id)
	}
	return err
}

func (r *groupCallRepository) AddParticipant(callID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO group_call_participants (call_id, user_id, joined_at)
		VALUES (?, ?, ?)`,
		callID, userID, time.Now().Format(time.RFC3339),
	)
	return err
}

func (r *groupCallRepository) RemoveParticipant(callID, userID string) error {
	_, err := r.db.Exec(
		`UPDATE group_call_participants SET left_at=? WHERE call_id=? AND user_id=?`,
		time.Now().Format(time.RFC3339), callID, userID,
	)
	return err
}

func (r *groupCallRepository) UpdateParticipantMute(callID, userID string, audioMuted, videoMuted bool) error {
	_, err := r.db.Exec(
		`UPDATE group_call_participants SET audio_muted=?, video_muted=? WHERE call_id=? AND user_id=?`,
		repository.BoolToInt(audioMuted), repository.BoolToInt(videoMuted), callID, userID,
	)
	return err
}

func (r *groupCallRepository) GetParticipants(callID string) ([]*calldomain.GroupCallParticipant, error) {
	rows, err := r.db.Query(
		`SELECT call_id, user_id, joined_at, left_at, audio_muted, video_muted
		FROM group_call_participants WHERE call_id = ?`, callID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*calldomain.GroupCallParticipant
	for rows.Next() {
		var p calldomain.GroupCallParticipant
		var joinedAt string
		var leftAt *string
		var audioMuted, videoMuted int
		if err := rows.Scan(&p.CallID, &p.UserID, &joinedAt, &leftAt, &audioMuted, &videoMuted); err != nil {
			return nil, err
		}
		p.JoinedAt = repository.ParseTime(joinedAt)
		if leftAt != nil {
			t := repository.ParseTime(*leftAt)
			p.LeftAt = &t
		}
		p.AudioMuted = audioMuted == 1
		p.VideoMuted = videoMuted == 1
		participants = append(participants, &p)
	}
	return participants, nil
}

func (r *groupCallRepository) FindByChatAndUser(chatID, userID string) ([]*calldomain.GroupCall, error) {
	rows, err := r.db.Query(
		`SELECT gc.id, gc.chat_id, gc.caller_id, gc.type, gc.status, gc.started_at, gc.ended_at
		FROM group_calls gc
		INNER JOIN group_call_participants gcp ON gcp.call_id = gc.id
		WHERE gc.chat_id = ? AND gcp.user_id = ?
		ORDER BY gc.started_at DESC`, chatID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var calls []*calldomain.GroupCall
	for rows.Next() {
		c, err := scanGroupCall(rows)
		if err != nil {
			return nil, err
		}
		calls = append(calls, c)
	}
	return calls, nil
}

func scanGroupCall(row repository.Scanner) (*calldomain.GroupCall, error) {
	var c calldomain.GroupCall
	var startedAt string
	var endedAt *string
	err := row.Scan(&c.ID, &c.ChatID, &c.CallerID, &c.Type, &c.Status, &startedAt, &endedAt)
	if err != nil {
		return nil, err
	}
	c.StartedAt = repository.ParseTime(startedAt)
	if endedAt != nil {
		t := repository.ParseTime(*endedAt)
		c.EndedAt = &t
	}
	return &c, nil
}
