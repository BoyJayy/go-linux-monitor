package handler

import (
	"encoding/json"
	"log"
	"monitoring/api"
	"net/http"
	"strconv"
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

func (h *HTTPHandlers) HandleGetLastMetricsByHostID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	host_id := r.URL.Query().Get("host_id")
	if host_id == "" {
		http.Error(w, "Host_ID in query is empty", http.StatusBadRequest)
		return
	}
	metrics, err := h.storage.GetLatestMetricsByHostID(r.Context(), host_id)
	if err != nil {
		log.Printf("failed to get latest metrics: %v", err)
		http.Error(w, "failed to get latest metrics: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metrics); err != nil {
		log.Printf("failed to encode latest metrics: %v", err)
	}
}

func (h *HTTPHandlers) HandleGetMetricsHistoryByHostID(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	hostID := r.URL.Query().Get("host_id")
	if hostID == "" {
		http.Error(w, "host_id is required", http.StatusBadRequest)
		return
	}
	limit := 100
	limitRaw := r.URL.Query().Get("limit")
	if limitRaw != "" {
		parsedLimit, err := strconv.Atoi(limitRaw)
		if err != nil {
			http.Error(w, "limit must be a number", http.StatusBadRequest)
			return
		}
		limit = parsedLimit
	}
	if limit <= 0 {
		http.Error(w, "limit must be positive", http.StatusBadRequest)
		return
	}
	if limit > 1000 {
		limit = 1000
	}
	metricsHistory, err := h.storage.GetMetricsHistoryByHostID(r.Context(), hostID, limit)
	if err != nil {
		log.Printf("failed to get metrics history: %v", err)
		http.Error(w, "failed to get metrics history: "+err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(metricsHistory); err != nil {
		log.Printf("failed to encode metrics history: %v", err)
	}
}
