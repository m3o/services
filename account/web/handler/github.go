package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/micro/go-micro/v2/logger"
	users "github.com/micro/services/users/service/proto"
)

// HandleGithubOauthLogin redirects the user to begin the oauth flow
func (h *Handler) HandleGithubOauthLogin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, h.github.Endpoint(), http.StatusFound)
}

// HandleGithubOauthVerify redirects the user to begin the oauth flow
func (h *Handler) HandleGithubOauthVerify(w http.ResponseWriter, req *http.Request) {
	// Consturct the requerst to get the access token
	data := url.Values{
		"client_id":     {h.github.Options().ClientID},
		"client_secret": {h.github.Options().ClientSecret},
		"redirect_uri":  {h.github.Redirect()},
		"code":          {req.FormValue("code")},
	}
	r, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	// Execute the request
	resp, err := client.Do(r)
	if err != nil {
		h.handleError(w, req, "Error getting access token from GitHub: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.handleError(w, req, "Error getting access token from GitHub. Status: %v", resp.Status)
		return
	}

	// Decode the token
	var result struct {
		Token string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	logger.Infof("TOKEN: %v", result.Token)

	// Use the token to get the users profile
	r, _ = http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Add("Authorization", "Bearer "+result.Token)
	resp, err = client.Do(r)
	if err != nil {
		h.handleError(w, req, "Error getting user from GitHub: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		bytes, _ := ioutil.ReadAll(resp.Body)
		h.handleError(w, req, "Error getting user from GitHub. Status: %v. Error: %v", resp.Status, string(bytes))
		return
	}

	// Decode the users profile
	var profile struct {
		ID       int64  `json:"id"`
		Name     string `json:"name"`
		Email    string `json:"email"`
		Username string `json:"login"`
	}
	json.NewDecoder(resp.Body).Decode(&profile)

	// Create the user in the users service
	uRsp, err := h.users.Create(req.Context(), &users.CreateRequest{
		User: &users.User{
			Id:       fmt.Sprintf("github_%v", profile.ID),
			Email:    profile.Email,
			Username: profile.Username,
		},
	})
	if err != nil {
		h.handleError(w, req, "Error creating account: %v", err)
		return
	}

	// Setup the roles
	roles := []string{"developer"}

	// Check to see if the user is part of the micro team
	url := fmt.Sprintf("https://api.github.com/teams/%v/memberships/%v", h.githubTeamID, profile.Username)
	r, _ = http.NewRequest("GET", url, nil)
	r.Header.Add("Authorization", "Bearer "+result.Token)
	resp, err = client.Do(r)
	logger.Info("GitHub url: ", url)
	logger.Info("GitHub response: ", resp.StatusCode)
	if err != nil {
		h.handleError(w, req, "Error getting user team membership from GitHub: %v", err)
		return
	}
	if resp.StatusCode == http.StatusOK {
		roles = append(roles, "collaborator")
	}

	// Logiin the user
	h.loginUser(w, req, uRsp.User, roles...)
}
