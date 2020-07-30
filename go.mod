module github.com/m3o/services

go 1.14

require (
	github.com/bradfitz/gomemcache v0.0.0-20190913173617-a41fca850d0b // indirect
	github.com/golang/protobuf v1.4.2
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.9.1-0.20200723075038-fbdf1f2c1c4c
	github.com/micro/micro/v3 v3.0.0-20200730101154-cc2a2ab5232b
	github.com/sethvargo/go-diceware v0.2.0
	github.com/slack-go/slack v0.6.5
	github.com/stretchr/testify v1.5.1
	github.com/stripe/stripe-go v70.15.0+incompatible
	github.com/teris-io/shortid v0.0.0-20171029131806-771a37caa5cf // indirect
	golang.org/x/crypto v0.0.0-20200709230013-948cd5f35899
	google.golang.org/protobuf v1.25.0
)

replace google.golang.org/grpc => google.golang.org/grpc v1.26.0
