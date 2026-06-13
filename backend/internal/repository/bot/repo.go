package botrepo

import (
	"database/sql"
	"time"

	botdomain "ChatServerGolang/backend/internal/domain/bot"
	"ChatServerGolang/backend/internal/repository"
)

type botRepository struct {
	db *sql.DB
}

func NewBotRepository(db *sql.DB) repository.BotRepository {
	return &botRepository{db: db}
}

func (r *botRepository) Create(bot *botdomain.Bot) error {
	_, err := r.db.Exec(
		`INSERT INTO bots (id, token, owner_id, name, avatar_url, webhook_url, created_at, active)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
		bot.ID, bot.Token, bot.OwnerID, bot.Name, bot.AvatarURL, bot.WebhookURL,
		bot.CreatedAt.Format(time.RFC3339), repository.BoolToInt(bot.Active),
	)
	return err
}

func (r *botRepository) FindByID(id string) (*botdomain.Bot, error) {
	row := r.db.QueryRow(
		`SELECT id, token, owner_id, name, avatar_url, webhook_url, created_at, active FROM bots WHERE id = ?`, id,
	)
	return scanBot(row)
}

func (r *botRepository) FindByOwnerID(ownerID string) ([]*botdomain.Bot, error) {
	rows, err := r.db.Query(
		`SELECT id, token, owner_id, name, avatar_url, webhook_url, created_at, active
		FROM bots WHERE owner_id = ? ORDER BY created_at DESC`, ownerID,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	bots := make([]*botdomain.Bot, 0)
	for rows.Next() {
		bot, err := scanBot(rows)
		if err != nil {
			return nil, err
		}
		bots = append(bots, bot)
	}
	return bots, nil
}

func (r *botRepository) Update(bot *botdomain.Bot) error {
	_, err := r.db.Exec(
		`UPDATE bots SET name=?, avatar_url=?, webhook_url=?, active=? WHERE id=?`,
		bot.Name, bot.AvatarURL, bot.WebhookURL, repository.BoolToInt(bot.Active), bot.ID,
	)
	return err
}

func (r *botRepository) RegenerateToken(id, token string) error {
	_, err := r.db.Exec(`UPDATE bots SET token=? WHERE id=?`, token, id)
	return err
}

func (r *botRepository) Delete(id string) error {
	_, err := r.db.Exec(`DELETE FROM bots WHERE id = ?`, id)
	return err
}

func (r *botRepository) FindByToken(token string) (*botdomain.Bot, error) {
	row := r.db.QueryRow(
		`SELECT id, token, owner_id, name, avatar_url, webhook_url, created_at, active
		FROM bots WHERE token = ?`, token,
	)
	return scanBot(row)
}

type botScanner interface {
	Scan(dest ...interface{}) error
}

func scanBot(row botScanner) (*botdomain.Bot, error) {
	var b botdomain.Bot
	var createdAt string
	var active int
	if err := row.Scan(&b.ID, &b.Token, &b.OwnerID, &b.Name, &b.AvatarURL, &b.WebhookURL, &createdAt, &active); err != nil {
		return nil, err
	}
	b.Active = active == 1
	b.CreatedAt = repository.ParseTime(createdAt)
	return &b, nil
}



