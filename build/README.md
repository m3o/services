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


Development Roadmap
-------------------

In `v1` the aim is to get basic functionality working as a stand-alone Micro service:
- [X] main.go
- [X] proto
- [X] handler
- [X] config from micro
- [X] login to Scaleway registry
- [X] build inside a container
- [X] push to registry
- [X] metrics
- [X] periodically clean up images
- [X] make sure that we catch the full output of docker commands
- [X] checkout the source repo with git before compiling (instead of relying on "go build" to do it for us)
- [X] add a streaming endpoint which stays open until a build finishes, periodically returning a status update:
    - state (enum)
    - output (from Docker commands etc)
    - errors
- [X] tests

 In `v2` the aim is to integrate with the micro platform and toolset:
- [ ] move the service proto to micro/micro
- [ ] move the build logic into go-micro
- [ ] potentially update the Build interface in go-micro (and drag the rest of the implementations with it?)
- [ ] build all of this into the runtime service, and use profiles to inject a build-service builder for the micro platform
- [ ] figure out how to override the resource-limits in the default deployment

 Future features:
- [ ] how to clean up all built images (but not build/base)?
- [ ] handle credentials for private git repos
- [ ] consider an option share Go cache (so the build container doesn't have to download so much cruft every time)


Example call
------------
* `micro call build Build.CreateImage '{"gitRepo": "https://github.com/micro/micro.git", "gitCommit": "master", "imageTag": "rg.fr-par.scw.cloud/micro/micro:latest"}'`

Links
-----
* https://stackoverflow.com/questions/38804313/build-docker-image-from-go-code
* https://medium.com/faun/how-to-build-docker-images-on-the-fly-2a1fd696c3fd
