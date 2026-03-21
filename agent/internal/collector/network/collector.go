package network

import (
	"errors"
	"time"
)

func Collect(interval time.Duration) (NetworkStat, error) {
	if interval <= 0 {
		return NetworkStat{}, errors.New("interval must be positive")
	}
	prevIfaces, err := ReadNetworkInterfaces()
	if err != nil {
		return NetworkStat{}, err
	}
	time.Sleep(interval)
	curIfaces, err := ReadNetworkInterfaces()
	if err != nil {
		return NetworkStat{}, err
	}
	prevByName := make(map[string]NetIfaceStat, len(prevIfaces))
	for _, iface := range prevIfaces {
		prevByName[iface.Name] = iface
	}
	seconds := interval.Seconds()
	newIfaces := make([]NetIfaceStat, 0, len(curIfaces))
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
		newIfaces = append(newIfaces, NetIfaceStat{
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
	return NetworkStat{
		RxBytesTotal: rxTotal,
		TxBytesTotal: txTotal,
		RxBpsTotal:   rxBpsTotal,
		TxBpsTotal:   txBpsTotal,
		Ifaces:       newIfaces,
	}, nil
}
