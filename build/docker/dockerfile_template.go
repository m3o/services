package docker

// DockerfileTemplate gets rendered to a buildable Dockerfile:
var DockerfileTemplate = `
# Build a service binary in a GoLang container:
FROM {{.BuildImage}} AS build
RUN go build -o /service {{.Source}}

# Copy the service binary into a lean Alpine container:
FROM {{.BaseImage}}
COPY --from=build /service /service
CMD ["/service"]
`
