package handler

import (
	"context"
	"encoding/json"

	"golang.org/x/crypto/bcrypt"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"

	pb "github.com/micro/services/login/service/proto/login"
	users "github.com/micro/services/users/service/proto"
)

// Handler implements the login service interface
type Handler struct {
	name  string
	store store.Store
}

// NewHandler returns an initialise handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:  srv.Name(),
		store: store.DefaultStore,
	}
}

// HandleUserEvent handles the events published by the uses service
func (h *Handler) HandleUserEvent(ctx context.Context, event *users.Event) error {
	switch event.Type {
	case users.EventType_UserDeleted:
		err := h.store.Delete(event.User.Email)
		if err != nil && err != store.ErrNotFound {
			return err
		}
	case users.EventType_UserUpdated:
		// TODO: If the email changed, read from the old email and write to the new one
	}

	return nil
}

// CreateLogin generates a set of credentials
func (h *Handler) CreateLogin(ctx context.Context, req *pb.CreateLoginRequest, rsp *pb.CreateLoginResponse) error {
	// Validate an email was provided
	if len(req.Email) == 0 {
		return errors.BadRequest(h.name, "Email required")
	}

	// Validate credentials don't exist for this email already
	if _, err := h.store.Read(req.Email); err != nil && err != store.ErrNotFound {
		return errors.InternalServerError(h.name, "Unable to read from store: %v", err)
	}

	// Validate the password
	if len(req.Password) < 6 {
		return errors.BadRequest(h.name, "Password must be at least 6 chars long")
	}

	// Hash the password
	pass, err := generatePassword(req.Password)
	if err != nil {
		return errors.BadRequest(h.name, "Password is invalid")
	}

	// Return at this point if the request is validate_only, since the caller will
	// at this point, create a user in the users-srv, then make this request again
	// with an id and validate_only = false
	if req.ValidateOnly {
		return nil
	}

	// Validate we have an ID
	if len(req.Id) == 0 {
		return errors.BadRequest(h.name, "ID required for non validate_only request")
	}

	// Serialize the user
	user := User{
		ID:           req.Id,
		PasswordHash: pass,
	}

	// Marshal to json
	bytes, err := json.Marshal(user)
	if err != nil {
		return errors.InternalServerError(h.name, "Unable to marshal JSON: %v", err)
	}

	// Write to the store
	err = h.store.Write(&store.Record{
		Key:   req.Email,
		Value: bytes,
	})
	if err != nil {
		return errors.InternalServerError(h.name, "Unable to write to store: %v", err)
	}

	return nil
}

// VerifyLogin validates a set of credentials
func (h *Handler) VerifyLogin(ctx context.Context, req *pb.VerifyLoginRequest, rsp *pb.VerifyLoginResponse) error {
	// Look up the user
	recs, err := h.store.Read(req.Email)
	if err == store.ErrNotFound || len(recs) != 1 {
		return errors.BadRequest(h.name, "Invalid Email")
	} else if err != nil {
		return errors.InternalServerError(h.name, "Unable to read from store: %v", err)
	}

	// Deserialize the user
	var user *User
	if err := json.Unmarshal(recs[0].Value, &user); err != nil {
		return errors.InternalServerError(h.name, "Unable to unmarshal JSON: %v", err)
	}

	// Compare the passwords
	if !comparePasswords(user.PasswordHash, req.Password) {
		return errors.BadRequest(h.name, "Incorrect Password")
	}

	// The user has been validates, return their ID
	rsp.Id = user.ID
	return nil
}

// UpdateEmail changes the users email, deleting the old key and writing to the new one
func (h *Handler) UpdateEmail(ctx context.Context, req *pb.UpdateEmailRequest, rsp *pb.UpdateEmailResponse) error {
	// Look up the user
	recs, err := h.store.Read(req.OldEmail)
	if err == store.ErrNotFound || len(recs) != 1 {
		return nil
	} else if err != nil {
		return errors.InternalServerError(h.name, "Unable to read from store: %v", err)
	}

	// Delete the old key
	if err := h.store.Delete(req.OldEmail); err != nil {
		return errors.InternalServerError(h.name, "Unable to delete from store: %v", err)
	}

	// Write the new key
	record := store.Record{Key: req.NewEmail, Value: recs[0].Value}
	if err := h.store.Write(&record); err != nil {
		return errors.InternalServerError(h.name, "Unable to write to store: %v", err)
	}

	return nil
}

// User is an object persisted in the store
type User struct {
	ID           string `json:"id"`
	PasswordHash string `json:"password"`
}

func generatePassword(s string) (string, error) {
	saltedBytes := []byte(s)
	hashedBytes, err := bcrypt.GenerateFromPassword(saltedBytes, bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	hash := string(hashedBytes[:])
	return hash, nil
}

func comparePasswords(hash string, s string) bool {
	incoming := []byte(s)
	existing := []byte(hash)
	return bcrypt.CompareHashAndPassword(existing, incoming) == nil
}
