package system

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/google/uuid"
)

const (
	envHostID     = "HOST_ID"
	envHostIDFile = "HOST_ID_FILE"

	machineIDPath       = "/etc/machine-id"
	defaultFallbackPath = "/var/lib/myagent/host_id"
)

func ReadTrimmedFromPath(path string) (string, error) {
	file, err := os.ReadFile(path)
	if err != nil {
		return "", err
	}
	value := strings.TrimSpace(string(file))
	if value == "" {
		return "", errors.New("file is empty " + path)
	}
	return value, nil
}

func ResolveHostId() (string, error) {
	if id := strings.TrimSpace(os.Getenv(envHostID)); id != "" {
		return id, nil
	}
	if id, err := ReadTrimmedFromPath(machineIDPath); err == nil {
		return id, nil
	}
	fallbackPath := strings.TrimSpace(os.Getenv(envHostIDFile))
	if fallbackPath == "" {
		fallbackPath = defaultFallbackPath
	}
	if id, err := ReadTrimmedFromPath(fallbackPath); err == nil {
		return id, nil
	}
	id := uuid.NewString()
	if err := os.MkdirAll(filepath.Dir(fallbackPath), 0755); err != nil {
		return "", fmt.Errorf("failed to create host id dir: %w", err)
	}
	if err := os.WriteFile(fallbackPath, []byte(id+"\n"), 0644); err != nil {
		return "", fmt.Errorf("failed to write fallback host id: %w", err)
	}
	return id, nil
}
