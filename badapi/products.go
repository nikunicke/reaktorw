package badapi

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
)

// Product represent a badapi product
type Product struct {
	ID           string   `json:"id"`
	Type         string   `json:"type"`
	Name         string   `json:"name"`
	Color        []string `json:"color"`
	Price        int32    `json:"price"`
	Manufacturer string   `json:"manufacturer"`
}

// ProductsService handles badapi products
type ProductsService struct {
	s *Service
}

// Products initiates a new ProductsService
func Products(s *Service) *ProductsService {
	return &ProductsService{s: s}
}

// List returns a list caller that requests a page at a time
func (s *ProductsService) List(category string) *ProductsListCall {
	return &ProductsListCall{s: s.s, ctg: category}
}

// ProductsListCall represents the list caller
type ProductsListCall struct {
	s *Service

	ctg    string
	header http.Header
}

// Do executes a list call
func (c *ProductsListCall) Do() (*ProductsListResponse, error) {
	urls := c.s.baseURL + "products/" + strings.ToLower(c.ctg)
	fmt.Println(urls)
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
	ret := &ProductsListResponse{
		ServerResponse: ServerResponse{
			Header: res.Header,
			Code:   res.StatusCode,
		},
	}
	ret.Products = nil
	target := &ret
	if err := json.NewDecoder(res.Body).Decode(&(*target).Products); err != nil {
		return nil, err
	}
	return ret, nil
}

// ProductsListResponse is the response from an executed list call
type ProductsListResponse struct {
	Products []*Product

	ServerResponse
}
