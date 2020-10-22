package main

import (
	"github.com/robfig/cron"
	"github.com/slack-go/slack"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/scaleway/scaleway-sdk-go/api/instance/v1"
	k8s "github.com/scaleway/scaleway-sdk-go/api/k8s/v1beta4"
	"github.com/scaleway/scaleway-sdk-go/api/lb/v1"
	"github.com/scaleway/scaleway-sdk-go/scw"

	pb "github.com/m3o/services/infrastructure/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	storeproto "github.com/micro/micro/v3/proto/store"
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

	s3Client  *minio.Client
	nsService nsproto.NamespacesService
	stService storeproto.StoreService
)

func main() {
	// Create the service
	svr := service.New(
		service.Name("infrastructure"),
		service.Version("latest"),
	)

	// Create a Slack client
	val, err := config.Get("micro.alert.slack_token")
	if err != nil {
		logger.Warnf("Error getting config: %v", err)
	}
	slackToken := val.String("")
	if len(slackToken) == 0 {
		logger.Fatal("Missing required config micro.alert.slack_token")
	}
	slackbot = slack.New(slackToken)

	accessKey := getConfig("access-key")
	secretKey := getConfig("secret-key")
	// Create a Scaleway client
	client, err := scw.NewClient(
		scw.WithDefaultOrganizationID(getConfig("org-id")),
		scw.WithAuth(accessKey, secretKey),
	)
	if err != nil {
		logger.Fatalf("Error creating scaleway client: %v", err)
	}

	minioOpts := &minio.Options{
		Secure: true,
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
	}
	s3client, err := minio.New(getConfig("s3-endpoint"), minioOpts)
	if err != nil {
		logger.Fatalf("Error creating object storage client: %v")
	}

	// Setup the clients
	lbAPI = lb.NewAPI(client)
	k8sAPI = k8s.NewAPI(client)
	inAPI = instance.NewAPI(client)
	s3Client = s3client

	nsService = nsproto.NewNamespacesService("namespaces", svr.Client())
	stService = storeproto.NewStoreService("store", svr.Client())

	// Check infra daily and report any wastage
	c := cron.New()
	c.AddFunc("0 9 * * *", checkInfraUsageCron)
	c.Start()

	// Register the RPC handler
	pb.RegisterInfrastructureHandler(svr.Server(), new(infrastructure))

	// Run the service
	if err := svr.Run(); err != nil {
		logger.Fatal(err)
	}
}

func getConfig(key string) string {
	val, err := config.Get("micro.infrastructure.scaleway." + key)
	if err != nil {
		logger.Warnf("Error getting config: %v", err)
	}
	if len(val.String("")) == 0 {
		logger.Fatalf("Missing required config: micro.infrastructure.scaleway.%v", key)
	}
	return val.String("")
}
