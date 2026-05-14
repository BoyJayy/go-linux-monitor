package sender

import (
	"net/http"
	"time"
)

type HTTPSender struct {
	endpoint    string
	client      *http.Client
	maxAttempts int
	baseBackoff time.Duration
}
