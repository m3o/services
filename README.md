# Services

This is a place for micro services

They currently reside at [github.com/microhq](https://github.com/microhq) but we'll move them soon.

## Overview

This repository serves as the monorepo for the M3O platform and as a reference architecture for others. 
Those invited to use the platform will be added to the Community team and have the ability to create 
and modify services here.

## Naming

Directories are the domain boundary for a specific concern e.g user, account, payment. They act as the 
alias for the otherwise fully qualified domain name "go.micro.srv.alias". Services should follow 
this naming convention and focus on single word naming where possible.

## Structure

Services should be generated using the `micro new` command using the alias e.g `micro new account`. 
The internal structure is defined by our new template generator. Extending this should follow 
a further convention as follows:

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

