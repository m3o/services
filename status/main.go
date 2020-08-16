package main

import (
	"github.com/m3o/services/metrics/prometheus"
	"github.com/m3o/services/metrics/wrapper"
	"github.com/m3o/services/status/handler"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
)

func main() {

	// Make a metrics.Reporter (should take a *service.Service but that is circular until we can embed metrics into the service code):
	reporter, err := prometheus.New(nil)
	if err != nil {
		log.Errorf("Error starting Prometheus metrics reporter")
	}

	metricsWrapper := wrapper.New(reporter)

	// New Service
	srv := service.New(
		service.Name("status"),
		service.Version("0.0.1"),
		service.WrapHandler(metricsWrapper.HandlerFunc),
	)

	// An example metric:
	reporter.Count("service_start", 1, nil)

	// grab services to monitor
	svcs := config.Get("micro", "status", "services").StringSlice(nil)
	log.Infof("Services to monitor %+v", svcs)

	// Register Handler
	srv.Handle(handler.NewStatusHandler(svcs))

	// Run service
	if err := srv.Run(); err != nil {
		log.Fatal(err)
	}
}
