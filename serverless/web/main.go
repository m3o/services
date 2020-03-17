package main

import (
	"log"
	"net/http"
	"strings"

	"serverless/web/handler/auth"

	"github.com/micro/go-micro/v2/web"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.serverless"),
	)

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	auth.RegisterHandlers(service)

	service.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		// Count is an ugly fix to serve urls containing micro service names ie. "go.micro.something"
		if strings.Contains(req.URL.Path, ".") && !strings.Contains(req.URL.Path, "go.micro") {
			http.ServeFile(w, req, "./app/dist/micro/"+req.URL.Path[1:])
			return
		}
		http.ServeFile(w, req, "./app/dist/micro/index.html")
	})
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
