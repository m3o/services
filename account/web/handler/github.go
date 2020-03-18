package handler

import (
	"encoding/json"
	"fmt"
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
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Redirect(w, req, "/account/error", http.StatusFound)
		fmt.Println(err)
		return
	}

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
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Redirect(w, req, "/account/error", http.StatusFound)
		fmt.Println(err)
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
		http.Redirect(w, req, "/account/error", http.StatusFound)
		fmt.Println(err)
		return
	}

	_ = uRsp

	// TODO:
	// - Add the developer role to the user
	// - Add the collaborator role to the user if they're part of the micro GH team

	// Set the cookie and redirect
	// http.SetCookie(w, &http.Cookie{
	// Name: auth.CookieName,
	// Value: uRsp.Token,
	// })
	http.Redirect(w, req, "/account", http.StatusFound)
}
