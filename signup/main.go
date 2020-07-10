package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/services/signup/handler"

	inviteproto "github.com/micro/services/account/invite/proto"
	k8sproto "github.com/micro/services/kubernetes/service/proto"
	paymentsproto "github.com/micro/services/payments/provider/proto"
	signup "github.com/micro/services/signup/proto/signup"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.service.signup"),
		micro.Version("latest"),
	)

	// Initialise service
	service.Init()

	// Register Handler
	signup.RegisterSignupHandler(service.Server(), handler.NewSignup(
		paymentsproto.NewProviderService("go.micro.service.payment.stripe", service.Options().Client),
		inviteproto.NewInviteService("go.micro.service.account.invite", service.Options().Client),
		k8sproto.NewKubernetesService("go.micro.service.kubernetes", service.Options().Client),
		service.Options().Store,
		service.Options().Config,
		service.Options().Auth,
	))

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
