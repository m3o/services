package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
	"web/handler"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Version("latest"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// register html handler
	service.HandleFunc("/", handler.Go)

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
