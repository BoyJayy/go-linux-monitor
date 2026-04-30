package storage

import "monitoring/api"

type Storage struct {
	metrics []api.Metrics
}

func NewStorage() *Storage {
	return &Storage{
		metrics: make([]api.Metrics, 0),
	}
}


func (st *Storage) GelLen() int {
	return len(st.metrics)
}
func (st *Storage) Save(m api.Metrics) {
	st.metrics = append(st.metrics, m)
}

func (st *Storage) GetLastMetrics() api.Metrics {
	if len(st.metrics) == 0 {
		return api.Metrics{}
	}
	return st.metrics[len(st.metrics)-1]
}
