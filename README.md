<p align="center">
  <img src="https://img.shields.io/badge/version-2.0.0-blue?style=for-the-badge" alt="Version">
  <img src="https://img.shields.io/badge/Go-1.21+-00ADD8?style=for-the-badge&logo=go" alt="Go">
  <img src="https://img.shields.io/badge/SQLite-003B57?style=for-the-badge&logo=sqlite" alt="SQLite">
  <img src="https://img.shields.io/badge/WebSocket-Real--Time-4A90D9?style=for-the-badge" alt="WebSocket">
  <img src="https://img.shields.io/badge/WebRTC-Voice%2FVideo-FF5722?style=for-the-badge" alt="WebRTC">
</p>

<h1 align="center">💬 Go Messenger Server</h1>

<p align="center">
  <b>High-performance chat messenger server</b><br>
  REST API + WebSocket real-time + WebRTC calls + Voice Chats + Stories + Channels + Polls + Stickers + Bots
</p>

<p align="center">
  <a href="#-quick-start">🚀 Quick Start</a> •
  <a href="#-features">✨ Features</a> •
  <a href="#-api-endpoints">📡 API</a> •
  <a href="#-websocket">🔌 WebSocket</a> •
  <a href="#-postman">📮 Postman</a> •
  <a href="#-architecture">🏗️ Architecture</a>
</p>

---

## 🚀 Quick Start

```bash
# 1. Clone and build
git clone https://github.com/JamikKhidirov/ChatServer.git
cd ChatServer
go build -o ChatServer.exe .

# 2. Run the server
./ChatServer.exe

# Server starts on http://localhost:8080
```

When the server starts, you'll see all available URLs:

```
╔══════════════════════════════════════════════════════════════════════╗
║                    CHAT MESSENGER SERVER v2.0                       ║
╠══════════════════════════════════════════════════════════════════════╣
║  Server:   http://localhost:8080                                     ║
║  Frontend: http://localhost:8080/app/                                ║
║  Swagger:  http://localhost:8080/swagger/index.html                  ║
║  Postman:  http://localhost:8080/postman                              ║
║  WebSocket: ws://localhost:8080/ws?token={jwt}                       ║
║  Health:   http://localhost:8080/health                              ║
╠══════════════════════════════════════════════════════════════════════╣
║  Import Postman: Download from /postman, import into Postman        ║
║  Auth: Register or Login → copy token → paste as Bearer token       ║
╚══════════════════════════════════════════════════════════════════════╝
```

### 📝 Test with curl

```bash
# Register
curl -s -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@mail.com","password":"secret123","displayName":"John"}'

# Login → copy the token
curl -s -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{"email":"john@mail.com","password":"secret123"}'

# Save token
$TOKEN="eyJhbGciOiJIUzI1NiIs..."

# Get profile
curl -s http://localhost:8080/api/users/profile -H "Authorization: Bearer $TOKEN"

# Create a chat
curl -s -X POST http://localhost:8080/api/chats \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"type":"private","participantIds":["USER_ID"]}'
```

---

## ✨ Features

<details open>
<summary><b>👤 User System</b></summary>

| Feature | Description |
|---------|-------------|
| Registration | Email + username + password |
| Login | Email/password, email code, phone code |
| JWT Auth | Bearer token with configurable TTL |
| Profiles | Display name, bio, avatar, status |
| Contacts | Sync phone contacts, search registered users |
| Blocks | Block/unblock users |
| Sessions | Manage active sessions |
</details>

<details>
<summary><b>💬 Messaging</b></summary>

| Feature | Description |
|---------|-------------|
| Text messages | Send and receive text |
| File uploads | Images, documents, voice messages |
| Video circles | Circular video messages |
| Location sharing | Send lat/lng coordinates |
| Message effects | confetti, fireworks, hearts, balloons, stars |
| Edit/Delete | Edit or delete sent messages |
| Reactions | Any emoji as reaction |
| Forward | Forward messages to other chats |
| Pin/Star | Pin to chat or star for personal bookmarks |
| Reply/Quote | Reply to specific messages |
| Read receipts | Track who read what |
| Scheduled | Schedule messages for later delivery |
| Self-destruct | Auto-delete after TTL |
| Search | Full-text search across all chats |
</details>

