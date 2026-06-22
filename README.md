# 💬 Go Messenger Server

High-performance chat messenger server with REST API + WebSocket real-time + WebRTC calls.

**Base URL:** `http://localhost:8080`  
**WebSocket:** `ws://localhost:8080/ws?token={jwt}`  
**Swagger:** `http://localhost:8080/swagger/index.html`  
**Postman:** `http://localhost:8080/postman`

---

## 📋 Содержание

- [Quick Start](#quick-start)
- [Авторизация (Auth Flow)](#авторизация-auth-flow)
- [Формат ответов](#формат-ответов)
- [API Endpoints](#api-endpoints)
  - [🔐 Auth](#-auth)
  - [👤 Users](#-users)
  - [💬 Chats](#-chats)
  - [✉️ Messages](#️-messages)
  - [📱 Contacts](#-contacts)
  - [📁 Folders](#-folders)
  - [🔗 Invite Links](#-invite-links)
  - [📊 Polls](#-polls)
  - [🎨 Stickers](#-stickers)
  - [🤖 Bots](#-bots)
  - [📞 Calls](#-calls)
  - [🎤 Voice Chats](#-voice-chats)
  - [📸 Stories](#-stories)
  - [📺 Channels](#-channels)
  - [📝 Drafts](#-drafts)
  - [⭐ Saved Messages](#-saved-messages)
  - [🎭 Custom Emojis](#-custom-emojis)
  - [🎨 GIFs](#-gifs)
  - [🖥️ Sessions](#️-sessions)
  - [🛡️ Verification](#️-verification)
  - [⚙️ Settings](#️-settings)
- [WebSocket](#websocket)
- [Postman Collection](#postman-collection)

---

## Quick Start

```bash
go build -o backend/ChatServer.exe ./backend/main.go
./backend/ChatServer.exe

# Server starts on http://localhost:8080
```

---

## Авторизация (Auth Flow)

### 1. Регистрация
```
POST /api/auth/register
Content-Type: application/json

{
  "username": "john",          // min:3, max:32
  "email": "john@mail.com",    // valid email
  "password": "secret123",     // min:6
  "display_name": "John Doe"   // min:1, max:64
}
```
**Ответ:** `201 Created`
```json
{
  "success": true,
  "data": {
    "token": "eyJhbGciOiJIUzI1NiIs...",
    "user": {
      "id": "uuid",
      "username": "john",
      "displayName": "John Doe",
      "email": "john@mail.com",
      "avatarUrl": "",
      "bio": "",
      "phone": "",
      "status": "online",
      "lastSeen": "2026-01-01T00:00:00Z"
    }
  }
}
```

### 2. Логин
```
POST /api/auth/login
Content-Type: application/json

{
  "email": "john@mail.com",
  "password": "secret123"
}
```
**Ответ:** `200 OK` — тот же формат, что и регистрация. Сохраните `token` — он понадобится для всех запросов.

### 3. Использование токена
Все защищённые эндпоинты требуют заголовок:
```
Authorization: Bearer eyJhbGciOiJIUzI1NiIs...
```
Токен живёт 24 часа (настраивается через `JWT_TTL`). По истечении вызовите:
```
GET /api/auth/refresh
Authorization: Bearer <token>
```

---

## Формат ответов

### ✅ Успех
```json
{
  "success": true,
  "data": { ... }
}
```

### 📄 С пагинацией
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
Параметры пагинации: `?offset=0&limit=50` (по умолчанию limit=50).

### ❌ Ошибка
```json
{
  "success": false,
  "error": "Описание ошибки",
  "code": "ERROR_CODE"
}
```

| HTTP Status | Code | Когда возникает |
|:-----------:|------|-----------------|
| 400 | `BAD_REQUEST` | Неверные данные запроса (валидация, JSON) |
| 401 | `UNAUTHORIZED` | Отсутствует/недействителен JWT токен |
| 403 | `FORBIDDEN` | Нет прав на действие (не участник чата, не владелец) |
| 404 | `NOT_FOUND` | Ресурс не найден (сообщение, чат, пользователь) |
| 400 | `VALIDATION_ERROR` | Поле не прошло валидацию |
| 401 | `UNAUTHORIZED` | Нет токена или неверный формат |
| 401 | `TOKEN_EXPIRED` | Токен истёк |
| 403 | `FORBIDDEN` | Нет прав (не админ, не владелец) |
| 404 | `NOT_FOUND` | Ресурс не найден |
| 409 | `DUPLICATE` | Уже существует (username/email заняты) |
| 429 | `RATE_LIMIT` | Слишком много запросов |
| 500 | `INTERNAL_ERROR` | Ошибка сервера |

### 📤 Загрузка файлов
Для загрузки файлов (аватар, фото чата, файлы в сообщениях) используйте `multipart/form-data`:
```
POST /api/chats/{id}/messages/file
Authorization: Bearer <token>
Content-Type: multipart/form-data

file: (binary)
caption: "Описание к файлу"
replyToId: "uuid" (опционально)
effect: "confetti" (опционально)
```

---

## API Endpoints

---

### 🔐 Auth

| Метод | Endpoint | Auth | Описание |
|-------|----------|:----:|----------|
| `POST` | `/api/auth/register` | 🔓 | Регистрация |
| `POST` | `/api/auth/admin/register` | 🔓 | Регистрация админа (+ `admin_secret`) |
| `POST` | `/api/auth/login` | 🔓 | Вход по email+password |
| `POST` | `/api/auth/login/email` | 🔓 | Отправить код на email |
| `POST` | `/api/auth/login/email/verify` | 🔓 | Подтвердить email-код |
| `POST` | `/api/auth/login/phone` | 🔓 | Отправить SMS-код |
| `POST` | `/api/auth/login/phone/verify` | 🔓 | Подтвердить SMS-код |
| `GET` | `/api/auth/refresh` | ✅ | Обновить токен |
| `PUT` | `/api/auth/change-password` | ✅ | Сменить пароль |

**Register** → см. [Авторизация](#авторизация-auth-flow)

**Login by email code:**
```
POST /api/auth/login/email
{ "email": "john@mail.com" }

→ Код отправлен на email

POST /api/auth/login/email/verify
{ "email": "john@mail.com", "code": "123456" }

→ { "success": true, "data": { "token": "...", "user": {...} } }
```

**Admin registration:**
```
POST /api/auth/admin/register
{
  "username": "admin",
  "email": "admin@mail.com",
  "password": "admin123",
  "display_name": "Admin",
  "admin_secret": "admin-secret-change-me"
}
```

**Change password:**
```
PUT /api/auth/change-password
Authorization: Bearer <token>
{ "oldPassword": "old123", "newPassword": "new456" }
```

---

### 👤 Users

#### GET /api/users/profile
Получить профиль текущего пользователя.
```json
// Response 200
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "john",
    "displayName": "John Doe",
    "email": "john@mail.com",
    "phone": "",
    "avatarUrl": "/uploads/avatar.jpg",
    "bio": "Hello!",
    "status": "online",
    "lastSeen": "2026-01-01T00:00:00Z",
    "isBot": false,
    "isAdmin": false
  }
}
```

#### PUT /api/users/profile
```json
{
  "displayName": "New Name",
  "bio": "New bio text",
  "username": "newusername"
}
```

#### GET /api/users/search?q=john&offset=0&limit=50
Поиск пользователей. `q` — обязательный параметр.

#### POST /api/users/avatar
`multipart/form-data` с полем `file`.

#### PUT /api/users/status
```json
{ "status": "online" }     // online | offline | away
```

#### PUT /api/users/push-token
```json
{ "token": "fcm_token_here", "provider": "fcm" }
```

#### POST /api/users/block
```json
{ "userId": "uuid" }
```

#### DELETE /api/users/block/{userId}
Разблокировать пользователя.

#### GET /api/users/blocked
Список заблокированных.

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/users/profile` | Профиль |
| `PUT` | `/api/users/profile` | Обновить профиль |
| `GET` | `/api/users/search?q=&offset=&limit=` | Поиск |
| `GET` | `/api/users/{id}` | По ID |
| `GET` | `/api/users/username/{username}` | По username |
| `GET` | `/api/users/{id}/last-seen` | Был(а) в сети |
| `PUT` | `/api/users/status` | Статус |
| `PUT` | `/api/users/push-token` | FCM токен |
| `POST` | `/api/users/avatar` | Аватар |
| `POST` | `/api/users/block` | Заблокировать |
| `DELETE` | `/api/users/block/{userId}` | Разблокировать |
| `GET` | `/api/users/blocked` | Список блокировок |

---

### 💬 Chats

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/chats` | Создать чат |
| `POST` | `/api/chats/start/{userId}` | Начать private chat |
| `GET` | `/api/chats` | Список чатов |
| `GET` | `/api/chats/{id}` | Инфо о чате |
| `PUT` | `/api/chats/{id}` | Обновить группу |
| `DELETE` | `/api/chats/{id}` | Удалить чат |
| `POST` | `/api/chats/{id}/read` | Прочитать |
| `POST` | `/api/chats/{id}/pin` | Закрепить |
| `DELETE` | `/api/chats/{id}/pin` | Открепить |
| `POST` | `/api/chats/{id}/archive` | Архивировать |
| `POST` | `/api/chats/{id}/unarchive` | Разархивировать |
| `POST` | `/api/chats/{id}/leave` | Выйти из группы |
| `POST` | `/api/chats/{id}/photo` | Фото чата |
| `POST` | `/api/chats/{id}/participants` | Добавить участника |
| `DELETE` | `/api/chats/{id}/participants/{userId}` | Удалить участника |
| `PUT` | `/api/chats/{id}/participants/{userId}/role` | Сменить роль |
| `GET` | `/api/chats/{id}/online` | Онлайн участники |
| `PUT` | `/api/chats/{id}/slow-mode` | Slow mode |
| `GET` | `/api/chats/archived` | Архив чатов |

#### POST /api/chats — Создать чат
```json
{
  "type": "private",         // "private" | "group" | "channel"
  "name": "My Group",        // обязательно для group/channel
  "participantIds": ["uuid1", "uuid2"],  // минимум 1 участник
  "description": "Описание"  // опционально
}
```
```json
// Response 201
{
  "success": true,
  "data": {
    "id": "uuid",
    "type": "group",
    "name": "My Group",
    "description": "Описание",
    "avatarUrl": "",
    "participants": [
      { "id": "uuid", "username": "john", "displayName": "John", "role": "owner" },
      { "id": "uuid2", "username": "jane", "displayName": "Jane", "role": "member" }
    ],
    "lastMessage": null,
    "unreadCount": 0,
    "pinned": false,
    "archived": false,
    "createdAt": "2026-01-01T00:00:00Z"
  }
}
```

#### POST /api/chats/start/{userId}
Создать или найти существующий private chat с пользователем.

#### GET /api/chats — Список чатов
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "type": "private",
      "name": "",
      "participants": [...],
      "lastMessage": {
        "id": "uuid",
        "content": "Hello!",
        "senderId": "uuid",
        "type": "text",
        "createdAt": "2026-01-01T00:00:00Z"
      },
      "unreadCount": 2,
      "pinned": false,
      "archived": false
    }
  ]
}
```

**Добавить участника:**
```
POST /api/chats/{id}/participants
{ "userId": "uuid" }
```

**Сменить роль:**
```
PUT /api/chats/{id}/participants/{userId}/role
{ "role": "admin" }
// role: "owner" | "admin" | "moderator" | "editor" | "member" | "read_only"
```

**Slow mode:**
```
PUT /api/chats/{id}/slow-mode
{ "intervalSeconds": 10 }  // 0 = выключено, 1-3600
```

---

### ✉️ Messages

Заголовки: `Authorization: Bearer <token>` + `Content-Type: application/json`

#### POST /api/chats/{chatId}/messages — Отправить сообщение
```json
{
  "content": "Привет!",       // текст сообщения (обязательно)
  "type": "text",            // text | image | file | gif | voice | video | audio | location | system
  "replyToId": "uuid",       // ID сообщения, на которое отвечаем (опц.)
  "latitude": 55.75,         // для type=location
  "longitude": 37.62,
  "locationTitle": "Москва",
  "effect": "confetti"       // confetti | fireworks | hearts | balloons | stars
}
```
```json
// Response 201
{
  "success": true,
  "data": {
    "id": "uuid",
    "chatId": "uuid",
    "senderId": "uuid",
    "content": "Привет!",
    "type": "text",
    "replyToId": null,
    "forwardFrom": null,
    "reactions": [],
    "readBy": [],
    "pinned": false,
    "createdAt": "2026-01-01T00:00:00Z",
    "updatedAt": "2026-01-01T00:00:00Z"
  }
}
```

#### GET /api/chats/{chatId}/messages?offset=0&limit=50
Получить сообщения чата (пагинация от новых к старым).

#### PUT /api/messages/{messageId} — Редактировать
```json
{ "content": "Новый текст" }
```

#### DELETE /api/messages/{messageId}
Удалить сообщение для всех.

#### POST /api/messages/{messageId}/reactions — Реакция
```json
{ "emoji": "👍" }
```
Популярные emoji: `👍` `❤️` `😆` `😮` `😢` `🙏`

#### DELETE /api/messages/{messageId}/reactions?emoji=👍
Удалить реакцию.

#### PUT /api/messages/{messageId}/pin
```json
{ "pin": true }
// true = закрепить, false = открепить
```

#### POST /api/messages/{messageId}/read
Отметить сообщение как прочитанное.

#### POST /api/messages/{messageId}/star / DELETE /api/messages/{messageId}/star
Добавить/удалить из избранного.

#### POST /api/messages/{messageId}/self-destruct
Установить таймер самоуничтожения сообщения. Сервер автоматически удалит сообщение через N секунд и разошлёт `message:deleted` через WebSocket всем участникам чата.

**Request (client → server):**
```json
{ "seconds": 60 }
// seconds: 1-86400 (1 секунда — 24 часа)
```

**Response (server → client):**
```json
{
  "success": true,
  "data": {
    "message": "self-destruct timer set"
  }
}
```

**WebSocket Event (server → client, через 60 сек):**
```json
{
  "type": "message:deleted",
  "payload": {
    "messageId": "uuid",
    "chatId": "uuid"
  }
}
```

**Механизм работы:** Фоновый процесс (горутина) каждые 15 секунд проверяет таблицу `message_self_destruct` и удаляет сообщения, у которых `delete_at <= now()`. После удаления рассылает `message:deleted` через WebSocket.

#### POST /api/messages/forward
```json
{
  "messageId": "uuid",
  "fromChatId": "uuid",
  "toChatId": "uuid"
}
```

#### GET /api/messages/search?q=text&offset=0&limit=50
Глобальный поиск по всем чатам.

#### POST /api/chats/{chatId}/messages/file
`multipart/form-data` — поле `file` (файл), `caption` (опц.), `replyToId` (опц.), `effect` (опц.)

#### POST /api/chats/{chatId}/messages/location
```json
{
  "latitude": 55.75,
  "longitude": 37.62,
  "title": "Москва",
  "replyToId": "uuid"
}
```

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/chats/{id}/messages` | Отправить сообщение |
| `GET` | `/api/chats/{id}/messages` | Список сообщений |
| `POST` | `/api/chats/{id}/messages/file` | Загрузить файл |
| `POST` | `/api/chats/{id}/messages/voice` | Голосовое сообщение |
| `POST` | `/api/chats/{id}/messages/location` | Геолокация |
| `POST` | `/api/chats/{id}/messages/video-circle` | Видео-кружок |
| `PUT` | `/api/messages/{id}` | Редактировать |
| `DELETE` | `/api/messages/{id}` | Удалить |
| `DELETE` | `/api/messages/{id}/for-me` | Удалить у себя |
| `POST` | `/api/messages/{id}/reactions` | Реакция |
| `DELETE` | `/api/messages/{id}/reactions` | Убрать реакцию |
| `PUT` | `/api/messages/{id}/pin` | Закрепить |
| `POST` | `/api/messages/{id}/star` | В избранное |
| `POST` | `/api/messages/{id}/read` | Прочитано |
| `POST` | `/api/messages/{id}/self-destruct` | Самоуничтожение |
| `POST` | `/api/messages/{id}/save` | Сохранить |
| `POST` | `/api/messages/{id}/report` | Пожаловаться |
| `GET` | `/api/messages/{id}/history` | История правок |
| `POST` | `/api/messages/forward` | Переслать |
| `POST` | `/api/messages/schedule` | Запланировать |
| `GET` | `/api/messages/search` | Поиск |
| `GET` | `/api/messages/starred` | Избранное |
| `GET` | `/api/messages/scheduled` | Запланированные |
| `GET` | `/api/chats/{id}/pinned` | Закреплённые |
| `GET` | `/api/chats/{id}/media` | Медиа-галерея |
| `GET` | `/api/chats/{id}/export` | Экспорт чата |

---

### 📱 Contacts

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/contacts/sync` | Синхронизировать контакты телефона |
| `GET` | `/api/contacts` | Список контактов |
| `GET` | `/api/contacts/search?q=` | Поиск по номеру |
| `GET` | `/api/contacts/registered` | Кто из контактов на платформе |
| `POST` | `/api/contacts/photo` | Фото контакта |

#### POST /api/contacts/sync
```json
{
  "contacts": [
    { "phone": "+79001234567", "name": "John", "lastName": "Doe" },
    { "phone": "+79007654321", "name": "Jane" }
  ]
}
```
```json
// Response 200 — список зарегистрированных контактов
{
  "success": true,
  "data": [
    {
      "userId": "uuid",
      "phone": "+79001234567",
      "name": "John",
      "displayName": "John Doe",
      "avatarUrl": ""
    }
  ]
}
```

---

### 📁 Folders

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/folders` | Список папок |
| `POST` | `/api/folders` | Создать папку |
| `PUT` | `/api/folders/{id}` | Обновить папку |
| `DELETE` | `/api/folders/{id}` | Удалить папку |

#### POST /api/folders
```json
{
  "name": "Work",
  "chatIds": ["uuid1", "uuid2"],
  "emoji": "💼",
  "order": 1
}
```

---

### 🔗 Invite Links

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/chats/{id}/invite-links` | Создать ссылку |
| `GET` | `/api/chats/{id}/invite-links` | Список ссылок |
| `DELETE` | `/api/chats/{id}/invite-links/{linkId}` | Удалить ссылку |
| `POST` | `/api/chats/join` | Присоединиться по ссылке |

#### POST /api/chats/{id}/invite-links
```json
{
  "expiresInMins": 1440,    // опц., 0 = без срока
  "usageLimit": 100          // опц., 0 = без лимита
}
```

#### POST /api/chats/join
```json
{ "code": "invite_link_code" }
```

---

### 📊 Polls

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/chats/{id}/polls` | Создать опрос |
| `GET` | `/api/chats/{id}/polls` | Список опросов |
| `POST` | `/api/polls/{pollId}/vote` | Проголосовать |
| `POST` | `/api/polls/{pollId}/close` | Закрыть опрос |

#### POST /api/chats/{id}/polls
```json
{
  "question": "Best programming language?",
  "options": ["Go", "Rust", "Python"],
  "isAnonymous": false,
  "multipleChoice": false,
  "expiresInMins": 1440      // опц., 0 = без срока
}
```

#### POST /api/polls/{pollId}/vote
```json
{
  "optionIndex": 0           // индекс выбранного варианта
}
```
Для `multipleChoice: true` можно передать массив:
```json
{ "optionIndex": [0, 2] }
```

---

### 🎨 Stickers

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/stickers/packs` | Все паки |
| `GET` | `/api/stickers/packs/my` | Мои паки |
| `POST` | `/api/stickers/packs` | Создать пак |
| `GET` | `/api/stickers/packs/{id}` | Пак по ID |
| `DELETE` | `/api/stickers/packs/{id}` | Удалить пак |
| `POST` | `/api/stickers/packs/{id}/stickers` | Добавить стикер |
| `POST` | `/api/stickers/library` | Добавить в библиотеку |
| `GET` | `/api/stickers/library` | Моя библиотека |

#### POST /api/stickers/packs
```json
{
  "name": "My Stickers",
  "emoji": "😎",
  "stickers": [
    { "fileUrl": "https://...", "emoji": "🔥" }
  ]
}
```

#### Добавить стикер в пак:
```
POST /api/stickers/packs/{id}/stickers
Content-Type: multipart/form-data

file: (image file, PNG/WEBP)
emoji: "🔥"
```

---

### 🤖 Bots

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/bots` | Создать бота |
| `GET` | `/api/bots` | Мои боты |
| `PUT` | `/api/bots/{id}` | Обновить |
| `DELETE` | `/api/bots/{id}` | Удалить |
| `POST` | `/api/bots/{id}/token` | Новый токен |

#### POST /api/bots
```json
{
  "username": "my_bot",
  "name": "My Bot",
  "description": "Bot description"
}
```
```json
// Response 201
{
  "success": true,
  "data": {
    "id": "uuid",
    "username": "my_bot",
    "name": "My Bot",
    "token": "bot_token_here",
    "description": "Bot description"
  }
}
```

---

### 📞 Calls

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/calls/initiate` | Начать звонок |
| `POST` | `/api/calls/{id}/respond` | Ответить |
| `POST` | `/api/calls/{id}/end` | Завершить |
| `GET` | `/api/calls/{id}` | Инфо о звонке |
| `GET` | `/api/chats/{chatId}/calls` | История звонков |

#### POST /api/calls/initiate
```json
{
  "chatId": "uuid",
  "type": "audio"     // "audio" | "video"
}
```

#### POST /api/calls/{id}/respond
```json
{ "action": "accept" }    // "accept" | "reject"
```

---

### 🎤 Voice Chats

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/chats/{id}/voice-chat` | Создать |
| `GET` | `/api/chats/{id}/voice-chats/active` | Активные |
| `GET` | `/api/chats/{id}/voice-chats/history` | История |
| `GET` | `/api/voice-chats/{id}` | Детали |
| `POST` | `/api/voice-chats/{id}/join` | Присоединиться |
| `POST` | `/api/voice-chats/{id}/leave` | Покинуть |
| `POST` | `/api/voice-chats/{id}/end` | Завершить |
| `POST` | `/api/voice-chats/{id}/mute` | Заглушить |

#### POST /api/chats/{id}/voice-chat
```json
{
  "title": "Evening Chat",
  "isPublic": true
}
```

#### POST /api/voice-chats/{id}/mute
```json
{
  "userId": "uuid",
  "muted": true
}
```

---

### 📸 Stories

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/stories` | Создать историю |
| `GET` | `/api/stories` | Истории подписок |
| `GET` | `/api/stories/my` | Мои истории |
| `GET` | `/api/stories/{id}` | Просмотр |
| `DELETE` | `/api/stories/{id}` | Удалить |
| `GET` | `/api/stories/{id}/views` | Просмотревшие |

#### POST /api/stories
```
Content-Type: multipart/form-data
Authorization: Bearer <token>

file: (image or video file)
caption: "Story text"
type: "photo"     // "photo" | "video"
```

---

### 📺 Channels

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/channels/subscribe` | Подписаться |
| `POST` | `/api/channels/unsubscribe` | Отписаться |
| `GET` | `/api/channels` | Мои каналы |
| `GET` | `/api/channels/{id}/subscribers` | Подписчики |
| `GET` | `/api/channels/{id}/subscribed` | Проверить подписку |
| `PUT` | `/api/channels/{id}/subscribers/{userId}/role` | Роль подписчика |

#### POST /api/channels/subscribe
```json
{ "channelId": "uuid" }
```

---

### 📝 Drafts

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/drafts` | Сохранить черновик |
| `GET` | `/api/drafts?chatId=uuid` | Получить черновик |
| `DELETE` | `/api/drafts/{id}` | Удалить |

#### POST /api/drafts
```json
{
  "chatId": "uuid",
  "content": "Draft message text",
  "replyToId": "uuid"    // опционально
}
```

---

### ⭐ Saved Messages

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/messages/{id}/save?chatId=uuid` | Сохранить сообщение |
| `GET` | `/api/saved-messages` | Сохранённые |
| `DELETE` | `/api/saved-messages/{id}` | Удалить |

---

### 🎭 Custom Emojis

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/emojis` | Загрузить emoji |
| `GET` | `/api/emojis` | Все публичные |
| `GET` | `/api/emojis/my` | Мои emoji |
| `DELETE` | `/api/emojis/{id}` | Удалить |

#### POST /api/emojis
```
Content-Type: multipart/form-data
Authorization: Bearer <token>

file: (image file)
shortcode: "my_emoji"    // как будет вызываться :my_emoji:
```

---

### 🎨 GIFs

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/gifs` | Сохранить GIF |
| `GET` | `/api/gifs` | Мои GIF |
| `DELETE` | `/api/gifs?url=https://...` | Удалить |

#### POST /api/gifs
```json
{
  "url": "https://media.giphy.com/media/.../giphy.gif",
  "title": "Funny cat"
}
```

---

### 🖥️ Sessions

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/sessions` | Список сессий |
| `DELETE` | `/api/sessions/{id}` | Завершить сессию |
| `DELETE` | `/api/sessions` | Все сессии |

#### GET /api/sessions
```json
// Response 200
{
  "success": true,
  "data": [
    {
      "id": "uuid",
      "device": "Chrome on Windows",
      "ip": "192.168.1.1",
      "lastActive": "2026-01-01T00:00:00Z",
      "createdAt": "2026-01-01T00:00:00Z"
    }
  ]
}
```

---

### 🛡️ Verification

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `POST` | `/api/verification/email/send` | Отправить код на email |
| `POST` | `/api/verification/email/verify` | Подтвердить email |
| `POST` | `/api/verification/phone/send` | Отправить SMS |
| `POST` | `/api/verification/phone/verify` | Подтвердить телефон |

```json
// Отправка кода
POST /api/verification/email/send
{ "email": "user@mail.com" }

// Подтверждение
POST /api/verification/email/verify
{ "email": "user@mail.com", "code": "123456" }
// → { "success": true, "data": { "verified": true, "message": "Email verified" } }
```

---

### ⚙️ Settings

| Метод | Endpoint | Описание |
|-------|----------|----------|
| `GET` | `/api/account/settings` | Настройки |
| `PUT` | `/api/account/settings` | Обновить |

```json
GET /api/account/settings
// Response
{
  "success": true,
  "data": {
    "language": "ru",
    "theme": "dark",
    "notificationsEnabled": true,
    "lastSeenEnabled": true,
    "onlineStatusEnabled": true
  }
}

PUT /api/account/settings
{
  "language": "en",
  "theme": "light",
  "notificationsEnabled": true
}
```

---

## WebSocket

WebSocket используется для real-time событий — сообщения приходят мгновенно, без polling.

### Подключение
```javascript
const ws = new WebSocket("ws://localhost:8080/ws?token=" + jwtToken);
```

### Формат всех сообщений
```json
{ "type": "название_события", "payload": { ... поля ... } }
```

### События от сервера к клиенту

Сервер автоматически рассылает эти события всем участникам чата.

#### Новое сообщение
```json
// Сервер отправляет всем участникам чата
{
  "type": "message:new",
  "payload": {
    "id": "uuid",
    "chatId": "uuid",
    "senderId": "uuid",
    "content": "Привет!",
    "type": "text",
    "reactions": [],
    "readBy": [],
    "pinned": false,
    "createdAt": "2026-01-01T00:00:00Z"
  }
}
```

#### Сообщение отредактировано
```json
{ "type": "message:edited", "payload": { ... Message ... } }
```

#### Сообщение удалено (вручную или самоуничтожение)
```json
{ "type": "message:deleted", "payload": { "messageId": "uuid", "chatId": "uuid" } }
```
Самоуничтожение: фоновый процесс проверяет таймеры каждые 15 секунд. Как только `delete_at` истекает — сообщение удаляется, и участники получают `message:deleted` через WebSocket.

#### Сообщение прочитано
```json
{ "type": "message:read", "payload": { "messageId": "uuid", "userId": "uuid", "chatId": "uuid" } }
```

#### Реакция
```json
{ "type": "message:reaction", "payload": { ... Message с reactions ... } }
```

#### Закреплено
```json
{ "type": "message:pinned", "payload": { ... Message ... } }
```

#### Пользователь онлайн/офлайн
```json
{ "type": "user:online", "payload": { "userId": "uuid", "online": true } }
{ "type": "user:offline", "payload": { "userId": "uuid", "online": false } }
```

#### Печатает
```json
{ "type": "user:typing", "payload": { "chatId": "uuid", "userId": "uuid" } }
{ "type": "user:stop_typing", "payload": { "chatId": "uuid", "userId": "uuid" } }
```

#### Чат создан/обновлён/удалён
```json
{ "type": "chat:created", "payload": { ... Chat ... } }
{ "type": "chat:updated", "payload": { ... Chat ... } }
{ "type": "chat:deleted", "payload": {} }
```

#### Звонки (WebRTC)
```json
{ "type": "call:offer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "offer_string" } }
{ "type": "call:answer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "answer_string" } }
{ "type": "call:ice", "payload": { "callId": "uuid", "candidate": "ice_candidate" } }
{ "type": "call:end", "payload": { "callId": "uuid", "userId": "uuid" } }
{ "type": "call:reject", "payload": { "callId": "uuid", "userId": "uuid" } }
```

### События от клиента к серверу

Клиент отправляет эти сообщения через WebSocket, чтобы выполнить действие.

#### Отправить сообщение через WebSocket
```json
{
  "type": "message:send",
  "payload": {
    "chatId": "uuid",
    "content": "Привет!",
    "type": "text",
    "replyToId": "uuid"   // опционально
  }
}
```
Типы сообщений: `text` | `image` | `file` | `gif` | `voice` | `video` | `audio` | `location` | `system`

#### Редактировать / Удалить
```json
{ "type": "message:edit", "payload": { "messageId": "uuid", "content": "Новый текст" } }
{ "type": "message:delete", "payload": { "messageId": "uuid", "chatId": "uuid" } }
```

#### Прочитано
```json
{ "type": "message:read", "payload": { "messageId": "uuid", "chatId": "uuid" } }
```

#### Реакции
```json
{ "type": "message:react", "payload": { "messageId": "uuid", "emoji": "👍" } }
{ "type": "message:unreact", "payload": { "messageId": "uuid", "emoji": "👍" } }
```

#### Закрепить / Избранное / Переслать
```json
{ "type": "message:pin", "payload": { "messageId": "uuid", "pin": true } }
{ "type": "message:star", "payload": { "messageId": "uuid" } }
{ "type": "message:forward", "payload": { "messageId": "uuid", "toChatId": "uuid" } }
```

#### Создать чат
```json
{
  "type": "chat:create",
  "payload": {
    "type": "group",
    "name": "Friends",
    "participantIds": ["uuid1", "uuid2"],
    "description": "Chat for friends"
  }
}
```

#### Обновить чат
```json
{ "type": "chat:update", "payload": { "chatId": "uuid", "name": "New Name", "description": "...", "avatarUrl": "..." } }
```

#### Управление участниками
```json
{ "type": "chat:add_participant", "payload": { "chatId": "uuid", "userId": "uuid" } }
{ "type": "chat:remove_participant", "payload": { "chatId": "uuid", "userId": "uuid" } }
{ "type": "chat:leave", "payload": { "chatId": "uuid" } }
```

#### Чат: закрепить / архив
```json
{ "type": "chat:pin", "payload": { "chatId": "uuid" } }
{ "type": "chat:archive", "payload": { "chatId": "uuid" } }
```

#### Печатание
```json
{ "type": "user:typing", "payload": { "chatId": "uuid" } }
{ "type": "user:stop_typing", "payload": { "chatId": "uuid" } }
```

#### Блокировки
```json
{ "type": "user:block", "payload": { "userId": "uuid" } }
{ "type": "user:unblock", "payload": { "userId": "uuid" } }
```

#### WebRTC звонки через WebSocket
```json
{ "type": "call:offer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "..." } }
{ "type": "call:answer", "payload": { "chatId": "uuid", "callId": "uuid", "sdp": "..." } }
{ "type": "call:ice", "payload": { "callId": "uuid", "candidate": "..." } }
```

---

## Postman Collection

Готовая коллекция Postman со всеми 168 эндпоинтами:
```
http://localhost:8080/postman
```

Или в репозитории: `docs/postman_collection.json`

**Импорт:** Postman → File → Import → выбрать файл → установить переменные:
- `base_url`: `http://localhost:8080`
- `token`: (вставится автоматически после Register/Login)

---

## ⚙️ Конфигурация

| Переменная | По умолчанию | Описание |
|-----------|-------------|----------|
| `SERVER_PORT` | `8080` | Порт сервера |
| `DATABASE_PATH` | `file:chat.db?mode=memory...` | SQLite |
| `JWT_SECRET` | `super-secret-change-me` | Ключ для JWT |
| `ADMIN_SECRET` | `admin-secret-change-me` | Секрет админа |
| `JWT_TTL` | `86400` | Время жизни токена (сек) |
| `CORS_ORIGIN` | `*` | CORS |
| `FCM_KEY_PATH` | — | Firebase ключ |
| `ENCRYPTION_KEY` | — | Ключ шифрования |
