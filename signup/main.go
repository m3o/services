package main

import (
	"github.com/m3o/services/signup/handler"
	log "github.com/micro/go-micro/v3/logger"

	inviteproto "github.com/m3o/services/invite/proto"
	k8sproto "github.com/m3o/services/kubernetes/proto"
	paymentsproto "github.com/m3o/services/payments/provider/proto"
	signup "github.com/m3o/services/signup/proto/signup"
	"github.com/micro/micro/v3/service"
	mauth "github.com/micro/micro/v3/service/auth/client"
)

func main() {
	// New Service
	srv := service.New(
		service.Name("go.micro.service.signup"),
	)

	// passing in auth because the DefaultAuth is the one used to set up the service
	auth := mauth.NewAuth()

	// Register Handler
	signup.RegisterSignupHandler(handler.NewSignup(
		paymentsproto.NewProviderService("go.micro.service.payment.stripe"),
		inviteproto.NewInviteService("go.micro.service.invite"),
		k8sproto.NewKubernetesService("go.micro.service.kubernetes"),
		auth,
	))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
