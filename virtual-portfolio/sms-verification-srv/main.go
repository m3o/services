package main

import (
	"fmt"
	"os"

	"github.com/kytra-app/helpers/sms"
	"github.com/kytra-app/sms-verification-srv/handler"
	proto "github.com/kytra-app/sms-verification-srv/proto"
	"github.com/kytra-app/sms-verification-srv/storage/postgres"
	"github.com/micro/go-micro"
	_ "github.com/micro/go-plugins/registry/kubernetes"
	"github.com/pkg/errors"
)

func main() {
	// Create The Service
	service := micro.NewService(
		micro.Name("kytra-srv-v1-sms-verification"),
		micro.Version("latest"),
	)
	service.Init()

	// Connect to the Database (Postgres)
	db, err := postgres.New(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the database"))
	}
	defer db.Close()

	// Connect to SMS service
	sms, err := sms.New(
		os.Getenv("SMS_ACCOUNT_SID"),
		os.Getenv("SMS_AUTH_TOKEN"),
	)
	if err != nil {
		panic(errors.Wrap(err, "Could not connect to the SMS service"))
	}

	// Register to Service Discovery
	hander := handler.New(db, sms)
	proto.RegisterSMSVerificationHandler(service.Server(), hander)

	if err := service.Run(); err != nil {
		fmt.Println(err)
	}
}