<details>
<summary><b>👥 Chats & Groups</b></summary>

| Feature | Description |
|---------|-------------|
| Private chats | 1-to-1 messaging |
| Group chats | Multi-participant with roles |
| Supergroup roles | Owner, admin, moderator, editor, member, read-only |
| Slow mode | Configurable interval between messages |
| Invite links | Generate with expiration and usage limits |
| Chat folders | Organize chats into custom folders |
| Archive/Hide | Archive or hide chats from list |
| Transfer ownership | Transfer group ownership |
| Promote/Demote | Change participant roles |
</details>

<details>
<summary><b>📺 Broadcast Channels</b></summary>

One-way broadcast channels with subscriber management and admin roles.
</details>

<details>
<summary><b>📸 Stories</b></summary>

Photo/video stories that disappear after 24 hours with view tracking.
</details>

<details>
<summary><b>📞 Voice/Video Calls</b></summary>

WebRTC-based 1-to-1 and group calls with signalling via REST + WebSocket.
</details>

<details>
<summary><b>🎤 Voice Chats</b></summary>

Persistent voice chat rooms in groups (join/leave/mute).
</details>

<details>
<summary><b>📊 Polls</b></summary>

Create polls with multiple options, vote, close, real-time updates via WebSocket.
</details>

<details>
<summary><b>🎨 Custom Emojis</b></summary>

Upload custom emoji images and use them via shortcodes.
</details>

<details>
<summary><b>🤖 Bots</b></summary>

Create and manage bot accounts with API tokens.
</details>

<details>
<summary><b>🔒 Security</b></summary>

| Feature | Description |
|---------|-------------|
| JWT authentication | HMAC-SHA256 signed tokens |
| Password hashing | bcrypt |
| E2E encryption | Support for encrypted messages |
| IP blocking | Block abusive IPs at application level |
| Rate limiting | Protect API from abuse |
| Admin panel | Admin registration, moderation, reporting |
</details>

<details>
<summary><b>🔔 Push Notifications</b></summary>

Firebase Cloud Messaging (FCM) integration for mobile push notifications.
</details>

---

## 📡 API Endpoints

All endpoints require `Authorization: Bearer <JWT_TOKEN>` unless marked with 🔓.

<details open>
<summary><b>🔐 Auth</b> — 9 endpoints</summary>

| Method | Endpoint | Auth | Description |
|--------|----------|:----:|-------------|
| `POST` | `/api/auth/register` | 🔓 | Register new user |
| `POST` | `/api/auth/admin/register` | 🔓 | Register admin (requires `admin_secret`) |
| `POST` | `/api/auth/login` | 🔓 | Login with email + password, returns JWT |
| `POST` | `/api/auth/login/email` | 🔓 | Send login code to email |
| `POST` | `/api/auth/login/email/verify` | 🔓 | Verify email login code and get JWT |
| `POST` | `/api/auth/login/phone` | 🔓 | Send login code via SMS |
| `POST` | `/api/auth/login/phone/verify` | 🔓 | Verify SMS login code and get JWT |
| `GET` | `/api/auth/refresh` | ✅ | Refresh JWT token |
| `PUT` | `/api/auth/change-password` | ✅ | Change password |
</details>

