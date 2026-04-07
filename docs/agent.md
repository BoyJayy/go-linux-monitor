# Monitoring System API

> Current API contract for the monitoring project.
> This document reflects the project state where:
> - the **agent** already builds unified snapshots and sends them over HTTP JSON
> - the **server** is being implemented with a first ingest endpoint

## 1. Purpose

The API is used to transfer a system snapshot from the agent to the backend.

Current transport:
- HTTP
- JSON
- request timeout configured on the agent side

## 2. Base idea

Flow:

`agent -> HTTP POST -> server`

The agent periodically:
1. builds a `Metrics` snapshot
2. serializes it to JSON
3. sends it to the backend endpoint

## 3. Endpoints

### 3.1 Health check

**Method:** `GET`  
**Path:** `/health`

#### Purpose
Used to verify that the backend process is alive.

#### Success response
- Status: `200 OK`

#### Example response
```json
{
  "status": "ok"
}
```

---

### 3.2 Metrics ingest

**Method:** `POST`  
**Path:** `/api/v1/metrics`

#### Purpose
Accepts one snapshot from the monitoring agent.

#### Request headers
```http
Content-Type: application/json
```

#### Request body
JSON object matching the `Metrics` structure.

## 4. Metrics payload

Current payload shape:

```json
{
  "timestamp": "2026-03-24T12:35:24.943217753Z",
  "host_id": "c5cd5bea-81fb-4f81-879f-cc5ca1a61e50",
  "cpu": {
    "user_pct": 0.0,
    "nice_pct": 0.0,
    "system_pct": 0.0623,
    "idle_pct": 99.8753,
    "iowait_pct": 0.0,
    "irq_pct": 0.0,
    "softirq_pct": 0.0623,
    "steal_pct": 0.0,
    "total_pct": 0.1246,
    "per_core_pct": {
      "cpu0": {
        "user_pct": 0.0,
        "nice_pct": 0.0,
        "system_pct": 0.0,
        "idle_pct": 100.0,
        "iowait_pct": 0.0,
        "irq_pct": 0.0,
        "softirq_pct": 0.0,
        "steal_pct": 0.0,
        "total_pct": 0.0
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

## 5. Field semantics

### Top-level fields

- `timestamp` — time when the snapshot was created by the agent
- `host_id` — stable host identifier resolved by the agent

### CPU

- `total_pct` — total CPU usage percent
- breakdown fields:
  - `user_pct`
  - `nice_pct`
  - `system_pct`
  - `idle_pct`
  - `iowait_pct`
  - `irq_pct`
  - `softirq_pct`
  - `steal_pct`
- `per_core_pct` — usage per CPU core

### Memory

- `total_bytes`
- `available_bytes`
- `used_bytes`
- `used_pct`

### Disk

Disk stats are collected **per mounted filesystem**, not per physical disk device.

Each entry contains:
- `mount`
- `total_bytes`
- `free_bytes`
- `used_bytes`
- `used_pct`

### Network

- `rx_bytes_total`
- `tx_bytes_total`
- `rx_bps_total`
- `tx_bps_total`
- `ifaces` — per-interface counters and rates

## 6. Validation rules for ingest

Current minimum validation for accepted snapshots:

- `host_id` must be non-empty
- `timestamp` must be present and non-zero

Further validation may be extended later.

## 7. Response codes

### `GET /health`
- `200 OK` — server is alive

### `POST /api/v1/metrics`
- `200 OK` — metrics accepted
- `400 Bad Request` — invalid JSON or invalid payload
- `405 Method Not Allowed` — wrong HTTP method
- `500 Internal Server Error` — internal processing error

## 8. Example curl requests

### Health
```bash
curl http://localhost:8080/health
```

### Metrics ingest
```bash
curl -X POST http://localhost:8080/api/v1/metrics \
  -H "Content-Type: application/json" \
  -d @metrics.json
```

## 9. Agent configuration related to the API

Current important runtime variables:

- `SERVER_URL` — full metrics ingest endpoint  
  Example: `http://localhost:8080/api/v1/metrics`
- `COLLECTION_INTERVAL` — snapshot collection interval  
  Example: `2s`
- `REQUEST_TIMEOUT` — maximum HTTP send timeout  
  Example: `3s`

## 10. Notes

- The sender currently uses synchronous HTTP sending.
- Retry policy is planned as a later improvement.
- Persistent storage is planned for the backend after the first in-memory ingest stage.
