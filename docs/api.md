
# Monitoring API (V1)

Base URL (local development):
http://192.168.0.111:8080

---

## 1. Health Check

### GET /health

Purpose:
Verify that the server is running and reachable.

Response:
200 OK
{
  "status": "ok"
}

---

## 2. Ingest Metrics

### POST /api/v1/metrics

Purpose:
Receive a metrics snapshot from an agent.

Headers:
Content-Type: application/json

Body:
MetricsSnapshot (see metrics-contract.md)

Response:
200 OK
{
  "status": "accepted"
}

Errors:
400 Bad Request – invalid JSON or missing required fields  
413 Payload Too Large – body exceeds limit  
415 Unsupported Media Type – wrong Content-Type  
500 Internal Server Error – unexpected server failure
