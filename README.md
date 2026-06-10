# Chat Messenger Server

**Сервер мессенджера на Go** с REST API, WebSocket real-time, WebRTC звонками, push-уведомлениями, опросами, стикерами, GIF, черновиками, отложенными сообщениями, сессиями, ботами, каналами, E2E шифрованием, admin панелью, капчей, верификацией email/SMS, закладками, жалобами, self-destruct сообщениями, историей редактирования, превью ссылок и IP-блокировкой.

---

## Быстрый старт

```bash
./ChatServer.exe

# Swagger UI: http://localhost:8080/swagger/index.html
# API Tester: http://localhost:8080/app/
# Health:     http://localhost:8080/health

# Register
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@mail.com","password":"secret123","displayName":"John"}'
```

---

## API Endpoints (100+ endpoints)

### Public endpoints

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/auth/register` | Регистрация (с поддержкой captcha) |
| `POST` | `/api/auth/login` | Вход по email + пароль |
| `POST` | `/api/auth/login/email` | Отправить код входа на email |
| `POST` | `/api/auth/login/email/verify` | Подтвердить вход по email коду |
| `POST` | `/api/auth/login/phone` | Отправить SMS код входа |
| `POST` | `/api/auth/login/phone/verify` | Подтвердить вход по SMS коду |
| `GET` | `/api/captcha/generate` | Сгенерировать captcha |
| `POST` | `/api/captcha/verify` | Проверить captcha |
| `GET` | `/api/preview?url=` | Получить link preview |

### Auth (authenticated)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/auth/refresh` | Обновить JWT |
| `PUT` | `/api/auth/change-password` | Сменить пароль |

### Users

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/users/profile` | Профиль |
| `PUT` | `/api/users/profile` | Обновить профиль |
| `POST` | `/api/users/avatar` | Загрузить аватарку |
| `PUT` | `/api/users/status` | Статус (Available/Busy/...) |
| `PUT` | `/api/users/push-token` | Push-токен |
| `POST` | `/api/users/push-test` | Тестовый push |
| `GET` | `/api/users/search?q=` | Поиск пользователей |
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
| `PUT` | `/api/account/settings` | Обновить (language, theme, notifications, sound, lastSeenMode) |

### Contacts

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/contacts/sync` | Синхронизировать телефонную книгу |
| `GET` | `/api/contacts` | Все контакты |
| `GET` | `/api/contacts/search?q=` | Поиск по номеру |
| `GET` | `/api/contacts/registered` | Найти пользователей сервера среди контактов |

### Chats

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/chats` | Список чатов |
| `GET` | `/api/chats/search?q=` | Поиск чатов по имени |
| `GET` | `/api/chats/archived` | Архивные чаты |
| `POST` | `/api/chats` | Создать (private/group/channel) |
| `GET` | `/api/chats/{id}` | Детали чата |
| `PUT` | `/api/chats/{id}` | Обновить группу |
| `DELETE` | `/api/chats/{id}` | Удалить чат |
| `POST` | `/api/chats/{id}/participants` | Добавить участника |
| `DELETE` | `/api/chats/{id}/participants/{userId}` | Удалить участника |
| `PUT` | `/api/chats/{id}/participants/{userId}/role` | Роль (admin/member) |
| `POST` | `/api/chats/{id}/leave` | Покинуть группу |
| `POST` | `/api/chats/{id}/read` | Отметить прочитанным |
| `POST` | `/api/chats/{id}/pin` | Закрепить чат |
| `DELETE` | `/api/chats/{id}/pin` | Открепить чат |
| `POST` | `/api/chats/{id}/archive` | Архивировать чат |
| `POST` | `/api/chats/{id}/unarchive` | Разархивировать чат |
| `POST` | `/api/chats/{id}/hide` | Скрыть чат |
| `POST` | `/api/chats/{id}/transfer-ownership` | Передать владельца |
| `PUT` | `/api/chats/{id}/notifications` | Mute/unmute |
| `GET` | `/api/chats/{id}/notifications` | Статус уведомлений |

### Messages

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/chats/{id}/messages?limit=&offset=` | Сообщения чата |
| `GET` | `/api/chats/{id}/messages/search?q=` | Поиск по чату |
| `GET` | `/api/chats/{id}/media?type=image|file|video|audio|gif` | Медиа чата |
| `POST` | `/api/chats/{id}/messages` | Отправить сообщение |
| `POST` | `/api/chats/{id}/messages/file` | Загрузить файл |
| `POST` | `/api/chats/{id}/messages/{msgId}/resend` | Переслать |
| `GET` | `/api/chats/{id}/pinned` | Закреплённые сообщения |
| `GET` | `/api/chats/{id}/export` | Экспорт чата (JSON) |
| `GET` | `/api/messages/search?q=` | Поиск по всем чатам |
| `GET` | `/api/messages/starred` | Избранные сообщения |
| `POST` | `/api/messages/forward` | Переслать в другой чат |
| `POST` | `/api/messages/schedule` | Отложить сообщение |
| `GET` | `/api/messages/scheduled` | Список отложенных |
| `DELETE` | `/api/messages/scheduled/{id}` | Отменить отложенное |
| `GET` | `/api/messages/{id}` | Сообщение по ID |
| `PUT` | `/api/messages/{id}` | Редактировать |
| `DELETE` | `/api/messages/{id}` | Удалить |
| `DELETE` | `/api/messages/{id}/for-me` | Удалить у себя |
| `POST` | `/api/messages/{id}/reactions` | Поставить реакцию |
| `DELETE` | `/api/messages/{id}/reactions?emoji=` | Убрать реакцию |
| `PUT` | `/api/messages/{id}/pin` | Закрепить/открепить |
| `POST` | `/api/messages/{id}/star` | Добавить в избранное |
| `DELETE` | `/api/messages/{id}/star` | Убрать из избранного |
| `POST` | `/api/messages/{id}/read` | Отметить прочитанным |

