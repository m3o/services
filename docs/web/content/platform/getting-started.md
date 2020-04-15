---
title: Get Started
keywords: platform, M3O
tags: [platform, M3O]
sidebar: home_sidebar
permalink: "/platform/getting-started.html"
summary: M3O is a serverless microservices platform
---

## Getting Access
M30 uses GitHub OAuth to login and verify users. While we're in a closed beta, you must be invited to the Micro community team on GitHub in order to gain access. If you don't have access to this team, please message us on the #platform channel on [Slack](micro-services.slack.com "Slack"). 

## Setup
Before starting let's ensure you have the latest version of Micro installed locally. To do this, run the following commands in your terminal:
```
rm $GOPATH/bin/micro
go get github.com/micro/micro/v2@master
```
Next, let's go to the [Micro website](https://micro.mu/platform/ "Micro") to get our login token. Press the "Sign in with GitHub button", if you are successfully authenticated, you'll be directed to the platform homepage where you'll see a list of running services and a button to login with CLI. Press this button and copy the command you're given, it will be in the format of:
```
micro login --token $token
```

If the login was successful, you will see the following message: `You have been logged in`.

## Writing your first service
As noted above, whilst M30 is in closed beta, the only services which can be deployed must be located within the github.com/micro/services repo. Let's close this repo, using the no checkout flag to speed up the process.
```
git clone https://github.com/micro/services && cd services
```

Next, let's create out first service (use your own first name for fun!) 
```
micro new foobar && cd foobar
```
At this point, you have a new micro service ready for deployment. All we need to do prior to deployment is build the proto buffer. We can use the Make command  to do this:
```
make build
```
## Deploying your first service
When you instruct M30 to run a service, it will pull the latest source code for the platform repo and run whatever service you specify. So firstly, let's push our changes to GitHub:
```
git add . && git commit -m "Initialising service" && git push
```
Next, let's use the `micro run` command to run the service. Note, foobar must be the directory of the service you wish to deploy.
```
micro run --platform foobar
```
If successful, you'll see the following message: `[Platform] Service foobar:latest created`. We can check on the progress of our deployment by running:
```
micro ps --platform
```
You'll now see a list of the services, including:
```
NAME		VERSION	SOURCE				STATUS		BUILD	METADATA
foobar		latest	github.com/micro/services	running		n/a	owner=n/a,group=n/a
```

## Interacting with your first service
Now we've deployed our first service, let's go and interact with it. We can do this via the [platform](https://micro.mu/platform "platform"). 
