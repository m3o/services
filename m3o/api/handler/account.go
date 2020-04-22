package handler

import (
	"context"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"

	pb "github.com/micro/services/m3o/api/proto"
	users "github.com/micro/services/users/service/proto"
)

// Account implments the M3O account service proto
type Account struct {
	name  string
	users users.UsersService
}

// NewAccount returns an initialised account handler
func NewAccount(service micro.Service) *Account {
	return &Account{
		name:  service.Name(),
		users: users.NewUsersService("go.micro.service.users", service.Client()),
	}
}

// Read the current users info
func (a *Account) Read(ctx context.Context, req *pb.ReadAccountRequest, rsp *pb.ReadAccountResponse) error {
	acc, err := auth.AccountFromContext(ctx)
	if err != nil {
		return err
	}

	uRsp, err := a.users.Read(ctx, &users.ReadRequest{Email: acc.ID})
	if err != nil {
		return err
	}

	rsp.User = serializeUser(uRsp.User)
	return nil
}

func serializeUser(u *users.User) *pb.User {
	return &pb.User{
		Id:                u.Id,
		FirstName:         u.FirstName,
		LastName:          u.LastName,
		Email:             u.Email,
		ProfilePictureUrl: u.ProfilePictureUrl,
	}
}
