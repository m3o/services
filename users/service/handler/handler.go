package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"regexp"
	"strings"
	"time"

	"github.com/google/uuid"

	"github.com/micro/go-micro/v2"
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
	store     store.Store
	publisher micro.Publisher
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) (*Handler, error) {
	// Return the initialised store
	return &Handler{
		store:     store.DefaultStore,
		publisher: micro.NewPublisher(srv.Name(), srv.Client()),
	}, nil
}

// Create inserts a new user into the the store
func (h *Handler) Create(ctx context.Context, req *pb.CreateRequest, rsp *pb.CreateResponse) error {
	// Validate the request
	if req.User == nil {
		return errors.BadRequest("go.micro.service.users", "User is missing")
	}

	// Check to see if the user already exists
	if user, err := h.findUser(req.User.Id); err == nil {
		rsp.User = user
		return nil
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
		return errors.InternalServerError("go.micro.service.users", "Coould not marshal user: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: user.Id, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.service.users", "Could not write to store: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserCreated,
		User: &user,
	})

	// Return the user and token in the response
	rsp.User = &user
	return nil
}

// Read retirves a user from the store
func (h *Handler) Read(ctx context.Context, req *pb.ReadRequest, rsp *pb.ReadResponse) error {
	var user *pb.User
	var err error

	if len(req.Email) > 0 {
		user, err = h.findUserByEmail(req.Email)
	} else {
		user, err = h.findUser(req.Id)
	}

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
		return errors.BadRequest("go.micro.service.users", "User is missing")
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
	user.Metadata = req.User.Metadata
	user.Updated = time.Now().Unix()
	user.ProfilePictureUrl = req.User.ProfilePictureUrl

	// Encode the updated user
	bytes, err := json.Marshal(user)
	if err != nil {
		return errors.InternalServerError("go.micro.service.users", "Coould not marshal user: %v", err)
	}

	// Write to the store
	if err := h.store.Write(&store.Record{Key: user.Id, Value: bytes}); err != nil {
		return errors.InternalServerError("go.micro.service.users", "Could not write to store: %v", err)
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
		return errors.InternalServerError("go.micro.service.users", "Could not write to store: %v", err)
	}

	// Publish the event
	go h.publisher.Publish(ctx, &pb.Event{
		Type: pb.EventType_UserDeleted,
		User: user,
	})

	return nil
}

// Search the users in th store, using full name
func (h *Handler) Search(ctx context.Context, req *pb.SearchRequest, rsp *pb.SearchResponse) error {
	// List all the records
	recs, err := h.store.Read("", store.ReadPrefix())
	if err != nil {
		return errors.InternalServerError("go.micro.service.users", "Could not read from store: %v", err)
	}

	// Decode the records
	users := make([]*pb.User, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &users[i]); err != nil {
			return errors.InternalServerError("go.micro.service.users", "Could not unmarshal user: %v", err)
		}
	}

	// Filter and return the users
	rsp.Users = make([]*pb.User, 0)
	for _, u := range users {
		fullname := fmt.Sprintf("%v %v", u.FirstName, u.LastName)
		if strings.Contains(fullname, req.Query) {
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
		return nil, errors.BadRequest("go.micro.service.users", "Missing ID")
	}

	// Get the records
	recs, err := h.store.Read(id)
	if err != nil {
		return nil, errors.InternalServerError("go.micro.service.users", "Could not read from store: %v", err)
	}
	if len(recs) == 0 {
		return nil, errors.NotFound("go.micro.service.users", "User not found")
	}
	if len(recs) > 1 {
		return nil, errors.InternalServerError("go.micro.service.users", "Store corrupted, %b records found for ID", len(recs))
	}

	// Decode the user
	var user *pb.User
	if err := json.Unmarshal(recs[0].Value, &user); err != nil {
		return nil, errors.InternalServerError("go.micro.service.users", "Could not unmarshal user: %v", err)
	}

	return user, nil
}

// findUserByEmail retrieves a user given an email
func (h *Handler) findUserByEmail(email string) (*pb.User, error) {
	// Validate the request
	if len(email) == 0 {
		return nil, errors.BadRequest("go.micro.service.users", "Missing Email")
	}

	// Get the records
	recs, err := h.store.Read("", store.ReadPrefix())
	if err != nil {
		return nil, errors.InternalServerError("go.micro.service.users", "Could not read from store: %v", err)
	}
	if len(recs) == 0 {
		return nil, errors.NotFound("go.micro.service.users", "User not found")
	}

	// Decode the users
	for _, r := range recs {
		var user *pb.User
		if err := json.Unmarshal(r.Value, &user); err != nil {
			return nil, errors.InternalServerError("go.micro.service.users", "Could not unmarshal user: %v", err)
		}
		if user.Email == email {
			return user, nil
		}
	}

	return nil, errors.NotFound("go.micro.service.users", "User not found")
}
