package handler

import (
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/helpers/sms"
	ledger "github.com/kytra-app/ledger-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	smsVer "github.com/kytra-app/sms-verification-srv/proto"
	user "github.com/kytra-app/users-srv/proto"
	users "github.com/kytra-app/users-srv/proto"
	proto "github.com/kytra-app/welcome-api/proto"

	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	auth       auth.Authenticator
	photos     photos.Service
	sms        sms.Service
	user       user.UsersService
	ledger     ledger.LedgerService
	smsVer     smsVer.SMSVerificationService
	portfolios portfolios.PortfoliosService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, pics photos.Service, sms sms.Service, client client.Client) Handler {
	return Handler{
		auth:       auth,
		photos:     pics,
		sms:        sms,
		user:       user.NewUsersService("kytra-srv-v1-users:8080", client),
		ledger:     ledger.NewLedgerService("kytra-srv-v1-ledger:8080", client),
		smsVer:     smsVer.NewSMSVerificationService("kytra-srv-v1-sms-verification:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
	}
}

func (h Handler) serializeUser(user *users.User) (*proto.User, error) {
	// Generate the JWT
	token, err := h.auth.EncodeUser(auth.User{UUID: user.Uuid})
	if err != nil {
		err = errors.InternalServerError("JWT", "An error occured generating the JWT")
		return &proto.User{}, err
	}

	// Serialize the user
	res := proto.User{
		Jwt:               token,
		Uuid:              user.Uuid,
		Email:             user.Email,
		FirstName:         user.FirstName,
		LastName:          user.LastName,
		Username:          user.Username,
		ProfilePictureUrl: h.photos.GetURL(user.ProfilePictureId),
	}

	return &res, nil
}