**Типы сообщений:** `text`, `image`, `file`, `gif`, `voice`, `video`, `audio`, `system`.

### Polls (опросы)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/chats/{id}/polls` | Создать опрос |
| `GET` | `/api/chats/{id}/polls` | Опросы чата |
| `POST` | `/api/polls/{pollId}/vote` | Проголосовать |
| `POST` | `/api/polls/{pollId}/close` | Закрыть опрос |

### Stickers

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/stickers/packs` | Все паки стикеров |
| `GET` | `/api/stickers/packs/my` | Мои паки |
| `POST` | `/api/stickers/packs` | Создать пак |
| `GET` | `/api/stickers/packs/{id}` | Пак со стикерами |
| `POST` | `/api/stickers/packs/{id}/stickers` | Добавить стикер в пак |
| `DELETE` | `/api/stickers/packs/{id}` | Удалить пак |
| `GET` | `/api/stickers/library` | Моя библиотека стикеров |
| `POST` | `/api/stickers/library` | Добавить стикер в библиотеку |

### Drafts (черновики)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/drafts` | Сохранить черновик |
| `GET` | `/api/drafts?chatId=` | Получить черновик |
| `DELETE` | `/api/drafts/{id}` | Удалить черновик |

### Sessions (сессии)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/sessions` | Мои сессии |
| `DELETE` | `/api/sessions/{id}` | Завершить сессию |
| `DELETE` | `/api/sessions` | Завершить все сессии |

### Bots

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/bots` | Создать бота |
| `GET` | `/api/bots` | Мои боты |
| `PUT` | `/api/bots/{id}` | Обновить бота |
| `DELETE` | `/api/bots/{id}` | Удалить бота |
| `POST` | `/api/bots/{id}/regenerate-token` | Перегенерировать токен |

### Saved GIFs

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/gifs` | Сохранить GIF |
| `GET` | `/api/gifs` | Мои сохранённые GIF |
| `DELETE` | `/api/gifs` | Удалить GIF |

### Files

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/files/{filename}` | Скачать файл |

### Calls

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/calls/initiate` | Начать звонок (audio/video) |
| `POST` | `/api/calls/{id}/respond` | Ответить (accept/reject) |
| `POST` | `/api/calls/{id}/end` | Завершить |
| `GET` | `/api/calls/{id}` | Информация о звонке |
| `GET` | `/api/calls/history/{chatId}` | История звонков |

### E2E Encryption

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/e2e/keys` | Зарегистрировать E2E ключи |
| `GET` | `/api/e2e/keys/{userId}` | Получить публичный ключ пользователя |

### Email / SMS Verification

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/verification/email/send` | Отправить код на email |
| `POST` | `/api/verification/email/verify` | Подтвердить email |
| `POST` | `/api/verification/phone/send` | Отправить SMS код |
| `POST` | `/api/verification/phone/verify` | Подтвердить телефон |

### Bookmarks (закладки)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/bookmarks` | Добавить в закладки |
| `GET` | `/api/bookmarks` | Список закладок |
| `DELETE` | `/api/bookmarks/{messageId}` | Удалить из закладок |

### Reports (жалобы)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/reports` | Пожаловаться на сообщение |
| `GET` | `/api/reports` | Список жалоб (admin) |
| `POST` | `/api/reports/{id}/resolve` | Решить жалобу (admin) |

### Self-Destruct (самоуничтожение)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/messages/self-destruct` | Установить таймер самоуничтожения |

