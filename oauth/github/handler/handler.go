package handler

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	util "github.com/micro/go-micro/v2/util/http"
	"github.com/micro/go-micro/v2/web"
	"golang.org/x/oauth2"

	gologinOauth "github.com/dghubble/gologin/v2/oauth2"
	gh "github.com/google/go-github/github"
	githubApi "github.com/google/go-github/v29/github"
	githubOAuth2 "golang.org/x/oauth2/github"
)

// Handler contains a reference to the users service
type Handler struct {
	auth auth.Auth

	githubOrgID  int64
	githubTeamID int64
}

// RegisterHandler adds the GitHub oauth handlers to the servie
func RegisterHandler(srv web.Service) {
	service := srv.Options().Service

	// Setup oauth2 config
	oauth2Config := &oauth2.Config{
		ClientID:     getConfig(service, "client_id"),
		ClientSecret: getConfig(service, "client_secret"),
		RedirectURL:  getConfig(service, "redirect"),
		Endpoint:     githubOAuth2.Endpoint,
		Scopes:       []string{"user:email", "read:org"},
	}

	h := Handler{
		auth: service.Options().Auth,
	}

	// Set GitHub Env Vars
	if id, err := strconv.ParseInt(getConfig(service, "team_id"), 10, 64); err != nil {
		log.Fatalf("Invalid team_id: %v", err)
	} else {
		h.githubTeamID = id
	}

	if id, err := strconv.ParseInt(getConfig(service, "org_id"), 10, 64); err != nil {
		log.Fatalf("Invalid org_id: %v", err)
	} else {
		h.githubOrgID = id
	}

	srv.HandleFunc("/*", func(w http.ResponseWriter, req *http.Request) {
		w.WriteHeader(200)
	})

	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	srv.Handle("/login", github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil)))
	srv.Handle("/verify", github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, func() http.Handler {
		return h.issueSession(srv)
	}(), nil)))
}

// issueSession issues a cookie session after successful GitHub login
func (h *Handler) issueSession(service web.Service) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()

		// Get the client and user from the  context
		client, err := clientFromContext(ctx)
		if err != nil {
			util.WriteInternalServerError(w, err)
			return
		}
		user, err := github.UserFromContext(ctx)
		if err != nil {
			util.WriteInternalServerError(w, err)
			return
		}

		// Check the user has access to the team
		membership, _, err := client.Teams.GetTeamMembership(ctx, h.githubTeamID, user.GetLogin())
		if err != nil || membership.GetState() != "active" {
			http.Redirect(w, req, "/not-invited", http.StatusFound)
			return
		}

		// get the primary email for the user
		if email, err := emailForGitHubUser(ctx, client); err != nil {
			util.WriteInternalServerError(w, err)
		} else {
			user.Email = &email
		}

		// Create the user
		md := h.metadataForUser(ctx, client, user)
		acc, err := h.auth.Generate(*user.Email, auth.Metadata(md))
		if err != nil {
			util.WriteInternalServerError(w, err)
			return
		}
		if acc == nil {
			util.WriteInternalServerError(w, errors.New("Account is empty"))
			return
		}

		// Write the token to cookies
		http.SetCookie(w, &http.Cookie{
			Name:    "micro_token",
			Value:   acc.Token,
			Expires: time.Now().Add(time.Hour * 24),
			Path:    "/",
		})

		http.Redirect(w, req, "/services", http.StatusFound)
	}

	return http.HandlerFunc(fn)
}

// metadataForUser gets the github metadata for the user
func (h *Handler) metadataForUser(ctx context.Context, client *githubApi.Client, user *gh.User) map[string]string {
	team, _, err := client.Teams.GetTeamByID(ctx, h.githubOrgID, h.githubTeamID)
	teamName := ""
	if err == nil {
		teamName = team.GetName()
	}

	org, _, err := client.Organizations.GetByID(ctx, h.githubOrgID)
	teamURL := ""
	if err == nil {
		teamURL = fmt.Sprintf("https://github.com/orgs/%v/teams/%v", org.GetLogin(), team.GetSlug())
	}

	return map[string]string{
		"email":                   *user.Email,
		"name":                    *user.Name,
		"login":                   *user.Login,
		"avatar_url":              user.GetAvatarURL(),
		"team_name":               teamName,
		"team_url":                teamURL,
		"organization_avatar_url": org.GetAvatarURL(),
	}
}

// getConfig loads a string variable from micro config
func getConfig(srv micro.Service, key string) string {
	path := []string{"micro", "oauth", "github", key}
	val := srv.Options().Config.Get(path...).String("")

	if len(val) == 0 {
		log.Fatalf("Missing Required Config: %v", strings.Join(path, "."))
	}
	return val
}

// emailForGitHubUser returns the github users primary email address
func emailForGitHubUser(ctx context.Context, client *githubApi.Client) (string, error) {
	emails, _, err := client.Users.ListEmails(ctx, nil)
	if err != nil {
		return "", err
	}

	for _, email := range emails {
		if email.GetPrimary() {
			return email.GetEmail(), nil
		}
	}

	return "", errors.New("User missing email")
}

// clientFromContext extracts the github client from the context
func clientFromContext(ctx context.Context) (*githubApi.Client, error) {
	oauthToken, err := gologinOauth.TokenFromContext(ctx)
	if err != nil {
		return nil, err
	}

	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: oauthToken.AccessToken},
	)

	tc := oauth2.NewClient(ctx, ts)
	return githubApi.NewClient(tc), nil
}
