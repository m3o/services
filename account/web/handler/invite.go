package handler

import (
	"net/http"
	"net/url"

	invites "github.com/micro/services/teams/invites/proto/invites"
)

// HandleInvite is the handler which gets called when a user clicks the
// link to join a team. The code is verified and then passed to the frontend.
func (h *Handler) HandleInvite(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	if len(code) == 0 {
		h.handleError(w, req, "Missing invite code")
		return
	}

	rsp, err := h.invites.Verify(req.Context(), &invites.VerifyRequest{Code: code})
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	params := url.Values{
		"inviteCode": {code},
		"teamName":   {rsp.TeamName},
		"email":      {rsp.Email},
	}
	http.Redirect(w, req, "/?"+params.Encode(), http.StatusFound)
}
