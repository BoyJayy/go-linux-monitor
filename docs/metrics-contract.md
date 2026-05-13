
# Контракт метрик распределённой системы мониторинга Linux-узлов вычислительного кластера

## 1. Назначение документа

Документ фиксирует структуру JSON-снимка, который агент отправляет на сервер. Контракт используется одновременно в `agent`, `server` и клиентских интерфейсах. На уровне Go он описан в модуле `api`.

Корневой тип:

```go
type Metrics struct {
    Timestamp time.Time      `json:"timestamp"`
    HostID    string         `json:"host_id"`
    CPU       UsageWithCores `json:"cpu"`
    Mem       MemoryUsage    `json:"mem"`
    Disk      []DiskUsage    `json:"disk"`
    Network   NetworkUsage   `json:"network"`
}
```

## 2. Общая структура JSON

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

## 3. Верхний уровень

| Поле | Тип | Обязательность | Описание |
| --- | --- | --- | --- |
| `timestamp` | string, RFC3339/RFC3339Nano | да | Время формирования снимка агентом. |
| `host_id` | string | да | Стабильный идентификатор узла. |
| `cpu` | object | да | Сводная и поядерная CPU-статистика. |
| `mem` | object | да | Статистика оперативной памяти. |
| `disk` | array | да | Список смонтированных файловых систем. |
| `network` | object | да | Сетевые счетчики и расчетные скорости. |

`host_id` определяется агентом в следующем порядке:

1. переменная окружения `HOST_ID`;
2. содержимое `/etc/machine-id`;
3. файл из `HOST_ID_FILE`;
4. резервный файл `/var/lib/myagent/host_id`;
5. новый UUID, сохраненный в резервный файл.

## 4. CPU

Источник данных: `/proc/stat`.

CPU-метрики рассчитываются на основании двух снимков счетчиков, разделенных интервалом `COLLECTION_INTERVAL`.

| Поле | Тип | Единица | Описание |
| --- | --- | --- | --- |
| `user_pct` | number | % | Доля времени CPU в пользовательском режиме. |
| `nice_pct` | number | % | Доля времени процессов с измененным nice-приоритетом. |
| `system_pct` | number | % | Доля времени в режиме ядра. |
| `idle_pct` | number | % | Доля простоя. |
| `iowait_pct` | number | % | Доля ожидания операций ввода-вывода. |
| `irq_pct` | number | % | Доля обработки аппаратных прерываний. |
| `softirq_pct` | number | % | Доля обработки программных прерываний. |
| `steal_pct` | number | % | Доля времени, забранная гипервизором. |
| `total_pct` | number | % | Суммарная занятость CPU без учета idle и iowait. |
| `per_core_pct` | object | % | Метрики по отдельным ядрам, ключи имеют вид `cpu0`, `cpu1`. |

Расчет `total_pct`:

```text
busy = delta_total - delta_idle
total_pct = busy / delta_total * 100
```

Где `delta_idle` включает `idle + iowait`.

## 5. Memory

Источник данных: `/proc/meminfo`.

| Поле | Тип | Единица | Описание |
| --- | --- | --- | --- |
| `total_bytes` | integer | bytes | Объем памяти из `MemTotal`. |
| `available_bytes` | integer | bytes | Доступная память из `MemAvailable`. |
| `used_bytes` | integer | bytes | `total_bytes - available_bytes`. |
| `used_pct` | number | % | Доля использованной памяти. |

Значения `/proc/meminfo` читаются в KiB и приводятся к байтам умножением на `1024`.

## 6. Disk

Источники данных:

- список mountpoint из `/proc/mounts`;
- значения файловой системы через `statfs`.

Метрики диска собираются по смонтированным файловым системам, а не напрямую по физическим устройствам.

| Поле | Тип | Единица | Описание |
| --- | --- | --- | --- |
| `mount` | string | - | Точка монтирования. |
| `total_bytes` | integer | bytes | Общий объем файловой системы. |
| `free_bytes` | integer | bytes | Свободное пространство по `statfs.Bfree`. |
| `used_bytes` | integer | bytes | `total_bytes - free_bytes`. |
| `used_pct` | number | % | Доля использованного пространства. |

Из сбора исключаются системные и временные файловые системы, включая `proc`, `sysfs`, `tmpfs`, `devtmpfs`, `cgroup`, `cgroup2`, `debugfs`, `tracefs`, `ramfs`, `autofs` и ряд служебных типов.

## 7. Network

Источник данных: `/proc/net/dev`.

Сетевые скорости рассчитываются по разнице счетчиков за интервал `COLLECTION_INTERVAL`.

| Поле | Тип | Единица | Описание |
| --- | --- | --- | --- |
| `rx_bytes_total` | integer | bytes | Суммарный текущий счетчик принятых байт по интерфейсам. |
| `tx_bytes_total` | integer | bytes | Суммарный текущий счетчик отправленных байт по интерфейсам. |
| `rx_bps_total` | number | bytes/s | Суммарная скорость приема. |
| `tx_bps_total` | number | bytes/s | Суммарная скорость передачи. |
| `ifaces` | array | - | Детализация по сетевым интерфейсам. |

Структура элемента `ifaces`:

| Поле | Тип | Единица | Описание |
| --- | --- | --- | --- |
| `name` | string | - | Имя интерфейса. |
| `rx_bytes` | integer | bytes | Текущий счетчик принятых байт. |
| `tx_bytes` | integer | bytes | Текущий счетчик отправленных байт. |
| `rx_bps` | number | bytes/s | Скорость приема по интерфейсу. |
| `tx_bps` | number | bytes/s | Скорость передачи по интерфейсу. |

Если интерфейс появился только во втором снимке или счетчик уменьшился, он пропускается в расчете текущей скорости.

## 8. Хранение на сервере

При успешном приеме сервер:

1. создает или обновляет запись в `devices` по `host_id`;
2. сохраняет снимок в `metric_snapshots`.

В `metric_snapshots` отдельно сохраняются:

- `collected_at`;
- агрегированные поля CPU, memory и network;
- JSONB-поля `cpu`, `memory`, `disk`, `network`;
- полный исходный снимок `raw`.

Такой подход позволяет одновременно выполнять простые аналитические выборки и сохранять полный контракт данных без потери детализации.

## 9. Совместимость и расширение

При расширении контракта необходимо соблюдать следующие правила:

- новые поля должны добавляться без удаления существующих;
- сервер должен сохранять полный `raw` JSON для обратной совместимости;
- клиенты должны устойчиво обрабатывать отсутствие новых полей;
- изменение семантики существующего поля требует отдельной версии API или миграционного периода.
