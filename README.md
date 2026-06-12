# Chat Messenger Server

High-performance chat messenger server built with Go. Features REST API, real-time messaging via WebSocket, WebRTC voice/video calls, voice chats, stories, broadcast channels, polls, stickers, bots, location sharing, message effects, custom emojis, saved messages, and more.

---

## Features

- **User System** — registration, authentication (JWT), profiles, avatars, contacts sync
- **Real-time Chat** — private & group chats, typing indicators, read receipts, slow mode
- **Messaging** — text, files, images, voice messages, video circles, mentions, reactions, forward, pin, star, schedule
- **Reply to Messages** — quote and reply to specific messages with inline citation
- **Supergroup Roles** — owner, admin, moderator, editor, member, read-only permission levels
- **Broadcast Channels** — one-way broadcast channels with subscribers and admins
- **Stories** — photo/video stories that disappear after 24 hours with view tracking
- **Saved Messages** — bookmark any message for later reference with paginated listing
- **Location Sharing** — share locations as messages with latitude/longitude
- **Message Effects** — send messages with animation effects: confetti, fireworks, hearts, balloons, stars
- **Custom Emojis** — upload custom emoji images and use them via shortcodes
- **Voice Chats** — persistent voice chat rooms in groups with join/leave/mute
- **WebSocket Events** — instant push for new messages, edits, deletes, calls, polls, stories, and more
- **Voice/Video Calls** — WebRTC-based signalling via REST + WS (1-to-1)
- **Group Calls** — multi-participant voice/video conference calls
- **Custom Reactions** — react with any emoji
- **Polls** — create, vote, close, real-time updates
- **Stickers** — packs, library, custom sticker sets
- **Bots** — create and manage bot accounts with tokens
- **Drafts** — save/load/delete message drafts per chat
- **Folders** — organize chats into custom folders
- **Invite Links** — generate and manage group invite links
- **Sessions** — manage active sessions, force logout
- **Self-Destructing Messages** — auto-delete after TTL
- **Message Export** — export chat history as CSV
- **Media Gallery** — filter by photo, video, audio, file, video_circle
- **Archived Chats** — archive/unarchive/hide chats
- **Global Search** — search messages and users across all chats
- **Push Notifications** — Firebase Cloud Messaging (FCM) integration
- **E2E Encryption** — support for end-to-end encrypted messages
- **IP Blocking** — block abusive IPs at application level
- **Admin Panel** — admin registration, moderation, reporting
- **Rate Limiting** — protect API from abuse

---

## Docker (Production)

```bash
# Build image
docker build -t chatserver .

# Run container
docker run -d --name chatserver -p 8080:8080 \
  -v chatserver-data:/app/data \
  -v chatserver-uploads:/app/uploads \
  -e DATABASE_PATH=/app/data/chatserver.db \
  -e JWT_SECRET=change-me-in-production \
  -e GIN_MODE=release \
  chatserver
```

## CI/CD

The project includes a GitHub Actions pipeline (`.github/workflows/deploy.yml`) that:
1. Runs on push to `master`/`main`
2. Checks out code, sets up Go
3. Downloads dependencies, builds, runs `go vet`
4. Validates Swagger documentation
5. Builds production binary and Docker image

## Quick Start

```bash
# Build & run
go build -o ChatServer.exe .
./ChatServer.exe

# Swagger UI:  http://localhost:8080/swagger/index.html
# API Tester:  http://localhost:8080/app/
# Health:      http://localhost:8080/health

# Register a new user
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@mail.com","password":"secret123","displayName":"John"}'

# Login
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@mail.com","password":"secret123"}'

# Admin registration (requires ADMIN_SECRET)
curl -X POST http://localhost:8080/api/auth/admin/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","email":"admin@mail.com","password":"admin123","displayName":"Admin","admin_secret":"admin-secret-change-me"}'

# WebSocket (replace JWT_TOKEN with your token)
wscat -c "ws://localhost:8080/ws?token=JWT_TOKEN"
```

