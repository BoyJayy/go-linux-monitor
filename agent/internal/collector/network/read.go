package network

import (
	"errors"
	//"fmt"
	"os"
	"strconv"
	"strings"
)

func ReadNetworkInterfaces() ([]NetIfaceStat, error) {
	file, err := os.ReadFile("/proc/net/dev")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(string(file), "\n")
	if len(lines) < 3 {
		return nil, errors.New("invalid /proc/net/dev format")
	}
	var ifaces []NetIfaceStat
	for _, line := range lines[2:] {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		name := strings.TrimSpace(parts[0])
		fields := strings.Fields(parts[1])
		if len(fields) < 16 {
			continue
		}
		rxBytes, err := strconv.ParseUint(fields[0], 10, 64)
		if err != nil {
			return nil, err
		}

		txBytes, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return nil, err
		}
		ifaces = append(ifaces, NetIfaceStat{
			Name:    name,
			RxBytes: rxBytes,
			TxBytes: txBytes,
		})
	}

	return ifaces, nil
}
