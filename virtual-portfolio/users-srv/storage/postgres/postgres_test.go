package postgres

import (
	"fmt"
	"testing"

	"github.com/kytra-app/users-srv/storage"
)

func TestMapWithoutBlank(t *testing.T) {
	tt := []struct {
		name  string
		user  storage.User
		count int
	}{
		{name: "Empty", user: storage.User{}, count: 2}, // created_at and updated_at are defaulted
		{name: "UUID", user: storage.User{UUID: "UUID"}, count: 3},
		{name: "Password", user: storage.User{Password: "Password"}, count: 3},
		{name: "Name", user: storage.User{FirstName: "John", LastName: "Doe"}, count: 4},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			r := mapWithoutBlank(tc.user)
			if len(r) != tc.count {
				fmt.Println(r)
				t.Errorf("Expected %v results, got %v", tc.count, len(r))
			}
		})
	}
}
