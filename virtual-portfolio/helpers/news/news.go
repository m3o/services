package news

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

// BaseURL for the IEX API
const BaseURL = "https://stocknewsapi.com/api/v1"

// ErrNetwork is returned when an unknown error occurred. It is safe to retry.
var ErrNetwork = errors.New("An unknown network error occured")

// ErrAuthentication is returned when a 401 or 403 status is returned by IEX
var ErrAuthentication = errors.New("The API key provided is not valid")

// ErrNotFound is returned when a 404 status is returned by IEX
var ErrNotFound = errors.New("Resource not found")

// ErrBadRequest is returned when a 400 status is returned by IEX
var ErrBadRequest = errors.New("Bad request")

// ErrUnknown is returned when an unexpected status code is returned by IEX
var ErrUnknown = errors.New("An unknown server error occured")

// Configuration is an option which can be passed when initializing the service
type Configuration struct {
	Name  string
	Value string
}

// New takes a News secret Token (normally prefixed with sk_), validates the
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

	// Ping the Tickers API to validate the Token
	_, err := h.Get("?tickers=AAPL&items=1")

	return h, err
}

// Service is a representation the of Stock News service.
type Service interface {
	Get(string) ([]byte, error)
	General() ([]Article, error)
	TopMentions() ([]Mention, error)
	Tickers(...string) ([]Article, error)
}

// Handler is an implementation of service
type Handler struct {
	Token   string
	BaseURL string
}

// Get performs a HTTP Get request on the IEX API
func (h Handler) Get(path string) ([]byte, error) {
	var url string
	if strings.Contains(path, "?") {
		url = fmt.Sprintf("%v/%v&token=%v", h.BaseURL, path, h.Token)
	} else {
		url = fmt.Sprintf("%v/%v?token=%v", h.BaseURL, path, h.Token)
	}
	fmt.Println(url)

	rsp, err := http.Get(url)
	if err != nil {
		return []byte{}, ErrNetwork
	}

	body, err := ioutil.ReadAll(rsp.Body)
	if err != nil {
		return []byte{}, err
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
		return body, err
	}
}

// Mention is returned by the TopMentions API
type Mention struct {
	TotalMentions    int    `json:"total_mentions"`    // Example: 1,
	PositiveMentions int    `json:"positive_mentions"` // Example: 0,
	NegativeMentions int    `json:"negative_mentions"` // Example: 0,
	NeutralMentions  int    `json:"neutral_mentions"`  // Example: 1,
	Ticker           string `json:"ticker"`            // Example: "COTY",
	Name             string `json:"name"`              // Example: "Coty Inc."
}

// Article is returned by the Tickers API
type Article struct {
	NewsURL   string   `json:"news_url"`    // Example: "https://www.geekwire.com/2019/amazon-double-seasonal-hiring-record-200k-workers-business-continues-boom/",
	ImageURL  string   `json:"image_url"`   // Example: "https://cdn.snapi.dev/images/v1/a/m/amzn-seasonal.jpg",
	Title     string   `json:"title"`       // Example: "Amazon to double seasonal hiring to record 200k workers as business continues to boom",
	Text      string   `json:"text"`        // Example: "Amazon plans to hire 200,000 seasonal workers this year, twice as many as last year, suggesting that it expects a strong holiday shopping season.",
	Source    string   `json:"source_name"` // Example: "GeekWire",
	Date      string   `json:"date"`        // Example: "Thu, 28 Nov 2019 00:32:52 -0500",
	Topics    []string `json:"topics"`      // Example: [],
	Sentiment string   `json:"sentiment"`   // Example: "Positive",
	Type      string   `json:"type"`        // Example: "Article",
	Tickers   []string `json:"tickers"`     // Example: ["AMZN"]
}

// CreatedAt is the time the article was published
func (a Article) CreatedAt() (time.Time, error) {
	t, err := time.Parse("Mon, 2 Jan 2006 15:04:05 -0700", a.Date)
	if err != nil {
		return time.Now(), err
	}

	// Reset to UTC
	return time.Unix(t.Unix(), 0), nil
}
