package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/dghubble/gologin/v2"
	"github.com/dghubble/gologin/v2/github"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/web"
	utils "github.com/micro/services/serverless/web/util"

	"golang.org/x/oauth2"

	gologinOauth "github.com/dghubble/gologin/v2/oauth2"
	githubApi "github.com/google/go-github/v29/github"
	githubOAuth2 "golang.org/x/oauth2/github"
)

// RegisterHandlers adds the GitHub oauth handlers to the servie
func RegisterHandlers(srv web.Service) error {
	conf := srv.Options().Service.Options().Config

	oauthClientID := conf.Get("micro", "serverless", "github", "oauth_client_id").String("")
	if len(oauthClientID) == 0 {
		return errors.New("No pauth client ID")
	}

	oauthClientSecret := conf.Get("micro", "serverless", "github", "oauth_client_secret").String("")
	if len(oauthClientSecret) == 0 {
		return errors.New("No oauth client secret")
	}

	redirectURL := conf.Get("micro", "serverless", "github", "oauth_redirect_url").String("")
	if len(redirectURL) == 0 {
		return errors.New("No oauth redirect url")
	}

	oauth2Config := &oauth2.Config{
		ClientID:     oauthClientID,
		ClientSecret: oauthClientSecret,
		RedirectURL:  redirectURL,
		Endpoint:     githubOAuth2.Endpoint,
		Scopes:       []string{"user:email", "read:org", "public_repo"},
	}

	// state param cookies require HTTPS by default; disable for localhost development
	stateConfig := gologin.DebugOnlyCookieConfig
	srv.HandleFunc("/v1/github/organisations", listOrgs(srv))
	srv.HandleFunc("/v1/github/repositories", listRepos(srv))
	srv.HandleFunc("/v1/github/folders", listFolders(srv))
	srv.Handle("/v1/github/login", github.StateHandler(stateConfig, github.LoginHandler(oauth2Config, nil)))
	srv.Handle("/v1/auth/verify", github.StateHandler(stateConfig, github.CallbackHandler(oauth2Config, func() http.Handler {
		return issueSession(srv)
	}(), nil)))
	srv.HandleFunc("/v1/user", userHandler(srv))

	return nil
}

// issueSession issues a cookie session after successful Github login
func issueSession(service web.Service) http.Handler {
	fn := func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		oauthToken, err := gologinOauth.TokenFromContext(ctx)
		if err != nil {
			utils.Write500(w, err)
			return
		}
		githubUser, err := github.UserFromContext(ctx)
		if err != nil {
			utils.Write500(w, err)
			return
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: oauthToken.AccessToken},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := githubApi.NewClient(tc)

		// Have to list the emails separately as users with a private email address
		// will not have an email in githubUser.Email
		emails, _, err := client.Users.ListEmails(ctx, nil)
		if err != nil {
			utils.Write500(w, err)
			return
		}
		primaryEmail := ""
		for _, email := range emails {
			if email.GetPrimary() {
				primaryEmail = email.GetEmail()
			}
		}
		githubUser.Email = &primaryEmail

		conf := service.Options().Service.Options().Config
		teamID, err := strconv.ParseInt(conf.Get("micro", "serverless", "github", "team_id").String(""), 10, 64)
		if err != nil {
			utils.Write500(w, err)
			return
		}

		frontendAddress := conf.Get("micro", "serverless", "frontend_address").String("")
		membership, _, err := client.Teams.GetTeamMembership(req.Context(), teamID, githubUser.GetLogin())
		if err != nil {
			http.Redirect(w, req, frontendAddress+"/not-invited", http.StatusFound)
			return
		}
		if membership.GetState() != "active" {
			http.Redirect(w, req, frontendAddress+"/not-invited", http.StatusFound)
			return
		}
		// gracefully degrading in case we have no ORG ID
		// ORG ID is only needed so we can read the team for teamname
		orgID, _ := strconv.ParseInt(conf.Get("micro", "serverless", "github", "org_id").String(""), 10, 64)
		team, _, err := client.Teams.GetTeamByID(req.Context(), orgID, teamID)
		teamName := ""
		if err == nil {
			teamName = team.GetName()
		}
		org, _, err := client.Organizations.GetByID(req.Context(), orgID)
		teamURL := ""
		if err == nil {
			teamURL = fmt.Sprintf("https://github.com/orgs/%v/teams/%v", org.GetLogin(), team.GetSlug())
		}
		acc, err := service.Options().Service.Options().Auth.Generate(*githubUser.Email, auth.Metadata(
			map[string]string{
				"email":                   *githubUser.Email,
				"name":                    *githubUser.Name,
				"login":                   *githubUser.Login,
				"avatar_url":              githubUser.GetAvatarURL(),
				"team_name":               teamName,
				"team_url":                teamURL,
				"organization_avatar_url": org.GetAvatarURL(),
				"github_access_token":     oauthToken.AccessToken,
			}))
		if err != nil {
			utils.Write500(w, err)
			return
		}
		if acc == nil {
			utils.Write500(w, errors.New("Account is empty"))
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:    "micro_token",
			Value:   acc.Token,
			Expires: acc.Expiry,
			Path:    "/",
		})

		http.Redirect(w, req, frontendAddress+"/", http.StatusFound)
	}
	return http.HandlerFunc(fn)
}

