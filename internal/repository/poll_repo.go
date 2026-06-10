package repository

import (
	"database/sql"
	"time"

	"ChatServerGolang/internal/domain"
)

type pollRepository struct {
	db *sql.DB
}

func NewPollRepository(db *sql.DB) PollRepository {
	return &pollRepository{db: db}
}

func (r *pollRepository) Create(poll *domain.Poll) error {
	_, err := r.db.Exec(
		`INSERT INTO polls (id, chat_id, creator_id, question, options, is_anonymous, multiple_choice, expires_at, created_at, closed)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		poll.ID, poll.ChatID, poll.CreatorID, poll.Question, poll.Options,
		boolToInt(poll.IsAnonymous), boolToInt(poll.MultipleChoice),
		poll.ExpiresAt, poll.CreatedAt.Format(time.RFC3339), boolToInt(poll.Closed),
	)
	return err
}

func (r *pollRepository) FindByID(id string) (*domain.Poll, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, creator_id, question, options, is_anonymous, multiple_choice, expires_at, created_at, closed
		FROM polls WHERE id = ?`, id,
	)
	var p domain.Poll
	var expiresAt, createdAt sql.NullString
	var isAnon, multi, closed int
	if err := row.Scan(&p.ID, &p.ChatID, &p.CreatorID, &p.Question, &p.Options,
		&isAnon, &multi, &expiresAt, &createdAt, &closed); err != nil {
		return nil, err
	}
	p.IsAnonymous = isAnon == 1
	p.MultipleChoice = multi == 1
	p.Closed = closed == 1
	if expiresAt.Valid {
		p.ExpiresAt = &expiresAt.String
	}
	if createdAt.Valid {
		p.CreatedAt = parseTime(createdAt.String)
	}
	return &p, nil
}

func (r *pollRepository) FindByChatID(chatID string) ([]*domain.Poll, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, creator_id, question, options, is_anonymous, multiple_choice, expires_at, created_at, closed
		FROM polls WHERE chat_id = ? ORDER BY created_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	polls := make([]*domain.Poll, 0)
	for rows.Next() {
		var p domain.Poll
		var expiresAt, createdAt sql.NullString
		var isAnon, multi, closed int
		if err := rows.Scan(&p.ID, &p.ChatID, &p.CreatorID, &p.Question, &p.Options,
			&isAnon, &multi, &expiresAt, &createdAt, &closed); err != nil {
			return nil, err
		}
		p.IsAnonymous = isAnon == 1
		p.MultipleChoice = multi == 1
		p.Closed = closed == 1
		if expiresAt.Valid {
			p.ExpiresAt = &expiresAt.String
		}
		if createdAt.Valid {
			p.CreatedAt = parseTime(createdAt.String)
		}
		polls = append(polls, &p)
	}
	return polls, nil
}

func (r *pollRepository) Update(poll *domain.Poll) error {
	_, err := r.db.Exec(
		`UPDATE polls SET question=?, options=?, closed=? WHERE id=?`,
		poll.Question, poll.Options, boolToInt(poll.Closed), poll.ID,
	)
	return err
}

func (r *pollRepository) AddVote(vote *domain.PollVote) error {
	_, err := r.db.Exec(
		`INSERT INTO poll_votes (poll_id, user_id, option_index, voted_at) VALUES (?, ?, ?, ?)`,
		vote.PollID, vote.UserID, vote.OptionIndex, vote.VotedAt.Format(time.RFC3339),
	)
	return err
}

func (r *pollRepository) HasVoted(pollID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM poll_votes WHERE poll_id = ? AND user_id = ?`, pollID, userID,
	).Scan(&count)
	return count > 0, err
}

func (r *pollRepository) GetVoteCount(pollID string, optionIndex int) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM poll_votes WHERE poll_id = ? AND option_index = ?`, pollID, optionIndex,
	).Scan(&count)
	return count, err
}

func (r *pollRepository) GetUserVote(pollID, userID string) (*domain.PollVote, error) {
	row := r.db.QueryRow(
		`SELECT poll_id, user_id, option_index, voted_at FROM poll_votes WHERE poll_id = ? AND user_id = ?`,
		pollID, userID,
	)
	var v domain.PollVote
	var votedAt string
	if err := row.Scan(&v.PollID, &v.UserID, &v.OptionIndex, &votedAt); err != nil {
		return nil, err
	}
	v.VotedAt = parseTime(votedAt)
	return &v, nil
}

func (r *pollRepository) GetTotalVotes(pollID string) (int, error) {
	var count int
	err := r.db.QueryRow(
		`SELECT COUNT(*) FROM poll_votes WHERE poll_id = ?`, pollID,
	).Scan(&count)
	return count, err
}

func (r *pollRepository) GetAllVotes(pollID string) ([]*domain.PollVote, error) {
	rows, err := r.db.Query(
		`SELECT poll_id, user_id, option_index, voted_at FROM poll_votes WHERE poll_id = ?`,
		pollID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	votes := make([]*domain.PollVote, 0)
	for rows.Next() {
		var v domain.PollVote
		var votedAt string
		if err := rows.Scan(&v.PollID, &v.UserID, &v.OptionIndex, &votedAt); err != nil {
			return nil, err
		}
		v.VotedAt = parseTime(votedAt)
		votes = append(votes, &v)
	}
	return votes, nil
}
