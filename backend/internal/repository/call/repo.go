package callrepo

import (
	"database/sql"
	"time"

	calldomain "ChatServerGolang/backend/internal/domain/call"
	"ChatServerGolang/backend/internal/repository"
)

type callRepository struct {
	db *sql.DB
}

func NewCallRepository(db *sql.DB) repository.CallRepository {
	return &callRepository{db: db}
}

func (r *callRepository) Create(call *calldomain.Call) error {
	var endedAt *string
	if call.EndedAt != nil {
		s := call.EndedAt.Format(time.RFC3339)
		endedAt = &s
	}
	_, err := r.db.Exec(
		`INSERT INTO calls (id, chat_id, caller_id, callee_id, call_type, status, started_at, ended_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		call.ID, call.ChatID, call.CallerID, call.CalleeID, call.Type, call.Status,
		call.StartedAt.Format(time.RFC3339), endedAt,
	)
	return err
}

func (r *callRepository) FindByID(id string) (*calldomain.Call, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, caller_id, callee_id, call_type, status, started_at, ended_at
		FROM calls WHERE id = ?`, id,
	)
	return scanCall(row)
}

func (r *callRepository) FindActiveByUser(userID string) (*calldomain.Call, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, caller_id, callee_id, call_type, status, started_at, ended_at
		FROM calls WHERE (caller_id = ? OR callee_id = ?) AND status IN ('initiated', 'ongoing')
		ORDER BY started_at DESC LIMIT 1`,
		userID, userID,
	)
	return scanCall(row)
}

func (r *callRepository) FindByChatAndUser(chatID, userID string) ([]*calldomain.Call, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, caller_id, callee_id, call_type, status, started_at, ended_at
		FROM calls WHERE chat_id = ? AND (caller_id = ? OR callee_id = ?)
		ORDER BY started_at DESC LIMIT 50`,
		chatID, userID, userID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	calls := make([]*calldomain.Call, 0)
	for rows.Next() {
		c, err := scanCall(rows)
		if err != nil {
			return nil, err
		}
		calls = append(calls, c)
	}
	return calls, nil
}

func (r *callRepository) UpdateStatus(id string, status calldomain.CallStatus) error {
	var endedAt *string
	if status == calldomain.CallEnded || status == calldomain.CallMissed || status == calldomain.CallRejected {
		s := time.Now().Format(time.RFC3339)
		endedAt = &s
	}
	_, err := r.db.Exec(
		`UPDATE calls SET status=?, ended_at=? WHERE id=?`,
		status, endedAt, id,
	)
	return err
}

type callScanner interface {
	Scan(dest ...interface{}) error
}

func scanCall(row callScanner) (*calldomain.Call, error) {
	var (
		c         calldomain.Call
		startedAt string
		endedAt   sql.NullString
		callType  string
	)
	err := row.Scan(&c.ID, &c.ChatID, &c.CallerID, &c.CalleeID, &callType, &c.Status, &startedAt, &endedAt)
	if err != nil {
		return nil, err
	}
	c.Type = calldomain.CallType(callType)
	c.StartedAt = repository.ParseTime(startedAt)
	if endedAt.Valid {
		t := repository.ParseTime(endedAt.String)
		c.EndedAt = &t
	}
	return &c, nil
}



