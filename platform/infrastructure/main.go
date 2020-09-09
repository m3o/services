package main

import (
	"github.com/robfig/cron"
	"github.com/slack-go/slack"

	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"

	pb "github.com/m3o/services/platform/infrastructure/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"
)

const (
	scalewayRegion = scw.RegionFrPar
	scalewayZone   = scw.ZoneFrPar1
)

var (
	// kubernetes api client
	k8sAPI *k8s.API
	// loadbalancer api client
	lbAPI *lb.API
	// instance api client
	inAPI *instance.API

	// slackbot client
	slackbot *slack.Client
)

func main() {
	// Create the service
	svr := service.New(
		service.Name("infrastructure"),
		service.Version("latest"),
	)

	// Create a Slack client
	slackToken := config.Get("micro", "alert", "slack_token").String("")
	if len(slackToken) == 0 {
		logger.Fatal("Missing required config micro.alert.slack_token")
	}
	slackbot = slack.New(slackToken)

	// Create a Scaleway client
	client, err := scw.NewClient(
		scw.WithDefaultOrganizationID(getConfig("org-id")),
		scw.WithAuth(getConfig("access-key"), getConfig("secret-key")),
	)
	if err != nil {
		logger.Fatalf("Error creating scaleway client: %v", err)
	}

	// Setup the clients
	lbAPI = lb.NewAPI(client)
	k8sAPI = k8s.NewAPI(client)
	inAPI = instance.NewAPI(client)

	// Check infra daily and report any wastage
	c := cron.New()
	c.AddFunc("0 9 * * *", checkInfraUsage)
	c.Start()

	checkInfraUsage()

	// Register the RPC handler
	pb.RegisterInfrastructureHandler(svr.Server(), new(infrastructure))

	// Run the service
	if err := svr.Run(); err != nil {
		logger.Fatal(err)
	}
}

func getConfig(key string) string {
	val := config.Get("micro", "infrastructure", "scaleway", key).String("")
	if len(val) == 0 {
		logger.Fatalf("Missing required config: micro.infrastructure.scaleway.%v", key)
	}
	return val
}
