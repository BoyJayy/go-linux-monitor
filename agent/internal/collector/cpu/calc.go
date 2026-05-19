package cpu

import (
	"errors"
	"monitoring/api"
)

func CPUCalculation(prev CPUCounters, cur CPUCounters) (api.CPUUsage, error) {
	curIdle, prevIdle := cur.Idle+cur.Iowait, prev.Idle+prev.Iowait
	if cur.Total < prev.Total || curIdle < prevIdle {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.User < prev.User {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.Nice < prev.Nice {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.System < prev.System {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.Idle < prev.Idle {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.Iowait < prev.Iowait {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.Irq < prev.Irq {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.Softirq < prev.Softirq {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if cur.Steal < prev.Steal {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	deltaTotal, deltaIdle := cur.Total-prev.Total, curIdle-prevIdle
	var busy float64 = float64(deltaTotal) - float64(deltaIdle)
	if busy < 0 {
		return api.CPUUsage{}, errors.New("Incorrect snapshot appeared")
	}
	if deltaTotal == 0 {
		return api.CPUUsage{}, nil
	}
	return api.CPUUsage{
		UserPct:    float64(cur.User-prev.User) / float64(deltaTotal) * 100,
		NicePct:    float64(cur.Nice-prev.Nice) / float64(deltaTotal) * 100,
		SystemPct:  float64(cur.System-prev.System) / float64(deltaTotal) * 100,
		IdlePct:    float64(cur.Idle-prev.Idle) / float64(deltaTotal) * 100,
		IowaitPct:  float64(cur.Iowait-prev.Iowait) / float64(deltaTotal) * 100,
		IrqPct:     float64(cur.Irq-prev.Irq) / float64(deltaTotal) * 100,
		SoftirqPct: float64(cur.Softirq-prev.Softirq) / float64(deltaTotal) * 100,
		StealPct:   float64(cur.Steal-prev.Steal) / float64(deltaTotal) * 100,
		TotalPct:   busy / float64(deltaTotal) * 100,
	}, nil
}

func CPUPerCoreCalculation(prevCores map[string]CPUCounters, curCores map[string]CPUCounters) (map[string]api.CPUUsage, error) {
	if len(prevCores) != len(curCores) {
		return nil, errors.New("Incorrect percore CPU snapshot appeared")
	}
	ans := make(map[string]api.CPUUsage)
	for k, prevSnap := range prevCores {
		curSnap, ok := curCores[k]
		if !ok {
			return nil, errors.New(k + "core does not exist in current snapshot")
		}
		val, err := CPUCalculation(prevSnap, curSnap)
		if err != nil {
			return nil, err
		}
		ans[k] = val
	}
	return ans, nil
}
