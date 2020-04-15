module github.com/micro/services/platform/web

go 1.13

require (
	github.com/micro/go-micro/v2 v2.4.1-0.20200409084150-c1ad6d6c7c11
	github.com/micro/micro/v2 v2.4.0
	github.com/micro/services/platform/service v0.0.0-20200313185528-4a795857eb73
)

replace github.com/micro/services/platform/service => ../service
