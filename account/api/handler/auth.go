package handler

import (
	"context"
	"time"

	"github.com/micro/go-micro/v2/client"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/go-micro/v2/auth"
	pb "github.com/micro/services/account/api/proto/account"
	invite "github.com/micro/services/account/invite/proto"
	payment "github.com/micro/services/payments/provider/proto"
	teamInvite "github.com/micro/services/teams/invites/proto/invites"
	teams "github.com/micro/services/teams/service/proto/teams"
	users "github.com/micro/services/users/service/proto"
)

// Token generates a new JWT using a RefreshToken token
func (h *Handler) Token(ctx context.Context, req *pb.TokenRequest, rsp *pb.TokenResponse) error {
	tok, err := h.auth.Token(auth.WithExpiry(time.Hour*24), auth.WithToken(req.RefreshToken))
	if err != nil {
		return err
	}
	rsp.Token = serializeToken(tok)
	return nil
}

// Login looks up an account using an email and password
func (h *Handler) Login(ctx context.Context, req *pb.LoginRequest, rsp *pb.LoginResponse) error {
	// Generate a token
	tok, err := h.auth.Token(auth.WithCredentials(req.Email, req.Password), auth.WithExpiry(time.Hour*24))
	if err != nil {
		return err
	}

	// Lookup the user
	uRsp, err := h.users.Read(ctx, &users.ReadRequest{Email: req.Email}, client.WithServiceToken())
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = serializeToken(tok)

	// Fetch the payment methods
	pRsp, err := h.payment.ListPaymentMethods(ctx, &payment.ListPaymentMethodsRequest{UserId: uRsp.User.Id})
	if err != nil {
		log.Infof("Error listing payment methods: %v", err)
		return nil
	}

	// Serialize the payment methods
	rsp.User.PaymentMethods = make([]*pb.PaymentMethod, len(pRsp.PaymentMethods))
	for i, p := range pRsp.PaymentMethods {
		rsp.User.PaymentMethods[i] = serializePaymentMethod(p)
	}

	// Fetch the subscriptions
	sRsp, err := h.payment.ListSubscriptions(ctx, &payment.ListSubscriptionsRequest{UserId: uRsp.User.Id})
	if err != nil {
		log.Infof("Error listing subscriptions: %v", err)
		return nil
	}

	// Serialize the subscriptions
	rsp.User.Subscriptions = make([]*pb.Subscription, len(sRsp.Subscriptions))
	for i, s := range sRsp.Subscriptions {
		rsp.User.Subscriptions[i] = serializeSubscription(s)
	}

	return nil
}

// Signup creates an account using an email and password
func (h *Handler) Signup(ctx context.Context, req *pb.SignupRequest, rsp *pb.SignupResponse) error {
	// Validate the invite code
	var err error
	if req.TeamInvite {
		// the invite code is from a team
		_, err = h.teamInvite.Verify(ctx, &teamInvite.VerifyRequest{Code: req.InviteCode})
	} else {
		// the invite code is from micro
		_, err = h.invite.Validate(ctx, &invite.ValidateRequest{Code: req.InviteCode})
	}
	if err != nil {
		return err
	}

	// Validate the user can be created
	_, err = h.users.Create(ctx, &users.CreateRequest{
		User:         &users.User{Email: req.Email},
		ValidateOnly: true,
	})
	if err != nil {
		return err
	}

	// Generate an account and token
	acc, err := h.auth.Generate(req.Email, auth.WithRoles("user"), auth.WithSecret(req.Password))
	if err != nil {
		return err
	}
	tok, err := h.auth.Token(auth.WithCredentials(acc.ID, acc.Secret), auth.WithExpiry(time.Hour*24))
	if err != nil {
		return err
	}

	// Create the user
	u := &users.User{Email: req.Email, InviteVerified: true}
	uRsp, err := h.users.Create(ctx, &users.CreateRequest{User: u})
	if err != nil {
		return err
	}

	// Assign the user to the team if they were invited
	if req.TeamInvite {
		if _, err = h.teamInvite.Redeem(ctx, &teamInvite.RedeemRequest{Code: req.InviteCode, UserId: uRsp.User.Id}); err != nil {
			return err
		}
	}

	// Get the users teams
	tRsp, err := h.teams.ListMemberships(ctx, &teams.ListMembershipsRequest{MemberId: uRsp.User.Id})
	if err != nil {
		return err
	}

	// Serialize the response
	rsp.User = serializeUser(uRsp.User)
	rsp.Token = serializeToken(tok)
	rsp.Teams = make([]*pb.Team, 0, len(tRsp.Teams))
	for _, t := range tRsp.Teams {
		rsp.Teams = append(rsp.Teams, serializeTeam(t))
	}

	return nil
}

func serializeTeam(t *teams.Team) *pb.Team {
	return &pb.Team{
		Id:        t.Id,
		Name:      t.Name,
		Namespace: t.Namespace,
	}
}
