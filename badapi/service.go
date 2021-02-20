package badapi

import (
	"net/http"

	"golang.org/x/xerrors"
)

// Service represents the badapi
type Service struct {
	c       *http.Client
	baseURL string
}

type ServerResponse struct {
	Code   int
	Header http.Header
}

// NewService initiates a new badapi service
func NewService() *Service {
	return &Service{
		c:       &http.Client{},
		baseURL: "https://bad-api-assignment.reaktor.com/v2/",
	}
}

// Do executes a http request and returns the response or an error. Exactly one
// return value will be non-nil.
func (s *Service) Do(req *http.Request) (*http.Response, error) {
	res, err := s.c.Do(req)
	if err != nil {
		return nil, err
	}
	if res.Header.Get("X-Error-Modes-Active") != "" {
		return nil, xerrors.New("badapi: failed to write data to response")
	}
	return res, nil
}
