package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	eventsproto "github.com/micro/micro/v3/service/events/proto"

	"github.com/micro/go-micro/v3/events"
	mcontext "github.com/micro/micro/v3/service/context"

	log "github.com/micro/go-micro/v3/logger"

	namespace "github.com/m3o/services/namespaces/proto"
	plproto "github.com/m3o/services/platform/proto"
	"github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
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
	streamService   eventsproto.StreamService
}

func New(plSvc plproto.PlatformService, streamService eventsproto.StreamService) *Namespaces {
	return &Namespaces{
		platformService: plSvc,
		streamService:   streamService,
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

	return n.eventPublish(nsTopic, NamespaceEvent{Namespace: *ns, Type: "namespaces.created"})

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
	if (request.Owner == "" && request.User == "") || (request.Owner != "" && request.User != "") {
		return errors.BadRequest("namespaces.list.validation", "Only one of Owner or User should be specified")
	}
	id := request.Owner
	prefix := prefixOwner
	if id == "" {
		id = request.User
		prefix = prefixUser
	}
	recs, err := mstore.Read(prefix+id+"/", store.ReadPrefix())
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
	return n.eventPublish(nsTopic,
		NamespaceEvent{Namespace: *ns, Type: "namespaces.adduser"},
		events.WithMetadata(map[string]string{"user": request.User}))
}

func (n Namespaces) RemoveUser(ctx context.Context, request *namespace.RemoveUserRequest, response *namespace.RemoveUserResponse) error {
	return errors.InternalServerError("notimplemented", "not implemented")
}

// TODO remove this and replace with publish from micro/micro
func (n Namespaces) eventPublish(topic string, msg interface{}, opts ...events.PublishOption) error {
	// parse the options
	options := events.PublishOptions{
		Timestamp: time.Now(),
	}
	for _, o := range opts {
		o(&options)
	}

	// encode the message if it's not already encoded
	var payload []byte
	if p, ok := msg.([]byte); ok {
		payload = p
	} else {
		p, err := json.Marshal(msg)
		if err != nil {
			return events.ErrEncodingMessage
		}
		payload = p
	}

	// execute the RPC
	_, err := n.streamService.Publish(mcontext.DefaultContext, &eventsproto.PublishRequest{
		Topic:     topic,
		Payload:   payload,
		Metadata:  options.Metadata,
		Timestamp: options.Timestamp.Unix(),
	}, client.WithAuthToken())

	return err
}
