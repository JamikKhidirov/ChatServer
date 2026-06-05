# ChatServerGolang

Полноценный сервер мессенджера на **Go** с поддержкой реального времени (WebSocket), голосовых/видеозвонков (WebRTC signaling) и push-уведомлений (FCM).

## Возможности

- Регистрация и авторизация (JWT)
- Личные и групповые чаты (создание, обновление, удаление, выход)
- Отправка, редактирование и удаление сообщений
- Статусы пользователя (как в WhatsApp: Available, Busy, At work, ...)
- Push-уведомления (Firebase Cloud Messaging) + тестовый эндпоинт
- Голосовые и видеозвонки (WebRTC signaling через WebSocket)
- Real-time обновления (онлайн/офлайн, печатает, новые сообщения)
- Поиск пользователей
- REST API + Swagger документация (браузер + терминал)
- Чистая архитектура (Clean Architecture)
- SQLite (не требует внешней БД)

## Быстрый старт

### Требования

- Go 1.21+
- Git (опционально)

### Запуск

```bash
# 1. Клонировать репозиторий
git clone <repo-url>
cd ChatServerGolang

# 2. Скачать зависимости
go mod tidy

# 3. Запустить сервер
go run main.go
```

Сервер запустится на `http://localhost:8080`.

### Переменные окружения

| Переменная | Значение по умолчанию | Описание |
|-----------|----------------------|----------|
| `SERVER_PORT` | `8080` | Порт сервера |
| `DATABASE_PATH` | `file:chat.db?cache=shared&mode=rwc` | Путь к SQLite БД |
| `JWT_SECRET` | `super-secret-key-change-in-production` | Секретный ключ JWT |
| `JWT_TTL` | `86400` | Время жизни токена (сек) |
| `ALLOW_ORIGINS` | `*` | CORS origins |
| `PUSH_ENABLED` | `false` | Включить push-уведомления |
| `FIREBASE_CREDENTIALS` | `` | Server Key Firebase Cloud Messaging |

## Архитектура

```
ChatServerGolang/
├── main.go                 # Точка входа
├── config/                 # Конфигурация
├── docs/                   # Swagger документация
├── internal/
│   ├── database/           # Слой базы данных (SQLite)
│   ├── domain/             # Модели данных
│   ├── handler/            # HTTP/WS хендлеры
│   ├── middleware/         # Middleware (CORS, JWT)
│   ├── repository/         # Репозитории (слой доступа к данным)
│   ├── service/            # Бизнес-логика
│   └── ws/                 # WebSocket hub и client
└── pkg/response/           # Утилиты для ответов API
```

## API Endpoints

### Аутентификация

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/auth/register` | Регистрация |
| POST | `/api/auth/login` | Вход |

### Пользователи (требуют JWT)

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/users/profile` | Профиль текущего пользователя |
| PUT | `/api/users/profile` | Обновить профиль |
| PUT | `/api/users/status` | Обновить статус (как в WhatsApp) |
| PUT | `/api/users/push-token` | Обновить push-токен |
| POST | `/api/users/push-test` | Отправить тестовое push-уведомление |
| GET | `/api/users/search?q=` | Поиск пользователей |
| GET | `/api/users/:id` | Получить пользователя по ID |
| GET | `/api/users/username/:username` | Получить пользователя по username |

### Чаты

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/chats` | Список чатов |
| POST | `/api/chats` | Создать чат (private/group) |
| GET | `/api/chats/:id` | Детали чата |
| PUT | `/api/chats/:id` | Обновить группу (имя, описание) |
| DELETE | `/api/chats/:id` | Удалить чат |
| POST | `/api/chats/:id/participants` | Добавить участника |
| DELETE | `/api/chats/:id/participants/:userId` | Удалить участника |
| PUT | `/api/chats/:id/participants/:userId/role` | Назначить admin/member |
| POST | `/api/chats/:id/leave` | Покинуть группу |
| POST | `/api/chats/:id/read` | Отметить как прочитанное |

### Сообщения

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/api/chats/:id/messages` | Сообщения чата |
| POST | `/api/chats/:id/messages` | Отправить сообщение |
| PUT | `/api/messages/:id` | Редактировать сообщение |
| DELETE | `/api/messages/:id` | Удалить сообщение |

### Звонки

| Метод | Путь | Описание |
|-------|------|----------|
| POST | `/api/calls/initiate` | Начать звонок |
| POST | `/api/calls/:id/respond` | Принять/отклонить |
| POST | `/api/calls/:id/end` | Завершить звонок |
| GET | `/api/calls/:id` | Информация о звонке |
| GET | `/api/calls/history/:chatId` | История звонков |

### Статус пользователя

Каждый пользователь имеет текстовый статус (как в WhatsApp):

- `Available` — по умолчанию при регистрации
- `Busy` — занят
- `At work` — на работе
- `Sleeping` — спит
- Любой кастомный текст (до 100 символов)

Изменить статус:
```bash
curl -X PUT http://localhost:8080/api/users/status \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"status": "At work"}'
```

Статус отображается в профиле, поиске, чатах и WebSocket-событиях.

### Другое

| Метод | Путь | Описание |
|-------|------|----------|
| GET | `/ws?token=` | WebSocket (real-time) |
| GET | `/health` | Healthcheck |
| GET | `/swagger/index.html` | Swagger UI |

## WebSocket

Подключение: `ws://localhost:8080/ws?token=<JWT_TOKEN>`

### События клиента

