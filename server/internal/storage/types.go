package storage

import (
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
)

type Device struct {
	ID          int64     `json:"id"`
	HostID      string    `json:"host_id"`
	FirstSeenAt time.Time `json:"first_seen_at"`
	LastSeenAt  time.Time `json:"last_seen_at"`
	Online      bool      `json:"online"`
}

type Storage struct {
	pool *pgxpool.Pool
}
