package test

import (
	"context"
	"fmt"
	"math/rand"

	"github.com/micro/micro/v3/service/auth"
)

func ContextWithAccount(issuer, id string) context.Context {
	return auth.ContextWithAccount(context.TODO(), &auth.Account{
		ID:       id,
		Type:     "",
		Issuer:   issuer,
		Metadata: nil,
		Scopes:   nil,
		Secret:   "",
		Name:     "",
	})
}

func TestEmail() string {
	return fmt.Sprintf("foo%d@bar.com", rand.Int())
}
