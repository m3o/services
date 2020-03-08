package main

import (
	"github.com/micro/go-micro/v2"
	log "github.com/micro/go-micro/v2/logger"

	"github.com/micro/services/distributed/api/handler/notes"
	"github.com/micro/services/distributed/api/handler/sprints"

	notesPb "github.com/micro/services/distributed/api/proto/notes"
	sprintsPb "github.com/micro/services/distributed/api/proto/sprints"
)

func main() {
	// New Service
	srv := micro.NewService(
		micro.Name("go.micro.api.distributed"),
		micro.Version("latest"),
	)

	// Initialise service
	srv.Init()

	// Register Handlers
	notesPb.RegisterDistributedNotesHandler(srv.Server(), notes.NewHandler(srv))
	sprintsPb.RegisterDistributedSprintsHandler(srv.Server(), sprints.NewHandler(srv))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
