# GitOps Service

GitOps is responsible for performing updates based on changes in git. The service is configured with a GitHub PAT, repo name and webhook URL. When it starts, it will register a webhook against the repo using the GitHub webhooks API for the given URL.

When a commit is added to master, the service will trigger an update of the micro services using the runtime. For example, when a commit is added to micro/micro, all the micro services (e.g. the router, registry) will be updated. This is intended to be used in a staging environment to enable continuous deployments.