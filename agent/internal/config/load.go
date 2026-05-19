package config

import (
	"errors"
	"fmt"
	"os"
	"time"
)

func Load() (Config, error) {
	serverURL, ok := os.LookupEnv("SERVER_URL")
	if !ok || serverURL == "" {
		return Config{}, errors.New("SERVER_URL is required")
	}
	intervalStr, ok := os.LookupEnv("COLLECTION_INTERVAL")
	if !ok || intervalStr == "" {
		intervalStr = "2s"
	}
	timeoutStr, ok := os.LookupEnv("REQUEST_TIMEOUT")
	if !ok || timeoutStr == "" {
		timeoutStr = "3s"
	}
	collectionInterval, err := time.ParseDuration(intervalStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid COLLECTION_INTERVAL: %w", err)
	}
	requestTimeout, err := time.ParseDuration(timeoutStr)
	if err != nil {
		return Config{}, fmt.Errorf("invalid REQUEST_TIMEOUT: %w", err)
	}
	return Config{
		ServerURL:          serverURL,
		CollectionInterval: collectionInterval,
		RequestTimeout:     requestTimeout,
	}, nil
}
