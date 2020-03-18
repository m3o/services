package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

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
	r1, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	r1.Header.Add("Accept", "application/json")
	r1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	// Execute the request
	resp, err := client.Do(r1)
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
	// Use the token to get the users profile
	r2, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r2.Header.Add("Authorization", "Bearer "+result.Token)
	resp, err = client.Do(r2)
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
		Username string `json:"login"`
		Name     string `json:"name"`
		Email    string `json:"email"`
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

	h.loginUser(w, req, uRsp.User, "developer")
}