<details>
<summary><b>👤 Users</b> — 12 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/users/profile` | Get current user profile |
| `PUT` | `/api/users/profile` | Update profile (displayName, bio, etc.) |
| `GET` | `/api/users/search?q=&offset=&limit=` | Search users (paginated) |
| `GET` | `/api/users/{id}` | Get user by ID |
| `GET` | `/api/users/username/{username}` | Get user by username |
| `GET` | `/api/users/{id}/last-seen` | Get user's last seen time |
| `PUT` | `/api/users/status` | Set online status (`online`/`offline`/`away`) |
| `PUT` | `/api/users/push-token` | Update FCM push token |
| `POST` | `/api/users/block` | Block a user |
| `DELETE` | `/api/users/block/{userId}` | Unblock a user |
| `GET` | `/api/users/blocked` | List blocked users |
| `POST` | `/api/users/avatar` | Upload profile avatar |
| `POST` | `/api/users/push-test` | Send test push notification |
</details>

<details>
<summary><b>💬 Chats</b> — 20+ endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats` | Create chat (private/group/channel) |
| `POST` | `/api/chats/start/{userId}` | Start or find private chat with user |
| `GET` | `/api/chats` | List my chats |
| `GET` | `/api/chats/search?q=` | Search chats by name |
| `GET` | `/api/chats/archived` | List archived chats |
| `GET` | `/api/chats/{id}` | Get chat details |
| `PUT` | `/api/chats/{id}` | Update group info |
| `DELETE` | `/api/chats/{id}` | Delete chat |
| `POST` | `/api/chats/{id}/read` | Mark as read (up to last message) |
| `POST` | `/api/chats/{id}/pin` | Pin chat to top |
| `DELETE` | `/api/chats/{id}/pin` | Unpin chat |
| `POST` | `/api/chats/{id}/archive` | Archive chat |
| `POST` | `/api/chats/{id}/unarchive` | Unarchive chat |
| `POST` | `/api/chats/{id}/hide` | Hide chat from list |
| `POST` | `/api/chats/{id}/leave` | Leave group chat |
| `POST` | `/api/chats/{id}/photo` | Upload chat photo |
| `POST` | `/api/chats/{id}/wallpaper` | Set chat wallpaper |
| `POST` | `/api/chats/{id}/transfer-ownership` | Transfer group ownership |
| `PUT` | `/api/chats/{id}/slow-mode` | Set slow mode interval (0-3600s) |
| `PUT` | `/api/chats/{id}/permissions` | Set chat permissions |
| `PUT` | `/api/chats/{id}/notifications` | Mute/unmute notifications |
| `GET` | `/api/chats/{id}/notifications` | Get notification settings |
| `POST` | `/api/chats/{id}/participants` | Add participant to group |
| `DELETE` | `/api/chats/{id}/participants/{userId}` | Remove participant |
| `PUT` | `/api/chats/{id}/participants/{userId}/role` | Change participant role |
| `POST` | `/api/chats/{id}/promote` | Promote to admin |
| `POST` | `/api/chats/{id}/demote` | Demote from admin |
| `GET` | `/api/chats/{id}/online` | Get online members |
</details>