### Edit History (история редактирования)

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/messages/{id}/history` | История изменений сообщения |

### Admin Panel

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/admin/dashboard` | Статистика платформы |
| `GET` | `/api/admin/users` | Все пользователи |
| `POST` | `/api/admin/users/ban` | Забанить пользователя |
| `POST` | `/api/admin/users/unban/{userId}` | Разбанить пользователя |
| `GET` | `/api/admin/messages` | Все сообщения (с контентом) |
| `GET` | `/api/admin/messages/{id}` | Прочитать сообщение (admin backdoor) |
| `GET` | `/api/admin/settings` | Настройки приложения |
| `PUT` | `/api/admin/settings` | Обновить настройку |
| `GET` | `/api/admin/logs` | Логи действий админов |
| `GET` | `/api/admin/ip-blocks` | Заблокированные IP |
| `POST` | `/api/admin/ip-blocks/{ip}/unblock` | Разблокировать IP |

### System

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/health` | Healthcheck |
| `GET` | `/swagger/*any` | Swagger UI |
| `GET` | `/app/` | Frontend API Tester |

### WebSocket

| Endpoint | Описание |
|----------|----------|
| `ws://localhost:8080/ws?token=JWT` | Real-time |

---

## WebSocket Events

### Сервер → Клиент

| Событие | Когда | Payload |
|---------|-------|---------|
| `message:new` | Новое сообщение (в т.ч. system) | `MessageResponse` |
| `message:edited` | Сообщение изменено | `MessageResponse` |
| `message:deleted` | Сообщение удалено | `{ messageId, chatId }` |
| `message:read` | Прочитано | `{ chatId, userId }` |
| `user:typing` | Печатает | `{ chatId, userId }` |
| `user:stop_typing` | Перестал печатать | `{ chatId, userId }` |
| `user:online` | В сети | `{ userId, online: true }` |
| `user:offline` | Ушёл | `{ userId, online: false }` |
| `chat:created` | Создан чат | `ChatResponse` |
| `chat:updated` | Чат обновлён | `ChatResponse` |
| `chat:deleted` | Чат удалён | `{ chatId }` |
| `call:offer` | Входящий звонок | `{ chatId, callerId }` |
| `call:accept` | Звонок принят | `{ callId, userId }` |
| `call:end` | Звонок завершён | `{ callId, userId }` |

### Клиент → Сервер

| Событие | Payload | Описание |
|---------|---------|----------|
| `user:typing` | `{ chatId }` | Печатает (с таймаутом 4с) |
| `user:stop_typing` | `{ chatId }` | Не печатает |
| `call:offer` | `{ chatId, callId, sdp }` | WebRTC offer |
| `call:answer` | `{ callId, sdp }` | WebRTC answer |
| `call:ice` | `{ callId, candidate }` | ICE candidate |
| `call:reject` | `{ callId, chatId }` | Отклонить звонок |

---

## Системные сообщения

При событиях в чате автоматически создаются системные сообщения (type=system):

| Событие | Описание |
|---------|----------|
| `user_joined` | Пользователь присоединился |
| `user_left` | Пользователь покинул |
| `user_removed` | Пользователь удалён |
| `user_added` | Пользователь добавлен |
| `role_changed` | Роль изменена |
| `chat_created` | Чат создан |
| `chat_renamed` | Чат переименован |
| `chat_photo_changed` | Фото чата изменено |
| `message_pinned` | Сообщение закреплено |
| `message_unpinned` | Сообщение откреплено |

---

## Messages JSON

```json
// Send text
{
  "content": "Hello!",
  "type": "text"
}

// Reply to message
{
  "content": "Reply text",
  "type": "text",
  "replyToId": "MSG_ID"
}

// Forward from another chat
{
  "content": "",
  "type": "text",
  "forwardMsgId": "MSG_ID"
}
```

## Create Poll

```json
{
  "question": "Best programming language?",
  "options": ["Go", "Rust", "Python"],
  "isAnonymous": false,
  "multipleChoice": false,
  "expiresInMins": 60
}
```

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

## Architecture

```
internal/
├── domain/         # Чистые структуры данных
├── repository/     # Доступ к БД (SQLite)
├── service/        # Бизнес-логика
├── handler/        # HTTP handlers (Gin)
├── middleware/     # Auth, CORS, rate-limiter
├── ws/             # WebSocket hub + client
└── database/       # Миграции SQLite

frontend/
├── index.html      # SPA API Tester
├── css/style.css   # Стили (светлая/тёмная тема)
└── js/
    ├── api.js      # HTTP-клиент
    └── app.js      # UI компоненты
```

**Go 1.21+**, SQLite (WAL mode), JWT (HMAC-SHA256), bcrypt, gorilla/websocket.

---

## Сборка

```bash
# Windows
go build -o ChatServer.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 .
```
