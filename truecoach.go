package truecoach

import (
	"net/http"
	"time"
)

const (
	defaultOrigin = "https://app.truecoach.co/proxy/api"
)

type Service struct {
	httpClient *http.Client
	origin     string
}

func NewService(token string) *Service {
	return &Service{
		httpClient: &http.Client{
			Timeout:   10 * time.Second,
			Transport: &transport{token},
		},
		origin: defaultOrigin,
	}
}

func (tc *Service) WithOrigin(origin string) *Service {
	tc.origin = origin
	return tc
}
