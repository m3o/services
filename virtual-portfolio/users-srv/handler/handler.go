package handler

import (
	"context"

	"github.com/kytra-app/helpers/microgorm"
	hash "github.com/kytra-app/helpers/passwordhasher"
	proto "github.com/kytra-app/users-srv/proto"
	"github.com/kytra-app/users-srv/storage"
	"github.com/micro/go-micro/errors"
)

// New returns an instance of Handler
func New(storage storage.Service) *Handler {
	return &Handler{db: storage}
}

// Handler is an object can process RPC requests
type Handler struct{ db storage.Service }

// Create creates a new User object
func (h *Handler) Create(ctx context.Context, req *proto.User, rsp *proto.User) error {
	if len(req.Password) < 8 {
		return errors.BadRequest("INVALID_PASSWORD", "Password must be at least 8 characters long")
	}

	params, err := h.reverseSerializeUser(*req)
	if err != nil {
		return err
	}

	u, err := h.db.Create(params)
	if err != nil {
		return err
	}

	*rsp = h.serializeUser(u)
	return nil
}

// Find looks up a User using the attributes set on the request
func (h *Handler) Find(ctx context.Context, req *proto.User, rsp *proto.User) error {
	if req.Uuid == "" && req.Email == "" && req.Username == "" && req.PhoneNumber == "" {
		return microgorm.ErrNotFound
	}

	u, err := h.db.Find(storage.User{
		UUID:        req.Uuid,
		Email:       req.Email,
		Username:    req.Username,
		PhoneNumber: req.PhoneNumber,
	})

	if err != nil {
		return err
	}

	*rsp = h.serializeUser(u)
	return nil
}

// Count returns the number of users in the database
func (h *Handler) Count(ctx context.Context, req *proto.CountRequest, rsp *proto.CountResponse) error {
	count, err := h.db.Count()
	if err != nil {
		return err
	}

	rsp.Count = count
	return nil
}

// ValidatePassword compares the password provided to that stored as a hash in the database
func (h *Handler) ValidatePassword(ctx context.Context, req *proto.User, rsp *proto.User) error {
	u, err := h.db.Find(storage.User{Email: req.Email})

	if err != nil || len(req.Email) == 0 {
		return errors.BadRequest("INVALID_EMAIL", "No user with this email has been found")
	}

	if !hash.Compare(u.Password, req.Password) {
		return errors.BadRequest("INCORRECT_PASSWORD", "The password does not match")
	}

	*rsp = h.serializeUser(u)
	return nil
}

// Update finds the User with the UUID provided and updates that object
func (h *Handler) Update(ctx context.Context, req *proto.User, rsp *proto.User) error {
	params, err := h.reverseSerializeUser(*req)
	if err != nil {
		return err
	}

	u, err := h.db.Update(params)
	if err != nil {
		return err
	}

	*rsp = h.serializeUser(u)
	return nil
}

// List returns an array of users which match the UUIDs or PhoneNumbers provided in the request
func (h *Handler) List(ctx context.Context, req *proto.ListRequest, rsp *proto.ListResponse) error {
	var users []storage.User
	var err error

	if len(req.Uuids) > 0 {
		users, err = h.db.List(req.Uuids)
	} else if len(req.PhoneNumbers) > 0 {
		users, err = h.db.ListByPhoneNumber(req.PhoneNumbers)
	}

	if err != nil {
		return err
	}

	rsp.Users = make([]*proto.User, len(users))
	for i, u := range users {
		serialized := h.serializeUser(u)
		rsp.Users[i] = &serialized
	}

	return nil
}

// All returns an array of all users
func (h *Handler) All(ctx context.Context, req *proto.AllRequest, rsp *proto.ListResponse) error {
	users, err := h.db.All()
	if err != nil {
		return err
	}

	rsp.Users = make([]*proto.User, len(users))
	for i, u := range users {
		serialized := h.serializeUser(u)
		rsp.Users[i] = &serialized
	}

	return nil
}

// Search returns all the users matching the query provided
func (h *Handler) Search(ctx context.Context, req *proto.SearchRequest, rsp *proto.ListResponse) error {
	if len(req.Query) == 0 {
		return errors.BadRequest("INVALID_QUERY", "Query can't be blank")
	}

	limit := req.Limit
	if limit == 0 {
		limit = 30
	}

	users, err := h.db.Query(req.Query, limit)
	if err != nil {
		return err
	}

	rsp.Users = make([]*proto.User, len(users))
	for i, u := range users {
		serialized := h.serializeUser(u)
		rsp.Users[i] = &serialized
	}

	return nil
}

func (h *Handler) reverseSerializeUser(req proto.User) (storage.User, error) {
	user := storage.User{
		UUID:             req.Uuid,
		FirstName:        req.FirstName,
		LastName:         req.LastName,
		PhoneNumber:      req.PhoneNumber,
		Email:            req.Email,
		Username:         req.Username,
		ProfilePictureID: req.ProfilePictureId,
	}

	if len(req.Password) == 0 {
		return user, nil
	}

	var err error
	user.Password, err = hash.Generate(req.Password)

	if err != nil {
		return user, errors.InternalServerError("HASHING_PASSWORD", "Error hashing password")
	}

	return user, nil
}

func (h *Handler) serializeUser(u storage.User) proto.User {
	return proto.User{
		Uuid:             u.UUID,
		FirstName:        u.FirstName,
		LastName:         u.LastName,
		Username:         u.Username,
		PhoneNumber:      u.PhoneNumber,
		Email:            u.Email,
		ProfilePictureId: u.ProfilePictureID,
		Admin:            u.Admin,
		CreatedAt:        u.CreatedAt.Unix(),
	}
}
