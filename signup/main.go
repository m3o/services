package main

import (
	aproto "github.com/m3o/services/alert/proto/alert"
	customersproto "github.com/m3o/services/customers/proto"
	inviteproto "github.com/m3o/services/invite/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	pproto "github.com/m3o/services/payments/provider/proto"
	"github.com/m3o/services/signup/handler"
	subproto "github.com/m3o/services/subscriptions/proto"
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
		inviteproto.NewInviteService("invite", srv.Client()),
		customersproto.NewCustomersService("customers", srv.Client()),
		nsproto.NewNamespacesService("namespaces", srv.Client()),
		subproto.NewSubscriptionsService("subscriptions", srv.Client()),
		auth,
		pproto.NewProviderService("payment.stripe", srv.Client()),
		aproto.NewAlertService("alert", srv.Client()),
	))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
