module github.com/micro/services/events/api

go 1.13

require (
	github.com/golang/protobuf v1.4.1
	github.com/micro/go-micro/v2 v2.9.1-0.20200630164038-dcf01ebbf033
	github.com/micro/services/events/service v0.0.0-00010101000000-000000000000
	github.com/micro/services/projects/environments v0.0.0-20200511093345-f9d4a9151fe3
	github.com/micro/services/projects/service v0.0.0-20200505095011-36eddd53bd2b
)

replace github.com/micro/services/events/service => ../service

replace github.com/micro/services/projects/service => ../../projects/service
