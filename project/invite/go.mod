module github.com/micro/services/project/invite

go 1.13

require (
	github.com/golang/protobuf v1.3.5
	github.com/google/uuid v1.1.1
	github.com/kr/text v0.2.0 // indirect
	github.com/micro/go-micro/v2 v2.5.1-0.20200422094434-e25ab9f4ca28
	github.com/micro/services/project/service v0.0.0-20200421073553-26a9ccb4988a
	github.com/micro/services/users/service v0.0.0-20200421073553-26a9ccb4988a
	github.com/niemeyer/pretty v0.0.0-20200227124842-a10e7caefd8e // indirect
	golang.org/x/net v0.0.0-20200301022130-244492dfa37a // indirect
	golang.org/x/sys v0.0.0-20200302150141-5c8b2ff67527 // indirect
	gopkg.in/check.v1 v1.0.0-20200227125254-8fa46927fb4f // indirect
	gopkg.in/yaml.v2 v2.2.4 // indirect
)

replace github.com/micro/services/project/service => ../service
