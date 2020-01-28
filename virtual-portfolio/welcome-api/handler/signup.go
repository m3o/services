package handler

import (
	"context"

	ledger "github.com/kytra-app/ledger-srv/proto"
	portfolios "github.com/kytra-app/portfolios-srv/proto"
	smsVer "github.com/kytra-app/sms-verification-srv/proto"
	user "github.com/kytra-app/users-srv/proto"
	"github.com/micro/go-micro/errors"
	proto "github.com/kytra-app/welcome-api/proto"
)

// CreateProfile creates a user object
func (h Handler) CreateProfile(ctx context.Context, req *proto.User, rsp *proto.User) error {
	// Fetch the JWT code
	if req.VerificationUuid == "" {
		return errors.BadRequest("INVALID_CODE", "A valid verification uuid is requiried")
	}
	ver, err := h.smsVer.Get(ctx, &smsVer.Verification{Uuid: req.VerificationUuid})
	if err != nil {
		return err
	}

	// Upload the profile picture
	profilePictureID, err := h.profilePictureIDFromBase64(req.ProfilePictureBase64)
	if err != nil {
		return err
	}

	// Create the User
	user, err := h.user.Create(ctx, &user.User{
		Email:            req.Email,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		Username:         req.Username,
		Password:         req.Password,
		PhoneNumber:      ver.PhoneNumber,
		ProfilePictureId: profilePictureID,
	})
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

	// Serialize the user and generate the JWT
	serializedUser, err := h.serializeUser(user)
	if err != nil {
		return err
	}

	*rsp = *serializedUser
	return nil
}

func (h Handler) profilePictureIDFromBase64(base64 string) (string, error) {
	if base64 == "" {
		return "", nil
	}

	return h.photos.Upload(base64)
}
