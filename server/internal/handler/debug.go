package handler

import (
	"encoding/json"
	"log"

	//"log"
	//"monitoring/api"
	"net/http"
)

func (h *HTTPHandlers) HandleDebugLast(w http.ResponseWriter, r *http.Request) {
	lastMetrics, err := h.storage.GetLastMetrics(r.Context())
	if err != nil {
		log.Printf("debug failed to get last metrics: %v", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	b, err := json.MarshalIndent(lastMetrics, "", "    ")
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
