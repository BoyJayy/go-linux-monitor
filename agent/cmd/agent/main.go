package main

import (
	"agent/internal/config"
	"agent/internal/sender"
	"agent/internal/snapshot"
	"log"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	sender := sender.NewHTTP(cfg.ServerURL, cfg.RequestTimeout)
	for {
		metrics, err := snapshot.BuildSnapshot(cfg.CollectionInterval)
		if err != nil {
			log.Printf("Error while reading metrics: %v", err)
			continue
		}
		//fmt.Printf("%+v\n", metrics)
		err = sender.Send(metrics)
		if err != nil {
			log.Printf("Error while sending metrics: %v", err)
			continue
		}
	}
}
