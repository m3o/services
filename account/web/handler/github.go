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
	// Get the token using the oauth code
	resp, err := http.PostForm("https://github.com/login/oauth/access_token", url.Values{
		"client_id":     {h.github.Options().ClientID},
		"client_secret": {h.github.Options().ClientSecret},
		"redirect_uri":  {h.github.Redirect()},
		"code":          {req.FormValue("code")},
	})
	if err != nil {
		h.handleError(w, req, "Error getting access token from GitHub: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.handleError(w, req, "Error getting access token from GitHub. Status: %v", resp.Status)
		return
	}

	bytes, _ := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))

	// Decode the token
	var oauthResult struct {
		Token string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&oauthResult)

	// Use the token to get the users profile
	req, err = http.NewRequest("GET", "https://api.github.com/user", nil)
	req.Header.Add("Authorization", "Bearer "+oauthResult.Token)
	client := &http.Client{}
	resp, err = client.Do(req)
	if err != nil {
		h.handleError(w, req, "Error getting user from GitHub: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.handleError(w, req, "Error getting user from GitHub. Status: %v", resp.Status)
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
