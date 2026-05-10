# Movie Night Planner Backend

REST API бэкенд для приложения "Movie Night Planner" (Планировщик кинопросмотра), написанный на Go с использованием PostgreSQL и интеграцией с TMDB API.

## Возможности

- **Управление киновечерами**: Создание, обновление, удаление киновечеров
- **Поиск фильмов**: Интеграция с TMDB API для поиска фильмов, просмотра постеров и трейлеров
- **Голосование**: Система голосования за фильмы (оценка 1-5 звезд)
- **Комментарии**: Возможность оставлять комментарии к киновечерам
- **JWT аутентификация**: Безопасная аутентификация пользователей
- **CORS**: Поддержка кросс-доменных запросов

## Требования

- Go 1.21+
- PostgreSQL 15+
- TMDB API ключ (бесплатно на [themoviedb.org](https://www.themoviedb.org/settings/api))

## Установка

### 1. Клонирование репозитория

```bash
git clone https://github.com/yourusername/movie-night-planner-backend.git
cd movie-night-planner-backend
```

### 2. Установка зависимостей

```bash
go mod download
```

### 3. Настройка окружения

Создайте файл `.env` на основе `.env.example`:

```bash
cp .env.example .env
```

Отредактируйте `.env` и добавьте ваш TMDB API ключ:

```env
TMDB_API_KEY=your-tmdb-api-key-here
JWT_SECRET=your-super-secret-key-change-in-production
```

### 4. Настройка базы данных

Создайте базу данных PostgreSQL:

```bash
createdb movie_night_planner
```

Или используйте Docker Compose (рекомендуется):

```bash
docker-compose -f docker/docker-compose.yml up -d
```

### 5. Запуск приложения

#### Локальный запуск

```bash
go run ./cmd/server
```

#### Использование Makefile

```bash
make run
```

#### Docker

```bash
make docker-build
make docker-run
```

## API Документация

### Базовый URL

```
http://localhost:8080/api/v1
```

### Аутентификация

#### Регистрация пользователя

```http
POST /api/v1/auth/register
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123",
  "username": "john_doe"
}
```

#### Вход пользователя

```http
POST /api/v1/auth/login
Content-Type: application/json

{
  "email": "user@example.com",
  "password": "securePassword123"
}
```

### Киновечера

#### Получить все киновечера

```http
GET /api/v1/evenings?page=1&limit=10
```

#### Получить киновечер по ID

```http
GET /api/v1/evenings/{id}
```

#### Создать киновечер

```http
POST /api/v1/evenings
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Пятница с пиццей",
  "description": "Смотрим классические комедии",
  "scheduled_at": "2026-05-15T20:00:00Z",
  "is_private": false
}
```

#### Обновить киновечер

```http
PUT /api/v1/evenings/{id}
Authorization: Bearer {token}
Content-Type: application/json

{
  "title": "Обновленное название",
  "is_private": true
}
```

#### Удалить киновечер

```http
DELETE /api/v1/evenings/{id}
Authorization: Bearer {token}
```

### Фильмы

#### Поиск фильмов

```http
GET /api/v1/movies/search?q=back+to+future&page=1
```

#### Получить детали фильма

```http
GET /api/v1/movies/{tmdbId}
```

#### Добавить фильм в киновечер

```http
POST /api/v1/evenings/{eveningId}/movies
Authorization: Bearer {token}
Content-Type: application/json

{
  "tmdb_id": 105
}
```

#### Удалить фильм из киновечера

```http
DELETE /api/v1/evenings/{eveningId}/movies/{tmdbId}
Authorization: Bearer {token}
```

### Голосования

#### Получить голоса для киновечера

```http
GET /api/v1/evenings/{eveningId}/votes
```

#### Голосовать за фильм

```http
POST /api/v1/evenings/{eveningId}/votes
Authorization: Bearer {token}
Content-Type: application/json

{
  "evening_film_id": "uuid",
  "value": 5
}
```

### Комментарии

#### Получить комментарии для киновечера

```http
GET /api/v1/evenings/{eveningId}/comments
```

#### Добавить комментарий

```http
POST /api/v1/evenings/{eveningId}/comments
Authorization: Bearer {token}
Content-Type: application/json

{
  "content": "Давайте смотреть в 20:00!"
}
```

## Структура проекта

```
movie-night-planner-backend/
├── cmd/
│   └── server/
│       └── main.go              # Точка входа приложения
├── internal/
│   ├── config/
│   │   └── config.go            # Конфигурация приложения
│   ├── database/
│   │   └── database.go          # Подключение к БД
│   ├── models/
│   │   └── models.go            # Модели данных
│   ├── repositories/
│   │   ├── user_repository.go   # Репозиторий пользователей
│   │   ├── evening_repository.go
│   │   ├── evening_film_repository.go
│   │   ├── vote_repository.go
│   │   └── comment_repository.go
│   ├── services/
│   │   ├── auth_service.go      # Сервис аутентификации
│   │   ├── evening_service.go
│   │   ├── movie_service.go
│   │   ├── vote_service.go
│   │   └── comment_service.go
│   ├── handlers/
│   │   ├── auth_handler.go      # HTTP handlers
│   │   ├── evening_handler.go
│   │   ├── movie_handler.go
│   │   ├── vote_handler.go
│   │   └── comment_handler.go
│   ├── middleware/
│   │   ├── auth.go              # JWT middleware
│   │   └── cors.go              # CORS middleware
│   ├── utils/
│   │   ├── jwt.go               # JWT утилиты
│   │   ├── password.go          # Хэширование паролей
│   │   └── errors.go            # Обработка ошибок
│   └── tmdb/
│       └── client.go            # TMDB API клиент
├── pkg/
│   └── response/
│       └── response.go          # Общие ответы API
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── .env.example
├── Makefile
├── go.mod
└── README.md
```

## 🐳 Docker

### Сборка образа

```bash
docker-compose -f docker/docker-compose.yml build
```

### Запуск

```bash
docker-compose -f docker/docker-compose.yml up
```

### Остановка

```bash
docker-compose -f docker/docker-compose.yml down
```

## Скрипты Makefile

```bash
make build       # Собрать приложение
make run         # Запустить приложение
make lint        # Запустить линтер
make clean       # Очистить артефакты сборки
make docker-build # Собрать Docker образ
make docker-run  # Запустить через Docker Compose
```
