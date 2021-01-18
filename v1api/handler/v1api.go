package handler

import (
	"context"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"time"

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

type V1api struct{}

const (
	storePrefixHashedKey = "hashed"
	storePrefixUserID    = "user"
	storePrefixKeyID     = "key"
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
}

// Generate generates an API key
func (e *V1api) Generate(ctx context.Context, req *v1api.GenerateRequest, rsp *v1api.GenerateResponse) error {
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
	// generate a new API key

	// are they allowed to generate with the requested scopes?
	if !checkGenerateScopes(acc, req.Scopes) {
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
	// TODO - this should really be the same account but with additional
	// TODO - what issuer should we use?
	authAcc, err := auth.Generate(
		uuid.New().String(),
		auth.WithSecret(apiKey),
		auth.WithIssuer("foobar"),
		auth.WithType("apikey"),
		auth.WithScopes(req.Scopes...),
	)
	if err != nil {
		log.Errorf("Error generating auth account %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}
	tok, err := auth.Token(
		auth.WithCredentials(authAcc.ID, apiKey),
		auth.WithTokenIssuer("foobar"),
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
	}
	if err := writeAPIRecord(&rec); err != nil {
		log.Errorf("Failed to write api record %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
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

func hashSecret(s string) (string, error) {
	h := sha256.New()
	h.Write([]byte(s))
	h.Sum(nil)
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

// checkGenerateScopes returns true if account has sufficient privileges for them to generate the requestedScopes.
// e.g. micro "admin" can generate whatever scopes they want
func checkGenerateScopes(account *auth.Account, requestedScopes []string) bool {
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
func (e *V1api) Endpoint(ctx context.Context, req *pb.Request, rsp *pb.Response) error {
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
				auth.WithTokenIssuer("foobar"),
				auth.WithExpiry(1*time.Hour))
			if err != nil {
				log.Errorf("Error generating token %s", err)
				return errors.InternalServerError("v1api.generate", "Failed to generate api key")
			}
			apiRec.Token = tok.AccessToken
			if err := writeAPIRecord(&apiRec); err != nil {
				log.Errorf("Error updating API record %s", err)
				return err
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
func (e *V1api) ListKeys(ctx context.Context, req *v1api.ListRequest, rsp *v1api.ListResponse) error {
	// Check account
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("v1api.listkeys", "Unauthorized call to listkeys")
	}
	recs, err := store.Read("", store.Prefix(fmt.Sprintf("%s:%s:%s:", storePrefixUserID, acc.Issuer, acc.ID)))
	if err != nil {
		if err == store.ErrNotFound {
			return nil
		}
		log.Errorf("Error reading keys %s", err)
		return errors.InternalServerError("v1api.listkeys", "")
	}
	rsp.ApiKeys = make([]*v1api.APIKey, len(recs))
	for i, rec := range recs {
		apiRec := &apiKeyRecord{}
		if err := json.Unmarshal(rec.Value, apiRec); err != nil {
			log.Errorf("Error unmarshalling key %s", err)
			return errors.InternalServerError("v1api.listkeys", "")
		}
		rsp.ApiKeys[i] = &v1api.APIKey{
			Id:          apiRec.ID,
			Description: apiRec.Description,
			CreatedTime: apiRec.Created,
		}
	}
	return nil
}

// Revoke revokes a given key
func (e *V1api) Revoke(ctx context.Context, req *v1api.RevokeRequest, rsp *v1api.RevokeResponse) error {
	// Check account
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("v1api.revoke", "Unauthorized call to revoke")
	}
	if len(req.Id) == 0 {
		return errors.BadRequest("v1api.revoke", "Missing ID field")
	}
	recs, err := store.Read(fmt.Sprintf("%s:%s:%s:%s", storePrefixKeyID, acc.Issuer, acc.ID, req.Id))
	if err != nil {
		if err == store.ErrNotFound {
			return errors.NotFound("v1api.revoke", "Key not found")
		}
		log.Errorf("Error reading key %s", err)
		return errors.InternalServerError("v1api.revoke", "")
	}
	apiRec := &apiKeyRecord{}
	if err := json.Unmarshal(recs[0].Value, apiRec); err != nil {
		log.Errorf("Error marshalling key %s", err)
		return errors.InternalServerError("v1api.revoke", "")
	}

	if err := deleteAPIRecord(apiRec); err != nil {
		log.Errorf("Error deleting record %s", err)
		return errors.InternalServerError("v1api.revoke", "Error while deleting record")
	}

	return nil
}
