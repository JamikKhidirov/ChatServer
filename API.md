# Chat Messenger Server API v2.0

Все API эндпоинты. Базовый URL: `http://localhost:8080`

---

## Содержание

- [1. Аутентификация (без JWT)](#1-аутентификация-без-jwt)
- [2. Авторизованные эндпоинты (JWT required)](#2-авторизованные-эндпоинты-jwt-required)
  - [2.1 Auth](#21-auth)
  - [2.2 Users / Profile](#22-users--profile)
  - [2.3 Account Settings](#23-account-settings)
  - [2.4 Contacts](#24-contacts)
  - [2.5 Chats](#25-chats)
  - [2.6 Invite Links](#26-invite-links)
  - [2.7 Chat Folders](#27-chat-folders)
  - [2.8 Messages](#28-messages)
  - [2.9 Polls](#29-polls)
  - [2.10 Stickers](#210-stickers)
  - [2.11 Drafts](#211-drafts)
  - [2.12 Sessions](#212-sessions)
  - [2.13 Bots](#213-bots)
  - [2.14 Saved GIFs](#214-saved-gifs)
  - [2.15 Stories](#215-stories)
  - [2.16 Calls](#216-calls)
  - [2.17 Group Calls](#217-group-calls)
  - [2.18 Channels](#218-channels)
  - [2.19 Saved Messages](#219-saved-messages)
  - [2.20 Custom Emojis](#220-custom-emojis)
  - [2.21 Voice Chats](#221-voice-chats)
  - [2.22 Verification](#222-verification)
- [3. Non-API Routes](#3-non-api-routes)
- [4. Типы запросов (Request Bodies)](#4-типы-запросов-request-bodies)
- [5. Типы ответов (Response Bodies)](#5-типы-ответов-response-bodies)
- [6. WebSocket API](#6-websocket-api)
  - [6.1 Подключение](#61-подключение)
  - [6.2 Формат сообщений](#62-формат-сообщений)
  - [6.3 Сервер → Клиент (входящие события)](#63-сервер--клиент-входящие-события)
  - [6.4 Клиент → Сервер (исходящие события)](#64-клиент--сервер-исходящие-события)
  - [6.5 Технические детали](#65-технические-детали)
  - [6.6 Примеры](#66-примеры)

---

## 1. Аутентификация (без JWT)

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/auth/register` | [`RegisterRequest`](#registerrequest) | [`AuthResponse`](#authresponse) |
| `POST` | `/api/auth/admin/register` | [`AdminRegisterRequest`](#adminregisterrequest) | [`AuthResponse`](#authresponse) |
| `POST` | `/api/auth/login` | [`LoginRequest`](#loginrequest) | [`AuthResponse`](#authresponse) |
| `POST` | `/api/auth/login/email` | `{ "email": "string" }` | `{ "message": "code sent" }` |
| `POST` | `/api/auth/login/email/verify` | `{ "email": "string", "code": "string" }` | `{ "token": "..." }` |
| `POST` | `/api/auth/login/phone` | `{ "phone": "string" }` | `{ "message": "code sent" }` |
| `POST` | `/api/auth/login/phone/verify` | `{ "phone": "string", "code": "string" }` | `{ "token": "..." }` |

### POST /api/auth/register
```json
{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "secret123",
    "display_name": "John Doe"
}
```
Ответ:
```json
{
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
        "id": "uuid",
        "username": "john_doe",
        "displayName": "John Doe",
        "avatarUrl": "",
        "bio": "",
        "status": "",
        "online": false,
        "lastSeen": "2026-01-01T00:00:00Z",
        "isAdmin": false
    }
}
```

### POST /api/auth/login
```json
{
    "email": "john@example.com",
    "password": "secret123"
}
```

---

## 2. Авторизованные эндпоинты (JWT required)

Заголовок: `Authorization: Bearer {token}`

### 2.1 Auth

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/auth/refresh` | — | `{ "token": "..." }` |
| `PUT` | `/api/auth/change-password` | [`ChangePasswordRequest`](#changepasswordrequest) | `{ "message": "password changed" }` |

### 2.2 Users / Profile

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/users/profile` | — | [`UserResponse`](#userresponse) |
| `PUT` | `/api/users/profile` | [`UpdateProfileRequest`](#updateprofilerequest) | [`UserResponse`](#userresponse) |
| `PUT` | `/api/users/username` | `{ "username": "..." }` | [`UserResponse`](#userresponse) |
| `PUT` | `/api/users/email` | `{ "email": "..." }` | [`UserResponse`](#userresponse) |
| `DELETE` | `/api/users/account` | — | `{ "message": "account deleted" }` |
| `PUT` | `/api/users/push-token` | [`UpdatePushTokenRequest`](#updatepushtokenrequest) | `{ "message": "push token updated" }` |
| `POST` | `/api/users/push-test` | `{ "title": "...", "body": "..." }` | `{ "message": "test push sent" }` |
| `PUT` | `/api/users/status` | [`UpdateStatusRequest`](#updatestatusrequest) | [`UserResponse`](#userresponse) |
| `GET` | `/api/users/search?q=&limit=&offset=` | — | [`APIResponse`](#apiresponse) paginated |
| `GET` | `/api/users/:id` | — | [`UserResponse`](#userresponse) |
| `GET` | `/api/users/username/:username` | — | [`UserResponse`](#userresponse) |
| `GET` | `/api/users/:id/last-seen` | — | `{ "userId": "...", "username": "...", "online": bool, "lastSeen": "..." }` |
| `POST` | `/api/users/block` | [`BlockUserRequest`](#blockuserrequest) | `{ "message": "user blocked" }` |
| `DELETE` | `/api/users/block/:userId` | — | `{ "message": "user unblocked" }` |
| `GET` | `/api/users/blocked` | — | `[UserResponse]` |
| `POST` | `/api/users/avatar` | `multipart: avatar (file)` | [`UserResponse`](#userresponse) |

### 2.3 Account Settings

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/account/settings` | — | [`AccountSetting`](#accountsetting) |
| `PUT` | `/api/account/settings` | [`UpdateAccountSettingRequest`](#updateaccountsettingrequest) | [`AccountSetting`](#accountsetting) |

### 2.4 Contacts

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/contacts/sync` | [`SyncContactsRequest`](#synccontactsrequest) | `{ "message": "contacts synced" }` |
| `GET` | `/api/contacts` | — | `[ContactResponse]` |
| `GET` | `/api/contacts/search?q=` | — | `[ContactResponse]` |
| `GET` | `/api/contacts/registered` | — | `[UserResponse]` |
| `POST` | `/api/contacts/photo` | [`UpdateContactPhotoRequest`](#updatecontactphotorequest) | `{ "message": "contact photo updated" }` |

### 2.5 Chats

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/chats` | — | `[ChatResponse]` |
| `GET` | `/api/chats/search?q=` | — | `[ChatResponse]` |
| `GET` | `/api/chats/archived` | — | `[ChatResponse]` |
| `POST` | `/api/chats` | [`CreateChatRequest`](#createchatrequest) | [`ChatResponse`](#chatresponse) |
| `POST` | `/api/chats/start/:userId` | — | [`ChatResponse`](#chatresponse) |
| `GET` | `/api/chats/:id` | — | [`ChatResponse`](#chatresponse) |
| `PUT` | `/api/chats/:id` | [`UpdateGroupRequest`](#updategrouprequest) | `{ "message": "group updated" }` |
| `DELETE` | `/api/chats/:id` | — | `{ "message": "chat deleted" }` |
| `POST` | `/api/chats/:id/participants` | [`AddParticipantRequest`](#addparticipantrequest) | `{ "message": "participant added" }` |
| `DELETE` | `/api/chats/:id/participants/:userId` | — | `{ "message": "participant removed" }` |
| `PUT` | `/api/chats/:id/participants/:userId/role` | `{ "role": "admin\|member" }` | `{ "message": "role updated to ..." }` |
| `POST` | `/api/chats/:id/leave` | — | `{ "message": "left the group" }` |
| `POST` | `/api/chats/:id/read` | — | `{ "message": "marked as read" }` |
| `POST` | `/api/chats/:id/pin` | — | `{ "message": "chat pinned" }` |
| `DELETE` | `/api/chats/:id/pin` | — | `{ "message": "chat unpinned" }` |
| `POST` | `/api/chats/:id/archive` | — | `{ "message": "chat archived" }` |
| `POST` | `/api/chats/:id/unarchive` | — | `{ "message": "chat unarchived" }` |
| `POST` | `/api/chats/:id/hide` | — | `{ "message": "chat hidden" }` |
| `POST` | `/api/chats/:id/transfer-ownership` | `{ "userId": "..." }` | `{ "message": "ownership transferred" }` |
| `POST` | `/api/chats/:id/promote` | `{ "userId": "..." }` | `{ "message": "user promoted to admin" }` |
| `POST` | `/api/chats/:id/demote` | `{ "userId": "..." }` | `{ "message": "user demoted to member" }` |
| `POST` | `/api/chats/:id/photo` | `multipart: photo (file)` | `{ "photoUrl": "..." }` |
| `POST` | `/api/chats/:id/wallpaper` | `multipart: wallpaper (file)` | `{ "wallpaperUrl": "..." }` |
| `PUT` | `/api/chats/:id/permissions` | `{ "whoCanSend": "...", "whoCanAdd": "..." }` | `{ "chatId": "...", "whoCanSend": "...", "whoCanAdd": "..." }` |
| `GET` | `/api/chats/:id/online` | — | `{ "userIds": [...] }` |
| `PUT` | `/api/chats/:id/notifications` | [`UpdateNotificationSettingRequest`](#updatenotificationsettingrequest) | `{ "message": "notifications muted\|unmuted" }` |
| `GET` | `/api/chats/:id/notifications` | — | `{ "muted": bool }` |
| `PUT` | `/api/chats/:id/slow-mode` | `{ "seconds": 0-3600 }` | `{ "message": "slow mode updated" }` |

### 2.6 Invite Links

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/chats/:id/invite-links` | [`CreateInviteLinkRequest`](#createinvitelinkrequest) | [`InviteLink`](#invitelink) |
| `GET` | `/api/chats/:id/invite-links` | — | `[InviteLink]` |
| `DELETE` | `/api/chats/:id/invite-links/:linkId` | — | `{ "message": "invite link deleted" }` |
| `POST` | `/api/chats/join` | `{ "code": "string" }` | `{ "message": "joined chat successfully" }` |

### 2.7 Chat Folders

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/folders` | — | `[ChatFolder]` |
| `POST` | `/api/folders` | [`CreateChatFolderRequest`](#createchatfolderrequest) | [`ChatFolder`](#chatfolder) |
| `PUT` | `/api/folders/:id` | [`UpdateChatFolderRequest`](#updatechatfolderrequest) | [`ChatFolder`](#chatfolder) |
| `DELETE` | `/api/folders/:id` | — | `{ "message": "folder deleted" }` |

### 2.8 Messages

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/chats/:id/messages?limit=&offset=` | — | [`APIResponse`](#apiresponse) paginated |
| `GET` | `/api/chats/:id/messages/search?q=&limit=&offset=` | — | [`APIResponse`](#apiresponse) paginated |
| `GET` | `/api/chats/:id/media?type=&limit=&offset=` | — | [`APIResponse`](#apiresponse) paginated |
| `POST` | `/api/chats/:id/messages` | [`SendMessageRequest`](#sendmessagerequest) | [`MessageResponse`](#messageresponse) |
| `POST` | `/api/chats/:id/messages/file` | `multipart: file (file), replyToId (string)` | [`MessageResponse`](#messageresponse) |
| `POST` | `/api/chats/:id/messages/voice` | `multipart: voice (file)` | [`MessageResponse`](#messageresponse) |
| `POST` | `/api/chats/:id/messages/:msgId/resend` | — | [`MessageResponse`](#messageresponse) |
| `GET` | `/api/chats/:id/pinned` | — | `[MessageResponse]` |
| `GET` | `/api/chats/:id/export` | — | `[MessageResponse]` |
| `GET` | `/api/messages/search?q=&limit=&offset=` | — | [`APIResponse`](#apiresponse) paginated |
| `GET` | `/api/messages/starred` | — | `[StarredMessageResponse]` |
| `POST` | `/api/messages/read/bulk` | `{ "messageIds": [...], "chatId": "..." }` | `{ "message": "...", "count": N }` |
| `DELETE` | `/api/messages/bulk` | `{ "messageIds": [...] }` | `{ "message": "...", "count": N }` |
| `POST` | `/api/messages/forward` | [`ForwardMessageRequest`](#forwardmessagerequest) | [`MessageResponse`](#messageresponse) |
| `POST` | `/api/messages/schedule` | [`ScheduleMessageRequest`](#dualmessagerequest) | [`ScheduledMessage`](#scheduledmessage) |
| `GET` | `/api/messages/scheduled` | — | `[ScheduledMessage]` |
| `DELETE` | `/api/messages/scheduled/:id` | — | `{ "message": "scheduled message cancelled" }` |
| `GET` | `/api/messages/:id` | — | [`MessageResponse`](#messageresponse) |
| `PUT` | `/api/messages/:id` | [`EditMessageRequest`](#editmessagerequest) | [`MessageResponse`](#messageresponse) |
| `DELETE` | `/api/messages/:id` | — | `{ "message": "message deleted" }` |
| `DELETE` | `/api/messages/:id/for-me` | — | `{ "message": "message deleted for you" }` |
| `POST` | `/api/messages/:id/reactions` | [`AddReactionRequest`](#addreactionrequest) | [`MessageResponse`](#messageresponse) |
| `DELETE` | `/api/messages/:id/reactions?emoji=` | — | [`MessageResponse`](#messageresponse) |
| `PUT` | `/api/messages/:id/pin` | [`PinMessageRequest`](#pinmessagerequest) | [`MessageResponse`](#messageresponse) |
| `POST` | `/api/messages/:id/star` | — | [`MessageResponse`](#messageresponse) |
| `DELETE` | `/api/messages/:id/star` | — | `{ "message": "message unstarred" }` |
| `POST` | `/api/messages/:id/self-destruct` | `{ "seconds": N }` | `{ "message": "self-destruct timer set" }` |
| `POST` | `/api/messages/:id/read` | — | `{ "message": "marked as read" }` |
| `POST` | `/api/messages/:id/report` | `{ "reason": "..." }` | `{ "messageId": "...", "reason": "...", "status": "reported" }` |
| `GET` | `/api/messages/:id/history` | — | [`MessageResponse`](#messageresponse) |
| `GET` | `/api/files/:filename` | — | Binary file stream |
| `POST` | `/api/chats/:id/messages/video-circle` | `multipart: video (file), caption (string)` | [`MessageResponse`](#messageresponse) |
| `POST` | `/api/chats/:id/messages/location` | [`SendLocationRequest`](#sendlocationrequest) | [`MessageResponse`](#messageresponse) |

### 2.9 Polls

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/chats/:id/polls` | [`CreatePollRequest`](#createpollrequest) | [`PollWithResults`](#pollwithresults) |
| `GET` | `/api/chats/:id/polls` | — | `[PollWithResults]` |
| `POST` | `/api/polls/:pollId/vote` | [`VotePollRequest`](#votepollrequest) | `{ "message": "vote recorded" }` |
| `POST` | `/api/polls/:pollId/close` | — | `{ "message": "poll closed" }` |

### 2.10 Stickers

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/stickers/packs` | — | `[StickerPack]` |
| `GET` | `/api/stickers/packs/my` | — | `[StickerPack]` |
| `POST` | `/api/stickers/packs` | [`CreateStickerPackRequest`](#createsticherpackrequest) | [`StickerPack`](#stickerpack) |
| `GET` | `/api/stickers/packs/:id` | — | [`StickerPackWithStickers`](#stickerpackwithstickers) |
| `POST` | `/api/stickers/packs/:id/stickers` | [`AddStickerRequest`](#addstickerrequest) | [`Sticker`](#sticker) |
| `DELETE` | `/api/stickers/packs/:id` | — | `{ "message": "pack deleted" }` |
| `GET` | `/api/stickers/library` | — | `[Sticker]` |
| `POST` | `/api/stickers/library` | `{ "stickerId": "..." }` | `{ "message": "sticker added to library" }` |

### 2.11 Drafts

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/drafts` | [`SaveDraftRequest`](#savedraftrequest) | [`Draft`](#draft) |
| `GET` | `/api/drafts?chatId=` | — | [`Draft`](#draft) |
| `DELETE` | `/api/drafts/:id` | — | `{ "message": "draft deleted" }` |

### 2.12 Sessions

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `GET` | `/api/sessions` | — | `[Session]` |
| `DELETE` | `/api/sessions/:id` | — | `{ "message": "session terminated" }` |
| `DELETE` | `/api/sessions` | — | `{ "message": "all other sessions terminated" }` |

### 2.13 Bots

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/bots` | [`CreateBotRequest`](#createbotrequest) | [`Bot`](#bot) |
| `GET` | `/api/bots` | — | `[Bot]` |
| `PUT` | `/api/bots/:id` | [`UpdateBotRequest`](#updatebotrequest) | `{ "message": "bot updated" }` |
| `DELETE` | `/api/bots/:id` | — | `{ "message": "bot deleted" }` |
| `POST` | `/api/bots/:id/regenerate-token` | — | `{ "message": "token regenerated" }` |

### 2.14 Saved GIFs

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/gifs` | `{ "url": "..." }` | `{ "message": "gif saved" }` |
| `GET` | `/api/gifs` | — | `[string]` (URLs) |
| `DELETE` | `/api/gifs` | `{ "url": "..." }` | `{ "message": "gif deleted" }` |

### 2.15 Stories

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/stories` | `multipart: type, caption, file` | [`StoryResponse`](#storyresponse) |
| `GET` | `/api/stories` | — | `[StoryResponse]` |
| `GET` | `/api/stories/my` | — | `[StoryResponse]` |
| `GET` | `/api/stories/:id` | — | [`StoryResponse`](#storyresponse) |
| `DELETE` | `/api/stories/:id` | — | `{ "message": "story deleted" }` |
| `GET` | `/api/stories/:id/views` | — | `[StoryView]` |

### 2.16 Calls

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/calls/initiate` | [`InitiateCallRequest`](#initiatecallrequest) | [`Call`](#call) |
| `POST` | `/api/calls/:id/respond` | [`RespondCallRequest`](#respondcallrequest) | `{ "message": "call accepted\|rejected" }` |
| `POST` | `/api/calls/:id/end` | — | `{ "message": "call ended" }` |
| `GET` | `/api/calls/:id` | — | [`CallResponse`](#callresponse) |
| `GET` | `/api/calls/history/:chatId` | — | `[Call]` |

### 2.17 Group Calls

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/calls/group/initiate` | [`GroupCallInitiateRequest`](#groupcallinitiaterequest) | [`GroupCallResponse`](#groupcallresponse) |
| `POST` | `/api/calls/group/respond` | [`GroupCallActionRequest`](#groupcallactionrequest) | `{ "message": "action ... completed" }` |
| `POST` | `/api/calls/group/:id/end` | — | `{ "message": "call ended" }` |
| `GET` | `/api/calls/group/:id` | — | [`GroupCallResponse`](#groupcallresponse) |
| `GET` | `/api/chats/:id/active-calls` | — | `[GroupCallResponse]` |

### 2.18 Channels

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/channels/subscribe` | `{ "channelId": "..." }` | `{ "message": "subscribed" }` |
| `POST` | `/api/channels/unsubscribe` | `{ "channelId": "..." }` | `{ "message": "unsubscribed" }` |
| `GET` | `/api/channels` | — | `[ChatResponse]` |
| `GET` | `/api/channels/:id/subscribers` | — | `[ChannelSubscriber]` |
| `GET` | `/api/channels/:id/subscribed` | — | `{ "subscribed": bool }` |
| `PUT` | `/api/channels/:id/subscribers/:userId/role` | `{ "role": "admin\|subscriber" }` | `{ "message": "role updated" }` |

### 2.19 Saved Messages

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/messages/:id/save?chatId=` | — | [`SavedMessageResponse`](#savedmessageresponse) |
| `GET` | `/api/saved-messages?limit=&offset=` | — | [`APIResponse`](#apiresponse) paginated |
| `DELETE` | `/api/saved-messages/:id` | — | `{ "message": "..." }` |

### 2.20 Custom Emojis

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/emojis` | `multipart: shortcode, emoji (file)` | [`CustomEmojiResponse`](#customemojiresponse) |
| `GET` | `/api/emojis` | — | `[CustomEmojiResponse]` |
| `GET` | `/api/emojis/my` | — | `[CustomEmojiResponse]` |
| `DELETE` | `/api/emojis/:id` | — | `{ "message": "..." }` |

### 2.21 Voice Chats

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/chats/:id/voice-chat` | [`CreateVoiceChatRequest`](#createvoicechatrequest) | [`VoiceChatResponse`](#voicechatresponse) |
| `GET` | `/api/chats/:id/voice-chats/active` | — | `[VoiceChatResponse]` |
| `GET` | `/api/chats/:id/voice-chats/history` | — | `[VoiceChatResponse]` |
| `GET` | `/api/voice-chats/:id` | — | [`VoiceChatResponse`](#voicechatresponse) |
| `POST` | `/api/voice-chats/:id/join` | — | `{ "message": "..." }` |
| `POST` | `/api/voice-chats/:id/leave` | — | `{ "message": "..." }` |
| `POST` | `/api/voice-chats/:id/end` | — | `{ "message": "..." }` |
| `POST` | `/api/voice-chats/:id/mute` | `{ "muted": bool }` | `{ "message": "..." }` |

### 2.22 Verification

| Method | Route | Request Body | Ответ |
|--------|-------|-------------|-------|
| `POST` | `/api/verification/email/send` | `{ "email": "..." }` | `{ "message": "code sent" }` |
| `POST` | `/api/verification/email/verify` | `{ "code": "..." }` | `{ "message": "email verified" }` |
| `POST` | `/api/verification/phone/send` | `{ "phone": "..." }` | `{ "message": "code sent" }` |
| `POST` | `/api/verification/phone/verify` | `{ "code": "..." }` | `{ "message": "phone verified" }` |

---

## 3. Non-API Routes

| Method | Route | Описание |
|--------|-------|----------|
| `GET` | `/ws?token={jwt}` | [WebSocket endpoint](#6-websocket-api) |
| `GET` | `/health` | `{ "status": "ok", "version": "2.0.0" }` |
| `GET` | `/swagger/*any` | Swagger UI |
| `GET` | `/postman` | Скачать Postman коллекцию |
| `GET` | `/` | Frontend SPA |
| `GET` | `/app` | Frontend SPA |

---

## 4. Типы запросов (Request Bodies)

### RegisterRequest
```json
{
    "username": "string (3-32 chars, required)",
    "email": "string (valid email, required)",
    "password": "string (min 6 chars, required)",
    "display_name": "string (1-64 chars, required)"
}
```

### AdminRegisterRequest
```json
{
    "username": "string (3-32 chars, required)",
    "email": "string (valid email, required)",
    "password": "string (min 6 chars, required)",
    "display_name": "string (1-64 chars, required)",
    "admin_secret": "string (required)"
}
```

### LoginRequest
```json
{
    "email": "string (required)",
    "password": "string (required)"
}
```

### ChangePasswordRequest
```json
{
    "oldPassword": "string (required)",
    "newPassword": "string (min 6 chars, required)"
}
```

### UpdateProfileRequest
```json
{
    "displayName": "string (optional, 1-64 chars)",
    "username": "string (optional, 3-32 chars)",
    "email": "string (optional, valid email)",
    "bio": "string (optional, max 256)",
    "phone": "string (optional, max 20)",
    "gender": "string (optional: male/female/other)",
    "dateOfBirth": "string (optional, max 10)",
    "avatarUrl": "string (optional)"
}
```

### UpdateStatusRequest
```json
{ "status": "string (max 100 chars)" }
```

### UpdatePushTokenRequest
```json
{
    "token": "string (required)",
    "provider": "string (required: fcm/apns)"
}
```

### BlockUserRequest
```json
{ "blockedId": "string (required)" }
```

### UpdateAccountSettingRequest
```json
{
    "language": "string (optional)",
    "theme": "string (optional)",
    "notifications": "bool (optional)",
    "soundEnabled": "bool (optional)",
    "lastSeenMode": "string (optional)"
}
```

### CreateChatRequest
```json
{
    "name": "string (optional, max 64, обяз. для group/channel)",
    "type": "string (required: private/group/channel)",
    "participantIds": ["uuid1", "uuid2", "..."],
    "description": "string (optional, max 512)"
}
```

### UpdateGroupRequest
```json
{
    "name": "string (optional, 1-64)",
    "description": "string (optional, max 512)",
    "avatarUrl": "string (optional)"
}
```

### AddParticipantRequest
```json
{ "userId": "string (required)" }
```

### CreateInviteLinkRequest
```json
{
    "expiresInMins": "int (optional)",
    "usageLimit": "int (optional)"
}
```

### CreateChatFolderRequest
```json
{
    "name": "string (required)",
    "emoji": "string (optional)",
    "chatIds": ["uuid1", "uuid2", "..."]
}
```

### UpdateChatFolderRequest
```json
{
    "name": "string (optional)",
    "emoji": "string (optional)",
    "order": "int (optional)",
    "chatIds": ["uuid1", "uuid2", "..."]
}
```

### SendMessageRequest
```json
{
    "content": "string (required)",
    "type": "string (required: text/image/file/gif/voice/video/audio/system/location)",
    "replyToId": "uuid (optional)",
    "forwardMsgId": "uuid (optional)",
    "latitude": "float (optional)",
    "longitude": "float (optional)",
    "locationTitle": "string (optional)",
    "effect": "string (optional: confetti/fireworks/hearts/balloons/stars/пусто)"
}
```

### SendLocationRequest
```json
{
    "latitude": "float (required)",
    "longitude": "float (required)",
    "title": "string (optional)",
    "replyToId": "uuid (optional)",
    "effect": "string (optional)"
}
```

### EditMessageRequest
```json
{ "content": "string (required)" }
```

### PinMessageRequest
```json
{ "pin": "bool (true = закрепить, false = открепить)" }
```

### ForwardMessageRequest
```json
{
    "messageId": "uuid (required)",
    "fromChatId": "uuid (required)",
    "toChatId": "uuid (required)"
}
```

### AddReactionRequest
```json
{ "emoji": "string (required: 👍/❤️/😆/😮/😢/🙏)" }
```

### CreatePollRequest
```json
{
    "chatId": "uuid (required)",
    "question": "string (required)",
    "options": ["option1", "option2", "..."],
    "isAnonymous": "bool (optional)",
    "multipleChoice": "bool (optional)",
    "expiresInMins": "int (optional)"
}
```

### VotePollRequest
```json
{ "optionIndex": "int (required)" }
```

### ScheduleMessageRequest
```json
{
    "chatId": "uuid (required)",
    "content": "string (required)",
    "type": "string (required: text/image/file/gif/voice/video)",
    "scheduledAt": "string (required, RFC3339)",
    "replyToId": "uuid (optional)"
}
```

### CreateBotRequest
```json
{
    "name": "string (required)",
    "webhookUrl": "string (optional)"
}
```

### UpdateBotRequest
```json
{
    "name": "string (optional)",
    "avatarUrl": "string (optional)",
    "webhookUrl": "string (optional)"
}
```

### SyncContactsRequest
```json
{
    "contacts": [
        { "phone": "+1234567890", "name": "Friend Name" }
    ]
}
```

### InitiateCallRequest
```json
{
    "chatId": "uuid (required)",
    "type": "string (required: audio/video)"
}
```

### RespondCallRequest
```json
{ "action": "string (required: accept/reject)" }
```

### CreateVoiceChatRequest
```json
{
    "title": "string (optional)",
    "scheduledInMins": "int (optional)"
}
```

---

## 5. Типы ответов (Response Bodies)

### APIResponse (стандартная обёртка)
```json
{
    "success": "bool",
    "data": "object (optional)",
    "error": "string (optional)",
    "code": "string (optional: BAD_REQUEST, UNAUTHORIZED и т.д.)",
    "meta": {
        "total": "int",
        "offset": "int",
        "limit": "int"
    }
}
```

### AuthResponse
```json
{
    "token": "string (JWT токен)",
    "user": {
        "id": "uuid",
        "username": "string",
        "displayName": "string",
        "avatarUrl": "string",
        "bio": "string",
        "status": "string",
        "online": "bool",
        "lastSeen": "datetime",
        "isAdmin": "bool"
    }
}
```

### UserResponse
```json
{
    "id": "uuid",
    "username": "string",
    "displayName": "string",
    "avatarUrl": "string",
    "bio": "string",
    "phone": "string (optional)",
    "gender": "string (optional)",
    "dateOfBirth": "string (optional)",
    "status": "string",
    "online": "bool",
    "lastSeen": "datetime",
    "isAdmin": "bool"
}
```

### ChatResponse
```json
{
    "id": "uuid",
    "name": "string",
    "description": "string",
    "avatarUrl": "string",
    "type": "string (private/group/channel)",
    "createdBy": "uuid",
    "participants": [
        { "UserResponse": "..." }
    ],
    "lastMessage": { "MessageResponse": "..." },
    "unreadCount": "int",
    "createdAt": "datetime"
}
```

### MessageResponse
```json
{
    "id": "uuid",
    "chatId": "uuid",
    "sender": { "UserResponse": "..." },
    "content": "string",
    "type": "string (text/image/file/gif/voice/video/video_circle/audio/system/location)",
    "replyTo": { "MessageResponse": "..." },
    "forwardFrom": { "UserResponse": "..." },
    "fileName": "string",
    "fileSize": "int64",
    "fileUrl": "string",
    "caption": "string",
    "mimeType": "string",
    "duration": "int",
    "width": "int",
    "height": "int",
    "latitude": "float",
    "longitude": "float",
    "locationTitle": "string",
    "effect": "string",
    "reactions": [
        {
            "messageId": "uuid",
            "userId": "uuid",
            "emoji": "string",
            "createdAt": "datetime",
            "user": { "UserResponse": "..." }
        }
    ],
    "pinned": "bool",
    "readBy": [{ "UserResponse": "..." }],
    "createdAt": "datetime",
    "updatedAt": "datetime",
    "edited": "bool",
    "deleted": "bool"
}
```

### PollWithResults
```json
{
    "id": "uuid",
    "chatId": "uuid",
    "creatorId": "uuid",
    "question": "string",
    "options": [
        { "text": "string", "votes": "int" }
    ],
    "isAnonymous": "bool",
    "multipleChoice": "bool",
    "expiresAt": "datetime (optional)",
    "createdAt": "datetime",
    "closed": "bool",
    "totalVotes": "int",
    "votedOption": "int (optional)"
}
```

### CallResponse
```json
{
    "id": "uuid",
    "chatId": "uuid",
    "caller": { "UserResponse": "..." },
    "callee": { "UserResponse": "..." },
    "type": "string (audio/video)",
    "status": "string (initiated/ongoing/ended/missed/rejected)",
    "startedAt": "datetime",
    "endedAt": "datetime (optional)",
    "duration": "int (optional)"
}
```

### VoiceChatResponse
```json
{
    "id": "uuid",
    "chatId": "uuid",
    "startedBy": "uuid",
    "title": "string (optional)",
    "status": "string (active/scheduled/ended)",
    "participantCount": "int",
    "participants": ["uuid1", "uuid2"],
    "scheduledAt": "datetime (optional)",
    "startedAt": "datetime (optional)",
    "endedAt": "datetime (optional)",
    "createdAt": "datetime"
}
```

### Error Codes
```
BAD_REQUEST, UNAUTHORIZED, FORBIDDEN, NOT_FOUND, INTERNAL_ERROR,
VALIDATION_ERROR, DUPLICATE, RATE_LIMIT, BLOCKED, ACCESS_DENIED,
USER_NOT_FOUND, CHAT_NOT_FOUND, MESSAGE_NOT_FOUND, INVALID_TOKEN,
TOKEN_EXPIRED, USERNAME_TAKEN, EMAIL_TAKEN, WEAK_PASSWORD
```

---

## 6. WebSocket API

### 6.1 Подключение

```
ws://localhost:8080/ws?token={JWT_TOKEN}
```

JWT токен получить через `POST /api/auth/login` или `POST /api/auth/register`.

### 6.2 Формат сообщений

Универсальный формат для всех сообщений:

**Входящие (клиент → сервер):**
```json
{
    "type": "название_события",
    "payload": { ... поля ... }
}
```

**Исходящие (сервер → клиент):**
```json
{
    "type": "название_события",
    "payload": { ... поля ... }
}
```

Поле `type` — строка, указывающая тип события.
Поле `payload` — объект с данными события.

### 6.3 Сервер → Клиент (входящие события)

Сервер автоматически рассылает эти события участникам чата в реальном времени.

#### 6.3.1 Сообщения

| Тип | Описание | Payload |
|-----|----------|---------|
| `message:new` | Новое сообщение | Полный [`MessageResponse`](#messageresponse) |
| `message:edited` | Сообщение отредактировано | Полный [`MessageResponse`](#messageresponse) |
| `message:deleted` | Сообщение удалено | `{ "messageId": "uuid", "chatId": "uuid" }` |
| `message:read` | Сообщение прочитано | `{ "messageId": "uuid", "userId": "uuid", "chatId": "uuid" }` |
| `message:reaction` | Реакция изменена | Полный [`MessageResponse`](#messageresponse) с `reactions` |
| `message:pinned` | Сообщение закреплено/откреплено | Полный [`MessageResponse`](#messageresponse) |
| `message:starred` | Добавлено в избранное | `{ "messageId": "uuid", "userId": "uuid" }` |
| `message:forward` | Сообщение переслано | Полный [`MessageResponse`](#messageresponse) |

**Пример `message:new`:**
```json
{
    "type": "message:new",
    "payload": {
        "id": "uuid",
        "chatId": "uuid",
        "senderId": "uuid",
        "content": "Привет!",
        "type": "text",
        "createdAt": "2026-01-01T00:00:00Z"
    }
}
```

**Пример `message:deleted`:**
```json
{
    "type": "message:deleted",
    "payload": {
        "messageId": "uuid",
        "chatId": "uuid"
    }
}
```

**Пример `message:read`:**
```json
{
    "type": "message:read",
    "payload": {
        "messageId": "uuid",
        "userId": "uuid",
        "chatId": "uuid"
    }
}
```

#### 6.3.2 Пользователи

| Тип | Описание | Payload |
|-----|----------|---------|
| `user:online` | Пользователь стал онлайн | `{ "userId": "uuid", "online": true }` |
| `user:offline` | Пользователь стал офлайн | `{ "userId": "uuid", "online": false }` |
| `user:typing` | Пользователь печатает | `{ "chatId": "uuid", "userId": "uuid" }` |
| `user:stop_typing` | Перестал печатать | `{ "chatId": "uuid", "userId": "uuid" }` |
| `user:keyboard_opened` | Клавиатура открыта (моб.) | `{ "chatId": "uuid", "userId": "uuid" }` |
| `user:keyboard_closed` | Клавиатура закрыта | `{ "chatId": "uuid", "userId": "uuid" }` |

**Пример `user:online`:**
```json
{
    "type": "user:online",
    "payload": {
        "userId": "uuid",
        "online": true
    }
}
```

**Пример `user:typing`:**
```json
{
    "type": "user:typing",
    "payload": {
        "chatId": "uuid",
        "userId": "uuid"
    }
}
```

#### 6.3.3 Чаты

| Тип | Описание | Payload |
|-----|----------|---------|
| `chat:created` | Создан новый чат | Полный [`ChatResponse`](#chatresponse) |
| `chat:updated` | Чат обновлён | Полный [`ChatResponse`](#chatresponse) |
| `chat:deleted` | Чат удалён | `{ "chatId": "uuid" }` |

#### 6.3.4 Звонки (WebRTC)

| Тип | Описание | Payload |
|-----|----------|---------|
| `call:offer` | Входящий звонок | `{ "chatId": "uuid", "callId": "uuid", "callerId": "uuid", "type": "audio\|video" }` |
| `call:answer` | Ответ на звонок | `{ "callId": "uuid", "userId": "uuid" }` |
| `call:ice` | ICE-кандидат | `{ "callId": "uuid", "candidate": "string" }` |
| `call:end` | Звонок завершён | `{ "callId": "uuid", "userId": "uuid" }` |
| `call:reject` | Звонок отклонён | `{ "callId": "uuid", "userId": "uuid" }` |
| `call:accept` | Звонок принят | `{ "callId": "uuid", "userId": "uuid" }` |
| `call:missed` | Пропущенный звонок | `{ "callId": "uuid", "userId": "uuid" }` |

### 6.4 Клиент → Сервер (исходящие события)

Отправляй эти события через WebSocket, чтобы выполнить действие.

#### 6.4.1 Отправка и управление сообщениями

**`message:send`** — Отправить сообщение в чат.
```json
{
    "type": "message:send",
    "payload": {
        "chatId": "uuid",
        "content": "Привет!",
        "type": "text",
        "replyToId": "uuid"
    }
}
```
| Поле | Обязательное | Описание |
|------|-------------|----------|
| chatId | ✅ | ID чата |
| content | ✅ | Текст / URL файла |
| type | ✅ | text, image, file, gif, voice, video, audio, location |
| replyToId | ❌ | ID сообщения, на которое отвечаем |

**`message:edit`** — Редактировать сообщение.
```json
{
    "type": "message:edit",
    "payload": {
        "messageId": "uuid",
        "content": "Новый текст"
    }
}
```

**`message:delete`** — Удалить сообщение.
```json
{
    "type": "message:delete",
    "payload": {
        "messageId": "uuid",
        "chatId": "uuid"
    }
}
```

**`message:read`** — Отметить сообщение как прочитанное.
```json
{
    "type": "message:read",
    "payload": {
        "messageId": "uuid",
        "chatId": "uuid"
    }
}
```

**`message:react`** — Добавить реакцию. emoji: 👍 ❤️ 😆 😮 😢 🙏
```json
{
    "type": "message:react",
    "payload": {
        "messageId": "uuid",
        "emoji": "👍"
    }
}
```

**`message:unreact`** — Удалить реакцию.
```json
{
    "type": "message:unreact",
    "payload": {
        "messageId": "uuid",
        "emoji": "👍"
    }
}
```

**`message:pin`** — Закрепить/открепить сообщение.
```json
{
    "type": "message:pin",
    "payload": {
        "messageId": "uuid",
        "pin": true
    }
}
```
`pin: true` — закрепить, `false` — открепить.

**`message:star`** — Добавить в избранное.
```json
{
    "type": "message:star",
    "payload": {
        "messageId": "uuid"
    }
}
```

**`message:unstar`** — Удалить из избранного.
```json
{
    "type": "message:unstar",
    "payload": {
        "messageId": "uuid"
    }
}
```

**`message:forward`** — Переслать сообщение.
```json
{
    "type": "message:forward",
    "payload": {
        "messageId": "uuid",
        "toChatId": "uuid"
    }
}
```

#### 6.4.2 Управление чатами

**`chat:create`** — Создать новый чат.
```json
{
    "type": "chat:create",
    "payload": {
        "type": "group",
        "name": "Friends",
        "participantIds": ["uuid1", "uuid2"],
        "description": "Чат для друзей"
    }
}
```
| Поле | Обязательное | Описание |
|------|-------------|----------|
| type | ✅ | private, group, channel |
| participantIds | ✅ | Список ID участников |
| name | ❌ | Название (обяз. для group/channel) |
| description | ❌ | Описание |
| avatarUrl | ❌ | Ссылка на аватар |

**`chat:update`** — Обновить чат.
```json
{
    "type": "chat:update",
    "payload": {
        "chatId": "uuid",
        "name": "Новое название",
        "description": "Описание",
        "avatarUrl": "https://..."
    }
}
```

**`chat:add_participant`** — Добавить участника.
```json
{
    "type": "chat:add_participant",
    "payload": {
        "chatId": "uuid",
        "userId": "uuid"
    }
}
```

**`chat:remove_participant`** — Удалить участника.
```json
{
    "type": "chat:remove_participant",
    "payload": {
        "chatId": "uuid",
        "userId": "uuid"
    }
}
```

**`chat:leave`** — Покинуть чат.
```json
{
    "type": "chat:leave",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`chat:pin`** — Закрепить чат в списке.
```json
{
    "type": "chat:pin",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`chat:unpin`** — Открепить чат.
```json
{
    "type": "chat:unpin",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`chat:archive`** — Архивировать чат.
```json
{
    "type": "chat:archive",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`chat:unarchive`** — Разархивировать чат.
```json
{
    "type": "chat:unarchive",
    "payload": {
        "chatId": "uuid"
    }
}
```

#### 6.4.3 Статус пользователя

**`user:typing`** — Индикатор печатания.
```json
{
    "type": "user:typing",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`user:stop_typing`** — Остановить индикатор.
```json
{
    "type": "user:stop_typing",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`user:keyboard_opened`** — Клавиатура открыта.
```json
{
    "type": "user:keyboard_opened",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`user:keyboard_closed`** — Клавиатура закрыта.
```json
{
    "type": "user:keyboard_closed",
    "payload": {
        "chatId": "uuid"
    }
}
```

**`user:block`** — Заблокировать пользователя.
```json
{
    "type": "user:block",
    "payload": {
        "userId": "uuid"
    }
}
```

**`user:unblock`** — Разблокировать пользователя.
```json
{
    "type": "user:unblock",
    "payload": {
        "userId": "uuid"
    }
}
```

#### 6.4.4 WebRTC звонки (сигналинг)

**`call:offer`** — Отправить WebRTC offer.
```json
{
    "type": "call:offer",
    "payload": {
        "chatId": "uuid",
        "callId": "uuid",
        "sdp": "offer_sdp_string"
    }
}
```

**`call:answer`** — Ответить на звонок.
```json
{
    "type": "call:answer",
    "payload": {
        "chatId": "uuid",
        "callId": "uuid",
        "sdp": "answer_sdp_string"
    }
}
```

**`call:ice`** — Отправить ICE-кандидат.
```json
{
    "type": "call:ice",
    "payload": {
        "callId": "uuid",
        "candidate": "ice_candidate_string"
    }
}
```

**`call:reject`** — Отклонить звонок.
```json
{
    "type": "call:reject",
    "payload": {
        "callId": "uuid"
    }
}
```

### 6.5 Технические детали

- **Пинг/Pong**: сервер шлёт ping каждые 54с, клиент отвечает pong (обрабатывается библиотекой websocket)
- **Таймаут**: если pong не получен 60с — соединение закрывается
- **Макс. размер сообщения**: 64 КБ
- **Rate limit**: 100 запросов/мин на API (REST + WebSocket)
- **Реконнект**: при обрыве клиент должен переподключаться с тем же JWT токеном
- **Автостатус**: при подключении пользователь автоматически становится `online`, при отключении — `offline`
- **Формат**: все сообщения в JSON

### 6.6 Примеры

**Пример 1: отправка сообщения**
1. Клиент А отправляет через WebSocket:
   ```json
   { "type": "message:send", "payload": { "chatId": "123", "content": "Привет!", "type": "text" } }
   ```
2. Сервер сохраняет в БД и рассылает всем участникам чата:
   ```json
   { "type": "message:new", "payload": { "id": "msg-uuid", "chatId": "123", "senderId": "userA-uuid", "content": "Привет!", "type": "text", "createdAt": "2026-01-01T00:00:00Z" } }
   ```
3. Клиент Б получает `message:new` и отображает сообщение

**Пример 2: звонок**
1. Клиент А отправляет `call:offer` → сервер шлёт `call:offer` всем участникам чата
2. Клиент Б отвечает `call:answer` → сервер шлёт `call:accept` всем
3. Клиенты обмениваются `call:ice` (ICE candidates) через сервер
4. Любая сторона шлёт `call:end` → сервер шлёт `call:end` всем
