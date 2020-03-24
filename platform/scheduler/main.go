package main

import (
	"scheduler/manager"

	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"
)

var (
	WorkflowFilename = "build.yml"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.platform.scheduler"),
	)

	// Initialise service
	service.Init()

	// start the scheduler
	if err := manager.Start(WorkflowFilename); err != nil {
		log.Fatal(err)
	}

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
