package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/micro/go-micro/v2/auth"
	users "github.com/micro/services/users/service/proto"
)

// HandleOauthLogin redirects the user to begin the oauth flow
func (h *Handler) HandleOauthLogin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, h.provider.Endpoint(), http.StatusFound)
}

// HandleOauthVerify redirects the user to begin the oauth flow
func (h *Handler) HandleOauthVerify(w http.ResponseWriter, req *http.Request) {
	// Get the token using the oauth code
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {h.provider.Options().ClientID},
		"client_secret": {h.provider.Options().ClientSecret},
		"redirect_uri":  {h.provider.Redirect()},
		"code":          {req.FormValue("code")},
		"grant_type":    {"authorization_code"},
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
	resp, err = http.Get("https://www.googleapis.com/oauth2/v1/userinfo?oauth_token=" + oauthResult.Token)
	if err != nil || resp.StatusCode != http.StatusOK {
		http.Redirect(w, req, "/account/error", http.StatusFound)
		fmt.Println(err)
		return
	}

	// Decode the users profile
	var profile struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
	}
	json.NewDecoder(resp.Body).Decode(&profile)

	// Create the user in the users service
	uRsp, err := h.users.Create(req.Context(), &users.CreateRequest{
		User: &users.User{
			Id:        fmt.Sprintf("google_%v", profile.ID),
			Email:     profile.Email,
			FirstName: profile.FirstName,
			LastName:  profile.LastName,
		},
	})
	if err != nil {
		http.Redirect(w, req, "/account/error", http.StatusFound)
		fmt.Println(err)
		return
	}

	// Set the cookie and redirect
	http.SetCookie(w, &http.Cookie{
		Name:  auth.CookieName,
		Value: uRsp.Token,
	})
	http.Redirect(w, req, "/account", http.StatusFound)
}
