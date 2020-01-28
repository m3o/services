package sms

import (
	"context"
	"net/url"

	twilio "github.com/kevinburke/twilio-go"
)

// Service is a module which allows the sending of SMS messages
type Service interface {
	Send(string, string) error
	Ping() error
}

// New returns an initialized instance of Service
func New(accountSID, authToken string) (Service, error) {
	client := twilio.NewClient(accountSID, authToken, nil)
	s := service{client: client}
	return s, s.Ping()
}

type service struct {
	client *twilio.Client
}

func (s service) Send(to, message string) error {
	data := url.Values{
		"To":   []string{to},
		"From": []string{"Kytra"},
		"Body": []string{message},
	}

	_, err := s.client.Messages.Create(context.TODO(), data)
	if err == nil {
		return nil
	}

	// Try using US number
	data["From"] = []string{"12023014378"}
	_, err = s.client.Messages.Create(context.TODO(), data)

	return err
}

func (s service) Ping() error {
	return s.Send("+447503196716", "Testing Twilio Credentials")
}
