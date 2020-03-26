package main

import (
	"log"

	"github.com/micro/go-micro/v2/web"
)

//go:generate ./hugo_build.sh

func main() {
	service := web.NewService(
		web.Name("go.micro.web.explore"),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
