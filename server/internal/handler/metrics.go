package handler

import (
	"encoding/json"
	"monitoring/api"
	"net/http"
	"time"
)

func (h *HTTPHandlers) HandlePostV1Metrics(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
	defer r.Body.Close()
	var metricsDTO api.Metrics
	if err := json.NewDecoder(r.Body).Decode(&metricsDTO); err != nil {
		errDTO := ErrorDTO{
			message: err.Error(),
			time:    time.Now(),
		}
		http.Error(w, errDTO.toString(), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusAccepted)
}
