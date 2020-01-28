package handler

import (
	"context"

	auth "github.com/kytra-app/helpers/authentication"
	ledger "github.com/kytra-app/ledger-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	proto "github.com/kytra-app/registration-api/proto"
	user "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/client"
	"github.com/micro/go-micro/errors"
)

// Handler is an object can process RPC requests
type Handler struct {
	user       user.UsersService
	auth       auth.Authenticator
	ledger     ledger.LedgerService
	portfolios portfolios.PortfoliosService
}

// New returns an instance of Handler
func New(auth auth.Authenticator, client client.Client) Handler {
	return Handler{
		auth:       auth,
		user:       user.NewUsersService("kytra-srv-v1-users:8080", client),
		ledger:     ledger.NewLedgerService("kytra-srv-v1-ledger:8080", client),
		portfolios: portfolios.NewPortfoliosService("kytra-srv-v1-portfolios:8080", client),
	}
}

// Count returns the number of users registered in the user-srv
func (h Handler) Count(ctx context.Context, req *proto.CountRequest, rsp *proto.CountResponse) error {
	countRsp, err := h.user.Count(ctx, &user.CountRequest{})
	if err != nil {
		return err
	}

	rsp.Count = countRsp.Count
	return nil
}

// Signup creates a user object
func (h Handler) Signup(ctx context.Context, req *proto.User, rsp *proto.Response) error {
	u := user.User{
		Email:            req.Email,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Username:         req.Username,
		ProfilePictureId: req.ProfilePictureId,
		Password:         req.Password,
	}

	// Create the User
	user, err := h.user.Create(ctx, &u)
	if err != nil {
		return err
	}

	// Create the simulated portfolio
	portfolio, err := h.portfolios.Create(ctx, &portfolios.Portfolio{UserUuid: user.Uuid})
	if err != nil {
		return err
	}

	// Insert a 100k deposit into the simulated portfolio
	transaction := ledger.Transaction{
		PortfolioUuid: portfolio.Uuid,
		Amount:        100000 * 100,
		Type:          ledger.TransactionType_DEPOSIT,
	}
	if _, err := h.ledger.CreateTransaction(ctx, &transaction); err != nil {
		return err
	}

	// Generate the JWT
	token, err := h.auth.EncodeUser(auth.User{UUID: user.Uuid})
	if err != nil {
		return errors.InternalServerError("JWT", "An error occured generating the JWT")
	}

	// Serialize the response
	rsp.Jwt = token
	rsp.User = &proto.User{
		Uuid:             user.Uuid,
		Email:            user.Email,
		FirstName:        user.FirstName,
		LastName:         user.LastName,
		Username:         user.Username,
		ProfilePictureId: user.ProfilePictureId,
		Admin:            user.Admin,
	}

	return nil
}
