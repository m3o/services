# Endtoend Service

This is the Endtoend service which performs end to end testing in live to ensure uptime.


## What does it do?
There is a cron job which fires the check every X minutes. The check will download micro via the m3o install script (https://install.m3o.com/micro) and run `micro signup`. This satisfies the requirement to run what the user runs. 

The email address used is from https://www.cloudmailin.com/ that allows us to receive emails on a webhook. This means we can pull out the OTP for email verification. 

The result of the check is recorded in the store and can be queried via the `check` endpoint. This can be called externally by uptime robot and the result will indicate whether the last check was successful and recent (within 5 mins). 