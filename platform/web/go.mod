module github.com/micro/services/platform/web

go 1.13

require (
	github.com/micro/go-micro/v2 v2.3.1-0.20200319095459-cbb958def528
	github.com/micro/micro/v2 v2.2.1-0.20200314171200-6192587db534
	github.com/micro/services/platform/service v0.0.0-20200313185528-4a795857eb73
	github.com/micro/services/users/service v0.0.0-20200319140645-20aa308d0728
)

replace github.com/micro/services/platform/service => ../service