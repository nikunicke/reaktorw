package badapi

import (
	"encoding/json"
	"net/http"
	"strings"
)

// Availability represent a badapi availability
type Availability struct {
	Code     int32 `json:"code"`
	Response []struct {
		ID          string `json:"id"`
		DataPayload string `json:"DATAPAYLOAD"`
	} `json:"response"`

	ServerResponse `json:"-"`
}

// AvailabilitiesService handles badapi products
type AvailabilitiesService struct {
	s *Service
}

// Availabilities initiates a new ProductsService
func Availabilities(s *Service) *AvailabilitiesService {
	return &AvailabilitiesService{s: s}
}

// Get return a get caller
func (s *AvailabilitiesService) Get(manufacturer string) *AvailabilitiesGetCall {
	return &AvailabilitiesGetCall{s: s.s, manufacturer: strings.ToLower(manufacturer)}
}

// AvailabilitiesGetCall represents an availability get caller
type AvailabilitiesGetCall struct {
	s *Service

	manufacturer string
	header       http.Header
}

// Do executes a get call request
func (c *AvailabilitiesGetCall) Do() (*Availability, error) {
	urls := c.s.baseURL + "availability/" + c.manufacturer
	req, err := http.NewRequest(http.MethodGet, urls, nil)
	if err != nil {
		return nil, err
	}
	res, err := c.s.Do(req)
	if err != nil {
		return nil, err
	}
	if err := checkResponse(res); err != nil {
		return nil, err
	}
	ret := &Availability{
		ServerResponse: ServerResponse{
			Header:     res.Header,
			StatusCode: res.StatusCode,
		},
	}
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(target); err != nil {
		return nil, err
	}
	return ret, nil
}
