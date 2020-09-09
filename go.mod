module github.com/m3o/services

go 1.14

require (
	github.com/alexellis/hmac v0.0.0-20180624211220-5c52ab81c0de
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.2
	github.com/micro/go-micro/v3 v3.0.0-beta.2
	github.com/micro/micro/v3 v3.0.0-beta.3.0.20200907201209-650d6e32a270
	github.com/patrickmn/go-cache v2.1.0+incompatible
	github.com/sethvargo/go-diceware v0.2.0
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go/v71 v71.28.0
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
