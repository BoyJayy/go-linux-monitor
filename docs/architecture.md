
# System Architecture (V1)

## Overview

The system consists of two main components:

1. Agent (Go binary)
   - Runs on Linux SBC (Orange Pi, Raspberry Pi, etc.)
   - Collects system metrics from /proc and /sys
   - Sends metrics to the backend server via HTTP

2. Server (Go binary)
   - Runs on MacBook (development) or Linux server (production)
   - Accepts metrics from agents
   - Validates and logs metrics (V1)
   - Will store metrics in PostgreSQL (V2)

---

## Data Flow

1. Agent samples system metrics at configurable interval.
2. Agent calculates derived values (CPU%, RX/TX rate).
3. Agent sends snapshot to:
   POST /api/v1/metrics
4. Server validates and processes the snapshot.

---

## Responsibilities

Agent:
- Metric collection
- Sampling logic
- Rate calculation
- Retry & backoff
- Timeout handling

Server:
- HTTP API
- Validation
- Logging (V1)
- Storage (future)

---

## Deployment (V1)

Agent:
- Single static binary
- Runs directly on device
- Later managed via systemd

Server:
- Runs locally on development machine
- Later deployable to remote Linux server
