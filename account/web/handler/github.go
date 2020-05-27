package handler

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/auth/provider"
	invite "github.com/micro/services/projects/invite/proto"
	users "github.com/micro/services/users/service/proto"
)

// HandleGithubOauthLogin redirects the user to begin the oauth flow
func (h *Handler) HandleGithubOauthLogin(w http.ResponseWriter, req *http.Request) {
	state, err := h.generateOauthState()
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// invite code is present if the user was redirected from a team invite
	inviteCode := req.URL.Query().Get("inviteCode")

	// record the invite code
	if len(inviteCode) > 0 {
		h.setInviteCode(state, inviteCode)
	}

	endpoint := h.github.Endpoint(provider.WithState(state))
	http.Redirect(w, req, endpoint, http.StatusFound)
}

// HandleGithubOauthVerify redirects the user to begin the oauth flow
func (h *Handler) HandleGithubOauthVerify(w http.ResponseWriter, req *http.Request) {
	// validate the oauth state
	valid, err := h.validateOauthState(req.FormValue("state"))
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	} else if !valid {
		h.handleError(w, req, "Invalid Oauth State")
		return
	}

	// perform the oauth call to exchange the code for an access token
	token, err := h.getGithubAccessToken(req.FormValue("code"))
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// get the users profile
	profile, err := h.getGithubProfile(token)
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// check to see if the user already exists
	if _, err := h.users.Read(req.Context(), &users.ReadRequest{Email: profile.Email}); err == nil {
		// user already exists, get the secret for their account
		secret, err := h.getAccountSecret(profile.Email)
		if err != nil {
			h.handleError(w, req, "Error getting auth secret: %v", err)
			return
		}

		// create a token
		tok, err := h.auth.Token(auth.WithCredentials(profile.Email, secret), auth.WithExpiry(time.Hour*24))
		if err != nil {
			h.handleError(w, req, err.Error())
			return
		}

		// Login the user (set the cookie and return)
		h.loginUser(w, req, tok)
		return
	}

	// Create the user in the users service
	uRsp, err := h.users.Create(req.Context(), &users.CreateRequest{
		User: &users.User{Email: profile.Email, ProfilePictureUrl: profile.Picture},
	})
	if err != nil {
		h.handleError(w, req, "Error creating account: %v", err)
		return
	}

	// Create an auth account
	acc, err := h.auth.Generate(profile.Email, auth.WithType("user"), auth.WithProvider("oauth/github"))
	if err != nil {
		h.handleError(w, req, "Error creating auth account: %v", err)
		return
	}
	if err := h.setAccountSecret(acc.ID, acc.Secret); err != nil {
		h.handleError(w, req, "Error storing auth secret: %v", err)
		return
	}

	// Generate a token
	tok, err := h.auth.Token(auth.WithCredentials(profile.Email, acc.Secret), auth.WithExpiry(time.Hour*24))
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	// Check to see if the user had an invite token, if they did, activate it
	if inv, err := h.getInviteCode(req.FormValue("state")); err == nil {
		// redeem the invite, adding the user to the team
		h.invite.Redeem(req.Context(), &invite.RedeemRequest{Code: inv, UserId: uRsp.User.Id})

		// update the user so the users account does not require
		// an invite token
		uRsp.User.InviteVerified = true
		h.users.Update(req.Context(), &users.UpdateRequest{User: uRsp.User})
	}

	// Login the user
	h.loginUser(w, req, tok)
}

func (h *Handler) getGithubAccessToken(code string) (string, error) {
	// Consturct the requerst to get the access token
	data := url.Values{
		"client_id":     {h.github.Options().ClientID},
		"client_secret": {h.github.Options().ClientSecret},
		"redirect_uri":  {h.github.Redirect()},
		"code":          {code},
	}
	r, _ := http.NewRequest("POST", "https://github.com/login/oauth/access_token", strings.NewReader(data.Encode()))
	r.Header.Add("Accept", "application/json")
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}

	// Execute the request
	resp, err := client.Do(r)
	if err != nil {
		return "", fmt.Errorf("Error getting access token from GitHub: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Error getting access token from GitHub. Status: %v", resp.Status)
	}

	// Decode the token
	var result struct {
		Token string `json:"access_token"`
	}
	json.NewDecoder(resp.Body).Decode(&result)
	return result.Token, nil
}

type githubProfile struct {
	ID       int64  `json:"id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Username string `json:"login"`
	Picture  string `json:"avatar_url"`
}

func (h *Handler) getGithubProfile(token string) (*githubProfile, error) {
	// Use the token to get the users profile
	client := &http.Client{}
	r, _ := http.NewRequest("GET", "https://api.github.com/user", nil)
	r.Header.Add("Authorization", "Bearer "+token)
	resp, err := client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Error getting user from GitHub: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error getting user from GitHub. Status: %v. Error: %v", resp.Status, string(bytes))
	}

	// Decode the users profile
	var profile *githubProfile
	json.NewDecoder(resp.Body).Decode(&profile)

	// Get the emails
	r, _ = http.NewRequest("GET", "https://api.github.com/user/emails", nil)
	r.Header.Add("Authorization", "Bearer "+token)
	resp, err = client.Do(r)
	if err != nil {
		return nil, fmt.Errorf("Error getting user from GitHub: %v", err)
	}
	if resp.StatusCode != http.StatusOK {
		bytes, _ := ioutil.ReadAll(resp.Body)
		return nil, fmt.Errorf("Error getting user emails from GitHub. Status: %v. Error: %v", resp.Status, string(bytes))
	}

	// Decode the response
	var emailsRsp []struct {
		Email   string `json:"email"`
		Primary bool   `json:"primary"`
	}
	json.NewDecoder(resp.Body).Decode(&emailsRsp)
	for _, r := range emailsRsp {
		if r.Primary == true {
			profile.Email = r.Email
		}
	}

	// Ensure an email was found
	if len(profile.Email) == 0 {
		return nil, fmt.Errorf("No primary email was found: %v", emailsRsp)
	}

	return profile, err
}
