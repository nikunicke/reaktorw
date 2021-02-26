package badapi

import (
	"net/http"
	"time"
)

// Service represents the badapi
type Service struct {
	c       *http.Client
	baseURL string
}

// ServerResponse includes response data. Included in all direct responses.
type ServerResponse struct {
	StatusCode int
	Header     http.Header
}

// NewService initiates a new badapi service. BaseURL default to
// "https://bad-api-assignment.reaktor.com/v2/"
func NewService() *Service {
	return &Service{
		c:       &http.Client{Timeout: 30 * time.Second},
		baseURL: "https://bad-api-assignment.reaktor.com/v2/",
	}
}

// URL sets a baseURL a badapi service
func (s *Service) URL(urls string) *Service {
	s.baseURL = urls
	return s
}

// Do executes a http request and returns the response or an error. Exactly one
// return value will be non-nil.
func (s *Service) Do(req *http.Request) (*http.Response, error) {
	res, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Header.Get("X-Error-Modes-Active") != "" {
		return nil, ErrModeActive
	}
	return res, nil
}