---

## Architecture

```
internal/
├── domain/              # Domain models (split by sub-package)
│   ├── auth/            #   login, register, captcha, admin
│   ├── bot/             #   bot accounts
│   ├── call/            #   WebRTC call models + group calls
│   ├── channel/         #   broadcast channel subscribers
│   ├── chat/            #   chats, participants, folders, invite links, starred, saved
│   ├── contact/         #   contacts sync
│   ├── draft/           #   message drafts + scheduled messages
│   ├── e2e/             #   end-to-end encryption keys
│   ├── emoji/           #   custom emojis
│   ├── ipblock/         #   IP block list
│   ├── link/            #   link previews
│   ├── message/         #   messages, reactions, read receipts, mentions,
│   │                   #   edit history, bookmarks, self-destruct, export
│   ├── notification/    #   notification settings
│   ├── poll/            #   polls & votes
│   ├── report/          #   abuse reports
│   ├── session/         #   user sessions
│   ├── sticker/         #   sticker packs & library
│   ├── user/            #   users, blocks, account settings, status
│   ├── verification/    #   email/phone verification codes
│   └── voicechat/       #   persistent voice chat rooms
├── repository/          # SQLite data access layer (interface + impl)
│   ├── account/         #   account settings
│   ├── bot/             #   bot accounts
│   ├── call/            #   call history
│   ├── channel/         #   channel subscribers
│   ├── chat/            #   chats, participants, slow mode
│   ├── contact/         #   contacts
│   ├── draft/           #   drafts + scheduled messages
│   ├── emoji/           #   custom emojis
│   ├── folder/          #   chat folders
│   ├── gif/             #   saved GIFs
│   ├── groupcall/       #   group call participants
│   ├── link/            #   invite links
│   ├── message/         #   messages, reactions, bookmarks, mentions
│   ├── poll/            #   polls & votes
│   ├── savedmsg/        #   saved messages
│   ├── schedmsg/        #   scheduled messages
│   ├── session/         #   user sessions
│   ├── sticker/         #   sticker packs & library
│   ├── story/           #   stories + views
│   ├── user/            #   users, blocks
│   ├── verification/    #   verification codes
│   └── voicechat/       #   voice chat rooms + participants
├── service/             # Business logic layer
│   ├── auth/            #   authentication service
│   ├── bot/             #   bot management
│   ├── call/            #   call signalling
│   ├── channel/         #   channel subscription logic
│   ├── chat/            #   chat CRUD, participants, pin, archive
│   ├── contact/         #   contact sync & search
│   ├── draft/           #   draft & scheduled message logic
│   ├── emoji/           #   custom emoji service
│   ├── folder/          #   folder CRUD
│   ├── gif/             #   saved GIFs
│   ├── groupcall/       #   group call management
│   ├── link/            #   invite links
│   ├── mention/         #   @mention parsing & notifications
│   ├── message/         #   message CRUD, reactions, forward, export
│   ├── poll/            #   poll CRUD & voting
│   ├── push/            #   FCM push notifications
│   ├── savedmsg/        #   saved messages service
│   ├── schedmsg/        #   scheduled message dispatch
│   ├── session/         #   session management
│   ├── sticker/         #   sticker packs & library
│   ├── story/           #   story CRUD + views
│   ├── systemmsg/       #   system-generated messages
│   ├── typing/          #   typing indicator broadcast
│   ├── user/            #   user profiles, blocks, status
│   ├── verification/    #   email/phone code verification
│   └── voicechat/       #   voice chat rooms
├── handler/             # HTTP handlers (Gin framework)
│   ├── auth/            #   register, login, refresh, change password
│   ├── bot/             #   bot CRUD
│   ├── call/            #   call init/respond/end
│   ├── channel/         #   channel subscribe/unsubscribe/roles
│   ├── chat/            #   chat CRUD, participants, notifications
│   ├── contact/         #   contact sync & search
│   ├── draft/           #   draft CRUD
│   ├── emoji/           #   custom emoji upload/CRUD
│   ├── folder/          #   folder CRUD
│   ├── gif/             #   saved GIFs
│   ├── groupcall/       #   group call handlers
│   ├── link/            #   invite links
│   ├── login/           #   login code (email/phone)
│   ├── message/         #   message CRUD, reactions, search, export, location
│   ├── poll/            #   poll CRUD & voting
│   ├── savedmsg/        #   saved messages CRUD
│   ├── schedmsg/        #   scheduled messages
│   ├── session/         #   session management
│   ├── sticker/         #   sticker packs & library
│   ├── story/           #   story CRUD
│   ├── user/            #   user profile, blocks, avatar
│   ├── verification/    #   verification endpoints
│   ├── voicechat/       #   voice chat rooms
│   └── ws/              #   WebSocket handler + event routing
├── middleware/          # JWT auth, admin auth, rate limiting
├── ws/                  # WebSocket hub, client, event bus
├── config/              # Environment-based configuration
└── ...                  # Router, DI container, main wiring

pkg/                     # Shared utilities
└── response/            #   JSON response helpers (Success, Paginated, Error)

docs/                    # API documentation
├── swagger.json         #   OpenAPI/Swagger spec
└── swagger.yaml         #   YAML format
```