<details>
<summary><b>✉️ Messages</b> — 20+ endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/{id}/messages` | Send a message |
| `GET` | `/api/chats/{id}/messages` | List messages (paginated) |
| `GET` | `/api/chats/{id}/messages/search?q=` | Search within chat |
| `POST` | `/api/chats/{id}/messages/file` | Upload file attachment |
| `POST` | `/api/chats/{id}/messages/voice` | Send voice message |
| `POST` | `/api/chats/{id}/messages/location` | Send location (lat/lng) |
| `POST` | `/api/chats/{id}/messages/video-circle` | Send circular video |
| `POST` | `/api/chats/{id}/messages/{msgId}/resend` | Re-send a message |
| `GET` | `/api/chats/{id}/pinned` | Get pinned messages |
| `GET` | `/api/chats/{id}/media?type=photo` | Media gallery |
| `GET` | `/api/chats/{id}/polls` | List polls in chat |
| `GET` | `/api/chats/{id}/export?format=csv` | Export chat history as CSV |
| `GET` | `/api/chats/{id}/calls` | Call history for chat |
| `GET` | `/api/chats/{id}/active-calls` | Active calls in chat |
| `GET` | `/api/chats/{id}/voice-chats/active` | Active voice chats |
| `GET` | `/api/chats/{id}/voice-chats/history` | Voice chat history |
| `GET` | `/api/messages/{id}` | Get message by ID |
| `PUT` | `/api/messages/{id}` | Edit message |
| `DELETE` | `/api/messages/{id}` | Delete message (for everyone) |
| `DELETE` | `/api/messages/{id}/for-me` | Delete message (for me only) |
| `GET` | `/api/messages/{id}/history` | Get edit history |
| `POST` | `/api/messages/{id}/reactions` | Add reaction |
| `DELETE` | `/api/messages/{id}/reactions` | Remove reaction |
| `PUT` | `/api/messages/{id}/pin` | Pin/unpin message |
| `POST` | `/api/messages/{id}/star` | Star message |
| `DELETE` | `/api/messages/{id}/star` | Unstar message |
| `POST` | `/api/messages/{id}/read` | Mark as read |
| `POST` | `/api/messages/{id}/save` | Save message |
| `POST` | `/api/messages/{id}/report` | Report message |
| `GET` | `/api/messages/search?q=&offset=&limit=` | Global search |
| `GET` | `/api/messages/starred` | List starred messages |
| `GET` | `/api/messages/scheduled` | List scheduled messages |
| `POST` | `/api/messages/forward` | Forward messages |
| `POST` | `/api/messages/schedule` | Schedule a message |
| `POST` | `/api/messages/read/bulk` | Bulk mark as read |
| `DELETE` | `/api/messages/bulk` | Bulk delete messages |
</details>

<details>
<summary><b>📱 Contacts</b> — 5 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/contacts/sync` | Sync phone contacts |
| `GET` | `/api/contacts` | List contacts |
| `GET` | `/api/contacts/search?q=` | Search contacts by phone |
| `GET` | `/api/contacts/registered` | Contacts on the platform |
| `POST` | `/api/contacts/photo` | Update contact photo |
</details>

<details>
<summary><b>📁 Folders</b> — 4 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/folders` | Create folder |
| `GET` | `/api/folders` | List folders |
| `PUT` | `/api/folders/{id}` | Update folder |
| `DELETE` | `/api/folders/{id}` | Delete folder |
</details>

<details>
<summary><b>🔗 Invite Links</b> — 4 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/{id}/invite-links` | Create invite link |
| `GET` | `/api/chats/{id}/invite-links` | List invite links |
| `DELETE` | `/api/chats/{id}/invite-links/{linkId}` | Revoke invite link |
| `POST` | `/api/chats/join` | Join via invite link |
</details>

<details>
<summary><b>📊 Polls</b> — 4 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/{id}/polls` | Create poll |
| `GET` | `/api/chats/{id}/polls` | List polls |
| `POST` | `/api/polls/{pollId}/vote` | Cast vote |
| `POST` | `/api/polls/{pollId}/close` | Close poll |
</details>

<details>
<summary><b>🎨 Stickers</b> — 8 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/stickers/packs` | Create sticker pack |
| `GET` | `/api/stickers/packs` | List all packs |
| `GET` | `/api/stickers/packs/my` | My sticker packs |
| `GET` | `/api/stickers/packs/{id}` | Get pack by ID |
| `DELETE` | `/api/stickers/packs/{id}` | Delete pack |
| `POST` | `/api/stickers/packs/{id}/stickers` | Add sticker to pack |
| `POST` | `/api/stickers/library` | Add sticker to library |
| `GET` | `/api/stickers/library` | My sticker library |
</details>

<details>
<summary><b>🤖 Bots</b> — 5 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/bots` | Create bot |
| `GET` | `/api/bots` | List my bots |
| `PUT` | `/api/bots/{id}` | Update bot |
| `DELETE` | `/api/bots/{id}` | Delete bot |
| `POST` | `/api/bots/{id}/token` | Regenerate bot token |
</details>

<details>
<summary><b>📞 Calls</b> — 5 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/calls` | Initiate WebRTC call |
| `POST` | `/api/calls/{id}/respond` | Accept/reject call |
| `POST` | `/api/calls/{id}/end` | End call |
| `GET` | `/api/calls/{id}` | Get call details |
| `GET` | `/api/chats/{chatId}/calls` | Call history |
</details>

