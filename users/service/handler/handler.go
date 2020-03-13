package handler

import (
	"context"
	"encoding/json"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/auth"
	"github.com/micro/go-micro/v2/errors"
	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/users/service/proto"
)

var (
	// URLSafeRegex is a function which returns true if a string is URL safe
	URLSafeRegex = regexp.MustCompile(`^[A-Za-z0-9_-].*?$`).MatchString
)

// Handler implements the users service interface
type Handler struct {
	auth      auth.Auth
	store     store.Store
	publisher micro.Publisher
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) (*Handler, error) {
	// create a new namespace in the default store
	s := store.DefaultStore
	if err := s.Init(store.Namespace(srv.Name())); err != nil {
		return nil, err
	}

	// Return the initialised store
	return &Handler{
		store:     s,
		auth:      srv.Options().Auth,
		publisher: micro.NewPublisher(srv.Name(), srv.Client()),
	}, nil
}

// Create inserts a new user into the the store
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// Validate the request
	if req.User == nil {
		return errors.BadRequest("go.micro.srv.users", "User is missing")
	}

	// Check to see if the user already exists
	if user, err := h.findUser(req.User.Id); err == nil {
		rsp.User = user

		if acc, err := h.auth.Generate(user.Id); err == nil {
			rsp.Token = acc.Token
		}

		return nil
	}

	// Validate the user
	if err := h.validateUser(req.User); err != nil {
		return err
	}

	// If validating only, return here
	if req.ValidateOnly {
		return nil
	}

	// Add the auto-generate fields
	var user pb.User = *req.User
	if len(user.Id) == 0 {
		// allow ID to be set for oauth providers
		user.Id = uuid.New().String()
	}
	user.Created = time.Now().Unix()
	user.Updated = time.Now().Unix()

	// Encode the user
	bytes, err := json.Marshal(user)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Coould not marshal user: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: user.Id, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Could not write to store: %v", err)
	}

	// Generate an auth account
	acc, err := h.auth.Generate(user.Id)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Could not generate auth account: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserCreated,
		User: &user,
	})

	// Return the user and token in the response
	rsp.User = &user
	rsp.Token = acc.Token
	return nil
}

// Read retirves a user from the store
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	user, err := h.findUser(req.Id)
	if err != nil {
		return err
	}

	rsp.User = user
	return nil
}

// Update modifies a user in the store
func (h *Handler) Update(ctx context.Context, req *pb.UpdateRequest, rsp *pb.UpdateResponse) error {
	// Validate the request
	if req.User == nil {
		return errors.BadRequest("go.micro.srv.users", "User is missing")
	}

	// Lookup the user
	user, err := h.findUser(req.User.Id)
	if err != nil {
		return err
	}

	// Update the user with the given attributes
	// TODO: Find a way which allows only updating a subset of attributes,
	// checking for blank values doesn't work since there needs to be a way
	// of unsetting attributes.
	user.FirstName = req.User.FirstName
	user.LastName = req.User.LastName
	user.Username = req.User.Username
	user.Email = req.User.Email
	user.Metadata = req.User.Metadata
	user.Updated = time.Now().Unix()

	// Validate the user
	if err := h.validateUser(req.User); err != nil {
		return err
	}

	// Encode the updated user
	bytes, err := json.Marshal(user)
	if err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Coould not marshal user: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: user.Id, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Could not write to store: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserUpdated,
		User: user,
	})

	// Return the user in the response
	rsp.User = user
	return nil
}

// Delete a user in the store
func (h *Handler) Delete(ctx context.Context, req *pb.DeleteRequest, rsp *pb.DeleteResponse) error {
	// Lookup the user
	user, err := h.findUser(req.Id)
	if err != nil {
		return err
	}

	// Delete from the store
	if err := h.store.Delete(user.Id); err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Could not write to store: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserDeleted,
		User: user,
	})

	return nil
}

// Search the users in th store, using username
func (h *Handler) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	// Validate the request
	if len(req.Username) == 0 {
		return errors.BadRequest("go.micro.srv.users", "Missing username")
	}

	// List all the records
	recs, err := h.store.Read("", store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.srv.users", "Could not read from store: %v", err)
	}

	// Decode the records
	users := make([]*pb.User, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &users[i]); err != nil {
			return errors.InternalServerError("go.micro.srv.users", "Could not unmarshal user: %v", err)
		}
	}

	// Filter and return the users
	rsp.Users = make([]*pb.User, 0)
	for _, u := range users {
		if strings.Contains(u.Username, req.Username) {
			rsp.Users = append(rsp.Users, u)
		}
	}

	return nil
}

// findUser retrieves a user given an ID. It is used by the Read, Update
// and Delete functions
func (h *Handler) findUser(id string) (*pb.User, error) {
	// Validate the request
	if len(id) == 0 {
		return nil, errors.BadRequest("go.micro.srv.users", "Missing ID")
	}

	// Get the records
	recs, err := h.store.Read(id, store.ReadPrefix())
	if err != nil {
		return nil, errors.InternalServerError("go.micro.srv.users", "Could not read from store: %v", err)
	}
	if len(recs) == 0 {
		return nil, errors.NotFound("go.micro.srv.users", "User not found")
	}
	if len(recs) > 1 {
		return nil, errors.InternalServerError("go.micro.srv.users", "Store corrupted, %b records found for ID", len(recs))
	}

	// Decode the user
	var user *pb.User
	if err := json.Unmarshal(recs[0].Value, &user); err != nil {
		return nil, errors.InternalServerError("go.micro.srv.users", "Could not unmarshal user: %v", err)
	}

	return user, nil
}

// usernameExists returns a bool if a user exists with this record,
// an error is also returned indicating there was an error reading
// from the store.
func (h *Handler) usernameExists(username string) (bool, error) {
	recs, err := h.store.Read(username, store.ReadSuffix())
	return len(recs) > 0, err
}

// validateUser performs some checks to ensure the validity of
// the data being written to the store. If the data is invalid
// a go-micro error is returned
func (h *Handler) validateUser(u *pb.User) error {
	if len(u.Username) == 0 {
		return nil
	}

	// Validate the username is url safe
	if safe := URLSafeRegex(u.Username); !safe {
		return errors.BadRequest("go.micro.srv.users", "Username is invalid, only a-Z, 0-9, dashes and underscores allowed")
	}

	// Ensure no other users with this username exist
	if exists, err := h.usernameExists(u.Username); err == nil && exists {
		return errors.BadRequest("go.micro.srv.users", "Username is taken")
	}

	return nil
}