---

## API Endpoints

All endpoints require `Authorization: Bearer <JWT_TOKEN>` unless noted.

### Auth

| Method | Endpoint | Auth | Description |
|--------|----------|------|-------------|
| `POST` | `/api/auth/register` | No | Register new user |
| `POST` | `/api/auth/admin/register` | No | Register admin (requires `admin_secret`) |
| `POST` | `/api/auth/login` | No | Login with email + password, returns JWT |
| `POST` | `/api/auth/login/email` | No | Send login code to email |
| `POST` | `/api/auth/login/email/verify` | No | Verify email login code |
| `POST` | `/api/auth/login/phone` | No | Send login code via SMS |
| `POST` | `/api/auth/login/phone/verify` | No | Verify SMS login code |
| `GET` | `/api/auth/refresh` | Yes | Refresh JWT token |
| `PUT` | `/api/auth/change-password` | Yes | Change password |

### Users

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users/profile` | Get current user profile |
| `PUT` | `/api/users/profile` | Update profile (displayName, bio, etc.) |
| `GET` | `/api/users/search?q=&offset=&limit=` | Search users (paginated) |
| `GET` | `/api/users/:id` | Get user by ID |
| `GET` | `/api/users/username/:username` | Get user by username |
| `PUT` | `/api/users/status` | Set online status (`online`/`offline`/`away`) |
| `PUT` | `/api/users/push-token` | Update FCM push token |
| `POST` | `/api/users/block` | Block a user |
| `DELETE` | `/api/users/block/:userId` | Unblock a user |
| `GET` | `/api/users/blocked` | List blocked users |
| `POST` | `/api/users/avatar` | Upload profile avatar |
| `POST` | `/api/users/push-test` | Send test push notification |

### Chats

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats` | Create chat (private or group) |
| `POST` | `/api/chats/start/:userId` | Start or find private chat with user |
| `GET` | `/api/chats` | List my chats |
| `GET` | `/api/chats/search?q=` | Search chats by name |
| `GET` | `/api/chats/archived` | List archived chats |
| `GET` | `/api/chats/:id` | Get chat details |
| `PUT` | `/api/chats/:id` | Update group info |
| `DELETE` | `/api/chats/:id` | Delete chat |
| `POST` | `/api/chats/:id/read` | Mark as read (up to last message) |
| `POST` | `/api/chats/:id/pin` | Pin chat to top |
| `DELETE` | `/api/chats/:id/pin` | Unpin chat |
| `POST` | `/api/chats/:id/archive` | Archive chat |
| `POST` | `/api/chats/:id/unarchive` | Unarchive chat |
| `POST` | `/api/chats/:id/hide` | Hide chat from list |
| `POST` | `/api/chats/:id/leave` | Leave group chat |
| `POST` | `/api/chats/:id/transfer-ownership` | Transfer group ownership |
| `PUT` | `/api/chats/:id/slow-mode` | Set slow mode interval (0-3600s) |
| `PUT` | `/api/chats/:id/notifications` | Mute/unmute notifications |
| `GET` | `/api/chats/:id/notifications` | Get notification settings |
| `POST` | `/api/chats/:id/participants` | Add participant to group |
| `DELETE` | `/api/chats/:id/participants/:userId` | Remove participant |
| `PUT` | `/api/chats/:id/participants/:userId/role` | Change participant role |

### Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/:id/messages` | Send a message |
| `GET` | `/api/chats/:id/messages` | List messages (paginated, newest first) |
| `GET` | `/api/chats/:id/messages/search?q=` | Search within chat |
| `POST` | `/api/chats/:id/messages/file` | Upload file attachment |
| `POST` | `/api/chats/:id/messages/:msgId/resend` | Re-send a message |
| `GET` | `/api/chats/:id/pinned` | Get pinned messages |
| `GET` | `/api/chats/:id/media?type=photo` | Media gallery (`photo`/`video`/`audio`/`file`) |
| `GET` | `/api/chats/:id/export?format=csv` | Export chat history as CSV |
| `GET` | `/api/messages/:id` | Get message by ID |
| `PUT` | `/api/messages/:id` | Edit message |
| `DELETE` | `/api/messages/:id` | Delete message (for everyone) |
| `DELETE` | `/api/messages/:id/for-me` | Delete message (for me only) |
| `POST` | `/api/messages/:id/reactions` | Add reaction (`{"emoji":"🔥"}`) |
| `DELETE` | `/api/messages/:id/reactions?emoji=` | Remove reaction |
| `PUT` | `/api/messages/:id/pin` | Pin/unpin message |
| `POST` | `/api/messages/:id/star` | Star message |
| `DELETE` | `/api/messages/:id/star` | Unstar message |
| `POST` | `/api/messages/:id/read` | Mark single message as read |
| `GET` | `/api/messages/search?q=&offset=&limit=` | Global message search |
| `GET` | `/api/messages/starred` | List starred messages |
| `POST` | `/api/messages/forward` | Forward messages to another chat |

### Contacts

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/contacts/sync` | Sync phone contacts |
| `GET` | `/api/contacts` | List contacts |
| `GET` | `/api/contacts/search?q=` | Search contacts by phone |
| `GET` | `/api/contacts/registered` | Contacts who are on the platform |
| `POST` | `/api/contacts/photo` | Update contact photo |

### Folders

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/folders` | Create folder |
| `GET` | `/api/folders` | List folders |
| `PUT` | `/api/folders/:id` | Update folder (name, order, chats) |
| `DELETE` | `/api/folders/:id` | Delete folder |

### Invite Links

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/:id/invite-links` | Create invite link |
| `GET` | `/api/chats/:id/invite-links` | List invite links for chat |
| `DELETE` | `/api/chats/:id/invite-links/:linkId` | Revoke invite link |
| `POST` | `/api/chats/join` | Join group via invite link |

### Polls

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/:id/polls` | Create poll |
| `GET` | `/api/chats/:id/polls` | List polls in chat |
| `POST` | `/api/polls/:pollId/vote` | Cast vote (multiple choice supported) |
| `POST` | `/api/polls/:pollId/close` | Close poll (admin only) |

