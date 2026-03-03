## Локальная разработка (Docker + Air)

Dev‑окружение поднимается в Docker и использует Air для live‑reload Go‑кода при изменениях файлов. [github](https://github.com/air-verse/air)

### Старт окружения

```bash
docker compose -f docker-compose.dev.yml up --build
```

После запуска backend‑сервис будет доступен по адресу и порту, указанным в `docker-compose.dev.yml` (например, `http://localhost:8080`). [oneuptime](https://oneuptime.com/blog/post/2026-01-07-go-hot-reloading-docker-air/view)

### Остановка окружения

```bash
docker compose -f docker-compose.dev.yml down
```

***

## Управление зависимостями Go

При добавлении новых импортов или модулей обновляй зависимости внутри контейнера. [medium.easyread](https://medium.easyread.co/today-i-learned-golang-live-reload-for-development-using-docker-compose-air-ecc688ee076)

### Обновить зависимости (go.mod / go.sum)

```bash
docker compose -f docker-compose.dev.yml exec svc go mod tidy
```

Где `svc` — имя сервиса с Go‑приложением в `docker-compose.dev.yml`. [stackoverflow](https://stackoverflow.com/questions/72546523/live-auto-reload-of-golang-apps-cosmtrek-air)

***

## Структура dev‑шаблона

Рекомендуемая структура проекта для работы с Air в Docker: [blog.logrocket](https://blog.logrocket.com/using-air-go-implement-live-reload/)

```text
.
├── .air.toml             # конфиг Air (watch + build + run)
├── Dockerfile.dev        # dev Dockerfile c установленным air
├── docker-compose.dev.yml
├── go.mod
├── go.sum
├── cmd/
│   └── app/main.go       # точка входа приложения
└── internal/ ...         # остальной код
```

***

## Пример конфигурации Air

В корне проекта создаётся `.air.toml`, который описывает, как собирать и запускать сервис. [blog.logrocket](https://blog.logrocket.com/using-air-go-implement-live-reload/)

```toml
# .air.toml
root = "."
tmp_dir = "tmp"

[build]
cmd = "go build -o ./tmp/app ./cmd/app"
bin = "./tmp/app"
include_ext = ["go"]
exclude_dir = ["tmp", "vendor"]

[log]
time = true
```

Air в контейнере следит за изменениями в смонтированном исходном коде и при каждом сохранении пересобирает и перезапускает бинарь. [courses.devopsdirective](https://courses.devopsdirective.com/docker-beginner-to-pro/lessons/11-development-workflow/01-hot-reloading)

***

## Пример Dockerfile.dev

```dockerfile
# Dockerfile.dev
FROM golang:1.22

WORKDIR /app

# зависимости
COPY go.mod go.sum ./
RUN go mod download

# установка air
RUN go install github.com/air-verse/air@latest

# исходники (в dev обычно перекрываются volume из docker-compose)
COPY . .

CMD ["air", "-c", ".air.toml"]
```

Этот Dockerfile устанавливает Air в dev‑образ и использует его как основной процесс контейнера. [stackoverflow](https://stackoverflow.com/questions/77973575/golang-hot-reloading-not-working-on-docker-with-compiledaemon-in-dev-mode)

***

## Пример docker-compose.dev.yml

```yaml
services:
  svc:
    build:
      context: .
      dockerfile: Dockerfile.dev
      target: dev
    ports:
      - "8080:8080"
    environment:
      - DATABASE_HOST=postgres
      - DATABASE_PORT=5432
    volumes:
      - .:/app:Z
      - go_cache:/app/tmp

volumes:
  go_cache:
```

Благодаря volume `.:/app:Z` Air внутри контейнера видит изменения локального кода и перезапускает приложение без пересборки самого контейнера. [betterprogramming](https://betterprogramming.pub/a-good-way-to-do-live-reload-for-go-b3707eb47336)
