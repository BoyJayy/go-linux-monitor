package handler

import (
	"encoding/json"
	"log"
	"monitoring/api"
	"net/http"
	"time"
)

func (h *HTTPHandlers) HandlePostV1Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
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
	if metricsDTO.HostID == "" {
		http.Error(w, "host_id is required", http.StatusBadRequest)
		return
	}

	if metricsDTO.Timestamp.IsZero() {
		http.Error(w, "timestamp is required", http.StatusBadRequest)
		return
	}

	if err := h.storage.SaveMetrics(r.Context(), metricsDTO); err != nil {
		log.Printf("failed to save metrics: %v", err)
		http.Error(w, "failed to save metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte("ok"))
}
