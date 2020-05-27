package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2/runtime"
)

const (
	// serviceRootFile is the file which a directory must contain to be considered a service
	serviceRootFile = "main.go"
	// source is the repository containing the source code
	source = "github.com/micro/services"
)

func main() {
	srv := micro.NewService()
	srv.Init()
	logger.Infof("Using %v runtime", srv.Options().Runtime)

	// setup an admin account for the service to use (required to run services in a custom namespace)
	// this is a temporaty solution until identity is setup then we'll need to pass a set of creds
	// as arguments.
	name := fmt.Sprintf("bootstrap-%v", srv.Options().Server.Options().Id)
	acc, err := srv.Options().Auth.Generate(name, auth.WithScopes("admin"))
	if err != nil {
		logger.Fatal(err)
	}
	tok, err := srv.Options().Auth.Token(auth.WithCredentials(acc.ID, acc.Secret))
	if err != nil {
		logger.Fatal(err)
	}
	srv.Options().Auth.Init(auth.ClientToken(tok))

	for _, name := range listServices() {
		logger.Infof("Creating %v", name)

		err := srv.Options().Runtime.Create(
			&runtime.Service{Name: name, Source: source},
			runtime.CreateNamespace("platform"),
			runtime.CreateImage("docker.pkg.github.com/micro/services/"+strings.ReplaceAll(name, "/", "-")),
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
