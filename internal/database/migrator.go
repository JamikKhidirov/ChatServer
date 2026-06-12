package database

import (
	"database/sql"
	"fmt"
	"log"
	"time"
)

type Migration struct {
	ID      int
	Name    string
	Up      string
}

var migrations = []Migration{
	{1, "initial_schema", initialSchema},
	{2, "add_description_to_chats", addDescriptionToChats},
	{3, "add_user_fields", addUserFields},
	{4, "add_forward_file_pinned", addForwardFilePinned},
	{5, "add_phone_gender_dob", addPhoneGenderDob},
	{6, "add_call_type", addCallType},
	{7, "add_pinned_archived_starred_deleted", addPinnedArchivedStarredDeleted},
	{8, "add_polls_stickers_drafts_scheduled", addPollsStickersDraftsScheduled},
	{9, "add_sessions_bots_mentions_gifs", addSessionsBotsMentionsGifs},
	{10, "add_verification_reports_bookmarks", addVerificationReportsBookmarks},
	{11, "add_e2e_keys_self_destruct_link_previews", addE2EKeysSelfDestructLinkPreviews},
	{12, "add_edit_history_captcha_ip_blocks", addEditHistoryCaptchaIPBlocks},
	{13, "add_admin_settings", addAdminSettings},
	{14, "add_login_codes", addLoginCodes},
	{15, "add_contact_photo", addContactPhoto},
	{16, "add_message_fields", addMessageFields},
	{17, "add_invite_links", addInviteLinks},
	{18, "add_chat_folders", addChatFolders},
	{19, "add_chat_slow_mode", addChatSlowMode},
	{20, "add_chat_themes", addChatThemes},
	{21, "add_is_admin_field", addIsAdminField},
}

