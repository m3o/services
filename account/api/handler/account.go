package handler

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"

	login "github.com/micro/services/login/service/proto/login"
	payment "github.com/micro/services/payments/provider/proto"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the account api proto interface
type Handler struct {
	name    string
	auth    auth.Auth
	users   users.UsersService
	login   login.LoginService
	payment payment.ProviderService
}

// NewHandler returns an initialised handle
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:    srv.Name(),
		auth:    srv.Options().Auth,
		users:   users.NewUsersService("go.micro.service.users", srv.Client()),
		login:   login.NewLoginService("go.micro.service.login", srv.Client()),
		payment: payment.NewProviderService("go.micro.service.payment.stripe", srv.Client()),
	}
}
