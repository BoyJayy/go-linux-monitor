package handler

import "monitoring/api"

func NewHTTPHandlers(m *api.Metrics) *HTTPHandlers {
	return &HTTPHandlers{
		metrics: m,
	}
}
