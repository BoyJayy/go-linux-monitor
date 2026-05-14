package main

import (
	"agent/internal/config"
	"agent/internal/sender"
	"agent/internal/snapshot"
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	metricSender := sender.NewHTTP(cfg.ServerURL, cfg.RequestTimeout)
	log.Printf(
		"agent started: server_url=%s collection_interval=%s request_timeout=%s",
		cfg.ServerURL,
		cfg.CollectionInterval,
		cfg.RequestTimeout,
	)
	for {
		select {
		case <-ctx.Done():
			log.Println("shutdown signal received")
			log.Println("agent stopped")
			return
		default:
		}
		metrics, err := snapshot.BuildSnapshot(cfg.CollectionInterval)
		if err != nil {
			log.Printf("collect metrics failed: error=%v", err)
			continue
		}
		select {
		case <-ctx.Done():
			log.Println("shutdown signal received")
			log.Println("agent stopped")
			return
		default:
		}
		err = metricSender.SendWithRetry(metrics)
		if err != nil {
			log.Printf("drop metrics snapshot: host_id=%s error=%v", metrics.HostID, err)
			continue
		}
		log.Printf("metrics sent: host_id=%s", metrics.HostID)
	}
}