<details>
<summary><b>👥 Group Calls</b> — 5 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/calls/group/initiate` | Start group call |
| `POST` | `/api/calls/group/respond` | Join/leave/mute |
| `POST` | `/api/calls/group/{id}/end` | End group call |
| `GET` | `/api/calls/group/{id}` | Get group call details |
| `GET` | `/api/chats/{chatId}/active-calls` | Active calls |
</details>

<details>
<summary><b>🎤 Voice Chats</b> — 8 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/chats/{id}/voice-chat` | Create voice chat |
| `GET` | `/api/chats/{id}/voice-chats/active` | Active voice chats |
| `GET` | `/api/chats/{id}/voice-chats/history` | Voice chat history |
| `GET` | `/api/voice-chats/{id}` | Get voice chat details |
| `POST` | `/api/voice-chats/{id}/join` | Join voice chat |
| `POST` | `/api/voice-chats/{id}/leave` | Leave voice chat |
| `POST` | `/api/voice-chats/{id}/end` | End voice chat |
| `POST` | `/api/voice-chats/{id}/mute` | Mute/unmute participant |
</details>

<details>
<summary><b>📸 Stories</b> — 6 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/stories` | Create story |
| `GET` | `/api/stories` | Following stories |
| `GET` | `/api/stories/my` | My active stories |
| `GET` | `/api/stories/{id}` | View story |
| `DELETE` | `/api/stories/{id}` | Delete story |
| `GET` | `/api/stories/{id}/views` | Get story viewers |
</details>

<details>
<summary><b>📺 Channels</b> — 6 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/channels/subscribe` | Subscribe to channel |
| `POST` | `/api/channels/unsubscribe` | Unsubscribe |
| `GET` | `/api/channels` | My channels |
| `GET` | `/api/channels/{id}/subscribers` | List subscribers |
| `GET` | `/api/channels/{id}/subscribed` | Check subscription |
| `PUT` | `/api/channels/{id}/subscribers/{userId}/role` | Set subscriber role |
</details>

<details>
<summary><b>📝 Drafts</b> — 3 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/drafts` | Save draft |
| `GET` | `/api/drafts?chatId=` | Get draft for chat |
| `DELETE` | `/api/drafts/{id}` | Delete draft |
</details>

<details>
<summary><b>⏰ Scheduled Messages</b> — 3 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/messages/schedule` | Schedule message |
| `GET` | `/api/messages/scheduled` | List scheduled |
| `DELETE` | `/api/messages/scheduled/{id}` | Cancel scheduled |
</details>

<details>
<summary><b>⭐ Saved Messages</b> — 3 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/messages/{id}/save` | Save message |
| `GET` | `/api/saved-messages` | List saved |
| `DELETE` | `/api/saved-messages/{id}` | Remove saved |
</details>

<details>
<summary><b>🎭 Custom Emojis</b> — 4 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/emojis` | Upload emoji |
| `GET` | `/api/emojis` | All public emojis |
| `GET` | `/api/emojis/my` | My uploaded emojis |
| `DELETE` | `/api/emojis/{id}` | Delete emoji |
</details>

<details>
<summary><b>🎨 GIFs</b> — 3 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `POST` | `/api/gifs` | Save GIF |
| `GET` | `/api/gifs` | List saved GIFs |
| `DELETE` | `/api/gifs` | Remove GIF |
</details>

<details>
<summary><b>🖥️ Sessions</b> — 3 endpoints</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `GET` | `/api/sessions` | List active sessions |
| `DELETE` | `/api/sessions/{id}` | Terminate session |
| `DELETE` | `/api/sessions` | Terminate all sessions |
</details>

<details>
<summary><b>🛡️ Security & Admin</b> — miscellaneous</summary>

| Method | Endpoint | Description |
|--------|----------|-------------|
| `PUT` | `/api/users/email` | Change email |
| `PUT` | `/api/users/username` | Change username |
| `DELETE` | `/api/account` | Delete account |
| `GET` | `/api/account/settings` | Get account settings |
| `PUT` | `/api/account/settings` | Update account settings |
| `POST` | `/api/verification/email` | Send email verification |
| `POST` | `/api/verification/email/verify` | Verify email code |
| `POST` | `/api/verification/phone` | Send phone verification |
| `POST` | `/api/verification/phone/verify` | Verify phone code |
</details>

