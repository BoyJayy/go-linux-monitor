package cpu

import (
	"errors"
	"os"
	"strconv"
	"strings"
)

func ReadCPUCounters() (total CPUCounters, cores map[string]CPUCounters, err error) {
	// здесь парс филдов
	cores = make(map[string]CPUCounters)
	proc, err := os.ReadFile("/proc/stat")
	if err != nil {
		return CPUCounters{}, nil, err
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
			return CPUCounters{}, nil, err
		}
		nice, err := strconv.ParseUint(fields[2], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		system, err := strconv.ParseUint(fields[3], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		idle, err := strconv.ParseUint(fields[4], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		iowait, err := strconv.ParseUint(fields[5], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		irq, err := strconv.ParseUint(fields[6], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		softirq, err := strconv.ParseUint(fields[7], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		steal, err := strconv.ParseUint(fields[8], 10, 64)
		if err != nil {
			return CPUCounters{}, nil, err
		}
		tot := user + nice + system + idle + iowait + irq + softirq + steal
		snap := CPUCounters{
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
		return CPUCounters{}, nil, errors.New("cpu total snapshot not found")
	}
	return total, cores, nil
}
