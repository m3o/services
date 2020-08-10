module github.com/m3o/services

go 1.14

require (
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v3 v3.0.0-alpha.0.20200809211157-65e6ee85667d
	github.com/micro/micro/v3 v3.0.0-alpha.0.20200809221604-9cc39c37378d
	github.com/sethvargo/go-diceware v0.2.0
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go v70.15.0+incompatible
	github.com/stripe/stripe-go/v71 v71.28.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
