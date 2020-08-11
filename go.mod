module github.com/m3o/services

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v3 v3.0.0-alpha.0.20200811171118-69a53e807013
	github.com/micro/micro/v3 v3.0.0-alpha.0.20200811140745-bc9bf56aeb2e
	github.com/sethvargo/go-diceware v0.2.0
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go v70.15.0+incompatible
	github.com/stripe/stripe-go/v71 v71.28.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
