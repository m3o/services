# Status Service

This is the Status service. It reports on the general status of the m3o platform. Used as an uptime "ping" type endpoint it will return various information about the status of core services

If things are OK you'll receive a `200 OK`, if not then you'll likely see `500 <some error>`

```
$ curl localhost:8080/status/call
{"statusCode":200,"body":"{\"go.micro.api\":\"OK\",\"go.micro.auth\":\"OK\",\"go.micro.broker\":\"OK\",\"go.micro.config\":\"OK\",\"go.micro.debug\":\"OK\",\"go.micro.network\":\"OK\",\"go.micro.proxy\":\"OK\",\"go.micro.registry\":\"OK\",\"go.micro.runtime\":\"OK\",\"go.micro.store\":\"OK\"}"}
``` 