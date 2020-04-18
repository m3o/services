package main

import (
	"net/http"
	"os"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	// create new web service
	service := web.NewService(
		web.Name("go.micro.web.m3o"),
	)

	// initialise service
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// Serve the web app
	service.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		path := "./app/build" + req.URL.Path

		// 404 to index.html since the frontend does dynamic
		// route generation client side
		if _, err := os.Stat(path); err != nil {
			path = "./app/build/index.html"
		}

		http.ServeFile(w, req, path)
	})

	// run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
