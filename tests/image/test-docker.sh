pushd micro
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build
popd
docker build -t micro -f services/tests/image/Dockerfile .
