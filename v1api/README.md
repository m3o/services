# V1api Service

This is the V1api service

In charge of
- routing requests to the correct service
- API keys lifecycle management
- blocking requests (if insufficient privilege or some other reason).

## Flow

### API key lifecycle
1. Create an API key with scopes. API key starts in blocked status so all requests will be blocked.
2. New auth account generated for API key under the user's namespace with required scopes. Separate account means that we can track the usage of this API key at a finer granularity.  
2. Event fired for API key creation
3. Something (Quota service) picks up and works out which endpoints should be unblocked, calls `UpdateAllowedPaths` to unblock the user for the relevant urls
4. Serve requests

### Request flow
1. API key passed in on request as `Authorization: Bearer <key>`
2. micro API service can't decode the key as a JWT so will forward to the micro namespace by default
3. v1api service verifies the key by looking up the hash and forwards on the request, adding a JWT token to the request auth

### Blocking requests
If requests for a client need to be blocked (insufficient privilege or exhausted quota) we call `UpdateAllowedPaths` and update the allowed paths for the user as a whole (which has the effect of just updating all their API keys). 

The same happens for unblocking requests

## Scopes
Scopes define what an API key can do. 

Proposal: By convention scope names could implicitly define which endpoints the client has access to. So `location:write` means they have access to all endpoints as defined by the `location` service. 
