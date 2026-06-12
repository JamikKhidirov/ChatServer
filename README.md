# Chat Messenger Server

Сервер мессенджера на Go: REST API + WebSocket real-time + WebRTC звонки.

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

# Admin Register (admin_secret из ADMIN_SECRET env)
curl -X POST http://localhost:8080/api/auth/admin/register \
  -H "Content-Type: application/json" \
  -d '{"username":"admin","email":"admin@mail.com","password":"admin123","displayName":"Admin","admin_secret":"admin-secret-change-me"}'
```

---

## Архитектура

```
/domain/         — модели данных (User, Chat, Message, Poll, ...)
/repository/     — интерфейсы + SQLite реализация по доменам
  /user/         — userrepo
  /chat/         — chatrepo
  /message/      — messagerepo
  ...
/service/        — бизнес-логика по доменам
  /auth/         — authservice
  /user/         — userservice
  /chat/         — chatterservice
  /message/      — messageservice
  ...
/handler/        — HTTP хендлеры (Gin) по доменам
  /auth/         — authhandler
  /user/         — userhandler
  /chat/         — chathandler
  /message/      — messagehandler
  /ws/           — wshandler (WebSocket)
  ...
/docs/           — swagger.json + swagger.yaml
/pkg/response/   — утилиты ответов (JSON, Paginated, Error)
/middleware/      — JWT auth middleware
/config/         — конфигурация из env
/internal/ws/    — WebSocket hub + client
```

---

## API Endpoints

### Auth
| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/auth/register` | Регистрация |
| `POST` | `/api/auth/admin/register` | Регистрация администратора (требует admin_secret) |
| `POST` | `/api/auth/login` | Вход по email+password |
| `POST` | `/api/auth/login/email` | Отправка кода на email |
| `POST` | `/api/auth/login/email/verify` | Верификация email кода |
| `POST` | `/api/auth/login/phone` | Отправка кода на телефон |
| `POST` | `/api/auth/login/phone/verify` | Верификация SMS кода |
| `GET` | `/api/auth/refresh` | Обновить JWT токен |
| `PUT` | `/api/auth/change-password` | Сменить пароль |

### Users
| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/users/profile` | Профиль |
| `PUT` | `/api/users/profile` | Обновить профиль |
| `GET` | `/api/users/search?q=` | Поиск пользователей (пагинация) |
| `GET` | `/api/users/:id` | По ID |
| `GET` | `/api/users/username/:username` | По username |
| `PUT` | `/api/users/status` | Статус (online/offline) |
| `PUT` | `/api/users/push-token` | Push-токен |
| `POST` | `/api/users/block` | Заблокировать |
| `DELETE` | `/api/users/block/:userId` | Разблокировать |
| `GET` | `/api/users/blocked` | Список заблокированных |
| `POST` | `/api/users/avatar` | Загрузить аватар |
| `POST` | `/api/users/push-test` | Тестовый push |

### Chats
| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/chats` | Создать чат (private/group) |
| `POST` | `/api/chats/start/:userId` | Начать приватный чат с пользователем |
| `GET` | `/api/chats` | Список чатов |
| `GET` | `/api/chats/search?q=` | Поиск чатов |
| `GET` | `/api/chats/archived` | Архивные чаты |
| `GET` | `/api/chats/:id` | Детали чата |
| `PUT` | `/api/chats/:id` | Обновить группу |
| `DELETE` | `/api/chats/:id` | Удалить чат |
| `POST` | `/api/chats/:id/read` | Отметить прочитанным |
| `POST` | `/api/chats/:id/pin` | Закрепить |
| `DELETE` | `/api/chats/:id/pin` | Открепить |
| `POST` | `/api/chats/:id/archive` | Архивировать |
| `POST` | `/api/chats/:id/unarchive` | Разархивировать |
| `POST` | `/api/chats/:id/hide` | Скрыть |
| `POST` | `/api/chats/:id/leave` | Выйти из группы |
| `POST` | `/api/chats/:id/transfer-ownership` | Передать права |
| `PUT` | `/api/chats/:id/slow-mode` | Slow mode (0-3600 сек) |
| `PUT` | `/api/chats/:id/notifications` | Mute/unmute |
| `GET` | `/api/chats/:id/notifications` | Статус mute |
| `POST` | `/api/chats/:id/participants` | Добавить участника |
| `DELETE` | `/api/chats/:id/participants/:userId` | Удалить участника |
| `PUT` | `/api/chats/:id/participants/:userId/role` | Сменить роль |

