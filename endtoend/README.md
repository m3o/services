# Endtoend Service

This is the Endtoend service which performs end to end testing in live to ensure uptime.


## What does it do?
There is a cron job which fires the check every X minutes. The check will download micro via the m3o install script (https://install.m3o.com/micro) and run `micro signup`. This satisfies the requirement to run what the user runs. 

The email address used is from https://www.cloudmailin.com/ that allows us to receive emails on a webhook. This means we can pull out the OTP for email verification. 

The result of the check is recorded in the store and can be queried via the `check` endpoint. This can be called externally by uptime robot and the result will indicate whether the last check was successful and recent (within 5 mins). 

## Setup
To setup this service for a new environment you should

#### Setup an email address 
See https://www.cloudmailin.com/ to sign up for a new email.

You should set it up to post email in `JSON normalized` format and set the target URL to `https://api.<domain>/endtoend/mailin`

#### Setup the config
You need to setup the email address that you've just signed up for
```
micro config set micro.endtoend.email <email address from previous step>
``` 

#### Setup auth rules
To enable external monitoring we need to allow public access to `/endtoend/check` through the API

```
micro auth create rule --resource="service:endtoend:Endtoend.Check" endtoendcheck
micro auth create rule --resource="service:endtoend:Endtoend.Mailin" endtoendmail
```

#### Deploy
```
micro run github.com/m3o/services/endtoend
```

#### Setup your external monitor
Setup your external monitor (uptime robot, etc) to ping `https://api.<domain>/endtoend/check`

## Manual checking
Checks can be triggered manually by hitting the `RunCheck` endpoint
```
micro endtoend RunCheck
``` 

and you can check the logs or the usual alert channels for success/failure. 
