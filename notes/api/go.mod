module github.com/micro/services/notes/api

go 1.13

require (
	github.com/golang/protobuf v1.3.2
	github.com/micro/go-micro/v2 v2.2.1-0.20200307205003-f01664a55156
	github.com/micro/services/notes/service v0.0.0-00010101000000-000000000000
)

replace github.com/micro/go-micro/v2 => ../../../go-micro

replace github.com/micro/services/notes/service => ../service
