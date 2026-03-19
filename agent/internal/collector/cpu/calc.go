package cpu

import (
	"errors"
)

func CPUCalculation(prev CPUCounters, cur CPUCounters) (float64, error) {
	curIdle, prevIdle := cur.Idle+cur.Iowait, prev.Idle+prev.Iowait
	if cur.Total < prev.Total || curIdle < prevIdle {
		return 0.0, errors.New("Incorrect snapshot appeared")
	}
	deltaTotal, deltaIdle := cur.Total-prev.Total, curIdle-prevIdle
	var busy float64 = float64(deltaTotal) - float64(deltaIdle)
	if busy < 0 {
		return 0.0, errors.New("Incorrect snapshot appeared")
	}
	if deltaTotal == 0 {
		return 0.0, nil
	}
	return busy / float64(deltaTotal) * 100, nil
}

func CPUPerCoreCalculation(prevCores map[string]CPUCounters, curCores map[string]CPUCounters) (map[string]float64, error) {
	if len(prevCores) != len(curCores) {
		return nil, errors.New("Incorrect percore CPU snapshot appeared")
	}
	ans := make(map[string]float64)
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
