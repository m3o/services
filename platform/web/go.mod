module github.com/micro/services/platform/web

go 1.13

require (
	github.com/micro/go-micro/v2 v2.7.1-0.20200523154723-bd049a51e637
	github.com/micro/micro/v2 v2.4.0
	github.com/micro/services/platform/service v0.0.0-20200313185528-4a795857eb73
)

replace github.com/micro/services/platform/service => ../service
