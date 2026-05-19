package memory

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func ReadMemCounters() (snap MemCounters, err error) {
	proc, err := os.ReadFile("/proc/meminfo")
	if err != nil {
		return MemCounters{}, err
	}
	strs := strings.Split(string(proc), "\n")
	var total, available uint64
	for i := 0; i < len(strs); i++ {
		//if strings.HasPrefix(strs[i])
		fields := strings.Fields(strs[i])
		if len(fields) < 2 {
			continue
		}
		if fields[0] == "MemTotal:" {
			total, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return MemCounters{}, err
			}
		}
		if fields[0] == "MemAvailable:" {
			available, err = strconv.ParseUint(fields[1], 10, 64)
			if err != nil {
				return MemCounters{}, err
			}
		}
		if total != 0 && available != 0 {
			break
		}
	}
	if total == 0 ||
		available == 0 {
		return MemCounters{}, errors.New("failed to read MemTotal or MemAvailable from /proc/meminfo")
	}
	if available > total {
		return MemCounters{}, errors.New("invalid value for memory (available > total)")
	}
	return MemCounters{
		Total:     total * 1024,
		Available: available * 1024,
		Used:      total*1024 - available*1024,
	}, nil
}
