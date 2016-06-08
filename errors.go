package plugapi

import (
	"errors"
	"fmt"
	"net/http"
)

// General API errors
var (
	ErrAuthentication         = errors.New("plugapi: authentication failed")
	ErrAuthenticationRequired = errors.New("plugapi: authentication details required")
	ErrUnknownData            = errors.New("plugapi: cannot structify request")
)

type ErrDataRequestError struct {
	Data     interface{}
	Endpoint string
}

func (e ErrDataRequestError) Error() string {
	return fmt.Sprintf("plugapi: invalid data from %s: %v", e.Endpoint, e.Data)
}

type ErrUnknownResponse struct {
	Response *http.Response
	Endpoint string
}

func (e ErrUnknownResponse) Error() string {
	return fmt.Sprintf("plugapi: bad reply. error %d from %s", e.Response.StatusCode, e.Endpoint)
}

func ErrIsUnknownResponse(err error) bool {
	switch err.(type) {
	case *ErrUnknownResponse:
		return true
	default:
		return false
	}
}

func ErrIsDataRequestError(err error) bool {
	switch err.(type) {
	case *ErrDataRequestError:
		return true
	default:
		return false
	}
}
