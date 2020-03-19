package auth

import (
	"fmt"
	"net/http"

	"github.com/micro/go-micro/v2/auth"

	"github.com/micro/go-micro/v2/web"
	utils "github.com/micro/services/platform/web/util"
	users "github.com/micro/services/users/service/proto"
)

// RegisterHandlers adds the GitHub oauth handlers to the servie
func RegisterHandlers(srv web.Service) error {
	h := UserHandler{
		users: users.NewUsersService("go.micro.service.users", srv.Options().Service.Client()),
	}

	srv.HandleFunc("/v1/user", h.ReadUser)
	return nil
}

type UserHandler struct {
	users users.UsersService
}

type User struct {
	Name                  string `json:"name"`
	Email                 string `json:"email"`
	AvatarURL             string `json:"avatarURL"`
	TeamName              string `json:"teamName"`
	TeamURL               string `json:"teamURL"`
	OrganizationAvatarURL string `json:"organizationAvatarURL"`
	Login                 string `json:"login"`
}

func (h *UserHandler) ReadUser(w http.ResponseWriter, req *http.Request) {
	utils.SetupResponse(&w, req)
	if (*req).Method == "OPTIONS" {
		return
	}

	acc, err := auth.AccountFromContext(req.Context())
	if err != nil {
		utils.Write500(w, err)
		return
	}

	if acc.Metadata == nil {
		acc.Metadata = make(map[string]string)
	}

	uRsp, err := h.users.Read(req.Context(), &users.ReadRequest{Id: acc.Id})
	if err != nil {
		utils.Write500(w, err)
		return
	}

	utils.WriteJSON(w, &User{
		Name:                  fmt.Sprintf("%v %v", uRsp.User.FirstName, uRsp.User.LastName),
		Email:                 uRsp.User.Email,
		Login:                 uRsp.User.Username,
		AvatarURL:             uRsp.User.ProfilePictureUrl,
		TeamName:              "Community",
		TeamURL:               "https://github.com/orgs/micro/teams/community",
		OrganizationAvatarURL: "https://avatars3.githubusercontent.com/u/5161210?v=4",
	})
}
