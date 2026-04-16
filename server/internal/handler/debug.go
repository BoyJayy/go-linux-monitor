package handler

import (
	"encoding/json"
	"log"

	//"log"
	//"monitoring/api"
	"net/http"
)

func (h *HTTPHandlers) HandleDebugLast(w http.ResponseWriter, r *http.Request) {
	/*var metricsDTO api.Metrics
	err := json.NewDecoder(r.Body).Decode(&metricsDTO)
	if err != nil {
		log.Printf("debug failed to decode")
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}*/
	last_metrics := h.storage.GetLastMetrics()
	b, err := json.MarshalIndent(last_metrics, "", "    ")
	if err != nil {
		log.Printf("debug failed to marshall")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	if _, err := w.Write(b); err != nil {
		log.Printf("debug failed to write http response")
		return
	}
}
