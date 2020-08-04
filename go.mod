module github.com/m3o/services

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v3 v3.0.0-alpha.0.20200804104301-07fef9fd33c2
	github.com/micro/micro/v3 v3.0.0-alpha.0.20200804151614-7b733e61add7
	github.com/sethvargo/go-diceware v0.2.0
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go v70.15.0+incompatible
	github.com/stripe/stripe-go/v71 v71.28.0
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
