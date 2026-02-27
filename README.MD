
# Linux Monitoring Service (V1)

A lightweight cross-platform monitoring system for Linux single-board computers.

## Features (V1)

- CPU usage and load average
- Memory statistics
- Disk usage
- Network traffic rate
- System uptime
- CPU temperature
- Process count
- Configurable sampling interval
- HTTP-based metric ingestion

## Architecture

Agent → HTTP → Server

- Agent collects metrics from /proc and /sys
- Server accepts and processes metrics

See:
- docs/metrics-contract.md
- docs/api.md
- docs/architecture.md

## Project Structure

monitoring/
  agent/
  server/
  docs/

## Roadmap

V1:
- Basic metric collection
- HTTP ingestion
- Logging

V2:
- PostgreSQL storage
- Device registry
- Historical metrics
- TUI client (btop-style)
