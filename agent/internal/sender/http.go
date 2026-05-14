package sender

import (
	"agent/internal/snapshot"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// my sender
// soon gonna rebuild for grpc but for now its just http
// already did json tags so we just need to implement sender

func NewHTTP(endpoint string, timeout time.Duration) *HTTPSender {
	return &HTTPSender{
		endpoint:    endpoint,
		client:      &http.Client{Timeout: timeout},
		maxAttempts: 3,
		baseBackoff: 500 * time.Millisecond,
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
	req.Header.Set("Accept", "application/json")
	resp, err := s.client.Do(req)
	if err != nil {
		return fmt.Errorf("send request: %w", err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(respBody))
	}
	return nil
}

func (s *HTTPSender) SendWithRetry(metrics snapshot.Metrics) error {
	var lastErr error
	backoff := s.baseBackoff

	for attempt := 1; attempt <= s.maxAttempts; attempt++ {
		err := s.Send(metrics)
		if err == nil {
			if attempt > 1 {
				fmt.Printf("metrics sent after retry: host_id=%s attempt=%d\n", metrics.HostID, attempt)
			}
			return nil
		}

		lastErr = err

		fmt.Printf(
			"send metrics failed: host_id=%s attempt=%d/%d error=%v\n",
			metrics.HostID,
			attempt,
			s.maxAttempts,
			err,
		)

		if attempt == s.maxAttempts {
			break
		}

		time.Sleep(backoff)
		backoff *= 2
	}

	return fmt.Errorf("send metrics failed after %d attempts: %w", s.maxAttempts, lastErr)
}
