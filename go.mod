module github.com/m3o/services

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/micro/go-micro/v3 v3.0.0-beta.2.0.20200910094249-53a32463765c
	github.com/micro/micro/v3 v3.0.0-beta.3.0.20200910111254-979f2fab21cd
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/robfig/cron v1.2.0
	github.com/scaleway/scaleway-sdk-go v1.0.0-beta.6
	github.com/sethvargo/go-diceware v0.2.0
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.6.1
	github.com/stripe/stripe-go/v71 v71.28.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
