package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/registry"
	"github.com/micro/micro/v3/service/server"
	"github.com/micro/micro/v3/service/store"

	"github.com/google/uuid"
)

type V1 struct {
}

const (
	storePrefixHashedKey = "hashed"
	storePrefixUserID    = "user"
	storePrefixKeyID     = "key"
	storePrefixAPI       = "api"
)

var (
	errUnauthorized = errors.Unauthorized("v1api", "Unauthorized")
	errInternal     = errors.InternalServerError("v1api", "Error processing request")
	errBlocked      = errors.Forbidden("v1api.blocked", "Client is blocked")
)

type apiKeyRecord struct {
	ID          string          `json:"id"`          // id of the key
	ApiKey      string          `json:"apiKey"`      // hashed api key
	Scopes      []string        `json:"scopes"`      // the scopes this key has granted
	UserID      string          `json:"userID"`      // the ID of the key's owner
	AccID       string          `json:"accID"`       // the ID of the service account
	Description string          `json:"description"` // optional description of the API key as given by user
	Namespace   string          `json:"namespace"`   // the namespace that this user belongs to (only because technically user IDs aren't globally unique)
	Token       string          `json:"token"`       // the short lived JWT token
	Created     int64           `json:"created"`     // creation time
	AllowList   map[string]bool `json:"allowList"`   // map of allowed path prefixes
}

