package handler

import (
	"server/internal/storage"
	"time"
)

type HTTPHandlers struct {
	storage *storage.Storage
}

type ErrorDTO struct {
	message string
	time    time.Time
}
