package handler

import (
	//"fmt"
	"net/http"
)

func (h *HTTPHandlers) HandleHeath(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte("ok"))
}