### Messages
| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/chats/:id/messages` | Отправить |
| `GET` | `/api/chats/:id/messages` | Список (пагинация) |
| `GET` | `/api/chats/:id/messages/search?q=` | Поиск (пагинация) |
| `POST` | `/api/chats/:id/messages/file` | Загрузить файл |
| `POST` | `/api/chats/:id/messages/:msgId/resend` | Переотправить |
| `GET` | `/api/chats/:id/pinned` | Закреплённые |
| `GET` | `/api/chats/:id/media?type=` | Медиа (фото/video/audio) |
| `GET` | `/api/chats/:id/export` | Экспорт чата |
| `GET` | `/api/messages/:id` | По ID |
| `PUT` | `/api/messages/:id` | Редактировать |
| `DELETE` | `/api/messages/:id` | Удалить |
| `DELETE` | `/api/messages/:id/for-me` | Удалить у себя |
| `POST` | `/api/messages/:id/reactions` | Добавить реакцию |
| `DELETE` | `/api/messages/:id/reactions?emoji=` | Убрать реакцию |
| `PUT` | `/api/messages/:id/pin` | Закрепить/открепить |
| `POST` | `/api/messages/:id/star` | В избранное |
| `DELETE` | `/api/messages/:id/star` | Из избранного |
| `POST` | `/api/messages/:id/read` | Отметить прочитанным |
| `GET` | `/api/messages/search?q=` | Глобальный поиск |
| `GET` | `/api/messages/starred` | Избранные |
| `POST` | `/api/messages/forward` | Переслать |

### Contacts, Folders, Invite Links
| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/contacts/sync` | Синхронизация контактов |
| `GET` | `/api/contacts` | Список контактов |
| `GET` | `/api/contacts/search?q=` | Поиск по телефону |
| `GET` | `/api/contacts/registered` | Зарегистрированные |
| `POST` | `/api/contacts/photo` | Обновить фото контакта |
| `POST` | `/api/folders` | Создать папку |
| `GET` | `/api/folders` | Список папок |
| `PUT` | `/api/folders/:id` | Обновить папку |
| `DELETE` | `/api/folders/:id` | Удалить папку |
| `POST` | `/api/chats/:id/invite-links` | Создать ссылку |
| `GET` | `/api/chats/:id/invite-links` | Список ссылок |
| `DELETE` | `/api/chats/:id/invite-links/:linkId` | Удалить ссылку |
| `POST` | `/api/chats/join` | Присоединиться по ссылке |

### Polls, Stickers, Calls
| Method | Endpoint | Описание |
|--------|----------|----------|
| `POST` | `/api/chats/:id/polls` | Создать опрос |
| `GET` | `/api/chats/:id/polls` | Опросы чата |
| `POST` | `/api/polls/:pollId/vote` | Голосовать |
| `POST` | `/api/polls/:pollId/close` | Закрыть опрос |
| `POST` | `/api/stickers/packs` | Создать паку |
| `GET` | `/api/stickers/packs` | Все паки |
| `GET` | `/api/stickers/packs/my` | Мои паки |
| `GET` | `/api/stickers/packs/:id` | По ID |
| `DELETE` | `/api/stickers/packs/:id` | Удалить паку |
| `POST` | `/api/stickers/library` | В библиотеку |
| `GET` | `/api/stickers/library` | Моя библиотека |
| `POST` | `/api/calls` | Инициировать звонок |
| `POST` | `/api/calls/:id/respond` | Ответить на звонок |
| `POST` | `/api/calls/:id/end` | Завершить |
| `GET` | `/api/calls/:id` | По ID |
| `GET` | `/api/chats/:chatId/calls` | История звонков |

### Sessions, Bots, Drafts, Gifs
| Method | Endpoint | Описание |
|--------|----------|----------|
| `GET` | `/api/sessions` | Активные сессии |
| `DELETE` | `/api/sessions/:id` | Завершить сессию |
| `DELETE` | `/api/sessions` | Завершить все |
| `POST` | `/api/bots` | Создать бота |
| `GET` | `/api/bots` | Мои боты |
| `PUT` | `/api/bots/:id` | Обновить |
| `DELETE` | `/api/bots/:id` | Удалить |
| `POST` | `/api/bots/:id/token` | Новый токен |
| `POST` | `/api/drafts` | Сохранить черновик |
| `GET` | `/api/drafts?chatId=` | Получить черновик |
| `DELETE` | `/api/drafts/:id` | Удалить черновик |
| `POST` | `/api/gifs` | Сохранить GIF |
| `GET` | `/api/gifs` | Мои GIF |
| `DELETE` | `/api/gifs` | Удалить GIF |
| `POST` | `/api/messages/schedule` | Запланировать сообщение |
| `GET` | `/api/messages/scheduled` | Запланированные |
| `DELETE` | `/api/messages/scheduled/:id` | Отменить |

---

## WebSocket Events

```
ws://localhost:8080/ws?token=JWT_TOKEN
```

