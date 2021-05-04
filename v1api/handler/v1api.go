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
	m3oauth "github.com/m3o/services/pkg/auth"
	publicapi "github.com/m3o/services/publicapi/proto"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service"
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
	papi publicapi.PublicapiService
}

const (
	storePrefixHashedKey = "hashed"
	storePrefixUserID    = "user"
	storePrefixKeyID     = "key"
)

var (
	errUnauthorized = errors.Unauthorized("v1api", "Unauthorized")
	errInternal     = errors.InternalServerError("v1api", "Error processing request")
	errBlocked      = errors.Forbidden("v1api.blocked", "Client is blocked")
)

type keyStatus string

const (
	keyStatusActive  = "active"
	keyStatusBlocked = "blocked" // blocked - probably because they've run out of money
)

type apiKeyRecord struct {
	ID          string   `json:"id"`          // id of the key
	ApiKey      string   `json:"apiKey"`      // hashed api key
	Scopes      []string `json:"scopes"`      // the scopes this key has granted
	UserID      string   `json:"userID"`      // the ID of the key's owner
	AccID       string   `json:"accID"`       // the ID of the service account
	Description string   `json:"description"` // optional description of the API key as given by user
	Namespace   string   `json:"namespace"`   // the namespace that this user belongs to (only because technically user IDs aren't globally unique)
	Token       string   `json:"token"`       // the short lived JWT token
	Created     int64    `json:"created"`     // creation time
	//AllowList   map[string]bool `json:"allowList"`   // map of allowed path prefixes
	Status keyStatus `json:"status"` // status of the key
}

func NewHandler(srv *service.Service) *V1 {
	return &V1{
		papi: publicapi.NewPublicapiService("publicapi", srv.Client()),
	}
}

// Generate generates an API key
func (e *V1) GenerateKey(ctx context.Context, req *v1api.GenerateKeyRequest, rsp *v1api.GenerateKeyResponse) error {
	if len(req.Scopes) == 0 {
		return errors.BadRequest("v1api.generate", "Missing scopes field")
	}
	if len(req.Description) == 0 {
		return errors.BadRequest("v1api.generate", "Missing description field")
	}

	acc, err := m3oauth.VerifyMicroCustomer(ctx, "v1api.generate")
	if err != nil {
		return err
	}
	// are they allowed to generate with the requested scopes?
	if !e.checkRequestedScopes(ctx, req.Scopes) {
		return errors.Forbidden("v1api.generate", "Not allowed to generate a key with requested scopes")
	}

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
		Status:      keyStatusBlocked,
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

func readAPIRecordByKeyID(ns, user, keyID string) (*apiKeyRecord, error) {
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
func (e *V1) checkRequestedScopes(ctx context.Context, requestedScopes []string) bool {
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

func readAPIRecordByAPIKey(authz string) (string, *apiKeyRecord, error) {
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
	// 1. Has the key been explicitly blocked? Happens if it's exhausted it's money for example
	// 2. Do the scopes of the token allow them to call the requested API? The name of the scopes correspond to the service

	if apiRec.Status == keyStatusBlocked {
		return errBlocked
	}

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

func parseContentType(ct string) string {
	if idx := strings.IndexRune(ct, ';'); idx >= 0 {
		ct = ct[:idx]
	}
	if len(ct) == 0 {
		ct = "application/json"
	}
	return ct
}

// Endpoint is a catch all for endpoints
func (e *V1) Endpoint(ctx context.Context, stream server.Stream) (retErr error) {
	// check api key
	defer stream.Close()

	md, ok := metadata.FromContext(ctx)
	if !ok {
		return errUnauthorized
	}

	key, apiRec, err := readAPIRecordByAPIKey(md["Authorization"])
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
		return serveStream(ctx, stream, service, endpoint, svcs, apiRec)
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
	publishEndpointEvent(reqURL, service, endpoint, apiRec)

	stream.Send(response)
	return nil

}

func publishEndpointEvent(reqURL, apiName, endpointName string, apiRec *apiKeyRecord) {
	if err := events.Publish("v1api", v1api.Event{Type: "Request",
		Request: &v1api.RequestEvent{
			UserId:       apiRec.UserID,
			Namespace:    apiRec.Namespace,
			ApiKeyId:     apiRec.ID,
			Url:          reqURL,
			ApiName:      apiName,
			EndpointName: endpointName,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}
}

// ListKeys lists all keys for a user
func (e *V1) ListKeys(ctx context.Context, req *v1api.ListRequest, rsp *v1api.ListResponse) error {
	// Check account
	acc, err := m3oauth.VerifyMicroCustomer(ctx, "v1api.listkeys")
	if err != nil {
		return err
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
	if len(request.Id) == 0 {
		return errors.BadRequest("v1api.Revoke", "Missing ID field")
	}
	acc, err := m3oauth.VerifyMicroCustomer(ctx, "v1api.Revoke")
	if err != nil {
		return err
	}

	rec, err := readAPIRecordByKeyID(acc.Issuer, acc.ID, request.Id)
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

func (e *V1) BlockKey(ctx context.Context, request *v1api.BlockKeyRequest, response *v1api.BlockKeyResponse) error {
	return e.updateKeyStatus(ctx, "v1api.BlockKey", request.Namespace, request.UserId, request.KeyId, keyStatusBlocked)
}

func (e *V1) UnblockKey(ctx context.Context, request *v1api.UnblockKeyRequest, response *v1api.UnblockKeyResponse) error {
	return e.updateKeyStatus(ctx, "v1api.UnblockKey", request.Namespace, request.UserId, request.KeyId, keyStatusActive)
}

func (e *V1) updateKeyStatus(ctx context.Context, methodName, ns, userID, keyID string, status keyStatus) error {

	if _, err := m3oauth.VerifyMicroAdmin(ctx, methodName); err != nil {
		return err
	}

	var keys []*apiKeyRecord
	if len(keyID) > 0 {
		rec, err := readAPIRecordByKeyID(ns, userID, keyID)
		if err != nil {
			log.Errorf("Error reading key %s", err)
			return errors.InternalServerError(methodName, "Error updating key")
		}
		keys = []*apiKeyRecord{rec}
	} else {
		recs, err := listKeysForUser(ns, userID)
		if err != nil {
			log.Errorf("Error listing keys %s", err)
			return errors.InternalServerError(methodName, "Error updating key")
		}
		keys = recs
	}
	for _, k := range keys {
		k.Status = status
		if err := writeAPIRecord(k); err != nil {
			log.Errorf("Error updating api key record %s", err)
			return errors.InternalServerError(methodName, "Error updating key")
		}
	}

	return nil
}

func (e *V1) listAPIs() ([]string, error) {
	rsp, err := e.papi.List(context.Background(), &publicapi.ListRequest{}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}

	ret := make([]string, len(rsp.Apis))
	for i, v := range rsp.Apis {
		ret[i] = v.Name
	}
	return ret, nil
}
