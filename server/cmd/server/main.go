package main

import (
	"context"
	"log"
	"os"
	"server/internal/handler"
	"server/internal/server"
	"server/internal/storage"

	//"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	databaseURL := os.Getenv("DATABASE_URL")
	if databaseURL == "" {
		log.Fatal("DATABASE_URL is required")
	}
	conn, err := pgxpool.New(ctx, databaseURL)
	if err != nil {
		log.Fatal("failed to create postgres pool: ", err)
	}
	defer conn.Close()
	if err := conn.Ping(ctx); err != nil {
		log.Fatal("failed to ping postgres: ", err)
	}
	//fmt.Println("successfully started db")
	store := storage.New(conn)
	httpHandlers := handler.NewHTTPHandlers(store)
	httpServer := server.NewHTTPServer(httpHandlers)
	httpServer.StartServer()
}
