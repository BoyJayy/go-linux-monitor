package cpu

import (
	"monitoring/api"
	"time"
)

func Collect(interval time.Duration) (api.UsageWithCores, error) {
	prevTot, prevCores, err := ReadCPUCounters()
	if err != nil {
		return api.UsageWithCores{}, err
	}
	time.Sleep(interval)
	curTot, curCores, err := ReadCPUCounters()
	if err != nil {
		return api.UsageWithCores{}, err
	}
	totCalc, err := CPUCalculation(prevTot, curTot)
	if err != nil {
		return api.UsageWithCores{}, err
	}
	coresCalc, err := CPUPerCoreCalculation(prevCores, curCores)
	if err != nil {
		return api.UsageWithCores{}, err
	}
	return api.UsageWithCores{
		TotalPct:   totCalc.TotalPct,
		UserPct:    totCalc.UserPct,
		NicePct:    totCalc.NicePct,
		SystemPct:  totCalc.SystemPct,
		IdlePct:    totCalc.IdlePct,
		IowaitPct:  totCalc.IowaitPct,
		IrqPct:     totCalc.IrqPct,
		SoftirqPct: totCalc.SoftirqPct,
		StealPct:   totCalc.StealPct,
		PerCorePct: coresCalc,
	}, nil
}
