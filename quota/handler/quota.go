package handler

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/micro/micro/v3/service/client"

	v1api "github.com/m3o/services/v1api/proto"

	log "github.com/micro/micro/v3/service/logger"

	"github.com/micro/micro/v3/service/store"

	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/errors"

	pb "github.com/m3o/services/quota/proto"
)

type Quota struct {
	v1Svc v1api.V1ApiService
}

type resetFrequency int

const (
	Never resetFrequency = iota
	Daily
	Monthly
)

const prefixQuotaID = "id"

func (r resetFrequency) String() string {
	return [...]string{"Never", "Daily", "Monthly"}[r]
}

type quota struct {
	id             string
	limit          int64
	resetFrequency resetFrequency
}

func New(client client.Client) *Quota {
	q := &Quota{
		v1Svc: v1api.NewV1ApiService("v1", client),
	}
	go q.consumeEvents()
	return q
}

func (q *Quota) Create(ctx context.Context, request *pb.CreateRequest, response *pb.CreateResponse) error {
	if err := verifyAdmin(ctx, "quota.Create"); err != nil {
		return err
	}
	if len(request.Id) == 0 {
		return errors.BadRequest("quota.Create", "Missing quota ID")
	}
	quot := &quota{
		id:             request.Id,
		limit:          request.Limit,
		resetFrequency: resetFrequency(request.ResetFrequency.Number()),
	}

	b, err := json.Marshal(quot)
	if err != nil {
		log.Errorf("Error marshalling json %s", err)
		return errors.InternalServerError("quota.Create", "Error creating quota")
	}
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s", prefixQuotaID, quot.id),
		Value: b,
	}); err != nil {
		log.Errorf("Error writing to store %s", err)
		return errors.InternalServerError("quota.Create", "Error creating quota")
	}
	return nil
}

func verifyAdmin(ctx context.Context, method string) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden(method, "Forbidden")
	}
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			return nil
		}
	}
	return errors.Forbidden(method, "Forbidden")
}

func (q *Quota) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest, response *pb.RegisterUserResponse) error {
	if err := verifyAdmin(ctx, "quota.RegisterUser"); err != nil {
		return err
	}

	// store association for each quota
	// update the v1api to unblock the user's api keys

	panic("implement me")
}

type quotaEntry struct {
	id     string
	userID string
	value  int64
}

func (q *Quota) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	panic("implement me")
}
