package main

import (
	"github.com/micro/services/oauth/github/handler"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.github"),
		web.Version("latest"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// register handler
	handler.RegisterHandler(service)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