### Stickers

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/stickers/packs` | Create sticker pack |
| `GET` | `/api/stickers/packs` | List all sticker packs |
| `GET` | `/api/stickers/packs/my` | My sticker packs |
| `GET` | `/api/stickers/packs/:id` | Get sticker pack by ID |
| `DELETE` | `/api/stickers/packs/:id` | Delete sticker pack |
| `POST` | `/api/stickers/library` | Add sticker to library |
| `GET` | `/api/stickers/library` | My sticker library |

### Calls

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/calls` | Initiate WebRTC call |
| `POST` | `/api/calls/:id/respond` | Accept/reject call |
| `POST` | `/api/calls/:id/end` | End call |
| `GET` | `/api/calls/:id` | Get call details |
| `GET` | `/api/chats/:chatId/calls` | Call history for chat |

### Sessions

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/sessions` | List active sessions |
| `DELETE` | `/api/sessions/:id` | Terminate specific session |
| `DELETE` | `/api/sessions` | Terminate all sessions (except current) |

### Bots

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/bots` | Create bot account |
| `GET` | `/api/bots` | List my bots |
| `PUT` | `/api/bots/:id` | Update bot |
| `DELETE` | `/api/bots/:id` | Delete bot |
| `POST` | `/api/bots/:id/token` | Generate new bot token |

### Drafts & Scheduled Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/drafts` | Save draft |
| `GET` | `/api/drafts?chatId=` | Get draft for chat |
| `DELETE` | `/api/drafts/:id` | Delete draft |
| `POST` | `/api/messages/schedule` | Schedule a message |
| `GET` | `/api/messages/scheduled` | List scheduled messages |
| `DELETE` | `/api/messages/scheduled/:id` | Cancel scheduled message |

### GIFs

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/gifs` | Save GIF to collection |
| `GET` | `/api/gifs` | List saved GIFs |
| `DELETE` | `/api/gifs` | Remove GIF |

### Stories

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/stories` | Create story (multipart: file + type + caption) |
| `GET` | `/api/stories` | Get stories from contacts and channels |
| `GET` | `/api/stories/my` | Get my active stories |
| `GET` | `/api/stories/:id` | View a story (marks as viewed) |
| `DELETE` | `/api/stories/:id` | Delete story (owner only) |
| `GET` | `/api/stories/:id/views` | Get story viewers (owner only) |

### Group Calls

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/calls/group/initiate` | Start a group call |
| `POST` | `/api/calls/group/respond` | Join/leave/mute in a group call |
| `POST` | `/api/calls/group/:id/end` | End a group call (caller only) |
| `GET` | `/api/calls/group/:id` | Get group call details |
| `GET` | `/api/chats/:chatId/active-calls` | Active calls in a chat |

### Broadcast Channels

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats` | Create channel (with `type: "channel"`) |
| `POST` | `/api/channels/subscribe` | Subscribe to a channel |
| `POST` | `/api/channels/unsubscribe` | Unsubscribe from a channel |
| `GET` | `/api/channels` | List my channels (owned + subscribed) |
| `GET` | `/api/channels/:id/subscribers` | List subscribers (admin only) |
| `GET` | `/api/channels/:id/subscribed` | Check if subscribed |
| `PUT` | `/api/channels/:id/subscribers/:userId/role` | Set subscriber role (admin) |

### Video Circle Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/:id/messages/video-circle` | Send a circular video message |

### Location Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/:id/messages/location` | Send a location message (lat/lng) |

### Message Effects

Available via `SendMessageRequest.effect` field: `confetti`, `fireworks`, `hearts`, `balloons`, `stars`.

### Saved Messages

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/messages/:id/save?chatId=` | Save a message for later |
| `GET` | `/api/saved-messages?limit=&offset=` | List saved messages (paginated) |
| `DELETE` | `/api/saved-messages/:id` | Remove a saved message |

### Custom Emojis

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/emojis` | Upload custom emoji (multipart: `shortcode` + `emoji` file) |
| `GET` | `/api/emojis` | List all public emojis |
| `GET` | `/api/emojis/my` | List my uploaded emojis |
| `DELETE` | `/api/emojis/:id` | Delete my emoji |

