package main

import (
	log "github.com/micro/go-micro/v2/logger"
	"github.com/micro/go-micro/v2"
	"scheduler/manager"
)

var (
	WorkflowID = "599918"
)

func main() {
	// New Service
	service := micro.NewService(
		micro.Name("go.micro.platform.scheduler"),
	)

	// Initialise service
	service.Init()

	// start the scheduler
	manager.Start(WorkflowID)

	// Run service
	if err := service.Run(); err != nil {
		log.Fatal(err)
	}
}
