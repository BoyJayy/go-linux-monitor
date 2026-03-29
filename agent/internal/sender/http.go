package sender

import (
	"agent/internal/snapshot"
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// my sender
// soon gonna rebuild for grpc but for now its just http
// already did json tags so we just need to implement sender

func NewHTTP(endpoint string, timeout time.Duration) *HTTPSender {
	return &HTTPSender{
		endpoint: endpoint,
		client:   &http.Client{Timeout: timeout},
	}

}

func (s *HTTPSender) Send(metrics snapshot.Metrics) error {
	body, err := json.Marshal(metrics)
	if err != nil {
		return fmt.Errorf("sender json packaging: %w", err)
	}
	req, err := http.NewRequest(http.MethodPost, s.endpoint, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("http request creation: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	return nil
}
