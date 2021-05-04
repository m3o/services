package auth

import (
	"context"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
)

const (
	CustomerNamespace = "micro"
)

func VerifyMicroCustomer(ctx context.Context, method string) (*auth.Account, error) {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != CustomerNamespace {
		return nil, errors.Forbidden(method, "Forbidden")
	}
	if acc.Type != "user" {
		return nil, errors.Forbidden(method, "Forbidden")
	}
	allowed := false
	for _, s := range acc.Scopes {
		if s == "customer" {
			allowed = true
			break
		}
	}
	if !allowed {
		return nil, errors.Forbidden(method, "Forbidden")
	}
	return acc, nil
}

func VerifyMicroAdmin(ctx context.Context, method string) (*auth.Account, error) {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return nil, errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != CustomerNamespace {
		return nil, errors.Forbidden(method, "Forbidden")
	}

	admin := false
	for _, s := range acc.Scopes {
		if (s == "admin" && acc.Type == "user") || (s == "service" && acc.Type == "service") {
			admin = true
			break
		}
	}
	if !admin {
		return nil, errors.Forbidden(method, "Forbidden")
	}
	return acc, nil
}
