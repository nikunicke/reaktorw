package badapi

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

// Availability represent a badapi availability
type Availability struct {
	Code     int32       `json:"code"`
	Response []*Response `json:"response"`

	ServerResponse `json:"-"`
}

type Response struct {
	ID          string `json:"id"`
	DataPayload string `json:"DATAPAYLOAD"`
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
	res, err := c.executeRequest(req)
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

func (c *AvailabilitiesGetCall) executeRequest(req *http.Request) (*http.Response, error) {
	attempts := 6
	resCh := make(chan *http.Response)
	errCh := make(chan error)
	ctx, cancelFn := context.WithCancel(context.Background())
	defer cancelFn()

	for i := 0; i < attempts; i++ {
		go func(context.Context) {
			res, err := c.s.Do(req)
			if err != nil {
				errCh <- err
				return
			}
			select {
			case <-ctx.Done():
			case resCh <- res:
			}
		}(ctx)
	}

	var allErr error
	for i := 0; i < attempts; i++ {
		select {
		case err := <-errCh:
			// allErr = multierror.Append(allErr, err)
			allErr = err
		case res := <-resCh:
			cancelFn()
			return res, nil
		}
	}
	fmt.Println(allErr)
	return nil, xerrors.New("All attempts to request availabilities failed")
}
