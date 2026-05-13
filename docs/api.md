
# HTTP API распределённой системы мониторинга Linux-узлов вычислительного кластера

## 1. Назначение

HTTP API используется для приема метрик от агентов и для выдачи данных клиентским интерфейсам. На текущем этапе API работает поверх HTTP и использует JSON для структурированных данных.

Основной поток записи:

```text
agent -> POST /api/v1/metrics -> server -> PostgreSQL
```

Основной поток чтения:

```text
web dashboard -> GET /api/v1/... -> server -> PostgreSQL
```

## 2. Базовый адрес

Для локального Docker Compose-стенда:

```text
http://localhost:8080
```

Внутри Docker-сети агенты обращаются к серверу по адресу:

```text
http://server:8080/api/v1/metrics
```

## 3. Health check

### `GET /health`

Проверяет, что процесс сервера запущен и принимает HTTP-запросы.

Успешный ответ:

```http
HTTP/1.1 200 OK
```

```text
ok
```

Пример:

```bash
curl http://localhost:8080/health
```

## 4. Прием метрик

### `POST /api/v1/metrics`

Принимает один снимок метрик от агента.

Заголовки запроса:

```http
Content-Type: application/json
Accept: application/json
```

Тело запроса должно соответствовать структуре `Metrics`.

Минимальный пример:

```json
{
  "timestamp": "2026-03-24T12:35:24.943217753Z",
  "host_id": "node-1",
  "cpu": {
    "user_pct": 0,
    "nice_pct": 0,
    "system_pct": 0.12,
    "idle_pct": 99.8,
    "iowait_pct": 0,
    "irq_pct": 0,
    "softirq_pct": 0.08,
    "steal_pct": 0,
    "total_pct": 0.2,
    "per_core_pct": {
      "cpu0": {
        "user_pct": 0,
        "nice_pct": 0,
        "system_pct": 0.1,
        "idle_pct": 99.9,
        "iowait_pct": 0,
        "irq_pct": 0,
        "softirq_pct": 0,
        "steal_pct": 0,
        "total_pct": 0.1
      }
    }
  },
  "mem": {
    "total_bytes": 8217731072,
    "available_bytes": 7573880832,
    "used_bytes": 643850240,
    "used_pct": 7.8349
  },
  "disk": [
    {
      "mount": "/",
      "total_bytes": 485473984512,
      "free_bytes": 481801867264,
      "used_bytes": 3672117248,
      "used_pct": 0.7564
    }
  ],
  "network": {
    "rx_bytes_total": 2118950,
    "tx_bytes_total": 26008,
    "rx_bps_total": 0,
    "tx_bps_total": 0,
    "ifaces": [
      {
        "name": "eth0",
        "rx_bytes": 2118950,
        "tx_bytes": 26008,
        "rx_bps": 0,
        "tx_bps": 0
      }
    ]
  }
}
```

Успешный ответ:

```http
HTTP/1.1 200 OK
```

```text
ok
```

Правила валидации:

- `host_id` обязателен и не должен быть пустым;
- `timestamp` обязателен и не должен быть нулевым;
- тело запроса должно быть валидным JSON.

Коды ответа:

- `200 OK` — снимок принят и сохранен;
- `400 Bad Request` — некорректный JSON или нарушение минимального контракта;
- `405 Method Not Allowed` — использован неподдерживаемый HTTP-метод;
- `500 Internal Server Error` — ошибка сохранения или внутренняя ошибка сервера.

Пример:

```bash
curl -X POST http://localhost:8080/api/v1/metrics \
  -H "Content-Type: application/json" \
  -d @metrics.json
```

## 5. Список устройств

### `GET /api/v1/devices`

Возвращает список зарегистрированных узлов, отсортированный по времени последней активности.

Успешный ответ:

```json
[
  {
    "id": 1,
    "host_id": "node-1",
    "first_seen_at": "2026-03-24T12:30:00Z",
    "last_seen_at": "2026-03-24T12:35:24Z",
    "online": true
  }
]
```

Семантика поля `online`:

