# Chat Messenger Server

**Сервер мессенджера на Go** с REST API, WebSocket real-time, WebRTC звонками, push-уведомлениями, загрузкой файлов, реакциями, закреплением сообщений, read-ресиптами, блокировками, rate limiting, refresh-токенами, настройками аккаунта, контактами телефона, аватарками и полным профилем пользователя.

---

## Быстрый старт

```bash
./ChatServer.exe

# Swagger UI: http://localhost:8080/swagger/index.html

# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@mail.com","password":"secret123","displayName":"John"}'
```

---

## API Endpoints (54 endpoints)

### Auth

| Method | Endpoint | Auth | Описание |
|--------|----------|------|----------|
| `POST` | `/api/auth/register` | ❌ | Регистрация |
| `POST` | `/api/auth/login` | ❌ | Вход |
| `GET` | `/api/auth/refresh` | ✅ | Обновить JWT |
| `PUT` | `/api/auth/change-password` | ✅ | Сменить пароль |

### Users

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/users/profile` | Профиль (phone, gender, dateOfBirth) |
| `PUT` | `/api/users/profile` | Обновить профиль |
| `POST` | `/api/users/avatar` | Загрузить аватарку (multipart) |
| `PUT` | `/api/users/status` | Статус (Available, Busy...) |
| `PUT` | `/api/users/push-token` | Push-токен |
| `POST` | `/api/users/push-test` | Тестовый push |
| `GET` | `/api/users/search?q=` | Поиск |
| `GET` | `/api/users/{id}` | По ID |
| `GET` | `/api/users/username/{username}` | По username |
| `POST` | `/api/users/block` | Заблокировать |
| `DELETE` | `/api/users/block/{userId}` | Разблокировать |
| `GET` | `/api/users/blocked` | Список заблокированных |
| `DELETE` | `/api/users/account` | Удалить аккаунт |

### Account Settings

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/account/settings` | Получить настройки |
| `PUT` | `/api/account/settings` | Обновить настройки |

Поля: `language` (en/ru), `theme` (light/dark), `notifications`, `soundEnabled`, `lastSeenMode` (everyone/nobody/contacts).

### Contacts

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/contacts/sync` | Синхронизировать телефонную книгу |
| `GET` | `/api/contacts` | Получить все контакты |
| `GET` | `/api/contacts/search?q=` | Поиск по номеру |
| `GET` | `/api/contacts/registered` | Найти пользователей сервера среди контактов |

При синхронизации контакты сохраняются в БД. Endpoint `/contacts/registered` ищет пользователей, чей номер телефона совпадает с сохранёнными контактами.

### Chats

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/chats` | Список чатов (скрытые не показываются) |
| `POST` | `/api/chats` | Создать (private/group) |
| `GET` | `/api/chats/{id}` | Детали |
| `PUT` | `/api/chats/{id}` | Обновить группу |
| `DELETE` | `/api/chats/{id}` | Удалить (только создатель) |
| `POST` | `/api/chats/{id}/participants` | Добавить участника |
| `DELETE` | `/api/chats/{id}/participants/{userId}` | Удалить участника |
| `PUT` | `/api/chats/{id}/participants/{userId}/role` | Роль (admin/member) |
| `POST` | `/api/chats/{id}/leave` | Покинуть группу |
| `POST` | `/api/chats/{id}/read` | Отметить прочитанным |
| `POST` | `/api/chats/{id}/hide` | Скрыть чат для себя |
| `PUT` | `/api/chats/{id}/notifications` | Mute/unmute |
| `GET` | `/api/chats/{id}/notifications` | Статус уведомлений |

**Hide chat**: скрывает чат только для текущего пользователя. Другие участники всё ещё видят чат. Чат не удаляется — его можно вернуть, создав новое сообщение.

