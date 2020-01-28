package main

import (
	"fmt"
	"os"

	auth "github.com/kytra-app/helpers/authentication"
	"github.com/kytra-app/post-enhancer-srv/handler"
	proto "github.com/kytra-app/post-enhancer-srv/proto"
	_ "github.com/micro/go-plugins/registry/kubernetes"

	"github.com/micro/go-micro"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-srv-v1-post-enhancer"),
		micro.Version("latest"),
	)
	service.Init()

	auth, err := auth.New(os.Getenv("JWT_SIGNING_KEY"))
	if err != nil {
		fmt.Printf("Could not initiate auth package: %v\n", err)
		os.Exit(2)
	}

	handler := handler.New(auth, service.Client())
	proto.RegisterPostEnhancerHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
