# Task Tracker Backend — микросервисная архитектура (NestJS)

Полноценный backend для таск‑трекера на NestJS с микросервисами и API Gateway.

## 🏗️ Архитектура

- API Gateway — HTTP/WebSocket шлюз (порт 3000), Swagger: `/api/docs`
- Auth Service — аутентификация и авторизация
- Users Service — управление пользователями
- Teams Service — управление командами
- Boards Service — доски и колонки
- Tasks Service — задачи, комментарии, логи активности

## 🛠️ Технологии

- NestJS + TypeScript
- TypeORM (PostgreSQL)
- Mongoose (MongoDB)
- RabbitMQ (межсервисное взаимодействие)
- JWT (авторизация)
- WebSocket (realtime)
- Docker Compose (инфраструктура)

## 📋 Требования

- Node.js 20+
- Docker и Docker Compose
- pnpm (или npm)

## 🚀 Быстрый старт (Docker)

1) Установите зависимости локально (опционально, для разработки):
```bash
cd backend
pnpm install
# или
npm install
```

2) Запустите инфраструктуру и сервисы:
```bash
docker-compose up -d
```

Будут подняты:
- PostgreSQL (5432)
- MongoDB (27017)
- RabbitMQ (5672, 15672 — management UI)
- Все микросервисы

3) Swagger UI:
- API Gateway: `http://localhost:3000/api/docs`
  - Ассеты Swagger обслуживаются локально (без Webpack и без внешнего CDN).

## 🔧 Локальная разработка (без Docker)

1) Поднимите локально PostgreSQL, MongoDB и RabbitMQ  
2) Настройте `.env` для сервисов  
3) Запустите сервисы:
```bash
# В отдельных терминалах
pnpm start:dev:gateway
pnpm start:dev:auth
pnpm start:dev:users
pnpm start:dev:teams
pnpm start:dev:boards
pnpm start:dev:tasks

# Или все сразу (нужен concurrently)
pnpm start:all
```

## 📡 Основные API

Auth (публичные):
- POST `/api/auth/register` — регистрация
- POST `/api/auth/login` — вход
- POST `/api/auth/refresh` — обновление токена

Users (JWT):
- GET `/api/users/me` — текущий пользователь
- GET `/api/users/:id` — пользователь по ID
- PATCH `/api/users/:id` — обновление профиля

Teams (JWT):
- GET `/api/teams` — список команд
- POST `/api/teams` — создать команду
- GET `/api/teams/:id` — команда по ID
- POST `/api/teams/:id/members` — добавить участника
- DELETE `/api/teams/:id/members/:userId` — удалить участника

Boards (JWT):
- GET `/api/boards` — доски пользователя
- POST `/api/boards` — создать доску
- GET `/api/boards/:id` — доска с колонками
- PUT `/api/boards/:id` — обновить доску
- DELETE `/api/boards/:id` — удалить доску

Columns (JWT):
- POST `/api/columns` — создать колонку
- GET `/api/columns/board/:boardId` — колонки доски
- PUT `/api/columns/:id` — обновить колонку
- DELETE `/api/columns/:id` — удалить колонку

Tasks (JWT):
- POST `/api/tasks` — создать задачу
- GET `/api/tasks/:id` — получить задачу
- GET `/api/tasks/board/:boardId` — задачи доски
- PUT `/api/tasks/:id` — обновить задачу
- DELETE `/api/tasks/:id` — удалить задачу
- PUT `/api/tasks/:id/move` — переместить задачу
- POST `/api/tasks/:id/comment` — добавить комментарий
- GET `/api/tasks/:id/comments` — комментарии задачи

## 🔌 WebSocket

Подключение: `ws://localhost:3000`

События для подписки:
- `join:board`, `leave:board`, `join:team`

События от сервера:
- `task.created`, `task.updated`, `task.deleted`
- `board.updated`
- `team.memberAdded`, `team.memberRemoved`

## 🗄️ Базы данных

PostgreSQL (TypeORM):
- пользователи, команды, доски, колонки, задачи, refresh‑токены

MongoDB (Mongoose):
- комментарии к задачам
- логи активности

## 🐰 RabbitMQ

RPC‑взаимодействие между сервисами.  
Management UI: `http://localhost:15672` (admin/admin123)

## 🧪 Тестирование

```bash
# Unit
pnpm test
# Coverage
pnpm test:cov
```

## 📦 Структура проекта

```
backend/
├── apps/
│   ├── api-gateway/      # HTTP/WebSocket Gateway
│   ├── auth-service/     # Аутентификация
│   ├── users-service/    # Пользователи
│   ├── teams-service/    # Команды
│   ├── boards-service/   # Доски и колонки
│   └── tasks-service/    # Задачи и комментарии
├── libs/
│   ├── common/           # Общие интерфейсы и DTO
│   └── database/         # TypeORM entities и утилиты БД
├── docker-compose.yml    # Инфраструктура
└── package.json
```

## 🔐 Безопасность

- Короткоживущие access‑токены (JWT) + refresh‑токены
- Хэширование паролей (bcrypt)
- Валидация входных данных (class-validator)

## 🐛 Отладка

```bash
# Логи сервисов
docker-compose logs -f api-gateway
docker-compose logs -f auth-service
# и т.д.
```

## 📚 Ссылки

- NestJS: https://docs.nestjs.com/
- RabbitMQ: https://www.rabbitmq.com/documentation.html
- MongoDB: https://docs.mongodb.com/

## 🤝 Вклад

Приветствуются pull requests и issues!
