package slack

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"net/url"
)

const baseURL = "https://slack.com/api/"

var (
	// ErrAuth is returned when there is an invalid token provided
	ErrAuth = errors.New("Invalid Token Provided")

	// ErrBadRequest is returned when an invalid request is made
	ErrBadRequest = errors.New("Bad Request")

	// ErrUnkown is returned when there is an unexected API error
	ErrUnkown = errors.New("Unknown Error")
)

// Client is an interface of the slack client
type Client interface {
	PostRequest(string, url.Values) (Response, error)
	Ping() error
	SendMessage(string, string) error
}

type client struct {
	token string
}

// Response is an instance of a response from the Slack API
type Response struct {
	Status int
	Body   map[string]interface{}
}

// NewClient takes a token and returns a slack client
func NewClient(token string) (Client, error) {
	c := client{token}

	// Ping the API on initialization to fail fast if invalid token
	if err := c.Ping(); err != nil {
		return c, err
	}

	return c, nil
}

// PostRequest executes a HTTP POST request on the Slack API
func (c client) PostRequest(path string, data url.Values) (Response, error) {
	data.Set("token", c.token)

	res, err := http.PostForm(baseURL+path, data)
	if err != nil {
		return Response{}, ErrUnkown
	}

	response := Response{Status: res.StatusCode}

	if res.StatusCode != http.StatusOK {
		switch res.StatusCode {
		case http.StatusBadRequest:
			return response, ErrBadRequest
		case http.StatusUnauthorized:
			return response, ErrAuth
		default:
			return response, ErrUnkown
		}
	}

	bodyBytes, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return response, ErrUnkown
	}

	if err = json.Unmarshal(bodyBytes, &response.Body); err != nil {
		return response, ErrUnkown
	}

	return response, nil
}

// Ping calls the test endpoint, verifying the token provided
func (c client) Ping() error {
	res, err := c.PostRequest("api.test", make(url.Values))
	if err != nil {
		return err
	}

	if res.Body["ok"] != true {
		return ErrAuth
	}

	return nil
}

// SendMessage posts a message to the given channel
func (c client) SendMessage(channel, msg string) error {
	res, err := c.PostRequest("chat.postMessage", url.Values{
		"channel": {channel},
		"text":    {msg},
	})

	if err != nil {
		return err
	}

	if res.Body["ok"] != true {
		return ErrBadRequest
	}

	return nil
}
