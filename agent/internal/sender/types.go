package sender

import (
	"net/http"
)

type HTTPSender struct {
	endpoint string
	client   *http.Client
}
