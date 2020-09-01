package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/micro/go-micro/v3/auth"

	namespace "github.com/m3o/services/namespaces/proto"
	plproto "github.com/m3o/services/platform/proto"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/events"
	log "github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	mevents "github.com/micro/micro/v3/service/events"
	mstore "github.com/micro/micro/v3/service/store"

	"github.com/sethvargo/go-diceware/diceware"
)

const (
	prefixNs    = "namespace/"
	prefixOwner = "owner/"
	prefixUser  = "user/"

	nsTopic = "namespaces"
)

type Namespaces struct {
	platformService plproto.PlatformService
}

func New(plSvc plproto.PlatformService) *Namespaces {
	return &Namespaces{
		platformService: plSvc,
	}
}

type NamespaceModel struct {
	ID      string
	Owners  []string
	Users   []string
	Created int64
}

func objToProto(ns *NamespaceModel) *namespace.Namespace {
	return &namespace.Namespace{
		Id:      ns.ID,
		Created: ns.Created,
		Owners:  ns.Owners,
		Users:   ns.Users,
	}
}

func (n Namespaces) Create(ctx context.Context, request *namespace.CreateRequest, response *namespace.CreateResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	if len(request.Owners) == 0 {
		return errors.BadRequest("namespaces.create.validation", "Owners is required")
	}

	id := request.Id
	if id == "" {
		list, err := diceware.Generate(3)
		if err != nil {
			return errors.InternalServerError("namespaces.create.name", "Error generating name for new namespace")
		}
		id = strings.Join(list, "-")
	}
	ns := &NamespaceModel{
		ID:      id,
		Owners:  request.Owners,
		Users:   request.Owners,
		Created: time.Now().Unix(),
	}
	_, err := n.platformService.CreateNamespace(ctx, &plproto.CreateNamespaceRequest{
		Name: ns.ID,
	}, client.WithRequestTimeout(10*time.Second), client.WithAuthToken())
	if err != nil {
		log.Errorf("Error creating namespace %s", err)
		return errors.InternalServerError("namespaces.create.creation", "Error creating namespace")
	}
	err = writeNamespace(ns)
	if err != nil {
		return err
	}
	response.Namespace = objToProto(ns)

	return mevents.Publish(nsTopic, NamespaceEvent{Namespace: *ns, Type: "namespaces.created"})

}

// writeNamespace writes to the store. We deliberately denormalise/duplicate across many indexes to optimise for reads
func writeNamespace(ns *NamespaceModel) error {
	b, err := json.Marshal(*ns)
	if err != nil {
		return err
	}
	if err := mstore.Write(&store.Record{
		Key:   prefixNs + ns.ID,
		Value: b,
	}); err != nil {
		return err
	}
	// index by owner
	for _, owner := range ns.Owners {
		if err := mstore.Write(&store.Record{
			Key:   prefixOwner + owner + "/" + ns.ID,
			Value: b,
		}); err != nil {
			return err
		}
	}
	// index by user
	for _, user := range ns.Users {
		if err := mstore.Write(&store.Record{
			Key:   prefixUser + user + "/" + ns.ID,
			Value: b,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (n Namespaces) Read(ctx context.Context, request *namespace.ReadRequest, response *namespace.ReadResponse) error {
	// TODO at some point we'll probably want to relax this
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	if request.Id == "" {
		return errors.BadRequest("namespaces.read.validation", "ID is required")
	}
	ns, err := readNamespace(request.Id)
	if err != nil {
		return err
	}
	response.Namespace = objToProto(ns)
	return nil
}

func readNamespace(id string) (*NamespaceModel, error) {
	recs, err := mstore.Read(prefixNs + id)
	if err != nil {
		return nil, err
	}
	if len(recs) != 1 {
		return nil, errors.InternalServerError("customers.read.toomanyrecords", "Cannot find record to update")
	}
	rec := recs[0]
	ns := &NamespaceModel{}
	if err := json.Unmarshal(rec.Value, ns); err != nil {
		return nil, err
	}
	return ns, nil
}

func (n Namespaces) Delete(ctx context.Context, request *namespace.DeleteRequest, response *namespace.DeleteResponse) error {
	return errors.InternalServerError("notimplemented", "not implemented")
}

func (n Namespaces) List(ctx context.Context, request *namespace.ListRequest, response *namespace.ListResponse) error {
	// TODO at some point we'll want to relax this
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	if request.Owner != "" && request.User != "" {
		return errors.BadRequest("namespaces.list.validation", "Only one of Owner or User should be specified")
	}
	key := ""
	switch {
	case request.Owner == "" && request.User == "":
		key = prefixNs
	case request.Owner != "":
		key = prefixOwner + request.Owner + "/"
	case request.User != "":
		key = prefixUser + request.User + "/"
	}

	recs, err := mstore.Read(key, store.ReadPrefix())
	if err != nil {
		return err
	}
	res := make([]*namespace.Namespace, len(recs))
	for i, rec := range recs {
		ns := &NamespaceModel{}
		if err := json.Unmarshal(rec.Value, ns); err != nil {
			return err
		}
		res[i] = objToProto(ns)
	}
	response.Namespaces = res
	return nil
}

func (n Namespaces) AddUser(ctx context.Context, request *namespace.AddUserRequest, response *namespace.AddUserResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	if request.Namespace == "" || request.User == "" {
		return errors.BadRequest("namespaces.adduser.validation", "User and Namespace are required")
	}
	ns, err := readNamespace(request.Namespace)
	if err != nil {
		return err
	}
	// quick check we haven't already added this user
	for _, user := range ns.Users {
		if user == request.User {
			// idempotent, just return success
			return nil
		}
	}
	ns.Users = append(ns.Users, request.User)
	// write it
	if err := writeNamespace(ns); err != nil {
		return err
	}
	// TODO anything else we need to do for adding a user to namespace?
	return mevents.Publish(nsTopic,
		NamespaceEvent{Namespace: *ns, Type: "namespaces.adduser"},
		events.WithMetadata(map[string]string{"user": request.User}))
}

func (n Namespaces) RemoveUser(ctx context.Context, request *namespace.RemoveUserRequest, response *namespace.RemoveUserResponse) error {
	return errors.InternalServerError("notimplemented", "not implemented")
}

func authorizeCall(ctx context.Context) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("namespaces", "Unauthorized request")
	}
	if account.Issuer != "micro" {
		return errors.Unauthorized("namespaces", "Unauthorized request")
	}
	return nil
}