---

## 🔌 WebSocket

Connect in real-time to receive instant events and send commands.

```bash
# Connect (replace TOKEN with your JWT)
wscat -c "ws://localhost:8080/ws?token=TOKEN"
```

### 📥 Server → Client Events

| Event | Payload | Description |
|-------|---------|-------------|
| `newMessage` | `{message}` | New message sent |
| `editMessage` | `{message}` | Message edited |
| `deleteMessage` | `{chatId, messageId}` | Message deleted |
| `readMessage` | `{chatId, userId, messageId}` | Message read |
| `pinMessage` | `{message}` | Pin toggled |
| `reaction` | `{message}` | Reaction added/removed |
| `chatCreated` | `{chat}` | Chat created |
| `chatUpdated` | `{chat}` | Chat updated |
| `chatDeleted` | `{chatId}` | Chat deleted |
| `callOffer` | `{chatId, callerId, type}` | Incoming call |
| `callAccept` | `{callId, userId}` | Call accepted |
| `callEnd` | `{callId, userId}` | Call ended |
| `online` | `{userId, online}` | User online status |
| `offline` | `{userId, online}` | User offline status |

### 📤 Client → Server Events

Send JSON over the WebSocket connection:

```json
{"type": "sendMessage", "payload": {"chatId": "...", "content": "Hello!"}}
{"type": "editMessage", "payload": {"messageId": "...", "content": "New text"}}
{"type": "deleteMessage", "payload": {"messageId": "...", "chatId": "..."}}
{"type": "createChat", "payload": {"type": "private", "participantIds": ["..."]}}
{"type": "addReaction", "payload": {"messageId": "...", "emoji": "🔥"}}
```

---

## 📮 Postman Collection

Download the complete Postman collection with all **168 endpoints**:

```
http://localhost:8080/postman
```

Or access it directly from the repository:
`docs/postman_collection.json`

<details>
<summary><b>How to import into Postman</b></summary>

1. Start the server: `./ChatServer.exe`
2. Open browser: `http://localhost:8080/postman`
3. Save the JSON file
4. Open Postman → **File** → **Import** → select the file
5. Set `base_url` variable to `http://localhost:8080`
6. Register/Login → copy JWT → set as `token` variable
</details>

---

## 📊 Response Format

### ✅ Success
```json
{ "success": true, "data": { ... } }
```

### 📄 Paginated
```json
{ "success": true, "data": [ ... ], "meta": { "total": 100, "offset": 0, "limit": 50 } }
```

### ❌ Error
```json
{ "success": false, "error": "Description", "code": "ERROR_CODE" }
```

| Code | Status | Description |
|------|:------:|-------------|
| `BAD_REQUEST` | 400 | Invalid request |
| `VALIDATION_ERROR` | 400 | Field validation failed |
| `UNAUTHORIZED` | 401 | Missing/invalid JWT |
| `TOKEN_EXPIRED` | 401 | JWT expired |
| `FORBIDDEN` | 403 | Insufficient permissions |
| `NOT_FOUND` | 404 | Resource not found |
| `DUPLICATE` | 409 | Already exists |
| `RATE_LIMIT` | 429 | Too many requests |
| `INTERNAL_ERROR` | 500 | Server error |

---

## 🏗️ Architecture

