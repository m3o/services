# Client

Client is a Micro services client to access anything from anywhere beyond Go and proto

## Usage

```
curl -XPOST \
-d '{"service": "go.micro.service.greeter", "endpoint": "Say.Hello"}' \
-H 'Content-Type: application/json' \
https://api.micro.mu/client
```
