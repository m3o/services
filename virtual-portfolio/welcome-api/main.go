package main

import (
	"fmt"
	"os"

	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/helpers/sms"
	"github.com/kytra-app/welcome-api/handler"
	proto "github.com/kytra-app/welcome-api/proto"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-api-v1-welcome"),
		micro.Version("latest"),
	)
	service.Init()

	auth, err := auth.New(os.Getenv("JWT_SIGNING_KEY"))
	if err != nil {
		fmt.Printf("Could not initiate auth package: %v\n", err)
		os.Exit(2)
	}

	// Connect to SMS service
	sms, err := sms.New(
		os.Getenv("SMS_ACCOUNT_SID"),
		os.Getenv("SMS_AUTH_TOKEN"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the SMS service"))
	}

	pics, err := photos.New(os.Getenv("PHOTOS_ADDRESS"))
	if err != nil {
		fmt.Printf("Could not initiate photos package: %v\n.", err)
		os.Exit(2)
	}

	handler := handler.New(auth, pics, sms, service.Client())
	proto.RegisterWelcomeHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
