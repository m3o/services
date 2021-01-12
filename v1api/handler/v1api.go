package handler

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"

	v1api "github.com/m3o/services/v1api/proto"
	pb "github.com/micro/micro/v3/proto/api"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

type V1api struct{}

type apiKeyRecord struct {
	ApiKey      string   // hashed api key
	Scopes      []string // the scopes this key has granted
	UserID      string   // the ID of the key's owner
	Description string   // optional description of the API key as given by user
	Namespace   string   // the namespace that this user belongs to (only because technically user IDs aren't globally unique)
}

// Generate generates an API key
func (e *V1api) Generate(ctx context.Context, req *v1api.GenerateRequest, rsp *v1api.GenerateResponse) error {
	if len(req.Scopes) == 0 {
		return errors.BadRequest("v1api.generate", "Scopes are mandatory")
	}
	if len(req.Description) == 0 {
		return errors.BadRequest("v1api.generate", "Description is mandatory")
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

	// hash API key and store with scopes
	rec := apiKeyRecord{
		ApiKey:      hashedKey,
		Scopes:      req.Scopes,
		UserID:      acc.ID,
		Namespace:   acc.Issuer,
		Description: req.Description,
	}
	b, err := json.Marshal(rec)
	if err != nil {
		log.Errorf("Error marshalling api key record %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}

	// store under hashed key for API usage
	if err := store.Write(&store.Record{
		Key:   hashedKey,
		Value: b,
	}); err != nil {
		log.Errorf("Error storing api key record %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}

	// store under the user ID for retrieval on dashboard
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s:%s", rec.Namespace, rec.UserID, hashedKey),
		Value: b,
	}); err != nil {
		log.Errorf("Error storing api key record %s", err)
		return errors.InternalServerError("v1api.generate", "Failed to generate api key")
	}

	rsp.ApiKey = apiKey

	return nil
}

func hashSecret(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func secretsMatch(hash string, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming) == nil
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
	log.Info("Received V1api.Call request")

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
	recs, err := store.Read(hashed)
	if err != nil {
		if err != store.ErrNotFound {
			log.Errorf("Error while looking up api key %s", err)
		}
		// not found == invalid (or even revoked)
		return errors.Unauthorized("v1api", "Unauthorized")
	}
	// rehydrate
	apiRec := apiKeyRecord{}
	if err := json.Unmarshal(recs[0].Value, &apiRec); err != nil {
		log.Errorf("Error while rehydrating api key record %s", err)
		return errors.Unauthorized("v1api", "Unauthorized")
	}

	// add scopes to context

	// send
	//client.Call(ctx, req, rsp)

	return nil
}
