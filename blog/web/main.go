package main

import (
	"net/http"

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

	// register html handler
	service.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "https://micro.mu/blog/", 302)
	})

	// run service
        if err := service.Run(); err != nil {
                log.Fatal(err)
        }
}
