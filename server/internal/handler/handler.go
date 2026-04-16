package handler

import (
	//"monitoring/api"
	"server/internal/storage"
)

func NewHTTPHandlers() *HTTPHandlers {
	return &HTTPHandlers{
		storage: storage.NewStorage(),
	}
}
