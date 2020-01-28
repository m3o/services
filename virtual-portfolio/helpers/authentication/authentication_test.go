package authetication

import "testing"

const testToken = "IB70H8C7lv9dEvCLCzj9HGeENMgHy3f91PqJzXH10L1k5SqYoJ"

func TestNew(t *testing.T) {
	tt := []struct {
		name  string
		token string
		err   error
	}{
		{name: "No Token", err: ErrInvalidSigningKey},
		{name: "Invalid Token", token: "DEMO", err: ErrInvalidSigningKey},
		{name: "Valid Token", token: testToken},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			if _, e := New(tc.token); e != tc.err {
				t.Errorf("Unexpected error, got %v, expected %v", e, tc.err)
			}
		})
	}
}

func TestEncodeUserDecodeToken(t *testing.T) {
	a, err := New(testToken)
	if err != nil {
		t.Fatalf("Unable to create a valid Authenticator: %v", err)
	}

	uuid := "MYUUID"

	var token string
	token, err = a.EncodeUser(User{UUID: uuid})
	if err != nil {
		t.Fatalf("Unable to encode user: %v", err)
	}

	var user User
	user, err = a.DecodeToken(token)
	if err != nil {
		t.Fatalf("Unable to decode token: %v", err)
	}

	if user.UUID != uuid {
		t.Errorf("Incorrect UUID returned, expected %v, got %v", uuid, user.UUID)
	}
}
