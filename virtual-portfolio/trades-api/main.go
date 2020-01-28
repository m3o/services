package main

import (
	"fmt"
	"os"

	auth "github.com/kytra-app/helpers/authentication"
	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/trades-api/handler"
	proto "github.com/kytra-app/trades-api/proto"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-api-v1-trades"),
		micro.Version("latest"),
	)
	service.Init()

	auth, err := auth.New(os.Getenv("JWT_SIGNING_KEY"))
	if err != nil {
		fmt.Printf("Could not initiate auth package: %v\n", err)
		os.Exit(2)
	}

	iex, err := iex.New(os.Getenv("IEX_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate iex package: %v\n.", err)
		os.Exit(2)
	}

	handler := handler.New(auth, iex, service.Client())
	proto.RegisterTradesHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
