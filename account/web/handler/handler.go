package handler

import (
	"log"
	"strings"

	"github.com/micro/go-micro/v2/auth/provider/oauth"

	"github.com/micro/go-micro/v2/auth/provider"

	"github.com/micro/go-micro/v2"
	login "github.com/micro/services/login/service/proto/login"
	users "github.com/micro/services/users/service/proto"
)

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	prov := oauth.NewProvider(
		provider.Credentials(
			getConfig(srv, "oauth", "google", "client_id"),
			getConfig(srv, "oauth", "google", "client_secret"),
		),
		provider.Redirect(
			getConfig(srv, "oauth", "google", "redirect"),
		),
		provider.Endpoint(
			getConfig(srv, "oauth", "google", "endpoint"),
		),
	)

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

func getConfig(srv micro.Service, keys ...string) string {
	path := append([]string{"micro", "account"}, keys...)
	val := srv.Options().Config.Get(path...).String("")
	if len(val) == 0 {
		log.Fatalf("Missing required config: %v", strings.Join(path, "."))
	}
	return val
}
