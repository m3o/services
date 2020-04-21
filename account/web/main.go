package main

import (
	"net/http"
	"os"

	"github.com/micro/services/account/web/handler"

	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/web"
)

func main() {
	// New Service
	service := web.NewService(
		web.Name("go.micro.web.account"),
		web.Version("latest"),
	)

	// Load the config (needed for the auth provider)
	if err := service.Init(); err != nil {
		log.Fatal(err)
	}

	// Handle the redirect from google when oauth completes
	h := handler.NewHandler(service.Options().Service)
	service.HandleFunc("/oauth/google/login", h.HandleGoogleOauthLogin)
	service.HandleFunc("/oauth/google/verify", h.HandleGoogleOauthVerify)
	service.HandleFunc("/oauth/github/login", h.HandleGithubOauthLogin)
	service.HandleFunc("/oauth/github/verify", h.HandleGithubOauthVerify)
	service.HandleFunc("/invite", h.HandleInvite)

	// Serve the web app
	service.HandleFunc("/", func(w http.ResponseWriter, req *http.Request) {
		path := "./app/build" + req.URL.Path

		// 404 to index.html since the frontend does dynamic
		// route generation client side
		if _, err := os.Stat(path); err != nil {
			path = "./app/build/index.html"
		}

		log.Logf(log.TraceLevel, "Serving file: %v", path)
		http.ServeFile(w, req, path)
	})

	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
