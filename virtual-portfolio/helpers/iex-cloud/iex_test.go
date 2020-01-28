package iex

import (
	"testing"
)

const TestToken = "Tsk_e369ac6651f543c8b84c5de83069a98e"
const TestBaseURL = "https://sandbox.iexapis.com/stable"

func TestNew(t *testing.T) {
	tt := []struct {
		Name  string
		Token string
		Error error
	}{
		{Name: "Valid token", Token: TestToken, Error: nil},
		{Name: "Invalid token", Token: "INVALID", Error: ErrAuthentication},
		{Name: "Missing token", Token: "", Error: ErrAuthentication},
	}

	urlConfig := Configuration{Name: "BaseURL", Value: TestBaseURL}

	for _, tc := range tt {
		t.Run(tc.Name, func(t *testing.T) {
			if _, err := New(tc.Token, urlConfig); err != tc.Error {
				t.Errorf("Incorrect error returned. Expected '%v', got '%v'", tc.Error, err)
			}
		})
	}
}

func validTestService() (Service, error) {
	urlConfig := Configuration{Name: "BaseURL", Value: TestBaseURL}
	return New(TestToken, urlConfig)
}
