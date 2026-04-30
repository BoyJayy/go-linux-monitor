package main

import (
	"context"
	"log"
	"server/internal/handler"
	"server/internal/server"

	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	conn, err := pgxpool.New(ctx, "postgres://postgres:537877@localhost:5432/monitoring service")
	if err != nil {
		log.Fatal("failed to create postgres pool: ", err)
	}
	defer conn.Close()
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("failed to ping postgres: ", err)
	}
	//fmt.Println("successfully started db")
	httpHandlers := handler.NewHTTPHandlers()
	httpServer := server.NewHTTPServer(httpHandlers)
	httpServer.StartServer()
}