### Server → Client (server отправляет)

| Event | Описание | Триггер |
|-------|----------|---------|
| `message:new` | Новое сообщение | POST /api/chats/:id/messages |
| `message:edited` | Сообщение изменено | PUT /api/messages/:id |
| `message:deleted` | Сообщение удалено | DELETE /api/messages/:id |
| `message:read` | Прочитано | POST /api/chats/:id/read |
| `message:pinned` | Закреплено/откреплено | PUT /api/messages/:id/pin |
| `message:starred` | В избранное | POST /api/messages/:id/star |
| `chat:created` | Создан чат | POST /api/chats |
| `chat:updated` | Чат обновлён | PUT /api/chats/:id |
| `chat:deleted` | Чат удалён | DELETE /api/chats/:id |
| `chat:slowmode` | Slow mode изменён | PUT /api/chats/:id/slow-mode |
| `chat:role` | Роль участника изменена | PUT /api/chats/:id/participants/:userId/role |
| `chat:ownership` | Права переданы | POST /api/chats/:id/transfer-ownership |
| `user:online` | Пользователь онлайн | WebSocket connect |
| `user:offline` | Пользователь офлайн | WebSocket disconnect |
| `user:typing` | Пользователь печатает | WS event |
| `user:stop_typing` | Перестал печатать | WS event |
| `user:keyboard_opened` | Открыл клавиатуру | WS event |
| `user:keyboard_closed` | Закрыл клавиатуру | WS event |
| `call:offer` | Входящий звонок | POST /api/calls |
| `call:accept` | Звонок принят | POST /api/calls/:id/respond |
| `call:end` | Звонок завершён | POST /api/calls/:id/end |
| `poll:created` | Новый опрос | POST /api/chats/:id/polls |
| `poll:vote` | Новый голос | POST /api/polls/:pollId/vote |
| `poll:closed` | Опрос закрыт | POST /api/polls/:pollId/close |
| `folder:created` | Папка создана | POST /api/folders |
| `folder:updated` | Папка обновлена | PUT /api/folders/:id |
| `folder:deleted` | Папка удалена | DELETE /api/folders/:id |
| `invite:created` | Приглашение создано | POST /api/chats/:id/invite-links |
| `invite:joined` | Присоединился по ссылке | POST /api/chats/join |
| `sticker:added` | Стикер добавлен | POST /api/stickers/library |
| `sticker:pack_created` | Пак создан | POST /api/stickers/packs |

### Client → Server (клиент отправляет)

| Event | Payload | Описание |
|-------|---------|----------|
| `user:typing` | `{"chatId":"..."}` | Я печатаю |
| `user:stop_typing` | `{"chatId":"..."}` | Я перестал печатать |
| `user:keyboard_opened` | `{"chatId":"..."}` | Я открыл клавиатуру |
| `user:keyboard_closed` | `{"chatId":"..."}` | Я закрыл клавиатуру |

---

## Формат ответов

### Успех
```json
{
  "success": true,
  "data": { ... }
}
```

### Пагинация
```json
{
  "success": true,
  "data": [ ... ],
  "meta": { "total": 100, "offset": 0, "limit": 50 }
}
```

### Ошибка
```json
{
  "success": false,
  "error": "description",
  "code": "ERROR_CODE"
}
```

### Коды ошибок
| Код | HTTP | Описание |
|-----|------|----------|
| BAD_REQUEST | 400 | Невалидный запрос |
| UNAUTHORIZED | 401 | Требуется авторизация |
| FORBIDDEN | 403 | Нет доступа |
| NOT_FOUND | 404 | Не найдено |
| VALIDATION_ERROR | 400 | Ошибка валидации |
| DUPLICATE | 409 | Уже существует |
| INTERNAL_ERROR | 500 | Внутренняя ошибка |
| RATE_LIMIT | 429 | Слишком много запросов |

---

## Технологии

Go 1.21+, SQLite (WAL mode, modernc.org/sqlite), JWT (HMAC-SHA256), bcrypt, gorilla/websocket, Gin, Swaggo.

## Сборка

```bash
# Windows
go build -o ChatServer.exe .

# Linux
GOOS=linux GOARCH=amd64 go build -o chat-server-linux-amd64 .
```

## Переменные окружения

| Переменная | По умолчанию | Описание |
|------------|-------------|----------|
| `SERVER_PORT` | `8080` | Порт |
| `DATABASE_PATH` | `file:chat.db?...` | Путь к SQLite |
| `JWT_SECRET` | `super-secret-...` | Секрет JWT |
| `ADMIN_SECRET` | `admin-secret-change-me` | Секрет для регистрации админа |
| `JWT_TTL` | `86400` | Время жизни токена (сек) |
