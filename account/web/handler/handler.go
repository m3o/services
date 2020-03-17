package handler

import (
	"log"
	"strings"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/auth/provider"
	"github.com/micro/go-micro/v2/auth/provider/oauth"

	login "github.com/micro/services/login/service/proto/login"
	users "github.com/micro/services/users/service/proto"
)

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	googleProv := oauth.NewProvider(
		provider.Credentials(
			getConfig(srv, "google", "client_id"),
			getConfig(srv, "google", "client_secret"),
		),
		provider.Redirect(
			getConfig(srv, "google", "redirect"),
		),
		provider.Endpoint(
			getConfig(srv, "google", "endpoint"),
		),
		provider.Scope(
			getConfig(srv, "google", "scope"),
		),
	)

	githubProv := oauth.NewProvider(
		provider.Credentials(
			getConfig(srv, "github", "client_id"),
			getConfig(srv, "github", "client_secret"),
		),
		provider.Redirect(
			getConfig(srv, "github", "redirect"),
		),
		provider.Endpoint(
			getConfig(srv, "github", "endpoint"),
		),
		provider.Scope(
			getConfig(srv, "github", "scope"),
		),
	)

	return &Handler{
		google: googleProv,
		github: githubProv,
		auth:   srv.Options().Auth,
		users:  users.NewUsersService("go.micro.srv.users", srv.Client()),
		login:  login.NewLoginService("go.micro.srv.login", srv.Client()),
	}
}

// Handler is used to handle oauth logic
type Handler struct {
	auth   auth.Auth
	users  users.UsersService
	login  login.LoginService
	google provider.Provider
	github provider.Provider
}

func getConfig(srv micro.Service, keys ...string) string {
	path := append([]string{"micro", "oauth"}, keys...)
	val := srv.Options().Config.Get(path...).String("")
	if len(val) == 0 {
		log.Fatalf("Missing required config: %v", strings.Join(path, "."))
	}
	return val
}
