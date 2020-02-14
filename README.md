# Services

This is a place for micro services

They currently reside at [github.com/microhq](https://github.com/microhq) but we'll move them soon.

## Overview

This repository serves as the monorepo for the M3O platform and as a reference architecture for others. 
Those invited to use the platform will be added to the Community team and have the ability to create 
and modify services here.

## Structure

Services should be generated using the `micro new` command.

Service and repo organisation should be as follows (example name used):

```
user/
    api/	# api routes
    web/	# web html
    client/	# generated clients
    service/	# core service types
    handler/	# request handlers
    subscriber/	# message subscribers
    proto/	# proto generated code
    main.go	# service main
    user.mu	# mu definition
    README.md	# readme
```

