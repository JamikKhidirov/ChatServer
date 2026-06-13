package voicechatrepo

import (
	"database/sql"
	"time"

	voicechatdomain "ChatServerGolang/backend/internal/domain/voicechat"
	"ChatServerGolang/backend/internal/repository"
)

type voiceChatRepository struct {
	db *sql.DB
}

func NewVoiceChatRepository(db *sql.DB) repository.VoiceChatRepository {
	return &voiceChatRepository{db: db}
}

func (r *voiceChatRepository) Create(vc *voicechatdomain.VoiceChat) error {
	_, err := r.db.Exec(
		`INSERT INTO voice_chats (id, chat_id, started_by, title, status, participant_count, scheduled_at, started_at, ended_at, created_at)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		vc.ID, vc.ChatID, vc.StartedBy, vc.Title, string(vc.Status), vc.ParticipantCount,
		nullTime(vc.ScheduledAt), nullTime(vc.StartedAt), nullTime(vc.EndedAt),
		vc.CreatedAt.Format(time.RFC3339),
	)
	return err
}

func (r *voiceChatRepository) FindByID(id string) (*voicechatdomain.VoiceChat, error) {
	row := r.db.QueryRow(
		`SELECT id, chat_id, started_by, title, status, participant_count, scheduled_at, started_at, ended_at, created_at
		FROM voice_chats WHERE id = ?`, id,
	)
	return scanVoiceChat(row)
}

func (r *voiceChatRepository) FindActiveByChatID(chatID string) ([]*voicechatdomain.VoiceChat, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, started_by, title, status, participant_count, scheduled_at, started_at, ended_at, created_at
		FROM voice_chats WHERE chat_id = ? AND status = 'active' ORDER BY created_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVoiceChats(rows)
}

func (r *voiceChatRepository) FindByChatID(chatID string) ([]*voicechatdomain.VoiceChat, error) {
	rows, err := r.db.Query(
		`SELECT id, chat_id, started_by, title, status, participant_count, scheduled_at, started_at, ended_at, created_at
		FROM voice_chats WHERE chat_id = ? ORDER BY created_at DESC`, chatID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	return scanVoiceChats(rows)
}

func (r *voiceChatRepository) UpdateStatus(id string, status voicechatdomain.VoiceChatStatus) error {
	_, err := r.db.Exec(`UPDATE voice_chats SET status = ? WHERE id = ?`, string(status), id)
	return err
}

func (r *voiceChatRepository) AddParticipant(vcID, userID string) error {
	_, err := r.db.Exec(
		`INSERT OR IGNORE INTO voice_chat_participants (voice_chat_id, user_id, joined_at, muted) VALUES (?, ?, ?, ?)`,
		vcID, userID, time.Now().Format(time.RFC3339), 0,
	)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE voice_chats SET participant_count = (SELECT COUNT(*) FROM voice_chat_participants WHERE voice_chat_id = ?) WHERE id = ?`, vcID, vcID)
	return err
}

func (r *voiceChatRepository) RemoveParticipant(vcID, userID string) error {
	_, err := r.db.Exec(`DELETE FROM voice_chat_participants WHERE voice_chat_id = ? AND user_id = ?`, vcID, userID)
	if err != nil {
		return err
	}
	_, err = r.db.Exec(`UPDATE voice_chats SET participant_count = (SELECT COUNT(*) FROM voice_chat_participants WHERE voice_chat_id = ?) WHERE id = ?`, vcID, vcID)
	return err
}

func (r *voiceChatRepository) IsParticipant(vcID, userID string) (bool, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM voice_chat_participants WHERE voice_chat_id = ? AND user_id = ?`, vcID, userID).Scan(&count)
	return count > 0, err
}

func (r *voiceChatRepository) GetParticipants(vcID string) ([]*voicechatdomain.VoiceChatParticipant, error) {
	rows, err := r.db.Query(
		`SELECT voice_chat_id, user_id, joined_at, left_at, muted FROM voice_chat_participants WHERE voice_chat_id = ? ORDER BY joined_at ASC`,
		vcID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var participants []*voicechatdomain.VoiceChatParticipant
	for rows.Next() {
		var p voicechatdomain.VoiceChatParticipant
		var joinedAt string
		var leftAt *string
		if err := rows.Scan(&p.VoiceChatID, &p.UserID, &joinedAt, &leftAt, &p.Muted); err != nil {
			return nil, err
		}
		p.JoinedAt = repository.ParseTime(joinedAt)
		if leftAt != nil {
			t := repository.ParseTime(*leftAt)
			p.LeftAt = &t
		}
		participants = append(participants, &p)
	}
	return participants, nil
}

func (r *voiceChatRepository) GetParticipantCount(vcID string) (int, error) {
	var count int
	err := r.db.QueryRow(`SELECT COUNT(*) FROM voice_chat_participants WHERE voice_chat_id = ?`, vcID).Scan(&count)
	return count, err
}

func (r *voiceChatRepository) SetParticipantMuted(vcID, userID string, muted bool) error {
	_, err := r.db.Exec(`UPDATE voice_chat_participants SET muted = ? WHERE voice_chat_id = ? AND user_id = ?`, repository.BoolToInt(muted), vcID, userID)
	return err
}

func scanVoiceChat(row repository.Scanner) (*voicechatdomain.VoiceChat, error) {
	var v voicechatdomain.VoiceChat
	var createdAt, scheduledAt, startedAt, endedAt *string
	err := row.Scan(&v.ID, &v.ChatID, &v.StartedBy, &v.Title, &v.Status, &v.ParticipantCount, &scheduledAt, &startedAt, &endedAt, &createdAt)
	if err != nil {
		return nil, err
	}
	if createdAt != nil {
		v.CreatedAt = repository.ParseTime(*createdAt)
	}
	if scheduledAt != nil {
		t := repository.ParseTime(*scheduledAt)
		v.ScheduledAt = &t
	}
	if startedAt != nil {
		t := repository.ParseTime(*startedAt)
		v.StartedAt = &t
	}
	if endedAt != nil {
		t := repository.ParseTime(*endedAt)
		v.EndedAt = &t
	}
	return &v, nil
}

func scanVoiceChats(rows *sql.Rows) ([]*voicechatdomain.VoiceChat, error) {
	var vcs []*voicechatdomain.VoiceChat
	for rows.Next() {
		v, err := scanVoiceChat(rows)
		if err != nil {
			return nil, err
		}
		vcs = append(vcs, v)
	}
	return vcs, nil
}

func nullTime(t *time.Time) *string {
	if t == nil {
		return nil
	}
	s := t.Format(time.RFC3339)
	return &s
}
