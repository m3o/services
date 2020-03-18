package handler

import (
	"fmt"
	"log"
	"net/http"
	"net/url"
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

	return &Handler{
		google:       googleProv,
		github:       githubProv,
		githubTeamID: getConfigInt(srv, "github", "team_id"),
		githubOrgID:  getConfigInt(srv, "github", "org_id"),
		auth:         srv.Options().Auth,
		users:        users.NewUsersService("go.micro.srv.users", srv.Client()),
		login:        login.NewLoginService("go.micro.srv.login", srv.Client()),
	}
}

// Handler is used to handle oauth logic
type Handler struct {
	githubTeamID int
	githubOrgID  int
	auth         auth.Auth
	users        users.UsersService
	login        login.LoginService
	google       provider.Provider
	github       provider.Provider
}

func getConfigString(srv micro.Service, keys ...string) string {
	path := append([]string{"micro", "oauth"}, keys...)
	val := srv.Options().Config.Get(path...).String("")
	if len(val) == 0 {
		log.Fatalf("Missing required config: %v", strings.Join(path, "."))
	}
	return val
}

func getConfigInt(srv micro.Service, keys ...string) int {
	path := append([]string{"micro", "oauth"}, keys...)
	val := srv.Options().Config.Get(path...).Int(0)
	if val == 0 {
		log.Fatalf("Missing required config: %v", strings.Join(path, "."))
	}
	return val
}

func (h *Handler) handleError(w http.ResponseWriter, req *http.Request, format string, args ...interface{}) {
	params := url.Values{"error": {fmt.Sprintf(format, args...)}}
	http.Redirect(w, req, "/account?"+params.Encode(), http.StatusFound)
}

func (h *Handler) loginUser(w http.ResponseWriter, req *http.Request, user *users.User, roleNames ...string) {
	// Determine the users roles
	var roles []*auth.Role
	for _, n := range roleNames {
		roles = append(roles, &auth.Role{Name: n})
	}

	// Create an auth token
	acc, err := h.auth.Generate(user.Id, auth.Roles(roles))
	if err != nil {
		h.handleError(w, req, "Error creating auth account: %v", err)
		return
	}

	// Set cookie and redirect
	http.SetCookie(w, &http.Cookie{
		Name:   auth.CookieName,
		Value:  acc.Token,
		Domain: "micro.mu",
		Path:   "/",
	})

	http.Redirect(w, req, "/account", http.StatusFound)
}
