# Infrastructure

The infrastructure service provides an RPC interface to access a summary of the current infrastructure. The service is currently setup to use the Scaleway API, so if the platform is swapped at a later time, this service will need updating.

To access the summary, run the following command: `micro infrastructure summary`

There is also a CRON job which runs daily at 9am checking for unused infra. If any used infra (e.g. servers) are found, a message a sent via slack to the team-important channel.