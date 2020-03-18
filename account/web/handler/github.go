package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	log "github.com/micro/go-micro/v2/logger"
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
	log.Infof("Token: %v", result.Token)

	// Use the token to get the users profile
	r, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Add("Authorization", "Bearer "+result.Token)
	resp, err = client.Do(r)
	if err != nil {
		h.handleError(w, r, "Error getting user from GitHub: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.handleError(w, r, "Error getting user from GitHub. Status: %v", resp.Status)
		return
	}

	// Decode the users profile
	var profile struct {
		ID        string `json:"id"`
		Username  string `json:"login"`
		Name      string `json:"name"`
		Email     string `json:"email"`
		FirstName string
		LastName  string
	}
	json.NewDecoder(resp.Body).Decode(&profile)

	nameComps := strings.Split(profile.Name, "")
	if len(nameComps) > 0 {
		profile.FirstName = nameComps[0]
	}
	if len(nameComps) > 1 {
		profile.LastName = strings.Join(nameComps[1:len(nameComps)-1], " ")
	}

	// Create the user in the users service
	uRsp, err := h.users.Create(req.Context(), &users.CreateRequest{
		User: &users.User{
			Id:        fmt.Sprintf("Github_%v", profile.ID),
			Email:     profile.Email,
			Username:  profile.Username,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
		},
	})
	if err != nil {
		h.handleError(w, req, "Error creating account: %v", err)
		return
	}

	h.loginUser(w, req, uRsp.User, "developer")
}
