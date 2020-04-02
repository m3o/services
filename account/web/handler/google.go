package handler

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strings"

	"github.com/micro/go-micro/v2/auth/provider"

	"github.com/micro/go-micro/v2/auth"
	users "github.com/micro/services/users/service/proto"
)

// HandleGoogleOauthLogin redirects the user to begin the oauth flow
func (h *Handler) HandleGoogleOauthLogin(w http.ResponseWriter, req *http.Request) {
	code := h.generateOauthState()
	http.Redirect(w, req, h.google.Endpoint(provider.WithState(code)), http.StatusFound)
}

// HandleGoogleOauthVerify redirects the user to begin the oauth flow
func (h *Handler) HandleGoogleOauthVerify(w http.ResponseWriter, req *http.Request) {
	// validate the oauth state
	if valid := h.validateOauthState(req.FormValue("state")); !valid {
		h.handleError(w, req, "Invalid Oauth State")
		return
	}

	// perform the oauth call to exchange the code for an access token
	token, err := h.getGoogleAccessToken(req.FormValue("code"))
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// Get the users profile
	profile, err := h.getGoogleProfile(token)
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// check to see if the user already exists
	if _, err := h.users.Read(req.Context(), &users.ReadRequest{Email: profile.Email}); err == nil {
		// user already exists, get the secret for their account
		secret, err := h.getAccountSecret(profile.Email)
		if err != nil {
			h.handleError(w, req, "Error storing auth secret: %v", err)
			return
		}

		// create a token
		tok, err := h.auth.Token(auth.WithCredentials(profile.Email, secret))
		if err != nil {
			h.handleError(w, req, err.Error())
			return
		}

		// Login the user (set the cookie and return)
		h.loginUser(w, req, tok)
		return
	}

	// Create the user in the users service
	_, err = h.users.Create(req.Context(), &users.CreateRequest{
		User: &users.User{
			Email:             profile.Email,
			FirstName:         profile.FirstName,
			LastName:          profile.LastName,
			ProfilePictureUrl: profile.Picture,
		},
	})
	if err != nil {
		h.handleError(w, req, "Error creating account: %v", err)
		return
	}

	// Setup the roles
	roles := []string{"user", "user.developer"}
	if strings.HasSuffix(profile.Email, "@micro.mu") {
		roles = append(roles, "admin", "user.collaborator")
	}

	// Create an auth account
	acc, err := h.auth.Generate(profile.Email, auth.WithRoles(roles...), auth.WithProvider("oauth/google"))
	if err != nil {
		h.handleError(w, req, "Error creating auth account: %v", err)
		return
	}
	if err := h.setAccountSecret(acc.ID, acc.Secret); err != nil {
		h.handleError(w, req, "Error storing auth secret: %v", err)
		return
	}

	// Generate a token
	tok, err := h.auth.Token(auth.WithCredentials(profile.Email, acc.Secret))
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// Login the user
	h.loginUser(w, req, tok)
}

func (h *Handler) getGoogleAccessToken(code string) (string, error) {
	// Get the token using the oauth code
	resp, err := http.PostForm("https://oauth2.googleapis.com/token", url.Values{
		"client_id":     {h.google.Options().ClientID},
		"client_secret": {h.google.Options().ClientSecret},
		"redirect_uri":  {h.google.Redirect()},
		"grant_type":    {"authorization_code"},
		"code":          {code},
	})
	if err != nil {
		return "", fmt.Errorf("Error getting access token from Google: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error getting access token from Google. Status: %v", resp.Status)
	}

	// Decode the token
	var result struct {
		Token string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Token, nil
}

type googleProfile struct {
	ID        string `json:"id"`
	Email     string `json:"email"`
	FirstName string `json:"given_name"`
	LastName  string `json:"family_name"`
	Picture   string `json:"picture"`
}

func (h *Handler) getGoogleProfile(token string) (*googleProfile, error) {
	resp, err := http.Get("https://www.googleapis.com/oauth2/v1/userinfo?oauth_token=" + token)
	if err != nil {
		return nil, fmt.Errorf("Error getting account from Google: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Error getting account from Google. Status: %v", resp.Status)
	}

	// Decode the users profile
	var profile *googleProfile
	json.NewDecoder(resp.Body).Decode(&profile)
	return profile, nil
}
