package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

func (h *HTTPHandlers) HandleListDevices(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	list, err := h.storage.GetListDevices(r.Context())
	if err != nil {
		log.Printf("Error receiving list of devices: %v", err.Error())
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	js, err := json.MarshalIndent(list, "", "    ")
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	w.Header().Set("Content-Type", "application-json")
	w.WriteHeader(http.StatusOK)
	w.Write(js)

}

