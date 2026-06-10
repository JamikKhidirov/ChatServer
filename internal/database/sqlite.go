package database

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

func NewDB(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite", path)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	if _, err := db.Exec("PRAGMA journal_mode=WAL"); err != nil {
		log.Printf("Warning: could not enable WAL mode: %v", err)
	}
	if _, err := db.Exec("PRAGMA foreign_keys=ON"); err != nil {
		log.Printf("Warning: could not enable foreign keys: %v", err)
	}
	if _, err := db.Exec("PRAGMA busy_timeout=5000"); err != nil {
		log.Printf("Warning: could not set busy timeout: %v", err)
	}

	if err := runMigrations(db); err != nil {
		return nil, err
	}

	log.Println("Database connected and migrated")
	return db, nil
}

func runMigrations(db *sql.DB) error {
	queries := []string{
		`CREATE TABLE IF NOT EXISTS users (
			id TEXT PRIMARY KEY,
			username TEXT UNIQUE NOT NULL,
			email TEXT UNIQUE NOT NULL,
			password_hash TEXT NOT NULL,
			display_name TEXT NOT NULL DEFAULT '',
			avatar_url TEXT NOT NULL DEFAULT '',
			bio TEXT NOT NULL DEFAULT '',
			user_status TEXT NOT NULL DEFAULT 'Available',
			push_token TEXT NOT NULL DEFAULT '',
			push_provider TEXT NOT NULL DEFAULT '',
			online INTEGER NOT NULL DEFAULT 0,
			deleted INTEGER NOT NULL DEFAULT 0,
			last_seen TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
		)`,
		`CREATE TABLE IF NOT EXISTS chats (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL DEFAULT '',
			description TEXT NOT NULL DEFAULT '',
			avatar_url TEXT NOT NULL DEFAULT '',
			type TEXT NOT NULL DEFAULT 'private',
			created_by TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (created_by) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS chat_participants (
			chat_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			role TEXT NOT NULL DEFAULT 'member',
			joined_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			last_read_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (chat_id, user_id),
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS messages (
			id TEXT PRIMARY KEY,
			chat_id TEXT NOT NULL,
			sender_id TEXT NOT NULL,
			content TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'text',
			reply_to_id TEXT,
			forward_from TEXT,
			file_name TEXT NOT NULL DEFAULT '',
			file_size INTEGER NOT NULL DEFAULT 0,
			file_path TEXT NOT NULL DEFAULT '',
			pinned INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			deleted_at TEXT,
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (sender_id) REFERENCES users(id),
			FOREIGN KEY (reply_to_id) REFERENCES messages(id)
		)`,
		`CREATE TABLE IF NOT EXISTS calls (
			id TEXT PRIMARY KEY,
			chat_id TEXT NOT NULL,
			caller_id TEXT NOT NULL,
			callee_id TEXT NOT NULL,
			status TEXT NOT NULL DEFAULT 'initiated',
			started_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			ended_at TEXT,
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (caller_id) REFERENCES users(id),
			FOREIGN KEY (callee_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS reactions (
			message_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			emoji TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (message_id, user_id, emoji),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS blocked_users (
			user_id TEXT NOT NULL,
			blocked_id TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, blocked_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (blocked_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS notification_settings (
			user_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			muted INTEGER NOT NULL DEFAULT 0,
			PRIMARY KEY (user_id, chat_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (chat_id) REFERENCES chats(id)
		)`,
		`CREATE TABLE IF NOT EXISTS read_receipts (
			message_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			read_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (message_id, user_id),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS account_settings (
			user_id TEXT PRIMARY KEY,
			language TEXT NOT NULL DEFAULT 'en',
			theme TEXT NOT NULL DEFAULT 'light',
			notifications INTEGER NOT NULL DEFAULT 1,
			sound_enabled INTEGER NOT NULL DEFAULT 1,
			last_seen_mode TEXT NOT NULL DEFAULT 'everyone',
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS contacts (
			user_id TEXT NOT NULL,
			phone TEXT NOT NULL,
			name TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, phone),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS hidden_chats (
			user_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			hidden_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, chat_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (chat_id) REFERENCES chats(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at)`,
		`CREATE INDEX IF NOT EXISTS idx_chat_participants_user_id ON chat_participants(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_calls_chat_id ON calls(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_calls_callee_id ON calls(callee_id)`,
		`CREATE INDEX IF NOT EXISTS idx_reactions_message_id ON reactions(message_id)`,
		`CREATE INDEX IF NOT EXISTS idx_blocked_users_user_id ON blocked_users(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_blocked_users_blocked_id ON blocked_users(blocked_id)`,
		`CREATE INDEX IF NOT EXISTS idx_read_receipts_message_id ON read_receipts(message_id)`,
		`CREATE INDEX IF NOT EXISTS idx_contacts_user_id ON contacts(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_hidden_chats_user_id ON hidden_chats(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_messages_chat_id_created_at ON messages(chat_id, created_at)`,
		`CREATE INDEX IF NOT EXISTS idx_notification_settings_user_id ON notification_settings(user_id, chat_id)`,
		`CREATE TABLE IF NOT EXISTS pinned_chats (
			user_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, chat_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (chat_id) REFERENCES chats(id)
		)`,
		`CREATE TABLE IF NOT EXISTS archived_chats (
			user_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			archived_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, chat_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (chat_id) REFERENCES chats(id)
		)`,
		`CREATE TABLE IF NOT EXISTS starred_messages (
			user_id TEXT NOT NULL,
			message_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, message_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (chat_id) REFERENCES chats(id)
		)`,
		`CREATE TABLE IF NOT EXISTS deleted_messages (
			user_id TEXT NOT NULL,
			message_id TEXT NOT NULL,
			deleted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, message_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (message_id) REFERENCES messages(id)
		)`,
		`CREATE INDEX IF NOT EXISTS idx_pinned_chats_user_id ON pinned_chats(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_archived_chats_user_id ON archived_chats(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_starred_messages_user_id ON starred_messages(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_deleted_messages_user_id ON deleted_messages(user_id)`,

		// Polls
		`CREATE TABLE IF NOT EXISTS polls (
			id TEXT PRIMARY KEY,
			chat_id TEXT NOT NULL,
			creator_id TEXT NOT NULL,
			question TEXT NOT NULL,
			options TEXT NOT NULL,
			is_anonymous INTEGER NOT NULL DEFAULT 0,
			multiple_choice INTEGER NOT NULL DEFAULT 0,
			expires_at TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			closed INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (creator_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS poll_votes (
			poll_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			option_index INTEGER NOT NULL,
			voted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (poll_id, user_id),
			FOREIGN KEY (poll_id) REFERENCES polls(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Sticker packs
		`CREATE TABLE IF NOT EXISTS sticker_packs (
			id TEXT PRIMARY KEY,
			name TEXT NOT NULL,
			creator_id TEXT NOT NULL,
			animated INTEGER NOT NULL DEFAULT 0,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (creator_id) REFERENCES users(id)
		)`,
		`CREATE TABLE IF NOT EXISTS stickers (
			id TEXT PRIMARY KEY,
			pack_id TEXT NOT NULL,
			emoji TEXT NOT NULL DEFAULT '',
			image_url TEXT NOT NULL DEFAULT '',
			file_path TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (pack_id) REFERENCES sticker_packs(id)
		)`,
		`CREATE TABLE IF NOT EXISTS user_stickers (
			user_id TEXT NOT NULL,
			sticker_id TEXT NOT NULL,
			PRIMARY KEY (user_id, sticker_id),
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (sticker_id) REFERENCES stickers(id)
		)`,

		// Drafts
		`CREATE TABLE IF NOT EXISTS drafts (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			chat_id TEXT NOT NULL,
			content TEXT NOT NULL DEFAULT '',
			reply_to_id TEXT,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id),
			FOREIGN KEY (chat_id) REFERENCES chats(id)
		)`,

		// Scheduled messages
		`CREATE TABLE IF NOT EXISTS scheduled_messages (
			id TEXT PRIMARY KEY,
			chat_id TEXT NOT NULL,
			sender_id TEXT NOT NULL,
			content TEXT NOT NULL,
			type TEXT NOT NULL DEFAULT 'text',
			reply_to_id TEXT,
			scheduled_at TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			sent INTEGER NOT NULL DEFAULT 0,
			FOREIGN KEY (chat_id) REFERENCES chats(id),
			FOREIGN KEY (sender_id) REFERENCES users(id)
		)`,

		// Sessions
		`CREATE TABLE IF NOT EXISTS sessions (
			id TEXT PRIMARY KEY,
			user_id TEXT NOT NULL,
			device_name TEXT NOT NULL DEFAULT '',
			ip_address TEXT NOT NULL DEFAULT '',
			last_active TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Bots
		`CREATE TABLE IF NOT EXISTS bots (
			id TEXT PRIMARY KEY,
			token TEXT NOT NULL,
			owner_id TEXT NOT NULL,
			name TEXT NOT NULL,
			avatar_url TEXT NOT NULL DEFAULT '',
			webhook_url TEXT NOT NULL DEFAULT '',
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			active INTEGER NOT NULL DEFAULT 1,
			FOREIGN KEY (owner_id) REFERENCES users(id)
		)`,

		// Mentions
		`CREATE TABLE IF NOT EXISTS mentions (
			message_id TEXT NOT NULL,
			user_id TEXT NOT NULL,
			username TEXT NOT NULL,
			PRIMARY KEY (message_id, user_id),
			FOREIGN KEY (message_id) REFERENCES messages(id),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Saved GIFs
		`CREATE TABLE IF NOT EXISTS saved_gifs (
			user_id TEXT NOT NULL,
			gif_url TEXT NOT NULL,
			created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
			PRIMARY KEY (user_id, gif_url),
			FOREIGN KEY (user_id) REFERENCES users(id)
		)`,

		// Indexes for new tables
		`CREATE INDEX IF NOT EXISTS idx_polls_chat_id ON polls(chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_poll_votes_poll_id ON poll_votes(poll_id)`,
		`CREATE INDEX IF NOT EXISTS idx_sticker_packs_creator ON sticker_packs(creator_id)`,
		`CREATE INDEX IF NOT EXISTS idx_stickers_pack_id ON stickers(pack_id)`,
		`CREATE INDEX IF NOT EXISTS idx_drafts_user_chat ON drafts(user_id, chat_id)`,
		`CREATE INDEX IF NOT EXISTS idx_scheduled_messages_sent ON scheduled_messages(sent, scheduled_at)`,
		`CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id)`,
		`CREATE INDEX IF NOT EXISTS idx_bots_owner_id ON bots(owner_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mentions_message_id ON mentions(message_id)`,
		`CREATE INDEX IF NOT EXISTS idx_mentions_user_id ON mentions(user_id)`,
	}

	for _, q := range queries {
		if _, err := db.Exec(q); err != nil {
			return err
		}
	}

	legacy := []string{
		`ALTER TABLE chats ADD COLUMN description TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN user_status TEXT NOT NULL DEFAULT 'Available'`,
		`ALTER TABLE messages ADD COLUMN forward_from TEXT`,
		`ALTER TABLE messages ADD COLUMN file_name TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE messages ADD COLUMN file_size INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE messages ADD COLUMN file_path TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE messages ADD COLUMN pinned INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN deleted INTEGER NOT NULL DEFAULT 0`,
		`ALTER TABLE users ADD COLUMN phone TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN gender TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE users ADD COLUMN date_of_birth TEXT NOT NULL DEFAULT ''`,
		`ALTER TABLE calls ADD COLUMN call_type TEXT NOT NULL DEFAULT 'audio'`,
	}
	for _, m := range legacy {
		db.Exec(m)
	}

	return nil
}
