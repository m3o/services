# Platform Web

This is the platform web dashboard

## Overview

The platform web dashboard is a consolidated UI for managing a platform.

## Usage

```
GITHUB_OAUTH_CLIENT_ID=[...] \
GITHUB_OAUTH_CLIENT_SECRET=[...] \
GITHUB_OAUTH_REDIRECT_URL=http://127.0.0.1:6060/v1/auth/verify \
MICRO_SERVER_ADDRESS=:6060 \
FRONTEND_ADDRESS=http://127.0.0.1:4200 \
go run main.go
```
