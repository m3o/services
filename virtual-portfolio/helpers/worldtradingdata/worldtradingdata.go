package worldtradingdata

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
)

// BaseURL for the WorldTradingData API
const BaseURL = "https://api.worldtradingdata.com/api/v1"

// ErrNetwork is returned when an unknown error occurred. It is safe to retry.
var ErrNetwork = errors.New("An unknown network error occured")

// ErrAuthentication is returned when a 401 or 403 status is returned by WorldTradingData
var ErrAuthentication = errors.New("The API key provided is not valid")

// ErrNotFound is returned when a 404 status is returned by WorldTradingData
var ErrNotFound = errors.New("Resource not found")

// ErrBadRequest is returned when a 400 status is returned by WorldTradingData
var ErrBadRequest = errors.New("Bad request")

// ErrUnknown is returned when an unexpected status code is returned by WorldTradingData
var ErrUnknown = errors.New("An unknown server error occured")

// Configuration is an option which can be passed when initializing the service
type Configuration struct {
	Name  string
	Value string
}

// New takes a WorldTradingData secret Token (normally prefixed with sk_), validates the
// Token, and returns a Service, and an error.
func New(token string, config ...Configuration) (Service, error) {
	if token == "" {
		return Handler{}, ErrAuthentication
	}

	// Initialize the handler with the default BaseURL
	h := Handler{Token: token, BaseURL: BaseURL}

	// Apply the configuration
	for _, c := range config {
		switch c.Name {
		case "BaseURL":
			h.BaseURL = c.Value
			break
		}
	}

	// Ping the Stock API to validate the Token
	_, err := h.Get("/stock?symbol=^IXIC")
	return h, err
}

// Service is a representation the of WorldTradingData Cloud service.
type Service interface {
	Get(string) ([]byte, error)
	Quote(string) (*Quote, error)
}

// Handler is an implementation of service
type Handler struct {
	Token   string
	BaseURL string
}

// Get performs a HTTP Get request on the WorldTradingData API
func (h Handler) Get(path string) ([]byte, error) {
	var url string
	if strings.Contains(path, "?") {
		url = fmt.Sprintf("%v/%v&api_token=%v", h.BaseURL, path, h.Token)
	} else {
		url = fmt.Sprintf("%v/%v?api_token=%v", h.BaseURL, path, h.Token)
	}

	rsp, err := http.Get(url)
	if err != nil {
		return []byte{}, ErrNetwork
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []byte{}, ErrUnknown
	}

	switch rsp.StatusCode {
	case 200, 201:
		return body, nil
	case 400:
		return body, ErrBadRequest
	case 401, 402, 403:
		return body, ErrAuthentication
	case 404:
		return body, ErrNotFound
	default:
		return body, ErrUnknown
	}
}
