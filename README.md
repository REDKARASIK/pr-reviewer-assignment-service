# pr-reviewer-assignment-service

PR Reviewer Assignment Service  
Тестовое задание, осень 2025 (Avito Trainee Backend).

На этом этапе настроена инфраструктура проекта:

- PostgreSQL в Docker
- миграции через `migrate/migrate` в отдельном контейнере
- backend-приложение в Docker
- управление через `Makefile`

Дальше поверх этого будет развиваться бизнес-логика сервиса распределения ревьюверов.

---

## Стек

- **Go** — backend-сервис (сборка в Docker)
- **PostgreSQL 17 (alpine)** — основная БД
- **golang-migrate** (образ `migrate/migrate`) — применение SQL-миграций
- **Docker + Docker Compose**
- **Makefile** для удобных команд

---

## Быстрый старт

### 1. Зависимости

Нужно установить:

- [Docker](https://www.docker.com/)
- [Docker Compose](https://docs.docker.com/compose/)
- `make` (обычно уже есть в Linux/macOS)

### 2. Поднять окружение

#### Основной запуск:

```bash
make up
```

#### Альтернативный запуск:

Если тебе не нужно разделять шаги (db → migrate → app), ты можешь поднять весь стек одной командой:

```bash
docker compose up --build
