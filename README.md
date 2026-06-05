# 🚀 Chat Messenger Server

**Сервер мессенджера на Go** с REST API, WebSocket, WebRTC звонками, push-уведомлениями и Swagger UI.

---

## 📋 Содержание

- [Быстрый старт](#-быстрый-старт)
- [Архитектура](#-архитектура)
- [API Endpoints](#-api-endpoints)
  - [Auth](#auth)
  - [Users](#users)
  - [Chats](#chats)
  - [Messages](#messages)
  - [Calls](#calls)
  - [WebSocket](#websocket)
  - [System](#system)
- [WebSocket Events](#-websocket-events)
- [Переменные окружения](#-переменные-окружения)
- [Push-уведомления (FCM)](#-push-уведомления-fcm)
- [Работа со Swagger](#-работа-со-swagger)
- [Примеры curl](#-примеры-curl)

---

## 🚀 Быстрый старт

### 1. Запуск бинарника

```bash
# Windows
./ChatServer.exe

# Linux
./chat-server-linux-amd64
```

### 2. Открой Swagger UI

```
http://localhost:8080/swagger/index.html
```

### 3. Зарегистрируй пользователя

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"john","email":"john@mail.com","password":"secret123","displayName":"John"}'
```

### 4. Получи JWT и авторизуйся в Swagger

Нажми **Authorize** → вставь `Bearer <токен>`.

---

## 🏗 Архитектура

```
├── main.go              # Точка входа, роутинг
├── config/              # Конфигурация (env vars)
├── docs/                # Swagger спецификация (swagger.json)
├── internal/
│   ├── domain/          # Модели данных (User, Chat, Message, Call)
│   ├── database/        # SQLite подключение + миграции
│   ├── repository/      # Запросы к БД (CRUD)
│   ├── service/         # Бизнес-логика
│   ├── handler/         # HTTP обработчики + WebSocket events
│   ├── middleware/      # JWT, CORS, Rate Limiter
│   └── ws/              # WebSocket hub + client
└── pkg/response/        # Утилиты для JSON-ответов
```

**Слои:** `domain → repository → service → handler → main`

---

## 📡 API Endpoints

### Auth

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/auth/register` | Регистрация нового пользователя |
| `POST` | `/api/auth/login` | Вход в систему |
| `GET` | `/api/auth/refresh` | Обновить JWT токен |
| `PUT` | `/api/auth/change-password` | Сменить пароль |

### Users

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/users/profile` | Профиль текущего пользователя |
| `PUT` | `/api/users/profile` | Обновить профиль |
| `PUT` | `/api/users/status` | Обновить статус (Available, Busy, ...) |
| `PUT` | `/api/users/push-token` | Обновить push-токен |
| `POST` | `/api/users/push-test` | Отправить тестовый push |
| `GET` | `/api/users/search?q=` | Поиск пользователей |
| `GET` | `/api/users/{id}` | Пользователь по ID |
| `GET` | `/api/users/username/{username}` | Пользователь по username |
| `POST` | `/api/users/block` | Заблокировать пользователя |
| `DELETE` | `/api/users/block/{userId}` | Разблокировать пользователя |
| `GET` | `/api/users/blocked` | Список заблокированных |
| `DELETE` | `/api/users/account` | Удалить аккаунт |

### Chats

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/chats` | Список чатов |
| `POST` | `/api/chats` | Создать чат (private/group) |
| `GET` | `/api/chats/{id}` | Детали чата |
| `PUT` | `/api/chats/{id}` | Обновить группу |
| `DELETE` | `/api/chats/{id}` | Удалить чат |
| `POST` | `/api/chats/{id}/participants` | Добавить участника |
| `DELETE` | `/api/chats/{id}/participants/{userId}` | Удалить участника |
| `PUT` | `/api/chats/{id}/participants/{userId}/role` | Назначить роль (admin/member) |
| `POST` | `/api/chats/{id}/leave` | Покинуть группу |
| `POST` | `/api/chats/{id}/read` | Отметить прочитанным |
| `PUT` | `/api/chats/{id}/notifications` | Включить/выключить уведомления |
| `GET` | `/api/chats/{id}/notifications` | Статус уведомлений |

### Messages

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/chats/{id}/messages` | Сообщения чата |
| `POST` | `/api/chats/{id}/messages` | Отправить сообщение |
| `GET` | `/api/chats/{id}/messages/search?q=` | Поиск по сообщениям |
| `POST` | `/api/chats/{id}/messages/file` | Загрузить файл (multipart) |
| `POST` | `/api/chats/{id}/messages/{msgId}/resend` | Переслать сообщение |
| `GET` | `/api/chats/{id}/pinned` | Закреплённые сообщения |
| `GET` | `/api/messages/{id}` | Сообщение по ID |
| `PUT` | `/api/messages/{id}` | Редактировать сообщение |
| `DELETE` | `/api/messages/{id}` | Удалить сообщение |
| `POST` | `/api/messages/{id}/reactions` | Добавить реакцию (emoji) |
| `DELETE` | `/api/messages/{id}/reactions?emoji=` | Удалить реакцию |
| `PUT` | `/api/messages/{id}/pin` | Закрепить/открепить |
| `POST` | `/api/messages/{id}/read` | Отметить прочитанным |

### Calls

| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/calls/initiate` | Начать звонок |
| `POST` | `/api/calls/{id}/respond` | Ответить на звонок (accept/reject) |
| `POST` | `/api/calls/{id}/end` | Завершить звонок |
| `GET` | `/api/calls/{id}` | Информация о звонке |
| `GET` | `/api/calls/history/{chatId}` | История звонков в чате |

### WebSocket

| Endpoint | Описание |
|----------|----------|
| `ws://localhost:8080/ws?token=JWT` | Real-time соединение |

### System

| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/health` | Healthcheck |

---

## 🔌 WebSocket Events

### События сервера (сервер → клиент)

| Событие | Описание | Payload |
|---------|----------|---------|
| `message:new` | Новое сообщение | `MessageResponse` |
| `message:edited` | Сообщение изменено | `MessageResponse` |
| `message:deleted` | Сообщение удалено | `{ messageId, chatId }` |
| `message:read` | Сообщение прочитано | `{ chatId, userId }` |
| `user:online` | Пользователь онлайн | `{ userId, online: true }` |
| `user:offline` | Пользователь офлайн | `{ userId, online: false }` |
| `user:typing` | Пользователь печатает | `{ chatId, userId }` |
| `user:stop_typing` | Перестал печатать | `{ chatId, userId }` |
| `chat:created` | Создан чат | `ChatResponse` |
| `chat:updated` | Чат обновлён | `ChatResponse` |
| `chat:deleted` | Чат удалён | `{ chatId }` |
| `call:offer` | Входящий звонок | `{ chatId, callerId }` |
| `call:accept` | Звонок принят | `{ callId, userId }` |
| `call:end` | Звонок завершён | `{ callId }` |

### Команды клиента (клиент → сервер)

```json
{ "type": "user:typing", "payload": { "chatId": "..." } }
{ "type": "user:stop_typing", "payload": { "chatId": "..." } }
```

---

## 🔧 Переменные окружения

| Переменная | По умолч. | Описание |
|------------|-----------|----------|
| `SERVER_PORT` | `8080` | Порт сервера |
| `DATABASE_PATH` | `file:chat.db?cache=shared&mode=rwc` | Путь к SQLite БД |
| `JWT_SECRET` | `super-secret-key-change-in-production` | Секрет для JWT |
| `JWT_TTL` | `86400` | Время жизни токена (сек) |
| `PUSH_ENABLED` | `false` | Включить push-уведомления |
| `FIREBASE_CREDENTIALS` | — | Server Key Firebase Cloud Messaging |

---

## 📬 Push-уведомления (FCM)

1. Получи **Server Key** из Firebase Console → Project Settings → Cloud Messaging
2. Установи переменные окружения:

```bash
set PUSH_ENABLED=true
set FIREBASE_CREDENTIALS=AIzaSy...
```

3. Обнови push-токен пользователя:

```bash
curl -X PUT http://localhost:8080/api/users/push-token \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"token":"<FCM_TOKEN>","provider":"fcm"}'
```

4. Отправь тестовое уведомление:

```bash
curl -X POST http://localhost:8080/api/users/push-test \
  -H "Authorization: Bearer <JWT>" \
  -H "Content-Type: application/json" \
  -d '{"title":"Привет!","body":"Тестовое уведомление"}'
```

---

## 📖 Работа со Swagger

Swagger UI доступен по адресу:

```
http://localhost:8080/swagger/index.html
```

**Авторизация:**
1. Зарегистрируйся через `POST /api/auth/register`
2. Скопируй `token` из ответа
3. Нажми **Authorize** в Swagger UI
4. Вставь `Bearer <токен>` → **Authorize**

---

## 📝 Примеры curl

<details>
<summary>Полный рабочий процесс</summary>

```bash
# 1. Регистрация
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{"username":"alice","email":"alice@mail.com","password":"123456","displayName":"Alice"}'

# Сохрани token из ответа
TOKEN="eyJhbG..."

# 2. Обновить статус
curl -X PUT http://localhost:8080/api/users/status \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"status":"Busy"}'

# 3. Поиск пользователей
curl -X GET "http://localhost:8080/api/users/search?q=bob" \
  -H "Authorization: Bearer $TOKEN"

# 4. Создать чат (private)
curl -X POST http://localhost:8080/api/chats \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"type":"private","participantIds":["USER_ID_2"]}'

# 5. Отправить сообщение
curl -X POST http://localhost:8080/api/chats/CHAT_ID/messages \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"content":"Привет!","type":"text"}'

# 6. Добавить реакцию
curl -X POST http://localhost:8080/api/messages/MSG_ID/reactions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"emoji":"👍"}'

# 7. Закрепить сообщение
curl -X PUT http://localhost:8080/api/messages/MSG_ID/pin \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"pin":true}'

# 8. Поиск по сообщениям
curl -X GET "http://localhost:8080/api/chats/CHAT_ID/messages/search?q=Привет" \
  -H "Authorization: Bearer $TOKEN"

# 9. Начать звонок
curl -X POST http://localhost:8080/api/calls/initiate \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d '{"chatId":"CHAT_ID"}'
```
</details>

---

## 🛠 Сборка из исходников

```bash
# Windows
go build -o ChatServer.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 .
```

Требования: Go 1.21+.

---

## 📄 Лицензия

MIT
