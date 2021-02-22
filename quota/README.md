# Quota Service

This is the Quota service

Records usage of a "thing" against a user. Enforces limits.

Flow:
1. Service creates a quota - quotaID, limit, frequency (when does it reset), description. locationsFree|locationsLevel1, 100, daily, free locations calls
2. Service records a relationship between a user and a quota. userID, thingID. user_1234, locationsFree.  
3. Quota service subscribes to event stream for endpoint invocation and records against the quotas
4. If a quota is breached it calls an endpoint on v1api to block that user. 
5. When a quota is reset (daily reset, someone gets upgraded, something else) we need to unblock all users by calling v1api again for each user


### Distributed Counting
We use Redis for our counter to keep track of quota usage. 

### Setup
Auth rules required 
```
micro auth create rule --resource="service:quota:Quota.List" --priority 1  quota-list-public   
micro auth create rule --resource="service:quota:Quota.ListUsage" --priority 1  quota-listusage-public
```
