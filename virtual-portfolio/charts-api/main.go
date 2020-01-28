package main

import (
	"fmt"
	"os"

	iex "github.com/kytra-app/helpers/iex-cloud"
	"github.com/kytra-app/charts-api/handler"
	proto "github.com/kytra-app/charts-api/proto"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
)

func main() {
	service := micro.NewService(
		micro.Name("kytra-api-v1-charts"),
		micro.Version("latest"),
	)
	service.Init()

	iex, err := iex.New(os.Getenv("IEX_TOKEN"))
	if err != nil {
		fmt.Printf("Could not initiate iex package: %v\n.", err)
		os.Exit(2)
	}

	handler := handler.New(iex, service.Client())
	proto.RegisterChartsHandler(service.Server(), handler)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
