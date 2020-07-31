package handler

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/m3o/services/account/web/provider"
	"github.com/m3o/services/account/web/provider/oauth"
	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	mconfig "github.com/micro/micro/v3/service/config"
	mstore "github.com/micro/micro/v3/service/store"

	login "github.com/m3o/services/login/service/proto/login"
	invite "github.com/m3o/services/projects/invite/proto"
	users "github.com/m3o/services/users/service/proto"
)

// NewHandler returns an initialised handler
func NewHandler() *Handler {
	googleProv := oauth.NewProvider(
		provider.Credentials(
			getConfigString("google", "client_id"),
			getConfigString("google", "client_secret"),
		),
		provider.Redirect(getConfigString("google", "redirect")),
		provider.Endpoint(getConfigString("google", "endpoint")),
		provider.Scope(getConfigString("google", "scope")),
	)

	githubProv := oauth.NewProvider(
		provider.Credentials(
			getConfigString("github", "client_id"),
			getConfigString("github", "client_secret"),
		),
		provider.Redirect(getConfigString("github", "redirect")),
		provider.Endpoint(getConfigString("github", "endpoint")),
		provider.Scope(getConfigString("github", "scope")),
	)

	return &Handler{
		google:       googleProv,
		github:       githubProv,
		githubOrgID:  getConfigInt("github", "org_id"),
		githubTeamID: getConfigInt("github", "team_id"),
		users:        users.NewUsersService("go.micro.service.users"),
		login:        login.NewLoginService("go.micro.service.login"),
		invite:       invite.NewInviteService("go.micro.service.projects.invite"),
	}
}

// Handler is used to handle oauth logic
type Handler struct {
	githubOrgID  int
	githubTeamID int
	users        users.UsersService
	login        login.LoginService
	invite       invite.InviteService
	google       provider.Provider
	github       provider.Provider
	store        store.Store
}

//
// Helper methods for persisting invite codes through oauth flows
//

const storePrefixInviteCode = "invite/"

func (h *Handler) setInviteCode(state, code string) error {
	return mstore.Write(&store.Record{
		Key:    storePrefixInviteCode + state,
		Expiry: time.Minute * 5,
		Value:  []byte(code),
	})
}

func (h *Handler) getInviteCode(state string) (string, error) {
	recs, err := mstore.Read(storePrefixInviteCode + state)
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
	return code, mstore.Write(record)
}

func (h *Handler) validateOauthState(code string) (bool, error) {
	_, err := mstore.Read(storePrefixOauthState + code)
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
	return mstore.Write(&store.Record{Key: key, Value: []byte(secret)})
}

func (h *Handler) getAccountSecret(id string) (string, error) {
	key := storePrefixAccountSecrets + id
	recs, err := mstore.Read(key)
	if err != nil {
		return "", err
	}
	fmt.Printf("getAccountSecret: %v = %v\n", id, recs[0].Value)
	return string(recs[0].Value), nil
}

//
// Helper methods for getting config
//

func getConfigString(keys ...string) string {
	path := append([]string{"micro", "oauth"}, keys...)
	return mconfig.Get(path...).String("")
}

func getConfigInt(keys ...string) int {
	path := append([]string{"micro", "oauth"}, keys...)
	return mconfig.Get(path...).Int(0)
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
		Name:   "micro-token",
		Domain: "micro.mu",
		Value:  tok.AccessToken,
		Path:   "/",
	})

	http.Redirect(w, req, "/", http.StatusFound)
}
