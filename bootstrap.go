package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v3/auth"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/runtime"
	"github.com/micro/micro/v3/service"
	mauth "github.com/micro/micro/v3/service/auth"
	mruntime "github.com/micro/micro/v3/service/runtime"
)

const (
	// serviceRootFile is the file which a directory must contain to be considered a service
	serviceRootFile = "main.go"
	// source is the repository containing the source code
	source = "github.com/m3o/services"
)

func main() {
	srv := service.New()
	logger.Infof("Using %v runtime", mruntime.DefaultRuntime)

	// setup an admin account for the service to use (required to run services in a custom namespace)
	// this is a temporaty solution until identity is setup then we'll need to pass a set of creds
	// as arguments.
	name := fmt.Sprintf("bootstrap-%v", srv.Server().Options().Id)
	acc, err := mauth.Generate(name, auth.WithScopes("admin"))
	if err != nil {
		logger.Fatal(err)
	}
	tok, err := mauth.Token(auth.WithCredentials(acc.ID, acc.Secret))
	if err != nil {
		logger.Fatal(err)
	}
	mauth.DefaultAuth.Init(auth.ClientToken(tok))

	for _, name := range listServices() {
		logger.Infof("Creating %v", name)

		err := mruntime.Create(
			&runtime.Service{Name: name, Source: source},
			runtime.CreateNamespace("platform"),
			runtime.CreateImage("docker.pkg.github.com/m3o/services/"+strings.ReplaceAll(name, "/", "-")),
		)

		if err != nil {
			logger.Fatal(err)
		}
	}
}

func listServices() []string {
	var services []string

	filepath.Walk(".", func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if filepath.Base(path) != serviceRootFile {
			return nil
		}
		services = append(services, filepath.Dir(path))
		return nil
	})

	return services
}
