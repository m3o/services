# Status Service

This is the Status service. It reports on the general status of the m3o platform. Used as an uptime "ping" type endpoint it will return various information about the status of core services

If things are OK you'll receive a `200 OK`, if not then you'll likely see `500 <some error>`

```
$ curl localhost:8080/status/call
{"statusCode":200,"body":"{\"go.micro.api\":\"OK\",\"go.micro.auth\":\"OK\",\"go.micro.broker\":\"OK\",\"go.micro.config\":\"OK\",\"go.micro.debug\":\"OK\",\"go.micro.network\":\"OK\",\"go.micro.proxy\":\"OK\",\"go.micro.registry\":\"OK\",\"go.micro.runtime\":\"OK\",\"go.micro.store\":\"OK\"}"}
``` 

## Config
The set of services to monitor can be loaded from config under the path `micro.status.services`. To set the list you can use the following call

```
micro call go.micro.config Config.Create '{"change":{"namespace":"micro","path":"micro.status.services","changeSet":{"format":"json","data":"go.micro.api,go.micro.auth,go.micro.broker,go.micro.config,go.micro.debug,go.micro.network,go.micro.proxy,go.micro.registry,go.micro.runtime,go.micro.server,go.micro.status,go.micro.store,go.micro.service.signup,go.micro.service.kubernetes,go.micro.service.invite,go.micro.service.payment"}}}'
```
