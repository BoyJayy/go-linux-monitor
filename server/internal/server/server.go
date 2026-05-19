package server

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"server/internal/handler"
	"syscall"
	"time"
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
	router.HandleFunc("/health", s.httpHandlers.HandleHeath).Methods("GET")
	router.HandleFunc("/api/v1/devices", s.httpHandlers.HandleListDevices).Methods("GET")
	router.HandleFunc("/api/v1/devices/latest", s.httpHandlers.HandleGetLastMetricsByHostID).Methods("GET")
	router.HandleFunc("/api/v1/devices/metrics", s.httpHandlers.HandleGetMetricsHistoryByHostID).Methods("GET")
	router.HandleFunc("/debug/api", s.httpHandlers.HandleDebugLast).Methods("GET")
	staticDir := "./web"
	router.PathPrefix("/").Handler(http.FileServer(http.Dir(staticDir)))
	httpServer := &http.Server{
		Addr:    ":8080",
		Handler: router,
	}
	serverErr := make(chan error, 1)
	go func() {
		log.Println("server started on :8080")
		if err := httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			serverErr <- err
			return
		}
		serverErr <- nil
	}()
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	select {
	case <-ctx.Done():
		log.Println("shutdown signal received")
	case err := <-serverErr:
		if err != nil {
			log.Fatalf("server start error: %v", err)
		}
		return
	}
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(shutdownCtx); err != nil {
		log.Fatalf("server shutdown error: %v", err)
	}
	log.Println("server stopped")
}