func RunMigrations(db *sql.DB) error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS schema_migrations (
		version INTEGER PRIMARY KEY,
		name TEXT NOT NULL,
		applied_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
	)`)
	if err != nil {
		return fmt.Errorf("failed to create migrations table: %w", err)
	}

	var currentVersion int
	err = db.QueryRow(`SELECT COALESCE(MAX(version), 0) FROM schema_migrations`).Scan(&currentVersion)
	if err != nil {
		return fmt.Errorf("failed to get current version: %w", err)
	}

	for _, m := range migrations {
		if m.ID <= currentVersion {
			continue
		}
		log.Printf("Running migration %d: %s", m.ID, m.Name)

		// Try to run migration in a transaction first, then fallback to direct execution for DDL
		tx, txErr := db.Begin()
		if txErr == nil {
			if _, execErr := tx.Exec(m.Up); execErr == nil {
				// Record migration in the same transaction
				if _, insErr := tx.Exec(
					`INSERT INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)`,
					m.ID, m.Name, time.Now().Format(time.RFC3339),
				); insErr == nil {
					if commitErr := tx.Commit(); commitErr == nil {
						log.Printf("Migration %d applied successfully", m.ID)
						continue
					}
				}
			}
			tx.Rollback()
		}

		// Fallback: try directly (for DDL that can't run in transaction)
		log.Printf("Migration %d: trying direct execution", m.ID)
		db.Exec(m.Up)

		// Record migration in a new transaction
		if recTx, err := db.Begin(); err == nil {
			if _, err := recTx.Exec(
				`INSERT OR IGNORE INTO schema_migrations (version, name, applied_at) VALUES (?, ?, ?)`,
				m.ID, m.Name, time.Now().Format(time.RFC3339),
			); err == nil {
				recTx.Commit()
				log.Printf("Migration %d recorded", m.ID)
			} else {
				recTx.Rollback()
			}
		}
	}

	return nil
}

const initialSchema = `
CREATE TABLE IF NOT EXISTS users (
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
);
CREATE TABLE IF NOT EXISTS chats (
	id TEXT PRIMARY KEY,
	name TEXT NOT NULL DEFAULT '',
	description TEXT NOT NULL DEFAULT '',
	avatar_url TEXT NOT NULL DEFAULT '',
	type TEXT NOT NULL DEFAULT 'private',
	created_by TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (created_by) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS chat_participants (
	chat_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'member',
	joined_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	last_read_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (chat_id, user_id),
	FOREIGN KEY (chat_id) REFERENCES chats(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS messages (
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
);
CREATE TABLE IF NOT EXISTS calls (
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
);
CREATE TABLE IF NOT EXISTS reactions (
	message_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	emoji TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (message_id, user_id, emoji),
	FOREIGN KEY (message_id) REFERENCES messages(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS blocked_users (
	user_id TEXT NOT NULL,
	blocked_id TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, blocked_id),
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (blocked_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS notification_settings (
	user_id TEXT NOT NULL,
	chat_id TEXT NOT NULL,
	muted INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (user_id, chat_id),
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS read_receipts (
	message_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	read_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (message_id, user_id),
	FOREIGN KEY (message_id) REFERENCES messages(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS account_settings (
	user_id TEXT PRIMARY KEY,
	language TEXT NOT NULL DEFAULT 'en',
	theme TEXT NOT NULL DEFAULT 'light',
	notifications INTEGER NOT NULL DEFAULT 1,
	sound_enabled INTEGER NOT NULL DEFAULT 1,
	last_seen_mode TEXT NOT NULL DEFAULT 'everyone',
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS contacts (
	user_id TEXT NOT NULL,
	phone TEXT NOT NULL,
	name TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, phone),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS hidden_chats (
	user_id TEXT NOT NULL,
	chat_id TEXT NOT NULL,
	hidden_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, chat_id),
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE INDEX IF NOT EXISTS idx_messages_chat_id ON messages(chat_id);
CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_deleted_at ON messages(deleted_at);
CREATE INDEX IF NOT EXISTS idx_chat_participants_user_id ON chat_participants(user_id);
CREATE INDEX IF NOT EXISTS idx_calls_chat_id ON calls(chat_id);
CREATE INDEX IF NOT EXISTS idx_calls_callee_id ON calls(callee_id);
CREATE INDEX IF NOT EXISTS idx_reactions_message_id ON reactions(message_id);
CREATE INDEX IF NOT EXISTS idx_blocked_users_user_id ON blocked_users(user_id);
CREATE INDEX IF NOT EXISTS idx_blocked_users_blocked_id ON blocked_users(blocked_id);
CREATE INDEX IF NOT EXISTS idx_read_receipts_message_id ON read_receipts(message_id);
CREATE INDEX IF NOT EXISTS idx_contacts_user_id ON contacts(user_id);
CREATE INDEX IF NOT EXISTS idx_hidden_chats_user_id ON hidden_chats(user_id);
CREATE INDEX IF NOT EXISTS idx_messages_chat_id_created_at ON messages(chat_id, created_at);
CREATE INDEX IF NOT EXISTS idx_notification_settings_user_id ON notification_settings(user_id, chat_id);
`

const addDescriptionToChats = `ALTER TABLE chats ADD COLUMN description TEXT NOT NULL DEFAULT ''`
const addUserFields = `ALTER TABLE users ADD COLUMN user_status TEXT NOT NULL DEFAULT 'Available'`
const addForwardFilePinned = `
ALTER TABLE messages ADD COLUMN forward_from TEXT;
ALTER TABLE messages ADD COLUMN file_name TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN file_size INTEGER NOT NULL DEFAULT 0;
ALTER TABLE messages ADD COLUMN file_path TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN pinned INTEGER NOT NULL DEFAULT 0;
`
const addPhoneGenderDob = `
ALTER TABLE users ADD COLUMN phone TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN gender TEXT NOT NULL DEFAULT '';
ALTER TABLE users ADD COLUMN date_of_birth TEXT NOT NULL DEFAULT '';
`
const addCallType = `ALTER TABLE calls ADD COLUMN call_type TEXT NOT NULL DEFAULT 'audio'`
const addPinnedArchivedStarredDeleted = `
CREATE TABLE IF NOT EXISTS pinned_chats (
	user_id TEXT NOT NULL, chat_id TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, chat_id), FOREIGN KEY (user_id) REFERENCES users(id), FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS archived_chats (
	user_id TEXT NOT NULL, chat_id TEXT NOT NULL, archived_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, chat_id), FOREIGN KEY (user_id) REFERENCES users(id), FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS starred_messages (
	user_id TEXT NOT NULL, message_id TEXT NOT NULL, chat_id TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, message_id), FOREIGN KEY (user_id) REFERENCES users(id), FOREIGN KEY (message_id) REFERENCES messages(id), FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS deleted_messages (
	user_id TEXT NOT NULL, message_id TEXT NOT NULL, deleted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (user_id, message_id), FOREIGN KEY (user_id) REFERENCES users(id), FOREIGN KEY (message_id) REFERENCES messages(id)
);
`
const addPollsStickersDraftsScheduled = `
CREATE TABLE IF NOT EXISTS polls (
	id TEXT PRIMARY KEY, chat_id TEXT NOT NULL, creator_id TEXT NOT NULL, question TEXT NOT NULL,
	options TEXT NOT NULL, is_anonymous INTEGER NOT NULL DEFAULT 0, multiple_choice INTEGER NOT NULL DEFAULT 0,
	expires_at TEXT, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, closed INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY (chat_id) REFERENCES chats(id), FOREIGN KEY (creator_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS poll_votes (
	poll_id TEXT NOT NULL, user_id TEXT NOT NULL, option_index INTEGER NOT NULL, voted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (poll_id, user_id), FOREIGN KEY (poll_id) REFERENCES polls(id), FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS sticker_packs (
	id TEXT PRIMARY KEY, name TEXT NOT NULL, creator_id TEXT NOT NULL, animated INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (creator_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS stickers (
	id TEXT PRIMARY KEY, pack_id TEXT NOT NULL, emoji TEXT NOT NULL DEFAULT '', image_url TEXT NOT NULL DEFAULT '',
	file_path TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (pack_id) REFERENCES sticker_packs(id)
);
CREATE TABLE IF NOT EXISTS user_stickers (user_id TEXT NOT NULL, sticker_id TEXT NOT NULL, PRIMARY KEY (user_id, sticker_id));
CREATE TABLE IF NOT EXISTS drafts (
	id TEXT PRIMARY KEY, user_id TEXT NOT NULL, chat_id TEXT NOT NULL, content TEXT NOT NULL DEFAULT '',
	reply_to_id TEXT, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id), FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS scheduled_messages (
	id TEXT PRIMARY KEY, chat_id TEXT NOT NULL, sender_id TEXT NOT NULL, content TEXT NOT NULL,
	type TEXT NOT NULL DEFAULT 'text', reply_to_id TEXT, scheduled_at TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, sent INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY (chat_id) REFERENCES chats(id), FOREIGN KEY (sender_id) REFERENCES users(id)
);
`
const addSessionsBotsMentionsGifs = `
CREATE TABLE IF NOT EXISTS sessions (
	id TEXT PRIMARY KEY, user_id TEXT NOT NULL, device_name TEXT NOT NULL DEFAULT '',
	ip_address TEXT NOT NULL DEFAULT '', last_active TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS bots (
	id TEXT PRIMARY KEY, token TEXT NOT NULL, owner_id TEXT NOT NULL, name TEXT NOT NULL,
	avatar_url TEXT NOT NULL DEFAULT '', webhook_url TEXT NOT NULL DEFAULT '',
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, active INTEGER NOT NULL DEFAULT 1,
	FOREIGN KEY (owner_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS mentions (message_id TEXT NOT NULL, user_id TEXT NOT NULL, username TEXT NOT NULL, PRIMARY KEY (message_id, user_id));
CREATE TABLE IF NOT EXISTS saved_gifs (user_id TEXT NOT NULL, gif_url TEXT NOT NULL, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, PRIMARY KEY (user_id, gif_url));
CREATE INDEX IF NOT EXISTS idx_polls_chat_id ON polls(chat_id);
CREATE INDEX IF NOT EXISTS idx_poll_votes_poll_id ON poll_votes(poll_id);
CREATE INDEX IF NOT EXISTS idx_sticker_packs_creator ON sticker_packs(creator_id);
CREATE INDEX IF NOT EXISTS idx_stickers_pack_id ON stickers(pack_id);
CREATE INDEX IF NOT EXISTS idx_drafts_user_chat ON drafts(user_id, chat_id);
CREATE INDEX IF NOT EXISTS idx_scheduled_messages_sent ON scheduled_messages(sent, scheduled_at);
CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_bots_owner_id ON bots(owner_id);
CREATE INDEX IF NOT EXISTS idx_mentions_message_id ON mentions(message_id);
CREATE INDEX IF NOT EXISTS idx_mentions_user_id ON mentions(user_id);
`
const addVerificationReportsBookmarks = `
CREATE TABLE IF NOT EXISTS email_verifications (
	id TEXT PRIMARY KEY, user_id TEXT NOT NULL, email TEXT NOT NULL, code TEXT NOT NULL,
	expires_at TEXT NOT NULL, verified INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS phone_verifications (
	id TEXT PRIMARY KEY, user_id TEXT NOT NULL, phone TEXT NOT NULL, code TEXT NOT NULL,
	expires_at TEXT NOT NULL, verified INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS message_reports (
	id TEXT PRIMARY KEY, message_id TEXT NOT NULL, reporter_id TEXT NOT NULL, reason TEXT NOT NULL,
	description TEXT NOT NULL DEFAULT '', status TEXT NOT NULL DEFAULT 'pending',
	resolved_by TEXT, resolved_at TEXT, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (message_id) REFERENCES messages(id), FOREIGN KEY (reporter_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS bookmarks (
	id TEXT PRIMARY KEY, user_id TEXT NOT NULL, message_id TEXT NOT NULL, chat_id TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id), FOREIGN KEY (message_id) REFERENCES messages(id)
);
`
const addE2EKeysSelfDestructLinkPreviews = `
CREATE TABLE IF NOT EXISTS e2e_keys (
	user_id TEXT PRIMARY KEY, public_key TEXT NOT NULL, private_key_encrypted TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS e2e_admin_key (
	id INTEGER PRIMARY KEY DEFAULT 1, admin_public_key TEXT NOT NULL, admin_private_key TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS message_self_destruct (
	message_id TEXT PRIMARY KEY, chat_id TEXT NOT NULL, delete_at TEXT NOT NULL,
	FOREIGN KEY (message_id) REFERENCES messages(id), FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS link_previews (
	url TEXT PRIMARY KEY, title TEXT NOT NULL DEFAULT '', description TEXT NOT NULL DEFAULT '',
	image_url TEXT NOT NULL DEFAULT '', site_name TEXT NOT NULL DEFAULT '', fetched_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
const addEditHistoryCaptchaIPBlocks = `
CREATE TABLE IF NOT EXISTS message_edit_history (
	id TEXT PRIMARY KEY, message_id TEXT NOT NULL, old_content TEXT NOT NULL,
	edited_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (message_id) REFERENCES messages(id)
);
CREATE TABLE IF NOT EXISTS captcha_tokens (
	token TEXT PRIMARY KEY, solution TEXT NOT NULL, expires_at TEXT NOT NULL, used INTEGER NOT NULL DEFAULT 0
);
CREATE TABLE IF NOT EXISTS ip_blocks (
	ip_address TEXT PRIMARY KEY, reason TEXT NOT NULL DEFAULT '', blocked_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	expires_at TEXT, attempts INTEGER NOT NULL DEFAULT 1
);
CREATE TABLE IF NOT EXISTS login_attempts (
	ip_address TEXT NOT NULL, attempted_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	email TEXT NOT NULL, success INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX IF NOT EXISTS idx_login_attempts_ip ON login_attempts(ip_address);
CREATE INDEX IF NOT EXISTS idx_message_self_destruct ON message_self_destruct(delete_at);
`
const addAdminSettings = `
CREATE TABLE IF NOT EXISTS admin_users (
	user_id TEXT PRIMARY KEY, role TEXT NOT NULL DEFAULT 'admin', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS admin_logs (
	id TEXT PRIMARY KEY, admin_id TEXT NOT NULL, action TEXT NOT NULL, target_type TEXT NOT NULL,
	target_id TEXT NOT NULL DEFAULT '', details TEXT NOT NULL DEFAULT '', created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (admin_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS app_settings (
	key TEXT PRIMARY KEY, value TEXT NOT NULL, updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
INSERT OR IGNORE INTO app_settings (key, value) VALUES ('captcha_enabled', 'false');
INSERT OR IGNORE INTO app_settings (key, value) VALUES ('max_login_attempts', '5');
INSERT OR IGNORE INTO app_settings (key, value) VALUES ('ip_block_duration_minutes', '30');
INSERT OR IGNORE INTO app_settings (key, value) VALUES ('self_destruct_default_minutes', '0');
INSERT OR IGNORE INTO app_settings (key, value) VALUES ('admin_can_read_messages', 'true');
`
const addLoginCodes = `
CREATE TABLE IF NOT EXISTS email_login_codes (
	id TEXT PRIMARY KEY, email TEXT NOT NULL, code TEXT NOT NULL,
	expires_at TEXT NOT NULL, verified INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
CREATE TABLE IF NOT EXISTS phone_login_codes (
	id TEXT PRIMARY KEY, phone TEXT NOT NULL, code TEXT NOT NULL,
	expires_at TEXT NOT NULL, verified INTEGER NOT NULL DEFAULT 0, created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP
);
`
const addContactPhoto = `
ALTER TABLE contacts ADD COLUMN photo_url TEXT NOT NULL DEFAULT '';
ALTER TABLE contacts ADD COLUMN user_id_ref TEXT NOT NULL DEFAULT '';
ALTER TABLE contacts ADD COLUMN updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP;
CREATE INDEX IF NOT EXISTS idx_contacts_user_id_ref ON contacts(user_id_ref);
`
const addMessageFields = `
ALTER TABLE messages ADD COLUMN caption TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN mime_type TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN duration INTEGER NOT NULL DEFAULT 0;
ALTER TABLE messages ADD COLUMN width INTEGER NOT NULL DEFAULT 0;
ALTER TABLE messages ADD COLUMN height INTEGER NOT NULL DEFAULT 0;
`
const addInviteLinks = `
CREATE TABLE IF NOT EXISTS invite_links (
	id TEXT PRIMARY KEY,
	chat_id TEXT NOT NULL,
	creator_id TEXT NOT NULL,
	code TEXT UNIQUE NOT NULL,
	expires_at TEXT,
	usage_limit INTEGER NOT NULL DEFAULT 0,
	usage_count INTEGER NOT NULL DEFAULT 0,
	active INTEGER NOT NULL DEFAULT 1,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (chat_id) REFERENCES chats(id),
	FOREIGN KEY (creator_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_invite_links_code ON invite_links(code);
CREATE INDEX IF NOT EXISTS idx_invite_links_chat_id ON invite_links(chat_id);
`
const addChatFolders = `
CREATE TABLE IF NOT EXISTS chat_folders (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	name TEXT NOT NULL,
	emoji TEXT NOT NULL DEFAULT '',
	folder_order INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS chat_folder_items (
	folder_id TEXT NOT NULL,
	chat_id TEXT NOT NULL,
	added_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (folder_id, chat_id),
	FOREIGN KEY (folder_id) REFERENCES chat_folders(id),
	FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE INDEX IF NOT EXISTS idx_chat_folder_items_folder ON chat_folder_items(folder_id);
`
const addChatSlowMode = `
ALTER TABLE chats ADD COLUMN slow_mode_seconds INTEGER NOT NULL DEFAULT 0;
`
const addChatThemes = `
CREATE TABLE IF NOT EXISTS chat_themes (
	chat_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	theme TEXT NOT NULL DEFAULT 'default',
	background_url TEXT NOT NULL DEFAULT '',
	updated_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (chat_id, user_id),
	FOREIGN KEY (chat_id) REFERENCES chats(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
`

const addIsAdminField = `
ALTER TABLE users ADD COLUMN is_admin INTEGER NOT NULL DEFAULT 0;
`

const addSavedMessagesEmojisVoiceChats = `
ALTER TABLE messages ADD COLUMN latitude REAL NOT NULL DEFAULT 0;
ALTER TABLE messages ADD COLUMN longitude REAL NOT NULL DEFAULT 0;
ALTER TABLE messages ADD COLUMN location_title TEXT NOT NULL DEFAULT '';
ALTER TABLE messages ADD COLUMN effect TEXT NOT NULL DEFAULT '';
CREATE TABLE IF NOT EXISTS saved_messages (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	message_id TEXT NOT NULL,
	chat_id TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id),
	FOREIGN KEY (message_id) REFERENCES messages(id),
	FOREIGN KEY (chat_id) REFERENCES chats(id)
);
CREATE TABLE IF NOT EXISTS custom_emojis (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	shortcode TEXT NOT NULL,
	file_url TEXT NOT NULL DEFAULT '',
	file_path TEXT NOT NULL DEFAULT '',
	animated INTEGER NOT NULL DEFAULT 0,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS voice_chats (
	id TEXT PRIMARY KEY,
	chat_id TEXT NOT NULL,
	started_by TEXT NOT NULL,
	title TEXT NOT NULL DEFAULT '',
	status TEXT NOT NULL DEFAULT 'active',
	participant_count INTEGER NOT NULL DEFAULT 0,
	scheduled_at TEXT,
	started_at TEXT,
	ended_at TEXT,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (chat_id) REFERENCES chats(id),
	FOREIGN KEY (started_by) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS voice_chat_participants (
	voice_chat_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	joined_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	left_at TEXT,
	muted INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (voice_chat_id, user_id),
	FOREIGN KEY (voice_chat_id) REFERENCES voice_chats(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_saved_messages_user_id ON saved_messages(user_id);
CREATE INDEX IF NOT EXISTS idx_custom_emojis_user_id ON custom_emojis(user_id);
CREATE INDEX IF NOT EXISTS idx_voice_chats_chat_id ON voice_chats(chat_id);
CREATE INDEX IF NOT EXISTS idx_voice_chat_participants_user ON voice_chat_participants(user_id);
`

const addStoriesGroupCallsChannels = `
CREATE TABLE IF NOT EXISTS stories (
	id TEXT PRIMARY KEY,
	user_id TEXT NOT NULL,
	file_path TEXT NOT NULL DEFAULT '',
	file_url TEXT NOT NULL DEFAULT '',
	type TEXT NOT NULL DEFAULT 'photo',
	caption TEXT NOT NULL DEFAULT '',
	expires_at TEXT NOT NULL,
	created_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS story_views (
	story_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	viewed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (story_id, user_id),
	FOREIGN KEY (story_id) REFERENCES stories(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS channel_subscribers (
	channel_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	role TEXT NOT NULL DEFAULT 'subscriber',
	subscribed_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	PRIMARY KEY (channel_id, user_id),
	FOREIGN KEY (channel_id) REFERENCES chats(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS group_calls (
	id TEXT PRIMARY KEY,
	chat_id TEXT NOT NULL,
	caller_id TEXT NOT NULL,
	type TEXT NOT NULL DEFAULT 'audio',
	status TEXT NOT NULL DEFAULT 'initiated',
	started_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	ended_at TEXT,
	FOREIGN KEY (chat_id) REFERENCES chats(id),
	FOREIGN KEY (caller_id) REFERENCES users(id)
);
CREATE TABLE IF NOT EXISTS group_call_participants (
	call_id TEXT NOT NULL,
	user_id TEXT NOT NULL,
	joined_at TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP,
	left_at TEXT,
	audio_muted INTEGER NOT NULL DEFAULT 0,
	video_muted INTEGER NOT NULL DEFAULT 0,
	PRIMARY KEY (call_id, user_id),
	FOREIGN KEY (call_id) REFERENCES group_calls(id),
	FOREIGN KEY (user_id) REFERENCES users(id)
);
CREATE INDEX IF NOT EXISTS idx_stories_user_id ON stories(user_id);
CREATE INDEX IF NOT EXISTS idx_stories_expires_at ON stories(expires_at);
CREATE INDEX IF NOT EXISTS idx_story_views_story_id ON story_views(story_id);
CREATE INDEX IF NOT EXISTS idx_channel_subscribers_user_id ON channel_subscribers(user_id);
CREATE INDEX IF NOT EXISTS idx_group_calls_chat_id ON group_calls(chat_id);
`

func init() {
	migrations = append(migrations, Migration{
		ID:   22,
		Name: "add_stories_group_calls_channels",
		Up:   addStoriesGroupCallsChannels,
	})
	migrations = append(migrations, Migration{
		ID:   23,
		Name: "add_saved_messages_emojis_voice_chats",
		Up:   addSavedMessagesEmojisVoiceChats,
	})
}
