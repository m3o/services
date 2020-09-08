package builder

// dockerfileTemplateRaw gets rendered to a buildable Dockerfile:
var dockerfileTemplateRaw = `# Build a service binary in a GoLang container:
FROM {{.BuildImage}} AS build
RUN mkdir -p /tmp/build
RUN git clone {{.SourceGitRepo}} /tmp/build
RUN cd /tmp/build && git checkout {{.SourceGitCommit}} && go build -o /service

# Copy the service binary into a lean Alpine container:
FROM {{.BaseImage}} AS service
COPY --from=build /service /service
CMD ["/service"]
`