```
internal/
├── domain/              # Domain models (auth, chat, message, user, ...)
│   ├── auth/            #   Register, Login, AdminRegister
│   ├── bot/             #   Bot accounts
│   ├── call/            #   WebRTC calls + group calls
│   ├── channel/         #   Broadcast channels
│   ├── chat/            #   Chats, folders, invite links
│   ├── contact/         #   Phone contacts
│   ├── draft/           #   Drafts + scheduled messages
│   ├── e2e/             #   E2E encryption keys
│   ├── emoji/           #   Custom emojis
│   ├── ipblock/         #   IP block list
│   ├── link/            #   Link previews
│   ├── message/         #   Messages, reactions, read receipts
│   ├── notification/    #   Notification settings
│   ├── poll/            #   Polls & votes
│   ├── report/          #   Abuse reports
│   ├── savedmsg/        #   Saved messages
│   ├── session/         #   User sessions
│   ├── sticker/         #   Sticker packs & library
│   ├── story/           #   Stories & views
│   ├── user/            #   Users, blocks, settings
│   ├── verification/    #   Email/phone verification
│   └── voicechat/       #   Voice chat rooms
├── repository/          # SQLite data access (interface + impl)
│   └──                 #   20+ sub-packages
├── service/             # Business logic layer
│   └──                 #   20+ sub-packages
├── handler/             # HTTP handlers (Gin)
│   ├── ws/              #   WebSocket handler + events
│   └──                 #   18+ sub-packages
├── middleware/          # JWT auth, admin, rate limiting
├── ws/                  # WebSocket hub, client, event bus
├── config/              # Environment configuration
└── database/            # Migrations + initializer

pkg/
└── response/            # JSON response helpers

docs/                    # API documentation
├── swagger.json         #   OpenAPI spec
├── swagger.yaml         #   YAML format
└── postman_collection.json  # Postman collection (168 endpoints)

frontend/                # React SPA (Vite + TypeScript)
└── dist/               #   Built frontend (served at /app/)
```

---

## ⚙️ Configuration

| Variable | Default | Description |
|----------|---------|-------------|
| `SERVER_PORT` | `8080` | HTTP server port |
| `DATABASE_PATH` | `file:chat.db?...` | SQLite connection string (WAL mode) |
| `JWT_SECRET` | `super-secret-change-me` | HMAC-SHA256 key for JWT |
| `ADMIN_SECRET` | `admin-secret-change-me` | Admin registration secret |
| `JWT_TTL` | `86400` | Token lifetime (seconds, 24h) |
| `CORS_ORIGIN` | `*` | Allowed CORS origin |
| `FCM_KEY_PATH` | — | Firebase service account JSON path |
| `ENCRYPTION_KEY` | — | E2E encryption key |

---

## 🐳 Docker

```bash
# Build
docker build -t chatserver .

# Run
docker run -d --name chatserver -p 8080:8080 \
  -v chatserver-data:/app/data \
  -v chatserver-uploads:/app/uploads \
  -e DATABASE_PATH=/app/data/chatserver.db \
  -e JWT_SECRET=change-me-in-production \
  -e GIN_MODE=release \
  chatserver
```

## 🔄 CI/CD

GitHub Actions pipeline (`.github/workflows/deploy.yml`):
1. ✅ Checkout + setup Go
2. ✅ Download dependencies + build
3. ✅ `go vet` + Swagger validation
4. ✅ Build Docker image
5. ✅ Push to container registry

## 🛠 Build

### Prerequisites
- Go 1.21+
- Node.js 18+ (for frontend)

### Commands
```bash
# Backend
go build -o ChatServer.exe .
./ChatServer.exe

# Frontend (optional — only if modifying)
cd frontend
npm install
npm run build    # production build
npm run dev      # dev server with HMR on :5173

# Cross-compile
GOOS=linux GOARCH=amd64 go build -o chat-server-linux .
GOOS=darwin GOARCH=amd64 go build -o chat-server-darwin .

# Tests
go test ./... -v
```

## 🧰 Tech Stack

| Component | Technology |
|-----------|-----------|
| **Language** | Go 1.21+ |
| **Database** | SQLite (WAL, modernc.org/sqlite — no CGO) |
| **HTTP** | Gin framework |
| **WebSocket** | gorilla/websocket |
| **Auth** | JWT (HMAC-SHA256) + bcrypt |
| **Push** | Firebase Cloud Messaging (FCM) |
| **Docs** | Swaggo (OpenAPI 2.0) |
| **Frontend** | React 19 + Vite + TypeScript |
| **Calls** | WebRTC signalling |
| **Real-time** | Custom WebSocket hub (per-client channels) |
