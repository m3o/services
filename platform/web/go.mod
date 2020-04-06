module github.com/micro/services/platform/web

go 1.13

require (
	github.com/micro/go-micro/v2 v2.4.1-0.20200406115547-bea7c3f7e720
	github.com/micro/services/platform/service v0.0.0-20200313185528-4a795857eb73
	golang.org/x/time v0.0.0-20191024005414-555d28b269f0 // indirect
)

replace github.com/micro/services/platform/service => ../service
