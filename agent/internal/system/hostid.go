package system

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/google/uuid"
)

const (
	machineIDPath = "/etc/machine-id"
	fallbackPath  = "/var/lib/myagent/host_id"
	hostIDDir     = "/var/lib/myagent"
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
	if id, err := ReadTrimmedFromPath(machineIDPath); err == nil {
		return id, nil
	}
	if id, err := ReadTrimmedFromPath(fallbackPath); err == nil {
		return id, nil
	}
	if err := os.MkdirAll(hostIDDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create host id dir: %w", err)
	}
	id := uuid.NewString()
	if err := os.WriteFile(fallbackPath, []byte(id), 0644); err != nil {
		return "", fmt.Errorf("failed to write fallback host id: %w", err)
	}
	return id, nil
}
