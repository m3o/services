module github.com/micro/services/account/web

go 1.13

require (
	github.com/google/uuid v1.1.1
	github.com/micro/go-micro/v2 v2.7.1-0.20200523154723-bd049a51e637
	github.com/micro/services/login/service v0.0.0-20200313083714-e72c0c76aa9a
	github.com/micro/services/projects/invite v0.0.0-20200421101014-4b009b48a425
	github.com/micro/services/users/service v0.0.0-20200501143857-056deed3461f
)

replace github.com/micro/services/projects/invite => ../../projects/invite

replace github.com/micro/services/projects/service => ../../projects/service
