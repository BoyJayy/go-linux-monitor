package network

import (
	"errors"
	"monitoring/api"
	"time"
)

func Collect(interval time.Duration) (api.NetworkUsage, error) {
	if interval <= 0 {
		return api.NetworkUsage{}, errors.New("interval must be positive")
	}
	prevIfaces, err := ReadNetworkInterfaces()
	if err != nil {
		return api.NetworkUsage{}, err
	}
	time.Sleep(interval)
	curIfaces, err := ReadNetworkInterfaces()
	if err != nil {
		return api.NetworkUsage{}, err
	}
	prevByName := make(map[string]api.NetIfaceStat, len(prevIfaces))
	for _, iface := range prevIfaces {
		prevByName[iface.Name] = iface
	}
	seconds := interval.Seconds()
	newIfaces := make([]api.NetIfaceStat, 0, len(curIfaces))
	var rxTotal, txTotal uint64
	var rxBpsTotal, txBpsTotal float64

	for _, curIface := range curIfaces {
		prevIface, ok := prevByName[curIface.Name]
		if !ok {
			continue
		}
		if curIface.RxBytes < prevIface.RxBytes || curIface.TxBytes < prevIface.TxBytes {
			continue
		}
		rxDelta := curIface.RxBytes - prevIface.RxBytes
		txDelta := curIface.TxBytes - prevIface.TxBytes
		rxBps := float64(rxDelta) / seconds //bytes per sec
		txBps := float64(txDelta) / seconds
		newIfaces = append(newIfaces, api.NetIfaceStat{
			Name:    curIface.Name,
			RxBytes: curIface.RxBytes,
			TxBytes: curIface.TxBytes,
			RxBps:   rxBps,
			TxBps:   txBps,
		})
		rxTotal += curIface.RxBytes
		txTotal += curIface.TxBytes
		rxBpsTotal += rxBps
		txBpsTotal += txBps
	}
	return api.NetworkUsage{
		RxBytesTotal: rxTotal,
		TxBytesTotal: txTotal,
		RxBpsTotal:   rxBpsTotal,
		TxBpsTotal:   txBpsTotal,
		Ifaces:       newIfaces,
	}, nil
}
