package server

import (
	"server/internal/handler"
)

type HTTPServer struct {
	httpHandlers *handler.HTTPHandlers
}
