package server

import (
	"log"
	"net/http"
	"server/internal/handler"

	"github.com/gorilla/mux"
)

func NewHTTPServer(h *handler.HTTPHandlers) *HTTPServer {
	return &HTTPServer{
		httpHandlers: h,
	}
}

func (s *HTTPServer) StartServer() {
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/metrics", s.httpHandlers.HandlePostV1Metrics).Methods("POST")
	router.HandleFunc("/debug/api", s.httpHandlers.HandleDebugLast).Methods("GET")
	if err := http.ListenAndServe(":8080", router); err != nil {
		log.Fatalf("Server start error: %v", err.Error())
	}
}
