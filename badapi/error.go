package badapi

import (
	"fmt"
	"io/ioutil"
	"net/http"

	"golang.org/x/xerrors"
)

// Error represents a badapi error
type Error struct {
	Code int
	Body string
}

var (
	ErrModeActive = xerrors.New("badapi: error mode active")
	ErrEmptyBody  = xerrors.New("badapi: empty response body")
)

func (e *Error) Error() string {
	return fmt.Sprintf("badapi: got http response status code %d with body %v", e.Code, e.Body)
}

func checkResponse(res *http.Response) error {
	if res.ContentLength <= 0 {
		return ErrEmptyBody
	}
	if res.StatusCode >= 200 && res.StatusCode <= 299 {
		return nil
	}
	errData, _ := ioutil.ReadAll(res.Body)
	return &Error{Code: res.StatusCode, Body: string(errData)}
}
