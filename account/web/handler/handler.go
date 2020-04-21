package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/auth/provider"
	"github.com/micro/go-micro/v2/auth/provider/oauth"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/store"

	login "github.com/micro/services/login/service/proto/login"
	invites "github.com/micro/services/teams/invites/proto/invites"
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
		githubOrgID:  getConfigInt(srv, "github", "org_id"),
		githubTeamID: getConfigInt(srv, "github", "team_id"),
		auth:         srv.Options().Auth,
		store:        store.DefaultStore,
		users:        users.NewUsersService("go.micro.service.users", srv.Client()),
		login:        login.NewLoginService("go.micro.service.login", srv.Client()),
		invites:      invites.NewInvitesService("go.micro.service.teams.invites", srv.Client()),
	}
}

// Handler is used to handle oauth logic
type Handler struct {
	githubOrgID  int
	githubTeamID int
	auth         auth.Auth
	users        users.UsersService
	login        login.LoginService
	invites      invites.InvitesService
	google       provider.Provider
	github       provider.Provider
	store        store.Store
}

//
// Helper methods for ensuring persisting invite codes through the oauth flow
//

const storePrefixInviteCode = "invite/"

func (h *Handler) setInviteCode(state, code string) error {
	return h.store.Write(&store.Record{
		Key:    storePrefixInviteCode + state,
		Expiry: time.Minute * 5,
		Value:  []byte(code),
	})
}

func (h *Handler) getInviteCode(state string) (string, error) {
	recs, err := h.store.Read(storePrefixInviteCode + state)
	if err != nil {
		return "", err
	}
	return string(recs[0].Value), nil
}

//
// Helper methods for ensuring oauth state is valid
//

const storePrefixOauthState = "state/"

func (h *Handler) generateOauthState() (string, error) {
	code := uuid.New().String()
	record := &store.Record{Key: storePrefixOauthState + code, Expiry: time.Minute * 5}
	return code, h.store.Write(record)
}

func (h *Handler) validateOauthState(code string) (bool, error) {
	_, err := h.store.Read(storePrefixOauthState + code)
	if err == nil {
		return true, nil
	} else if err == store.ErrNotFound {
		return false, nil
	}
	return false, err
}

//
// Helper methods for recording account secrets
//

const storePrefixAccountSecrets = "secrets/"

func (h *Handler) setAccountSecret(id, secret string) error {
	key := storePrefixAccountSecrets + id
	fmt.Printf("setAccountSecret: %v = %v\n", id, secret)
	return h.store.Write(&store.Record{Key: key, Value: []byte(secret)})
}

func (h *Handler) getAccountSecret(id string) (string, error) {
	key := storePrefixAccountSecrets + id
	recs, err := h.store.Read(key)
	if err != nil {
		return "", err
	}
	fmt.Printf("getAccountSecret: %v = %v\n", id, recs[0].Value)
	return string(recs[0].Value), nil
}

//
// Helper methods for getting config
//

func getConfigString(srv micro.Service, keys ...string) string {
	path := append([]string{"micro", "oauth"}, keys...)
	return srv.Options().Config.Get(path...).String("")
}

func getConfigInt(srv micro.Service, keys ...string) int {
	path := append([]string{"micro", "oauth"}, keys...)
	return srv.Options().Config.Get(path...).Int(0)
}

//
// Helper methods handling errors and setting cookies
//

func (h *Handler) handleError(w http.ResponseWriter, req *http.Request, format string, args ...interface{}) {
	logger.Errorf(format, args...)
	params := url.Values{"error": {fmt.Sprintf(format, args...)}}
	http.Redirect(w, req, "/?"+params.Encode(), http.StatusFound)
}

func (h *Handler) loginUser(w http.ResponseWriter, req *http.Request, tok *auth.Token) {
	http.SetCookie(w, &http.Cookie{
		Name:   auth.TokenCookieName,
		Value:  tok.AccessToken,
		Domain: "micro.mu",
		Path:   "/",
	})

	http.Redirect(w, req, "/", http.StatusFound)
}
