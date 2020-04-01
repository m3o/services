package handler

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"

	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	name    string
	auth    auth.Auth
	users   users.UsersService
	payment payment.ProviderService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:    srv.Name(),
		auth:    srv.Options().Auth,
		users:   users.NewUsersService("go.micro.service.users", srv.Client()),
		payment: payment.NewProviderService("go.micro.service.payment.stripe", srv.Client()),
	}
}