### Voice Chats

Persistent voice chat rooms in groups (not group calls).

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/:id/voice-chat` | Create/start a voice chat in a group |
| `GET` | `/api/chats/:id/voice-chats/active` | List active voice chats |
| `GET` | `/api/chats/:id/voice-chats/history` | Voice chat history |
| `GET` | `/api/voice-chats/:id` | Get voice chat details |
| `POST` | `/api/voice-chats/:id/join` | Join a voice chat |
| `POST` | `/api/voice-chats/:id/leave` | Leave a voice chat |
| `POST` | `/api/voice-chats/:id/end` | End a voice chat |
| `POST` | `/api/voice-chats/:id/mute` | Mute/unmute in voice chat |

---

## WebSocket Events

```
ws://localhost:8080/ws?token=JWT_TOKEN
```

### Server to Client

Events pushed from server to connected clients in real time.

| Event | Payload | Trigger |
|-------|---------|---------|
| `message:new` | `{"message":{...}}` | Message sent |
| `message:edited` | `{"message":{...}}` | Message edited |
| `message:deleted` | `{"chatId":"...","messageId":"..."}` | Message deleted |
| `message:read` | `{"chatId":"...","userId":"...","readUpTo":"..."}` | Messages read |
| `message:pinned` | `{"chatId":"...","messageId":"...","pinned":bool}` | Message pin toggled |
| `message:starred` | `{"messageId":"...","starred":bool}` | Message star toggled |
| `message:scheduled` | `{"message":{...}}` | Scheduled message created |
| `chat:created` | `{"chat":{...}}` | Chat created |
| `chat:updated` | `{"chat":{...}}` | Chat updated |
| `chat:deleted` | `{"chatId":"..."}` | Chat deleted |
| `chat:slowmode` | `{"chatId":"...","slowMode":int}` | Slow mode changed |
| `chat:role` | `{"chatId":"...","userId":"...","role":"..."}` | Participant role changed |
| `chat:ownership` | `{"chatId":"...","newOwnerId":"..."}` | Ownership transferred |
| `chat:participant_added` | `{"chatId":"...","user":{...}}` | Participant added |
| `chat:participant_removed` | `{"chatId":"...","userId":"..."}` | Participant removed |
| `user:online` | `{"userId":"...","username":"..."}` | User connected via WS |
| `user:offline` | `{"userId":"...","username":"..."}` | User disconnected |
| `user:typing` | `{"chatId":"...","userId":"..."}` | User started typing |
| `user:stop_typing` | `{"chatId":"...","userId":"..."}` | User stopped typing |
| `user:keyboard_opened` | `{"chatId":"...","userId":"..."}` | User opened keyboard |
| `user:keyboard_closed` | `{"chatId":"...","userId":"..."}` | User closed keyboard |
| `user:status` | `{"userId":"...","status":"..."}` | User status changed |
| `call:offer` | `{"call":{...}}` | Incoming call |
| `call:accept` | `{"callId":"...","accepted":bool}` | Call accepted |
| `call:end` | `{"callId":"...","ended":bool}` | Call ended |
| `poll:created` | `{"poll":{...}}` | Poll created |
| `poll:vote` | `{"pollId":"...","optionId":"..."}` | Vote cast |
| `poll:closed` | `{"pollId":"...","closed":bool}` | Poll closed |
| `folder:created` | `{"folder":{...}}` | Folder created |
| `folder:updated` | `{"folder":{...}}` | Folder updated |
| `folder:deleted` | `{"folderId":"..."}` | Folder deleted |
| `invite:created` | `{"link":{...}}` | Invite link created |
| `invite:joined` | `{"chatId":"...","userId":"..."}` | User joined via invite |
| `sticker:added` | `{"sticker":{...}}` | Sticker added to library |
| `sticker:pack_created` | `{"pack":{...}}` | Sticker pack created |
| `story:new` | `{"story":{...}}` | Story created |
| `story:deleted` | `{"storyId":"..."}` | Story deleted |
| `story:viewed` | `{"storyId":"...","userId":"..."}` | Story viewed |
| `call:group_started` | `{"call":{...}}` | Group call started |
| `call:group_joined` | `{"callId":"...","userId":"..."}` | User joined group call |
| `call:group_left` | `{"callId":"...","userId":"..."}` | User left group call |
| `call:group_mute` | `{"callId":"...","userId":"...","audioMuted":bool}` | Mute status changed |
| `call:group_ended` | `{"callId":"..."}` | Group call ended |

### Client to Server

Send these JSON events over the WebSocket connection.

| Event | Payload | Description |
|-------|---------|-------------|
| `user:typing` | `{"chatId":"..."}` | Notify chat that you are typing |
| `user:stop_typing` | `{"chatId":"..."}` | Notify chat you stopped typing |
| `user:keyboard_opened` | `{"chatId":"..."}` | Notify chat you opened keyboard |
| `user:keyboard_closed` | `{"chatId":"..."}` | Notify chat you closed keyboard |

---

## Response Format

### Success

```json
{
  "success": true,
  "data": { ... }
}
```

### Paginated List

```json
{
  "success": true,
  "data": [ ... ],
  "meta": {
    "total": 100,
    "offset": 0,
    "limit": 50
  }
}
```

### Error

```json
{
  "success": false,
  "error": "Human-readable description",
  "code": "ERROR_CODE"
}
```

### Error Codes

| Code | HTTP Status | Description |
|------|-------------|-------------|
| `BAD_REQUEST` | 400 | Invalid request body or parameters |
| `VALIDATION_ERROR` | 400 | Field validation failed |
| `UNAUTHORIZED` | 401 | Missing or invalid JWT token |
| `TOKEN_EXPIRED` | 401 | JWT token has expired |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `DUPLICATE` | 409 | Resource already exists |
| `RATE_LIMIT` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Internal server error |

---

## Configuration

All configuration is via environment variables.

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `file:chat.db?cache=shared&_journal_mode=WAL` | SQLite connection string |
| `JWT_SECRET` | `super-secret-change-me-in-production` | HMAC-SHA256 key for JWT |
| `ADMIN_SECRET` | `admin-secret-change-me` | Secret key for admin registration |
| `JWT_TTL` | `86400` | JWT token lifetime in seconds (24h) |
| `CORS_ORIGIN` | `*` | Allowed CORS origin |
| `FCM_KEY_PATH` | — | Path to Firebase service account JSON |
| `ENCRYPTION_KEY` | — | Key for E2E encryption |

---

## Build

### Prerequisites

- Go 1.21 or later
- GCC/MinGW (for SQLite via modernc.org/sqlite — pure Go, no CGO required)

### Commands

```bash
# Build for Windows
go build -o ChatServer.exe .

# Build for Linux
GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 .

# Build for macOS
GOOS=darwin GOARCH=amd64 go build -o chat-server-darwin-amd64 .

# Run tests
go test ./... -v

# Run integration tests
$env:INTEGRATION=1; go test ./... -v -run Integration
```

---

## Tech Stack

- **Language:** Go 1.21+
- **Database:** SQLite (WAL mode, via modernc.org/sqlite — pure Go, no CGO)
- **HTTP Framework:** Gin
- **WebSocket:** gorilla/websocket
- **Auth:** JWT (HMAC-SHA256), bcrypt
- **Push Notifications:** Firebase Cloud Messaging (FCM)
- **API Documentation:** Swaggo (swagger.json)
- **Real-time:** Custom WebSocket hub with per-client channels
- **Architecture:** Domain-driven sub-packages (domain/repository/service/handler)
