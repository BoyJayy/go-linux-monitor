package snapshot

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func ReadCPUSnapshot() (total CpuSnapshot, cores map[string]CpuSnapshot, err error) {
	// здесь парс филдов
	cores = make(map[string]CpuSnapshot)
	proc, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CpuSnapshot{}, nil, err
	}
	strs := strings.Split(string(proc), "\n")
	for i := 0; i < len(strs); i++ {
		fields := strings.Fields(strs[i])
		if len(fields) == 0 {
			continue
		}
		name := fields[0]
		if !strings.HasPrefix(name, "cpu") {
			continue
		}
		if len(fields) < 9 {
			continue
		}
		user, err := strconv.ParseUint(fields[1], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		nice, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		system, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		idle, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		iowait, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		irq, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		softirq, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		steal, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return CpuSnapshot{}, nil, err
		}
		tot := user + nice + system + idle + iowait + irq + softirq + steal
		snap := CpuSnapshot{
			User:    user,
			Nice:    nice,
			System:  system,
			Idle:    idle,
			Iowait:  iowait,
			Irq:     irq,
			Softirq: softirq,
			Steal:   steal,
			Total:   tot,
		}
		if name == "cpu" {
			total = snap
			continue
		}
		cores[name] = snap
	}
	if total.Total == 0 {
		return CpuSnapshot{}, nil, errors.New("cpu total snapshot not found")
	}
	return total, cores, nil
}
