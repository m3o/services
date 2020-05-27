module github.com/micro/services/projects/invite

go 1.13

require (
	github.com/golang/protobuf v1.4.0
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.7.1-0.20200523154723-bd049a51e637
	github.com/micro/services/projects/service v0.0.0-20200421073553-26a9ccb4988a
	github.com/micro/services/users/service v0.0.0-20200421073553-26a9ccb4988a
)

replace github.com/micro/services/projects/service => ../service
