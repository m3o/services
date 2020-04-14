# Update API

The update API is a system update api for micro

## Overview

The update api is pinged by both the docker and github webhooks to keep track of whats changing. 
The micro runtime polls this api for changes to know when it should update. 
