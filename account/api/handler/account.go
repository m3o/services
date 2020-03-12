package handler

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	log "github.com/micro/go-micro/v2/logger"

	login "github.com/micro/services/login/service/proto/login"
	"github.com/micro/services/payments/provider"
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
	pay, err := provider.NewProvider("stripe", srv.Client())
	if err != nil {
		log.Fatalf("Error setting up payment provider: %v", err)
	}

	return &Handler{
		payment: pay,
		name:    srv.Name(),
		auth:    srv.Options().Auth,
		users:   users.NewUsersService("go.micro.srv.users", srv.Client()),
		login:   login.NewLoginService("go.micro.srv.login", srv.Client()),
	}
}
