package badapi

import (
	"fmt"
	"io/ioutil"
	"net/http"
)

// Error represents a badapi error
type Error struct {
	Code int
	Body string
}

func (e *Error) Error() string {
	return fmt.Sprintf("badapi: got http response status code %d with body %v", e.Code, e.Body)
}

func checkResponse(res *http.Response) error {
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	errData, _ := ioutil.ReadAll(res.Body)
	return &Error{Code: res.StatusCode, Body: string(errData)}
}
