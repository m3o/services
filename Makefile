# test runs the tests for all the services in the repo,
# ensuring at the very least they can build
testall:
	find . -name "main.go" | xargs -n 1 go test