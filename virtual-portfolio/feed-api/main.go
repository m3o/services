package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/feed-api/handler"
	proto "github.com/kytra-app/feed-api/proto"
	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/helpers/photos"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra.api.v1.feed"),
		micro.Version("latest"),
	)
	service.Init()

	auth, err := auth.New(os.Getenv("JWT_SIGNING_KEY"))
	if err != nil {
		fmt.Printf("Could not initiate auth package: %v\n", err)
		os.Exit(2)
	}

	pics, err := photos.New(os.Getenv("PHOTOS_ADDRESS"))
	if err != nil {
		fmt.Printf("Could not initiate photos package: %v\n.", err)
		os.Exit(2)
	}

	handler := handler.New(auth, pics, service.Client())
	proto.RegisterFeedHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
