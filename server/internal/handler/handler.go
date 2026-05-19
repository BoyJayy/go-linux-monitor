package handler

import (
	//"monitoring/api"
	"server/internal/storage"
)

func NewHTTPHandlers(storage *storage.Storage) *HTTPHandlers {
	return &HTTPHandlers{
		storage: storage,
	}
}
