Build Service
=============


Features
--------
* Takes a Build() request
    - Git repo
    - Docker registry URL
* Builds a go binary inside a GoLang alpine container
* Builds a runtime image using a configured base image
* Logs in to the docker registry
* Pushes to the specific container registry


Questions
---------
* Should we use the hosts docker daemon, or an ephemeral one?
* How do we manage storage usage?
* Should we use a cache?
* Should we use a go-mod proxy?


Config
------
* Base image URL
* Docker credentials


Todo
----
- [X] main.go
- [X] proto
- [X] handler
- [X] config from micro
- [X] login to Scaleway registry
- [X] build inside a container
- [X] push to registry
- [X] metrics
- [X] periodically clean up images
- [ ] how to clean up all built images (but not build/base)?
- [ ] handle credentials for private git repos
- [ ] share Go cache (so the build container doesn't have to download so much cruft every time)
- [ ] store the state of each build, and provide an RPC to query this
- [ ] make sure that we catch the full output of docker commands
- [X] checkout the source repo with git before compiling (instead of relying on "go build" to do it for us)

Test
----
* `micro call build Build.CreateImage '{"gitRepo": "github.com/micro/micro", "dockerRegistry": "rg.fr-par.scw.cloud/micro", "imageTag": "rg.fr-par.scw.cloud/micro/micro:latest"}'`

Links
-----
* https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
* https://medium.com/faun/how-to-build-docker-images-on-the-fly-2a1fd696c3fd
