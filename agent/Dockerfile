FROM golang:1.25.5 AS builder

WORKDIR /src

COPY api ./api
COPY agent ./agent
RUN go work init ./api ./agent

WORKDIR /src/agent

RUN go mod download
RUN go build -o /bin/monitoring-agent ./cmd/agent
FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /bin/monitoring-agent /app/monitoring-agent

CMD ["/app/monitoring-agent"]