### Messages

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/chats/{id}/messages` | Сообщения чата |
| `POST` | `/api/chats/{id}/messages` | Отправить |
| `GET` | `/api/chats/{id}/messages/search?q=` | Поиск |
| `POST` | `/api/chats/{id}/messages/file` | Загрузить файл (multipart) |
| `POST` | `/api/chats/{id}/messages/{msgId}/resend` | Переслать |
| `GET` | `/api/chats/{id}/pinned` | Закреплённые |
| `GET` | `/api/messages/{id}` | По ID |
| `PUT` | `/api/messages/{id}` | Редактировать |
| `DELETE` | `/api/messages/{id}` | Удалить |
| `POST` | `/api/messages/{id}/reactions` | Реакция (emoji) |
| `DELETE` | `/api/messages/{id}/reactions?emoji=` | Удалить реакцию |
| `PUT` | `/api/messages/{id}/pin` | Закрепить/открепить |
| `POST` | `/api/messages/{id}/read` | Прочитано |

**Типы сообщений:** `text`, `image`, `file`, `gif`, `voice`, `video`, `system`.

### Files

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/files/{filename}` | Скачать файл |

### Calls

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/calls/initiate` | Начать звонок |
| `POST` | `/api/calls/{id}/respond` | Ответить (accept/reject) |
| `POST` | `/api/calls/{id}/end` | Завершить |
| `GET` | `/api/calls/{id}` | Информация |
| `GET` | `/api/calls/history/{chatId}` | История |

### WebSocket

| Endpoint | Описание |
|----------|----------|
| `ws://localhost:8080/ws?token=JWT` | Real-time |

