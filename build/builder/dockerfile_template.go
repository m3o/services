package builder

// dockerfileTemplateRaw gets rendered to a buildable Dockerfile:
var dockerfileTemplateRaw = `# Build a service binary in a GoLang container:
FROM {{.BuildImage}} AS build
RUN go get {{.SourceGitRepo}}
RUN go build -o /service {{.SourceGitRepo}}

# Copy the service binary into a lean Alpine container:
FROM {{.BaseImage}} AS service
COPY --from=build /service /service
CMD ["/service"]
`
