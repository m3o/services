package handler

import (
	"log"

	"github.com/micro/go-micro/v2/auth/provider"

	"github.com/micro/go-micro/v2"
	login "github.com/micro/services/login/service/proto/login"
	users "github.com/micro/services/users/service/proto"
)

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	prov := srv.Options().Auth.Options().Provider
	if prov == nil {
		log.Fatal("Micro Auth provider required")
	}

	return &Handler{
		provider: prov,
		users:    users.NewUsersService("go.micro.srv.users", srv.Client()),
		login:    login.NewLoginService("go.micro.srv.login", srv.Client()),
	}
}

// Handler is used to handle oauth logic
type Handler struct {
	users    users.UsersService
	login    login.LoginService
	provider provider.Provider
}
