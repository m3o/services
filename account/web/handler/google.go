package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	users "github.com/micro/services/users/service/proto"
)

// HandleGoogleOauthLogin redirects the user to begin the oauth flow
func (h *Handler) HandleGoogleOauthLogin(w http.ResponseWriter, req *http.Request) {
	http.Redirect(w, req, h.google.Endpoint(), http.StatusFound)
}

// HandleGoogleOauthVerify redirects the user to begin the oauth flow
func (h *Handler) HandleGoogleOauthVerify(w http.ResponseWriter, req *http.Request) {
	// Get the token using the oauth code
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {h.google.Options().ClientID},
		"client_secret": {h.google.Options().ClientSecret},
		"redirect_uri":  {h.google.Redirect()},
		"code":          {req.FormValue("code")},
		"grant_type":    {"authorization_code"},
	})
	if err != nil {
		h.handleError(w, req, "Error getting access token from Google: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.handleError(w, req, "Error getting access token from Google. Status: %v", resp.Status)
		return
	}

	// Decode the token
	var oauthResult struct {
		Token string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&oauthResult)

	// Use the token to get the users profile
	resp, err = http.Get("https://www.googleapis.com/oauth2/v1/userinfo?oauth_token=" + oauthResult.Token)
	if err != nil {
		h.handleError(w, req, "Error getting account from Google: %v", err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		h.handleError(w, req, "Error getting account from Google. Status: %v", resp.Status)
		return
	}

	// Decode the users profile
	var profile struct {
		ID        string `json:"id"`
		Email     string `json:"email"`
		FirstName string `json:"given_name"`
		LastName  string `json:"family_name"`
		Picture   string `json:"picture"`
	}
	json.NewDecoder(resp.Body).Decode(&profile)
	fmt.Println(profile)

	// Create the user in the users service
	uRsp, err := h.users.Create(req.Context(), &users.CreateRequest{
		User: &users.User{
			Id:                fmt.Sprintf("google_%v", profile.ID),
			Email:             profile.Email,
			FirstName:         profile.FirstName,
			LastName:          profile.LastName,
			ProfilePictureUrl: profile.Picture,
		},
	})
	if err != nil {
		h.handleError(w, req, "Error creating user account: %v", err)
		return
	}

	var roles []string
	if strings.HasSuffix(profile.Email, "@micro.mu") {
		roles = append(roles, "admin", "developer", "collaborator")
	}
	h.loginUser(w, req, uRsp.User, roles...)
}
