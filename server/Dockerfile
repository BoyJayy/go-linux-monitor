FROM golang:1.25.5 AS builder
WORKDIR /src

COPY api ./api
COPY server ./server
RUN go work init ./api ./server
WORKDIR /src/server

RUN go mod download
RUN go build -o /bin/monitoring-server ./cmd/server


FROM debian:bookworm-slim

WORKDIR /app

COPY --from=builder /bin/monitoring-server /app/monitoring-server
COPY server/web /app/web

EXPOSE 8080

CMD ["/app/monitoring-server"]