// Generate generates an API key
func (e *V1) GenerateKey(ctx context.Context, req *v1api.GenerateKeyRequest, rsp *v1api.GenerateKeyResponse) error {
	if len(req.Scopes) == 0 {
		return errors.BadRequest("v1api.generate", "Missing scopes field")
	}
	if len(req.Description) == 0 {
		return errors.BadRequest("v1api.generate", "Missing description field")
	}
	// Check account
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("v1api.generate", "Unauthorized call to generate")
	}
	// only namespace admins can generate a key
	admin := false
	for _, s := range acc.Scopes {
		if s == "admin" {
			admin = true
			break
		}
	}
	if !admin {
		return errors.Forbidden("v1api.generate", "Forbidden")
	}
	// generate a new API key

	// are they allowed to generate with the requested scopes?
	if !e.checkRequestedScopes(acc, req.Scopes) {
		return errors.Forbidden("v1api.generate", "Not allowed to generate a key with requested scopes")
	}

	// generate API key
	id, err := uuid.NewRandom()
	if err != nil {
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}

	apiKey := base64.StdEncoding.EncodeToString([]byte(id.String()))
	hashedKey, err := hashSecret(apiKey)
	if err != nil {
		log.Errorf("Error hashing api key %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}

	// api key is the secret for a new account
	// generate the new account + short lived access token for it
	authAcc, err := auth.Generate(
		uuid.New().String(),
		auth.WithSecret(apiKey),
		auth.WithIssuer(acc.Issuer),
		auth.WithType("apikey"),
		auth.WithScopes(req.Scopes...),
		auth.WithMetadata(map[string]string{"apikey_owner": acc.ID}),
	)
	if err != nil {
		log.Errorf("Error generating auth account %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}
	tok, err := auth.Token(
		auth.WithCredentials(authAcc.ID, apiKey),
		auth.WithTokenIssuer(acc.Issuer),
		auth.WithExpiry(1*time.Hour))
	if err != nil {
		log.Errorf("Error generating token %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}
	// hash API key and store with scopes
	rec := apiKeyRecord{
		ID:          uuid.New().String(),
		ApiKey:      hashedKey,
		Scopes:      req.Scopes,
		UserID:      acc.ID,
		Namespace:   acc.Issuer,
		Description: req.Description,
		AccID:       authAcc.ID,
		Token:       tok.AccessToken,
		Created:     time.Now().Unix(),
		AllowList:   map[string]bool{},
	}
	if err := writeAPIRecord(&rec); err != nil {
		log.Errorf("Failed to write api record %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}

	if err := events.Publish("v1api", v1api.Event{Type: "APIKeyCreate",
		ApiKeyCreate: &v1api.APIKeyCreateEvent{
			UserId:    rec.UserID,
			Namespace: rec.Namespace,
			ApiKeyId:  rec.ID,
			Scopes:    rec.Scopes,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}
	// return the unhashed key
	rsp.ApiKey = apiKey
	return nil
}

func writeAPIRecord(rec *apiKeyRecord) error {
	b, err := json.Marshal(rec)
	if err != nil {
		return err
	}

	// store under hashed key for API usage
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s", storePrefixHashedKey, rec.ApiKey),
		Value: b,
	}); err != nil {
		return err
	}

	// store under the user ID for retrieval on dashboard
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s:%s:%s", storePrefixUserID, rec.Namespace, rec.UserID, rec.ApiKey),
		Value: b,
	}); err != nil {
		return err
	}

	// store under the key's ID for deletion
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s:%s:%s", storePrefixKeyID, rec.Namespace, rec.UserID, rec.ID),
		Value: b,
	}); err != nil {
		return err
	}

	return nil
}

func deleteAPIRecord(rec *apiKeyRecord) error {
	// store under hashed key for API usage
	if err := store.Delete(fmt.Sprintf("%s:%s", storePrefixHashedKey, rec.ApiKey)); err != nil {
		return err
	}

	// store under the user ID for retrieval on dashboard
	if err := store.Delete(fmt.Sprintf("%s:%s:%s:%s", storePrefixUserID, rec.Namespace, rec.UserID, rec.ApiKey)); err != nil {
		return err
	}

	// store under the key's ID for deletion
	if err := store.Delete(fmt.Sprintf("%s:%s:%s:%s", storePrefixKeyID, rec.Namespace, rec.UserID, rec.ID)); err != nil {
		return err
	}

	return nil
}

func readAPIRecord(ns, user, keyID string) (*apiKeyRecord, error) {
	recs, err := store.Read(fmt.Sprintf("%s:%s:%s:%s", storePrefixKeyID, ns, user, keyID))
	if err != nil {
		return nil, err
	}

	rec := recs[0]
	keyRec := &apiKeyRecord{}
	if err := json.Unmarshal(rec.Value, keyRec); err != nil {
		return nil, err
	}
	return keyRec, nil
}

func hashSecret(s string) (string, error) {
	h := sha256.New()
	h.Write([]byte(s))
	h.Sum(nil)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// checkRequestedScopes returns true if account has sufficient privileges for them to generate the requestedScopes.
// e.g. micro "admin" can generate whatever scopes they want
func (e *V1) checkRequestedScopes(account *auth.Account, requestedScopes []string) bool {
	allowedScopes, err := e.listAPIs()
	if err != nil {
		// fail closed
		return false
	}
	// add wildcard
	allowedScopes = append(allowedScopes, "*")
	for _, requested := range requestedScopes {
		found := false
		for _, allowed := range allowedScopes {
			if allowed == requested {
				found = true
				break
			}
		}
		if !found {
			return false
		}
	}
	return true
}

func loadAPIRec(authz string) (string, *apiKeyRecord, error) {
	if len(authz) == 0 || !strings.HasPrefix(authz, "Bearer ") {
		return "", nil, errUnauthorized
	}

	// do lookup on hash of key
	key := authz[7:]
	hashed, err := hashSecret(key)
	if err != nil {
		return "", nil, errUnauthorized
	}
	recs, err := store.Read(fmt.Sprintf("%s:%s", storePrefixHashedKey, hashed))
	if err != nil {
		if err != store.ErrNotFound {
			log.Errorf("Error while looking up api key %s", err)
		}
		// not found == invalid (or even revoked)
		log.Infof("Authz not found %+v", hashed)
		return "", nil, errUnauthorized
	}
	// rehydrate
	apiRec := apiKeyRecord{}
	if err := json.Unmarshal(recs[0].Value, &apiRec); err != nil {
		log.Errorf("Error while rehydrating api key record %s", err)
		return "", nil, errUnauthorized
	}
	return key, &apiRec, nil
}

func verifyCallAllowed(apiRec *apiKeyRecord, reqURL string) error {
	// checks
	// We do 2 types of checks
	// 1. Do the scopes of the token allow them to call the requested API? The name of the scopes correspond to the service
	// 2. Has the key been explicitly blocked? Happens if it's exhausted it's quota for example

	// Type 1 check. This is a bit belt and braces to be honest, the type 2 check should probably be enough
	scopeMatch := false
	for _, s := range apiRec.Scopes {
		if strings.HasPrefix(reqURL, fmt.Sprintf("/v1/%s/", s)) {
			scopeMatch = true
			break
		}
	}
	if !scopeMatch {
		// TODO better error please
		return errBlocked
	}

	// Type 2 check
	allowed := false
	for prefix := range apiRec.AllowList {
		if strings.HasPrefix(reqURL[3:], prefix) {
			allowed = true
		}
	}
	if !allowed {
		// TODO better error please
		return errBlocked
	}
	return nil
}

func refreshToken(apiRec *apiKeyRecord, key string) error {
	// do we need to refresh the token?
	tok, _, err := new(jwt.Parser).ParseUnverified(apiRec.Token, jwt.MapClaims{})
	if err != nil {
		log.Errorf("Error parsing existing jwt %s", err)
		return errUnauthorized
	}
	if claims, ok := tok.Claims.(jwt.MapClaims); ok {
		if !claims.VerifyExpiresAt(time.Now().Unix(), true) {
			// needs a refresh
			tok, err := auth.Token(
				auth.WithCredentials(apiRec.AccID, key),
				auth.WithTokenIssuer(apiRec.Namespace),
				auth.WithExpiry(1*time.Hour))
			if err != nil {
				log.Errorf("Error refreshing token %s", err)
				return errInternal
			}
			apiRec.Token = tok.AccessToken
			if err := writeAPIRecord(apiRec); err != nil {
				log.Errorf("Error updating API record %s", err)
				return errInternal
			}
		}
	} else {
		log.Errorf("Error parsing existing jwt claims %s", err)
		return errUnauthorized
	}
	return nil
}

func getRequestedService(reqURL string) (string, string, []*registry.Service, error) {
	trimmedPath := strings.TrimPrefix(reqURL, "/v1/")
	parts := strings.Split(trimmedPath, "/")
	if len(parts) < 2 {
		// can't work out service and method
		return "", "", nil, errors.NotFound("v1api", "")
	}

	service := parts[0]
	svcs, err := registry.GetService(service)
	if err != nil {
		if err == registry.ErrNotFound {
			return "", "", nil, errors.NotFound("v1api", "No such API")
		}
		log.Errorf("Error looking up service %s", err)
		return "", "", nil, errInternal
	}
	endpoint := ""
	if len(parts) == 2 {
		// /v1/helloworld/call -> helloworld Helloworld.Call
		endpoint = fmt.Sprintf("%s.%s", strings.Title(parts[0]), strings.Title(parts[1]))
	} else {
		// /v1/hello/world/call -> hello World.Call
		endpoint = fmt.Sprintf("%s.%s", strings.Title(parts[1]), strings.Title(parts[2]))
	}

	return service, endpoint, svcs, nil
}

// Endpoint is a catch all for endpoints
func (e *V1) Endpoint(ctx context.Context, stream server.Stream) error {
	// check api key
	log.Infof("Received request")
	defer stream.Close()

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errUnauthorized
	}

	key, apiRec, err := loadAPIRec(md["Authorization"])
	if err != nil {
		return err
	}

	reqURL, ok := md.Get("url")
	if !ok {
		log.Errorf("Requested URL not found")
		return errInternal
	}

	if err := verifyCallAllowed(apiRec, reqURL); err != nil {
		return err
	}

	if err := refreshToken(apiRec, key); err != nil {
		return err
	}
	// set the auth
	ctx = metadata.Set(ctx, "Authorization", fmt.Sprintf("Bearer %s", apiRec.Token))

	// assume application/json for now
	ct := "application/json"
	service, endpoint, svcs, err := getRequestedService(reqURL)
	if err != nil {
		return err
	}

	if isStream(endpoint, svcs) {
		log.Infof("It's a stream!!!")
		return serveStream(ctx, stream, service, endpoint, svcs)
	}

	// forward the request
	var payload json.RawMessage
	if err := stream.Recv(&payload); err != nil {
		log.Errorf("Error receiving from stream %s", err)
		return errInternal
	}

	request := client.DefaultClient.NewRequest(
		service,
		endpoint,
		&payload,
		client.WithContentType(ct),
	)

	// create request/response
	var response json.RawMessage
	// make the call
	if err := client.Call(ctx, request, &response); err != nil {
		return err
	}

	if err := events.Publish("v1api", v1api.Event{Type: "Request",
		Request: &v1api.RequestEvent{
			UserId:    apiRec.UserID,
			Namespace: apiRec.Namespace,
			ApiKeyId:  apiRec.ID,
			Url:       reqURL,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	stream.Send(response)
	return nil

}

// ListKeys lists all keys for a user
func (e *V1) ListKeys(ctx context.Context, req *v1api.ListRequest, rsp *v1api.ListResponse) error {
	// Check account
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("v1api.listkeys", "Unauthorized call to listkeys")
	}
	recs, err := listKeysForUser(acc.Issuer, acc.ID)
	if err != nil {
		log.Errorf("Error listing keys %s", err)
		return errors.InternalServerError("v1aapi.listkeys", "Error listing keys")
	}
	rsp.ApiKeys = make([]*v1api.APIKey, len(recs))
	for i, apiRec := range recs {
		rsp.ApiKeys[i] = &v1api.APIKey{
			Id:          apiRec.ID,
			Description: apiRec.Description,
			CreatedTime: apiRec.Created,
			Scopes:      apiRec.Scopes,
		}
	}
	return nil
}

func listKeysForUser(ns, userID string) ([]*apiKeyRecord, error) {
	recs, err := store.Read(fmt.Sprintf("%s:%s:%s:", storePrefixUserID, ns, userID), store.ReadPrefix())
	if err != nil {
		if err == store.ErrNotFound {
			return nil, nil
		}
		return nil, err
	}
	ret := make([]*apiKeyRecord, len(recs))
	for i, rec := range recs {
		apiRec := &apiKeyRecord{}
		if err := json.Unmarshal(rec.Value, apiRec); err != nil {
			return nil, err
		}
		ret[i] = apiRec
	}
	return ret, nil
}

func (e *V1) RevokeKey(ctx context.Context, request *v1api.RevokeRequest, response *v1api.RevokeResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("v1api.Revoke", "Unauthorized call to revoke")
	}
	if len(request.Id) == 0 {
		return errors.BadRequest("v1api.Revoke", "Missing ID field")
	}

	rec, err := readAPIRecord(acc.Issuer, acc.ID, request.Id)
	if err != nil {
		if err == store.ErrNotFound {
			return errors.NotFound("v1api.Revoke", "Not found")
		}
		log.Errorf("Error reading API key record %s", err)
		return errors.InternalServerError("v1pi.Revoke", "Error revoking key")
	}
	if err := deleteAPIRecord(rec); err != nil {
		log.Errorf("Error deleting API key record %s", err)
		return errors.InternalServerError("v1pi.Revoke", "Error revoking key")
	}

	if err := events.Publish("v1api", v1api.Event{Type: "APIKeyRevoke",
		ApiKeyRevoke: &v1api.APIKeyRevokeEvent{
			UserId:    acc.ID,
			Namespace: acc.Issuer,
			ApiKeyId:  request.Id,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	return nil
}

func (e *V1) UpdateAllowedPaths(ctx context.Context, request *v1api.UpdateAllowedPathsRequest, response *v1api.UpdateAllowedPathsResponse) error {
	if err := verifyMicroAdmin(ctx, "v1api.UpdateAllowedPaths"); err != nil {
		return err
	}

	var keys []*apiKeyRecord
	if len(request.KeyId) > 0 {
		rec, err := readAPIRecord(request.Namespace, request.UserId, request.KeyId)
		if err != nil {
			log.Errorf("Error reading key %s", err)
			return errors.InternalServerError("v1api.UpdateAllowedPaths", "Error updating user")
		}
		keys = []*apiKeyRecord{rec}
	} else {
		recs, err := listKeysForUser(request.Namespace, request.UserId)
		if err != nil {
			log.Errorf("Error listing keys %s", err)
			return errors.InternalServerError("v1api.UpdateAllowedPaths", "Error updating user")
		}
		keys = recs
	}
	update := func(key *apiKeyRecord, allow, block []string) error {
		for _, a := range allow {
			for _, s := range key.Scopes {
				if strings.HasPrefix(a, fmt.Sprintf("/%s/", s)) {
					key.AllowList[a] = true
					break
				}
			}
		}
		for _, b := range block {
			delete(key.AllowList, b)
		}
		return writeAPIRecord(key)
	}
	for _, k := range keys {
		if err := update(k, request.Allowed, request.Blocked); err != nil {
			log.Errorf("Error updating key api key record %s", err)
			return errors.InternalServerError("v1api.UpdateAllowedPaths", "Error updating allowed paths")
		}
	}

	return nil
}

func verifyMicroAdmin(ctx context.Context, method string) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden(method, "Forbidden")
	}
	admin := false
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			admin = true
			break
		}
	}
	if !admin {
		return errors.Forbidden(method, "Forbidden")
	}
	return nil
}

type apiEntry struct {
	Name string
}

func (e *V1) EnableAPI(ctx context.Context, request *v1api.EnableAPIRequest, response *v1api.EnableAPIResponse) error {
	if err := verifyMicroAdmin(ctx, "v1api.EnableAPI"); err != nil {
		return err
	}
	if len(request.Name) == 0 {
		return errors.BadRequest("v1api.EnableAPI", "Missing Name field")
	}
	// store
	ae := &apiEntry{
		Name: request.Name,
	}
	b, err := json.Marshal(ae)
	if err != nil {
		log.Errorf("Error marshalling API %s", err)
		return errors.InternalServerError("v1api.EnableAPI", "Error enabling API")
	}
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s", storePrefixAPI, request.Name),
		Value: b,
	}); err != nil {
		log.Errorf("Error persisting API %s", err)
		return errors.InternalServerError("v1api.EnableAPI", "Error enabling API")
	}
	// rules
	// we need two rules, one for the /v1/foo/bar from public internet and one for v1api->foo
	//micro auth create rule --resource="service:v1.helloworld:*" --priority 1 helloworld-v1
	//micro auth create rule --resource="service:helloworld:*" --priority 1 --scope '+' helloworld-internal
	if err := auth.Grant(&auth.Rule{
		ID:    fmt.Sprintf("%s-v1", request.Name),
		Scope: "",
		Resource: &auth.Resource{
			Name:     fmt.Sprintf("v1.%s", request.Name),
			Type:     "service",
			Endpoint: "*",
		},
		Access:   auth.AccessGranted,
		Priority: 1,
	}); err != nil {
		log.Errorf("Error adding rule %s", err)
		return errors.InternalServerError("v1api.EnableAPI", "Error enabling API")
	}

	if err := auth.Grant(&auth.Rule{
		ID:    fmt.Sprintf("%s-internal", request.Name),
		Scope: "+",
		Resource: &auth.Resource{
			Name:     request.Name,
			Type:     "service",
			Endpoint: "*",
		},
		Access:   auth.AccessGranted,
		Priority: 1,
	}); err != nil {
		log.Errorf("Error adding rule %s", err)
		return errors.InternalServerError("v1api.EnableAPI", "Error enabling API")
	}

	// event
	if err := events.Publish("v1api", v1api.Event{Type: "APIEnable",
		ApiEnable: &v1api.APIEnableEvent{
			Name: request.Name,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	return nil
}

func (e *V1) DisableAPI(ctx context.Context, request *v1api.DisableAPIRequest, response *v1api.DisableAPIResponse) error {
	if err := verifyMicroAdmin(ctx, "v1api.DisableAPI"); err != nil {
		return err
	}
	if len(request.Name) == 0 {
		return errors.BadRequest("v1api.DisableAPI", "Missing Name field")
	}

	// delete from store
	if err := store.Delete(fmt.Sprintf("%s:%s", storePrefixAPI, request.Name)); err != nil {
		log.Errorf("Error deleting store entry %s", err)
		return errors.InternalServerError("v1pi.Disable", "Error disabling API")
	}
	// delete rules
	if err := auth.Revoke(&auth.Rule{ID: fmt.Sprintf("%s-internal", request.Name)}); err != nil {
		log.Errorf("Error deleting rule %s", err)
		return errors.InternalServerError("v1pi.Disable", "Error disabling API")
	}
	if err := auth.Revoke(&auth.Rule{ID: fmt.Sprintf("%s-v1", request.Name)}); err != nil {
		log.Errorf("Error deleting rule %s", err)
		return errors.InternalServerError("v1pi.Disable", "Error disabling API")
	}
	// event
	if err := events.Publish("v1api", v1api.Event{Type: "APIDisable",
		ApiDisable: &v1api.APIDisableEvent{
			Name: request.Name,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	return nil
}

func (e *V1) ListAPIs(ctx context.Context, request *v1api.ListAPIsRequest, response *v1api.ListAPIsResponse) error {
	ret, err := e.listAPIs()
	if err != nil {
		log.Errorf("Error listing APIs %s", err)
		return errors.InternalServerError("v1pi.ListAPIs", "Error listing APIs")
	}
	response.Names = ret
	return nil
}

func (e *V1) listAPIs() ([]string, error) {
	recs, err := store.Read(fmt.Sprintf("%s:", storePrefixAPI), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		return nil, err
	}
	ret := make([]string, len(recs))
	for i, v := range recs {
		ae := &apiEntry{}
		if err := json.Unmarshal(v.Value, ae); err != nil {
			return nil, err
		}
		ret[i] = ae.Name
	}
	return ret, nil
}
