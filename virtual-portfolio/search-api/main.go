package main

import (
	"fmt"
	"os"

	auth "github.com/kytra-app/helpers/authentication"
	photos "github.com/kytra-app/helpers/photos"
	"github.com/kytra-app/search-api/handler"
	proto "github.com/kytra-app/search-api/proto"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-api-v1-search"),
		micro.Version("latest"),
	)
	service.Init()

	pics, err := photos.New(os.Getenv("PHOTOS_ADDRESS"))
	if err != nil {
		fmt.Printf("Could not initiate photos package: %v\n.", err)
		os.Exit(2)
	}

	auth, err := auth.New(os.Getenv("JWT_SIGNING_KEY"))
	if err != nil {
		fmt.Printf("Could not initiate auth package: %v\n", err)
		os.Exit(2)
	}


	handler := handler.New(service.Client(), auth, pics)
	proto.RegisterSearchHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
