package handler

import (
	"net/http"
	"net/url"

	invite "github.com/micro/services/project/invite/proto"
)

// HandleInvite is the handler which gets called when a user clicks the
// link to join a project. The code is verified and then passed to the frontend.
func (h *Handler) HandleInvite(w http.ResponseWriter, req *http.Request) {
	code := req.URL.Query().Get("code")
	if len(code) == 0 {
		h.handleError(w, req, "Missing invite code")
		return
	}

	rsp, err := h.invite.Verify(req.Context(), &invite.VerifyRequest{Code: code})
	if err != nil {
		h.handleError(w, req, err.Error())
		return
	}

	params := url.Values{
		"inviteCode":  {code},
		"projectName": {rsp.ProjectName},
		"email":       {rsp.Email},
	}
	http.Redirect(w, req, "/?"+params.Encode(), http.StatusFound)
}
