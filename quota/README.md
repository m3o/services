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
Redis.

Counting quickly and reliably is hard. For our requirements it's OK to be accurate-ish. We don't want to overestimate how much a user has used since that is unfair. Under estimating means that we potentially take a loss if the operation is expensive but that's better than the alternative of having an accurate but slow count. Accurate up to 1 second. (in case of node failure you lose up to 1 second of counts).

Background loop which every second syncs to the DB
