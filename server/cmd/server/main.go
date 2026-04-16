package main

import (
	"server/internal/handler"
	"server/internal/server"
)

func main() {
	httpHandlers := handler.NewHTTPHandlers()
	httpServer := server.NewHTTPServer(httpHandlers)
	httpServer.StartServer()
}
