module github.com/micro/services/account/web

go 1.13

require (
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.4.1-0.20200403120726-ed6fe67880a4
	github.com/micro/services/login/service v0.0.0-20200313083714-e72c0c76aa9a
	github.com/micro/services/users/service v0.0.0-20200319140645-20aa308d0728
)

replace github.com/micro/services/users/service => ../../users/service