```json
{ "type": "user:typing", "payload": { "chatId": "..." } }
{ "type": "user:stop_typing", "payload": { "chatId": "..." } }
```

### События сервера

```json
{ "type": "message:new", "payload": { /* MessageResponse */ } }
{ "type": "message:edited", "payload": { /* MessageResponse */ } }
{ "type": "message:deleted", "payload": { "messageId": "..." } }
{ "type": "message:read", "payload": { "chatId": "...", "userId": "..." } }
{ "type": "user:typing", "payload": { "chatId": "...", "userId": "..." } }
{ "type": "user:stop_typing", "payload": { "chatId": "...", "userId": "..." } }
{ "type": "user:online", "payload": { "userId": "..." } }
{ "type": "user:offline", "payload": { "userId": "..." } }
{ "type": "chat:created", "payload": { /* ChatResponse */ } }
{ "type": "chat:updated", "payload": { /* ChatResponse */ } }
{ "type": "chat:deleted", "payload": { "chatId": "..." } }
{ "type": "call:offer", "payload": { "callId", "chatId", "sdp" } }
{ "type": "call:answer", "payload": { "callId", "sdp" } }
{ "type": "call:ice", "payload": { "callId", "candidate", "sdpMLineIdx" } }
{ "type": "call:end", "payload": { "callId" } }
```

## Push-уведомления

### Настройка FCM

Для включения push-уведомлений через Firebase Cloud Messaging:

1. Получите **Server Key** из Firebase Console
2. Запустите сервер с настройками:

```bash
set PUSH_ENABLED=true
set FIREBASE_CREDENTIALS=<YOUR_FCM_SERVER_KEY>
go run main.go
```

### Регистрация устройства

Клиент отправляет свой push-токен на сервер:

```bash
curl -X PUT http://localhost:8080/api/users/push-token \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"token": "<FCM_DEVICE_TOKEN>", "provider": "fcm"}'
```

### Тестирование без FCM

Если `FIREBASE_CREDENTIALS` не задан, push-уведомления логируются в консоль:
```
[FCM] Would send push (no credentials configured): title=..., body=...
```

### Тестовый эндпоинт

Для проверки интеграции клиента:

```bash
curl -X POST http://localhost:8080/api/users/push-test \
  -H "Authorization: Bearer <TOKEN>" \
  -H "Content-Type: application/json" \
  -d '{"title": "Hello", "body": "This is a test push"}'
```

### Когда отправляются пуши

- **Новое сообщение** — всем участникам чата (кроме отправителя)
- **Входящий звонок** — всем участникам чата (кроме инициатора)
- **Тестовый пуш** — только текущему пользователю

## Звонки (WebRTC)

Звонки реализованы через REST API + WebSocket signaling:

1. **POST** `/api/calls/initiate` — инициировать звонок  
   → Сервер создает запись `Call` и шлет уведомление через WS
2. Клиент B получает `call:offer` через WebSocket
3. **POST** `/api/calls/:id/respond` — `{ "action": "accept" }`  
   → Статус меняется на `ongoing`
4. Обмен ICE кандидатами через WebSocket (`call:ice`)
5. **POST** `/api/calls/:id/end` — завершить звонок

### Пример запроса

```bash
curl -X POST http://localhost:8080/api/calls/initiate \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{"chatId": "<CHAT_ID>"}'
```

## Примеры запросов

### Регистрация

```bash
curl -X POST http://localhost:8080/api/auth/register \
  -H "Content-Type: application/json" \
  -d '{
    "username": "john_doe",
    "email": "john@example.com",
    "password": "password123",
    "displayName": "John Doe"
  }'
```

### Вход

```bash
curl -X POST http://localhost:8080/api/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "john@example.com",
    "password": "password123"
  }'
```

### Создать чат

```bash
curl -X POST http://localhost:8080/api/chats \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "name": "Group Chat",
    "type": "group",
    "participantIds": ["user-id-2"]
  }'
```

### Отправить сообщение

```bash
curl -X POST http://localhost:8080/api/chats/<CHAT_ID>/messages \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <TOKEN>" \
  -d '{
    "content": "Hello World!",
    "type": "text"
  }'
```

## Swagger (интерактивная документация)

### Способ 1: Встроенный Swagger UI (рекомендуемый)

Запустите сервер и откройте в браузере:
```
http://localhost:8080/swagger/index.html
```

**Как пользоваться:**
1. Откройте ссылку в браузере — загрузится Swagger UI
2. Нажмите **Authorize** (🔓) в правом верхнем углу
3. Введите токен: `Bearer <JWT_TOKEN>`
4. Нажимайте на любой endpoint → **Try it out** → заполните параметры → **Execute**
5. Для регистрации/входа токен не нужен (эти ручки без авторизации)

**Порядок тестирования в Swagger:**
1. `/api/auth/register` — создать двух пользователей, скопировать token
2. Authorize → вставить токен первого пользователя
3. `/api/users/search?q=` — найти ID второго пользователя
4. `/api/chats` (POST) — создать чат с этим пользователем
5. `/api/chats/{id}/messages` (POST) — отправить сообщение
6. И так далее по всем ручкам

### Способ 2: Файл `docs/swagger.yaml`

Откройте `docs/swagger.yaml` в [Swagger Editor](https://editor.swagger.io/).  
Или используйте curl для локального просмотра:
```bash
curl http://localhost:8080/swagger.yaml
```

### Способ 3: Raw JSON (для интеграции)
```bash
curl http://localhost:8080/swagger/doc.json
```
