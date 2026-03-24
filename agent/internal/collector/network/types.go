package network

type NetIfaceStat struct {
	Name    string
	RxBytes uint64
	TxBytes uint64
	RxBps   float64
	TxBps   float64
}

type NetworkUsage struct {
	RxBytesTotal uint64
	TxBytesTotal uint64
	RxBpsTotal   float64
	TxBpsTotal   float64
	Ifaces       []NetIfaceStat
}
