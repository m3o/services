package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/micro/micro/v3/service/events"

	"github.com/micro/micro/v3/service/client"

	"github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"

	v1api "github.com/m3o/services/v1api/proto"
	pb "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/context/metadata"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

type V1Api struct{}

const (
	storePrefixHashedKey = "hashed"
	storePrefixUserID    = "user"
	storePrefixKeyID     = "key"
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
func (e *V1Api) GenerateKey(ctx context.Context, req *v1api.GenerateKeyRequest, rsp *v1api.GenerateKeyResponse) error {
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
	if !checkRequestedScopes(acc, req.Scopes) {
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
func checkRequestedScopes(account *auth.Account, requestedScopes []string) bool {
	if account.Issuer == "micro" {
		for _, scope := range account.Scopes {
			if scope == "admin" {
				return true
			}
		}
	}
	// TODO take from config
	allowedScopes := map[string]bool{
		"location:write": true,
		"location:read":  true,
	}

	for _, scope := range requestedScopes {
		if !allowedScopes[scope] {
			return false
		}
	}
	return true
}

// Endpoint is a catch all for endpoints
func (e *V1Api) Endpoint(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
	// check api key
	authz := req.Header["Authorization"]
	if authz == nil || len(authz.Values) == 0 {
		return errors.Unauthorized("v1api", "Unauthorized")
	}

	// do lookup on hash of key
	key := authz.Values[0]
	if !strings.HasPrefix(key, "Bearer ") {
		return errors.Unauthorized("v1api", "Unauthorized")
	}
	key = key[7:]
	hashed, err := hashSecret(key)
	if err != nil {
		return errors.Unauthorized("v1api", "Unauthorized")
	}
	recs, err := store.Read(fmt.Sprintf("%s:%s", storePrefixHashedKey, hashed))
	if err != nil {
		if err != store.ErrNotFound {
			log.Errorf("Error while looking up api key %s", err)
		}
		// not found == invalid (or even revoked)
		log.Infof("Authz not found %+v", hashed)
		return errors.Unauthorized("v1api", "Unauthorized")
	}
	// rehydrate
	apiRec := apiKeyRecord{}
	if err := json.Unmarshal(recs[0].Value, &apiRec); err != nil {
		log.Errorf("Error while rehydrating api key record %s", err)
		return errors.Unauthorized("v1api", "Unauthorized")
	}

	// allowed? this *doesn't* check based on scopes, but whether the user has been specifically allowed (likely due to quota)
	allowed := false
	for prefix := range apiRec.AllowList {
		if strings.HasPrefix(req.Url, prefix) {
			allowed = true
		}
	}
	if !allowed {
		// TODO better error please
		return errors.Forbidden("v1api.blocked", "Client is blocked")
	}

	// do we need to refresh the token?
	tok, _, err := new(jwt.Parser).ParseUnverified(apiRec.Token, jwt.MapClaims{})
	if err != nil {
		log.Errorf("Error parsing existing jwt %s", err)
		return errors.Unauthorized("v1api", "Unauthorized")
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
				return errors.InternalServerError("v1api", "Failed to refresh api key")
			}
			apiRec.Token = tok.AccessToken
			if err := writeAPIRecord(&apiRec); err != nil {
				log.Errorf("Error updating API record %s", err)
				return errors.InternalServerError("v1api", "Failed to refresh api key")
			}
		}
	} else {
		log.Errorf("Error parsing existing jwt claims %s", err)
		return errors.Unauthorized("v1api", "Unauthorized")
	}

	// assume application/json for now
	ct := "application/json"

	// forward the request
	var payload json.RawMessage
	if len(req.Body) > 0 {
		payload = json.RawMessage(req.Body)
	}

	trimmedPath := strings.TrimPrefix(req.Path, "/v1/")
	parts := strings.Split(trimmedPath, "/")
	if len(parts) < 2 {
		// can't work out service and method
		return errors.NotFound("v1api", "")
	}

	service := parts[0]
	endpoint := fmt.Sprintf("%s.%s", strings.Title(parts[0]), strings.Title(parts[1]))
	request := client.DefaultClient.NewRequest(
		service,
		endpoint,
		&payload,
		client.WithContentType(ct),
	)

	// set the auth
	ctx = metadata.Set(ctx, "Authorization", fmt.Sprintf("Bearer %s", apiRec.Token))

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
			Url:       req.Url,
		}}); err != nil {
		log.Errorf("Error publishing event %s", err)
	}

	// marshal response
	// TODO implement errors
	b, err := response.MarshalJSON()
	if err != nil {
		return err
	}
	rsp.Body = string(b)
	return nil

}

// ListKeys lists all keys for a user
func (e *V1Api) ListKeys(ctx context.Context, req *v1api.ListRequest, rsp *v1api.ListResponse) error {
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
		}
	}
	return nil
}

func listKeysForUser(ns, userID string) ([]*apiKeyRecord, error) {
	recs, err := store.Read("", store.Prefix(fmt.Sprintf("%s:%s:%s:", storePrefixUserID, ns, userID)))
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

func (e *V1Api) RevokeKey(ctx context.Context, request *v1api.RevokeRequest, response *v1api.RevokeResponse) error {
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
	return nil
}

func (e *V1Api) UpdateAllowedPaths(ctx context.Context, request *v1api.UpdateAllowedPathsRequest, response *v1api.UpdateAllowedPathsResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("v1api.UpdateAllowedPaths", "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden("v1api.UpdateAllowedPaths", "Forbidden")
	}
	admin := false
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			admin = true
			break
		}
	}
	if !admin {
		return errors.Forbidden("v1api.UpdateAllowedPaths", "Forbidden")
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
			key.AllowList[a] = true
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
