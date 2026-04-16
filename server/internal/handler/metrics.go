package handler

import (
	"encoding/json"
	"log"
	"monitoring/api"
	"net/http"
	"time"
	//"server/internal/storage"
)

func (h *HTTPHandlers) HandlePostV1Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return;
	}
	defer r.Body.Close()
	log.Println("HandleMetricsPost called")
	var metricsDTO api.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricsDTO); err != nil {
		errDTO := ErrorDTO{
			message: err.Error(),
			time:    time.Now(),
		}
		http.Error(w, errDTO.toString(), http.StatusBadRequest)
		return
	}
	h.storage.Save(metricsDTO);
	log.Printf("v1 receiver method: %+v\n", time.Now())
	log.Printf("received metrics: %+v\n", metricsDTO)
	w.WriteHeader(http.StatusAccepted)
	log.Printf("metrics saved, total=%d\n", h.storage.GelLen())
}
