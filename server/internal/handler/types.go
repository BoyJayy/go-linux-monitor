package handler

import (
	"monitoring/api"
	"time"
)

type HTTPHandlers struct {
	metrics *api.Metrics
}

type ErrorDTO struct {
	message string
	time    time.Time
}