### System

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/health` | Healthcheck |
| `GET` | `/swagger/*any` | Swagger UI |

---

## WebSocket Events

Подключение: `ws://localhost:8080/ws?token=JWT`

Сервер отмечает пользователя **online** при подключении и **offline** при отключении.

### Сервер → Клиент (сервер отправляет)

| Событие | Когда | Payload |
|---------|-------|---------|
| `message:new` | Новое сообщение | `MessageResponse` |
| `message:edited` | Сообщение изменено | `MessageResponse` |
| `message:deleted` | Сообщение удалено | `{ messageId, chatId }` |
| `message:read` | Прочитано | `{ chatId, userId }` |
| `user:typing` | Пользователь печатает | `{ chatId, userId }` |
| `user:stop_typing` | Перестал печатать | `{ chatId, userId }` |
| `user:online` | Появился в сети | `{ userId, online: true }` |
| `user:offline` | Ушёл | `{ userId, online: false }` |
| `chat:created` | Создан чат | `ChatResponse` |
| `chat:updated` | Чат обновлён | `ChatResponse` |
| `chat:deleted` | Чат удалён | `{ chatId }` |
| `call:offer` | Входящий звонок | `{ chatId, callerId }` |
| `call:accept` | Звонок принят | `{ callId, userId }` |
| `call:end` | Звонок завершён | `{ callId, userId }` |

### Клиент → Сервер (клиент отправляет)

| Событие | Payload | Описание |
|---------|---------|----------|
| `user:typing` | `{ chatId }` | Печатает (рассылается всем, кроме отправителя) |
| `user:stop_typing` | `{ chatId }` | Не печатает |
| `call:offer` | `{ chatId, callId, sdp }` | WebRTC offer (ретранслируется) |
| `call:answer` | `{ callId, sdp }` | WebRTC answer (ретранслируется) |
| `call:ice` | `{ callId, candidate, sdpMLineIdx }` | ICE candidate (ретранслируется) |
| `call:reject` | `{ callId, chatId }` | Отклонить звонок |

**Важно:** typing-события приходят **всем участникам чата, кроме отправителя**. WebRTC события ретранслируются через сервер к другому участнику (signaling relay).

---

## Настройки профиля

```
PUT /api/users/profile
```

```json
{
  "displayName": "New Name",
  "bio": "About me",
  "phone": "+79001234567",
  "gender": "female",
  "dateOfBirth": "1995-06-15",
  "avatarUrl": "/uploads/avatars/userid.jpg"
}
```

- `gender`: `male`, `female`, `other`
- `dateOfBirth`: формат `YYYY-MM-DD`
- `phone`: любой формат номера

---

## Контакты телефона

Пользователь может загрузить свою телефонную книгу на сервер:

```bash
# Sync contacts
curl -X POST http://localhost:8080/api/contacts/sync \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"contacts":[{"phone":"+79001234567","name":"Alice"},{"phone":"+79876543210","name":"Bob"}]}'

# Find registered users
curl http://localhost:8080/api/contacts/registered \
  -H "Authorization: Bearer $TOKEN"
```

Endpoint `/contacts/registered` возвращает пользователей сервера, чьи номера телефонов есть в контактах (поле `phone` в профиле).

---

## Типы сообщений

- `text` — текстовое сообщение
- `image` — изображение
- `file` — файл
- `gif` — GIF-анимация
- `voice` — голосовое сообщение
- `video` — видеосообщение/кружок
- `system` — системное (серверное) сообщение

---

## Переменные окружения

| Переменная | По умолч. | Описание |
|------------|-----------|----------|
| `SERVER_PORT` | `8080` | Порт |
| `DATABASE_PATH` | `file:chat.db?cache=shared&mode=rwc` | SQLite |
| `JWT_SECRET` | `super-secret-key-...` | Секрет JWT |
| `JWT_TTL` | `86400` | TTL токена (сек) |
| `ALLOW_ORIGINS` | `*` | CORS origin |
| `PUSH_ENABLED` | `false` | Push-уведомления |
| `FIREBASE_CREDENTIALS` | — | Server Key FCM |

---

## Push-уведомления (FCM)

```bash
set PUSH_ENABLED=true
set FIREBASE_CREDENTIALS=AIzaSy...

curl -X PUT http://localhost:8080/api/users/push-token \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"token":"<FCM_TOKEN>","provider":"fcm"}'

curl -X POST http://localhost:8080/api/users/push-test \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Hello!","body":"Test"}'
```

---

## Optimizations & Scaling

### Database
- **WAL mode** — параллельные чтения не блокируют запись
- **Connection pool** — 25 max open, 5 idle
- **Busy timeout** — 5 секунд ожидания блокировки
- **Индексы** — на все ключевые поля

### Stability
- **Graceful shutdown** — `server.Shutdown(ctx)` с таймаутом 10s
- **Rate limiter** — 100 requests/min per IP с `stopCh`
- **Auto-miss call** — горутина `time.Sleep(30s)` для автоматического пропуска
- **WebSocket ping/pong** — 60s timeout, автоматическое переподключение

### WebSocket
- **Typing events** — рассылаются всем, кроме отправителя
- **Online/offline** — при подключении/отключении WS
- **WebRTC relay** — сервер ретранслирует SDP/ICE между участниками

### Безопасность
- **Проверка блокировки** при отправке сообщений (в обе стороны)
- **Доступ к чату** только для участников
- **Редактирование/удаление** только своих сообщений (или админом)

---

## Примеры curl

```bash
# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@mail.com","password":"123456","displayName":"Alice"}'

# Update profile with phone, gender, date of birth
curl -X PUT http://localhost:8080/api/users/profile \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"phone":"+79001234567","gender":"female","dateOfBirth":"1995-06-15"}'

# Upload avatar
curl -X POST http://localhost:8080/api/users/avatar \
  -H "Authorization: Bearer $TOKEN" \
  -F "avatar=@photo.jpg"

# Sync contacts
curl -X POST http://localhost:8080/api/contacts/sync \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"contacts":[{"phone":"+79001234567","name":"Friend"}]}'

# Find registered contacts
curl http://localhost:8080/api/contacts/registered \
  -H "Authorization: Bearer $TOKEN"

# Send GIF
curl -X POST http://localhost:8080/api/chats/CHAT_ID/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"funny.gif","type":"gif"}'

# Send voice
curl -X POST http://localhost:8080/api/chats/CHAT_ID/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"voice.ogg","type":"voice"}'

# Hide chat
curl -X POST http://localhost:8080/api/chats/CHAT_ID/hide \
  -H "Authorization: Bearer $TOKEN"
```

---

## Сборка

```bash
# Windows
go build -o ChatServer.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 .
```

**Go 1.21+**

---

## Лицензия

MIT
