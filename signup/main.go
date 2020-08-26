package main

import (
	customersproto "github.com/m3o/services/customers/proto"
	inviteproto "github.com/m3o/services/invite/proto"
	paymentsproto "github.com/m3o/services/payments/provider/proto"
	plproto "github.com/m3o/services/platform/proto"
	"github.com/m3o/services/signup/handler"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	mauth "github.com/micro/micro/v3/service/auth/client"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("signup"),
	)

	// passing in auth because the DefaultAuth is the one used to set up the service
	auth := mauth.NewAuth()

	// Register Handler
	srv.Handle(handler.NewSignup(
		paymentsproto.NewProviderService("payment.stripe", srv.Client()),
		inviteproto.NewInviteService("invite", srv.Client()),
		plproto.NewPlatformService("platform", srv.Client()),
		customersproto.NewCustomersService("customers", srv.Client()),
		auth,
	))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
