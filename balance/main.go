package main

import (
	"github.com/m3o/services/balance/handler"
	pb "github.com/m3o/services/balance/proto"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/util/opentelemetry/jaeger"
	"github.com/opentracing/opentracing-go"

	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/util/opentelemetry"
)

func main() {
	// Create service
	srv := service.New(
		service.Name("balance"),
		service.Version("latest"),
	)

	// Register handler
	pb.RegisterBalanceHandler(srv.Server(), handler.NewHandler(srv))

	c, _ := config.Get("jaegeraddress")

	openTracer, tracecloser, err := jaeger.New(
		opentelemetry.WithServiceName("balance"),
		opentelemetry.WithTraceReporterAddress(c.String("localhost:6831")),
	)
	defer tracecloser.Close()
	if err != nil {
		logger.Fatalf("Error configuring opentracing: %v", err)
	}
	// Set the global default opentracing tracer:
	opentracing.SetGlobalTracer(openTracer)
	opentelemetry.DefaultOpenTracer = openTracer
	logger.Infof("Configured jaeger to %s", c.String("localhost:6831"))

	// Run service
	if err := srv.Run(); err != nil {
		logger.Fatal(err)
	}
}
