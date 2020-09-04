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
- [ ] build inside a container
- [ ] push to registry
- [ ] metrics
