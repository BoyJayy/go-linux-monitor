package main

import (
	"context"
	"server/internal/handler"
	"server/internal/server"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()
	conn, err := pgx.Connect(ctx, "postgres://postgres:537877@localhost:5432/monitoring service")
	if err != nil {
		panic(err)
	}
	defer conn.Close(ctx)
	//fmt.Println("successfully started db")
	httpHandlers := handler.NewHTTPHandlers()
	httpServer := server.NewHTTPServer(httpHandlers)
	httpServer.StartServer()
}