type User struct {
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	AvatarURL             string `json:"avatarURL"`
	TeamName              string `json:"teamName"`
	TeamURL               string `json:"teamURL"`
	OrganizationAvatarURL string `json:"organizationAvatarURL"`
	Login                 string `json:"login"`
}

func userHandler(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		token := req.URL.Query().Get("token")
		if len(token) == 0 {
			utils.Write400(w, errors.New("Token missing"))
			return
		}

		acc, err := service.Options().Service.Options().Auth.Verify(token)
		if err != nil {
			utils.Write400(w, err)
			return
		}
		if acc == nil {
			utils.Write400(w, errors.New("Not found"))
			return
		}

		if acc.Metadata == nil {
			utils.Write400(w, errors.New("Metadata not found"))
			return
		}

		utils.WriteJSON(w, &User{
			Name:                  acc.Metadata["name"],
			Email:                 acc.Metadata["email"],
			AvatarURL:             acc.Metadata["avatar_url"],
			TeamName:              acc.Metadata["team_name"],
			TeamURL:               acc.Metadata["team_url"],
			OrganizationAvatarURL: acc.Metadata["organization_avatar_url"],
			Login:                 acc.Metadata["login"],
		})
	}
}

func listOrgs(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		token := req.URL.Query().Get("token")
		if len(token) == 0 {
			utils.Write400(w, errors.New("Token missing"))
			return
		}

		acc, err := service.Options().Service.Options().Auth.Verify(token)
		if err != nil {
			utils.Write400(w, err)
			return
		}
		if acc == nil {
			utils.Write400(w, errors.New("Not found"))
			return
		}

		if acc.Metadata == nil {
			utils.Write400(w, errors.New("Metadata not found"))
			return
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: acc.Metadata["github_access_token"]},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := githubApi.NewClient(tc)

		orgs, _, err := client.Organizations.List(ctx, acc.Metadata["login"], nil)
		if err != nil {
			utils.Write500(w, err)
			return
		}
		utils.WriteJSON(w, orgs)
	}
}

func listRepos(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		token := req.URL.Query().Get("token")
		if len(token) == 0 {
			utils.Write400(w, errors.New("Token missing"))
			return
		}

		acc, err := service.Options().Service.Options().Auth.Verify(token)
		if err != nil {
			utils.Write400(w, err)
			return
		}
		if acc == nil {
			utils.Write400(w, errors.New("Not found"))
			return
		}

		if acc.Metadata == nil {
			utils.Write400(w, errors.New("Metadata not found"))
			return
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: acc.Metadata["github_access_token"]},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := githubApi.NewClient(tc)

		org := req.URL.Query().Get("organisation")
		if len(org) == 0 {
			utils.Write400(w, errors.New("Organization missing"))
			return
		}

		repos, _, err := client.Repositories.ListByOrg(ctx, org, nil)
		if err != nil {
			utils.Write500(w, err)
			return
		}
		utils.WriteJSON(w, repos)
	}
}

func listFolders(service web.Service) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		utils.SetupResponse(&w, req)
		if (*req).Method == "OPTIONS" {
			return
		}
		token := req.URL.Query().Get("token")
		if len(token) == 0 {
			utils.Write400(w, errors.New("Token missing"))
			return
		}

		acc, err := service.Options().Service.Options().Auth.Verify(token)
		if err != nil {
			utils.Write400(w, err)
			return
		}
		if acc == nil {
			utils.Write400(w, errors.New("Not found"))
			return
		}

		if acc.Metadata == nil {
			utils.Write400(w, errors.New("Metadata not found"))
			return
		}

		ts := oauth2.StaticTokenSource(
			&oauth2.Token{AccessToken: acc.Metadata["github_access_token"]},
		)
		tc := oauth2.NewClient(ctx, ts)
		client := githubApi.NewClient(tc)

		org := req.URL.Query().Get("organisation")
		if len(org) == 0 {
			utils.Write400(w, errors.New("Organization missing"))
			return
		}

		repo := req.URL.Query().Get("repository")
		if len(org) == 0 {
			utils.Write400(w, errors.New("Repository missing"))
			return
		}

		path := req.URL.Query().Get("path")
		if len(org) == 0 {
			utils.Write400(w, errors.New("Repository missing"))
			return
		}

		repoParts := strings.Split(repo, "/")
		if len(repoParts) > 1 {
			repo = repoParts[1]
		}
		_, dirs, _, err := client.Repositories.GetContents(ctx, org, repo, path, nil)
		if err != nil {
			utils.Write500(w, err)
			return
		}
		utils.WriteJSON(w, dirs)
	}
}
