package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
	"time"

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
			getConfigString(srv, "google", "client_id"),
			getConfigString(srv, "google", "client_secret"),
		),
		provider.Redirect(getConfigString(srv, "google", "redirect")),
		provider.Endpoint(getConfigString(srv, "google", "endpoint")),
		provider.Scope(getConfigString(srv, "google", "scope")),
	)

	githubProv := oauth.NewProvider(
		provider.Credentials(
			getConfigString(srv, "github", "client_id"),
			getConfigString(srv, "github", "client_secret"),
		),
		provider.Redirect(getConfigString(srv, "github", "redirect")),
		provider.Endpoint(getConfigString(srv, "github", "endpoint")),
		provider.Scope(getConfigString(srv, "github", "scope")),
	)

	account, err := srv.Options().Auth.Generate(srv.Name(),
		auth.WithRoles("service", fmt.Sprintf("service.%v", srv.Name())),
	)
	if err != nil {
		log.Fatalf("Unable to generate service auth account: %v", err)
	}
	token, err := srv.Options().Auth.Refresh(account.Secret.Token)
	if err != nil {
		log.Fatalf("Unable to generate service auth token: %v", err)
	}

	return &Handler{
		google:       googleProv,
		github:       githubProv,
		authToken:    token.Token,
		githubOrgID:  getConfigInt(srv, "github", "org_id"),
		githubTeamID: getConfigInt(srv, "github", "team_id"),
		auth:         srv.Options().Auth,
		users:        users.NewUsersService("go.micro.service.users", srv.Client()),
		login:        login.NewLoginService("go.micro.service.login", srv.Client()),
	}
}

// Handler is used to handle oauth logic
type Handler struct {
	githubOrgID  int
	githubTeamID int
	authToken    string
	auth         auth.Auth
	users        users.UsersService
	login        login.LoginService
	google       provider.Provider
	github       provider.Provider
}

func getConfigString(srv micro.Service, keys ...string) string {
	path := append([]string{"micro", "oauth"}, keys...)
	return srv.Options().Config.Get(path...).String("")
}

func getConfigInt(srv micro.Service, keys ...string) int {
	path := append([]string{"micro", "oauth"}, keys...)
	return srv.Options().Config.Get(path...).Int(0)
}

func (h *Handler) handleError(w http.ResponseWriter, req *http.Request, format string, args ...interface{}) {
	params := url.Values{"error": {fmt.Sprintf(format, args...)}}
	http.Redirect(w, req, "/account?"+params.Encode(), http.StatusFound)
}

func (h *Handler) loginUser(w http.ResponseWriter, req *http.Request, user *users.User, roles ...string) {
	// Create an auth account
	acc, err := h.auth.Generate(user.Id, auth.WithRoles(roles...))
	if err != nil {
		h.handleError(w, req, "Error creating auth account: %v", err)
		return
	}

	// Create an auth token
	tok, err := h.auth.Refresh(acc.Secret.Token, auth.WithTokenExpiry(time.Hour*24))
	if err != nil {
		h.handleError(w, req, "Error creating auth token: %v", err)
		return
	}

	// Set cookie and redirect
	http.SetCookie(w, &http.Cookie{
		Name:   auth.TokenCookieName,
		Value:  tok.Token,
		Domain: "micro.mu",
		Path:   "/",
	})

	http.Redirect(w, req, "/account", http.StatusFound)
}
