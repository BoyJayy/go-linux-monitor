# Запуск распределённой системы мониторинга Linux-узлов вычислительного кластера

## 1. Назначение

Документ описывает запуск текущего демонстрационного стенда и базовые проверки работоспособности. Основной поддерживаемый сценарий на данном этапе — Docker Compose.

## 2. Требования

Для запуска стенда необходимы:

- Docker;
- Docker Compose;
- свободный локальный порт `8080`;
- файл `.env`, подготовленный на основе `.env.example`.

Для локальной разработки без контейнеров дополнительно требуется Go `1.25.5`, указанный в `go.mod` и `go.work`.

## 3. Подготовка окружения

Создать файл `.env`:

```bash
cp .env.example .env
```

Проверить значения:

```text
POSTGRES_USER=postgres
POSTGRES_PASSWORD=change_me
POSTGRES_DB=monitoring service
DATABASE_URL=postgres://postgres:change_me@postgres:5432/monitoring%20service?sslmode=disable
SERVER_URL=http://localhost:8080
COLLECTION_INTERVAL=2s
REQUEST_TIMEOUT=3s
```

Важно: если имя базы содержит пробел, в `DATABASE_URL` оно должно быть URL-encoded как `monitoring%20service`.

## 4. Запуск Docker Compose

```bash
docker compose up --build
```

Состав стенда:

- `monitoring-postgres`;
- `monitoring-server`;
- `monitoring-agent-node-1`;
- `monitoring-agent-node-2`;
- `monitoring-agent-node-3`.

Сервер публикуется на локальный порт:

```text
http://localhost:8080
```

## 5. Проверка после запуска

Проверить эндпоинт состояния:

```bash
curl http://localhost:8080/health
```

Ожидаемый ответ:

```text
ok
```

Проверить список устройств:

```bash
curl http://localhost:8080/api/v1/devices
```

Ожидаемый результат — JSON-массив с узлами `node-1`, `node-2`, `node-3` после того, как агенты успели отправить первые снимки.

Проверить последний снимок:

```bash
curl "http://localhost:8080/api/v1/devices/latest?host_id=node-1"
```

Проверить историю:

```bash
curl "http://localhost:8080/api/v1/devices/metrics?host_id=node-1&limit=10"
```

## 6. Веб-интерфейс

Веб-панель доступна по адресу:

```text
http://localhost:8080
```

Панель отображает:

- список узлов;
- количество online/offline;
- последний снимок выбранного узла;
- графики CPU, памяти, RX/TX;
- таблицу последних измерений.

Автообновление выполняется каждые 2 секунды.

## 7. Остановка стенда

Остановить контейнеры:

```bash
docker compose down
```

Остановить контейнеры и удалить volume PostgreSQL:

```bash
docker compose down -v
```

Удаление volume приводит к потере сохраненной истории и повторному применению SQL-файлов из `server/migrations` при следующем запуске.

## 8. Локальный запуск сервера без Docker

Для локального запуска сервера требуется доступный PostgreSQL и переменная `DATABASE_URL`.

```bash
cd server
DATABASE_URL="postgres://postgres:change_me@localhost:5432/monitoring%20service?sslmode=disable" \
go run ./cmd/server
```

Сервер будет слушать порт `8080`.

## 9. Локальный запуск агента без Docker

Агент должен получить полный URL эндпоинта приема метрик:

```bash
cd agent
SERVER_URL="http://localhost:8080/api/v1/metrics" \
COLLECTION_INTERVAL="2s" \
REQUEST_TIMEOUT="3s" \
HOST_ID="local-node" \
go run ./cmd/agent
```

При запуске на Linux-узле агент собирает метрики той системы, где он выполняется.

## 10. Особенности миграций

В Docker Compose SQL-файлы из `server/migrations` монтируются в `/docker-entrypoint-initdb.d`. PostgreSQL выполняет эти файлы только при первичной инициализации пустого volume.

Если volume уже существует, новые SQL-файлы автоматически не применяются. Для демонстрационного стенда допустимо удалить volume командой:

```bash
docker compose down -v
```

Для промышленной эксплуатации требуется отдельный управляемый механизм миграций.

## 11. Диагностика

Проверить конфигурацию Compose:

```bash
docker compose config
```

Проверить логи сервера:

```bash
docker compose logs -f server
```

Проверить логи агента:

```bash
docker compose logs -f agent-node-1
```

Проверить доступность PostgreSQL внутри Compose-сети можно по healthcheck контейнера `postgres`.
