package builder

// Builder does the actual building of docker images:
type Builder interface {
	Build(sourceGitRepo, targetImageTag string) error
}

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
