package main

import (
	"net/http"
	"os"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	service := web.NewService(
		web.Name("go.micro.web.notes"),
		web.Version("latest"),
	)

	service.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		path := "./app/build" + req.URL.Path

		// 404 to index.html since the frontend does dynamic
		// route generation client side
		if _, err := os.Stat(path); err != nil {
			path = "./app/build/index.html"
		}

		log.Logf(log.InfoLevel, "Serving file: %v", path)
		http.ServeFile(w, req, path)
	})

	if err := service.Init(); err != nil {
		log.Fatal(err)
	}
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