- `true` — последний снимок получен не более 30 секунд назад;
- `false` — последний снимок старше 30 секунд.

Коды ответа:

- `200 OK` — список успешно получен;
- `405 Method Not Allowed` — использован неподдерживаемый метод;
- `500 Internal Server Error` — ошибка чтения из хранилища.

Пример:

```bash
curl http://localhost:8080/api/v1/devices
```

## 6. Последний снимок по узлу

### `GET /api/v1/devices/latest?host_id=<id>`

Возвращает последний сохраненный снимок метрик для указанного `host_id`.

Параметры запроса:

- `host_id` — обязательный идентификатор узла.

Успешный ответ:

```json
{
  "timestamp": "2026-03-24T12:35:24Z",
  "host_id": "node-1",
  "cpu": {
    "total_pct": 0.2,
    "per_core_pct": {}
  },
  "mem": {
    "total_bytes": 8217731072,
    "available_bytes": 7573880832,
    "used_bytes": 643850240,
    "used_pct": 7.8349
  },
  "disk": [],
  "network": {
    "rx_bytes_total": 2118950,
    "tx_bytes_total": 26008,
    "rx_bps_total": 0,
    "tx_bps_total": 0,
    "ifaces": []
  }
}
```

Коды ответа:

- `200 OK` — снимок найден;
- `400 Bad Request` — параметр `host_id` отсутствует;
- `405 Method Not Allowed` — использован неподдерживаемый метод;
- `500 Internal Server Error` — снимок не найден или произошла ошибка чтения.

Пример:

```bash
curl "http://localhost:8080/api/v1/devices/latest?host_id=node-1"
```

## 7. История метрик по узлу

### `GET /api/v1/devices/metrics?host_id=<id>&limit=<n>`

Возвращает историю снимков для указанного `host_id`. Снимки отсортированы от новых к старым.

Параметры запроса:

- `host_id` — обязательный идентификатор узла;
- `limit` — необязательное количество снимков, по умолчанию `100`.

Правила для `limit`:

- значение должно быть числом;
- значение должно быть больше нуля;
- значения больше `1000` ограничиваются сервером до `1000`.

Успешный ответ:

```json
[
  {
    "timestamp": "2026-03-24T12:35:24Z",
    "host_id": "node-1",
    "cpu": {
      "total_pct": 0.2,
      "per_core_pct": {}
    },
    "mem": {
      "total_bytes": 8217731072,
      "available_bytes": 7573880832,
      "used_bytes": 643850240,
      "used_pct": 7.8349
    },
    "disk": [],
    "network": {
      "rx_bytes_total": 2118950,
      "tx_bytes_total": 26008,
      "rx_bps_total": 0,
      "tx_bps_total": 0,
      "ifaces": []
    }
  }
]
```

Коды ответа:

- `200 OK` — история получена;
- `400 Bad Request` — отсутствует `host_id`, `limit` не является числом или `limit <= 0`;
- `405 Method Not Allowed` — использован неподдерживаемый метод;
- `500 Internal Server Error` — ошибка чтения из хранилища.

Пример:

```bash
curl "http://localhost:8080/api/v1/devices/metrics?host_id=node-1&limit=80"
```

## 8. Диагностический эндпоинт

### `GET /debug/api`

Возвращает последний сохраненный снимок метрик независимо от `host_id`. Эндпоинт предназначен для локальной отладки и не должен рассматриваться как стабильный публичный контракт.

Пример:

```bash
curl http://localhost:8080/debug/api
```

## 9. Требования к клиентам

Клиент, отправляющий метрики, должен:

- передавать `Content-Type: application/json`;
- использовать полный URL эндпоинта приема метрик;
- обеспечивать непустой `host_id`;
- передавать ненулевой `timestamp`;
- учитывать, что доставка на текущем этапе является синхронной и не содержит серверной очереди.

Клиент, читающий данные, должен:

- указывать `host_id` для запросов последнего снимка и истории;
- учитывать сортировку истории от новых снимков к старым;
- обрабатывать временное отсутствие данных по узлу как штатную ситуацию начального запуска.
