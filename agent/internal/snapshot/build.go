package snapshot

import (
	"agent/internal/collector/cpu"
	"agent/internal/collector/disk"
	"agent/internal/collector/memory"
	"agent/internal/collector/network"
	"agent/internal/system"
	"fmt"
	"time"
)

func BuildSnapshot(interval time.Duration) (Metrics, error) {
	cpuCh := make(chan CPUResult, 1)
	memCh := make(chan MemResult, 1)
	netCh := make(chan NetworkResult, 1)
	diskCh := make(chan DiskResult, 1)
	timeCh := make(chan time.Time, 1)
	/*var cpuclc CPUResult
	var memclc MemResult
	var diskclc DiskResult
	var netclc NetworkResult*/
	go func() {
		stat, err := cpu.Collect(interval)
		cpuCh <- CPUResult{stat: stat, err: err}
	}()
	go func() {
		stat, err := memory.Collect()
		memCh <- MemResult{stat: stat, err: err}
	}()
	go func() {
		stat, err := network.Collect(interval)
		netCh <- NetworkResult{stat: stat, err: err}
	}()
	go func() {
		stat, err := disk.Collect()
		diskCh <- DiskResult{stat: stat, err: err}
	}()
	go func() {
		timeCh <- system.CollectTimestamp(interval)
	}()
	cpuclc := <-cpuCh
	memclc := <-memCh
	netclc := <-netCh
	diskclc := <-diskCh
	timeclc := <-timeCh
	hostid, err := system.ResolveHostId()

	if err != nil {
		return Metrics{}, fmt.Errorf("host ID collect: %w", err)
	}
	if cpuclc.err != nil {
		return Metrics{}, fmt.Errorf("cpu collect: %w", cpuclc.err)
	}
	if memclc.err != nil {
		return Metrics{}, fmt.Errorf("memory collect: %w", memclc.err)
	}
	m := Metrics{HostID: hostid, Timestamp: timeclc, CPU: cpuclc.stat, Mem: memclc.stat}
	if netclc.err == nil {
		m.Network = netclc.stat
	}
	if diskclc.err == nil {
		m.Disk = diskclc.stat
	}
	return m, nil
}
