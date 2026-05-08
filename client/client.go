package client

import (
	"net/http"
	"time"
)

type uaTransport struct {
	userAgent string
	base      http.RoundTripper
}

func (tr *uaTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req2 := req.Clone(req.Context())
	req2.Header.Set("User-Agent", tr.userAgent)

	// delegate the actual transport to the base
	return tr.base.RoundTrip(req2)
}

// New Client that can be designate your custom User-Agent.
func NewUAClient(ua string) *http.Client {
	return &http.Client{
		Timeout: 60 * time.Second,
		Transport: &uaTransport{
			userAgent: ua,
			base:      http.DefaultTransport,
		},
	}
}
