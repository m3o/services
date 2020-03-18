package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/micro/go-micro/v2/logger"

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
		logger.Errorf("Error fetching google account: %v", err)
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
		logger.Errorf("Error creating user account: %v", err)
		return
	}

	var roles []string
	if strings.HasSuffix(profile.Email, "@micro.mu") {
		roles = append(roles, "admin", "developer", "collaborator")
	}
	h.loginUser(w, req, uRsp.User, roles...)
}
