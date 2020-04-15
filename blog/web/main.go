package main

import (
        log "github.com/micro/go-micro/v2/logger"
        "github.com/micro/go-micro/v2/web"
)

func main() {
	// create new web service
        service := web.NewService(
                web.Name("go.micro.web.blog"),
        )

	// initialise service
        if err := service.Init(); err != nil {
                log.Fatal(err)
        }

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
