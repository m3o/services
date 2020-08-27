package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	customer "github.com/m3o/services/customers/proto"
	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/events"
	"github.com/micro/go-micro/v3/store"
	mcontext "github.com/micro/micro/v3/service/context"
	eventsproto "github.com/micro/micro/v3/service/events/proto"
	mstore "github.com/micro/micro/v3/service/store"
)

type Customers struct {
	streamService eventsproto.StreamService
}

const (
	statusUnverified = "unverified"
	statusVerified   = "verified"
	statusActive     = "active"
	statusDeleted    = "deleted"

	prefixCustomer = "customers/"
	custTopic      = "customers"
)

type CustomerModel struct {
	ID      string
	Status  string
	Created int64
}

func New(streamService eventsproto.StreamService) *Customers {
	return &Customers{
		streamService: streamService,
	}
}

func objToProto(cust *CustomerModel) *customer.Customer {
	return &customer.Customer{
		Id:      cust.ID,
		Status:  cust.Status,
		Created: cust.Created,
	}
}

func (c *Customers) Create(ctx context.Context, request *customer.CreateRequest, response *customer.CreateResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customers.create", "ID is required")
	}
	cust := &CustomerModel{
		ID:      request.Id,
		Status:  statusUnverified,
		Created: time.Now().Unix(),
	}
	b, err := json.Marshal(*cust)
	if err != nil {
		return err
	}
	if err := mstore.Write(&store.Record{
		Key:   prefixCustomer + cust.ID,
		Value: b,
	}); err != nil {
		return err
	}
	response.Customer = objToProto(cust)

	return c.eventPublish(custTopic, CustomerEvent{Customer: *cust, Type: "customers.created"})
	//return nil
}

func (c *Customers) MarkVerified(ctx context.Context, request *customer.MarkVerifiedRequest, response *customer.MarkVerifiedResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customers.create", "ID is required")
	}
	cust, err := updateCustomerStatus(request.Id, statusVerified)
	if err != nil {
		return err
	}
	return c.eventPublish(custTopic, CustomerEvent{Customer: *cust, Type: "customers.verified"})
}

func readCustomer(customerID string) (*CustomerModel, error) {
	recs, err := mstore.Read(prefixCustomer + customerID)
	if err != nil {
		return nil, err
	}
	if len(recs) != 1 {
		return nil, errors.InternalServerError("customers.read.toomanyrecords", "Cannot find record to update")
	}
	rec := recs[0]
	cust := &CustomerModel{}
	if err := json.Unmarshal(rec.Value, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

func (c *Customers) Read(ctx context.Context, request *customer.ReadRequest, response *customer.ReadResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customers.create", "ID is required")
	}
	cust, err := readCustomer(request.Id)
	if err != nil {
		return err
	}
	response.Customer = objToProto(cust)
	// TODO fill out subscription and namespaces
	return nil
}

func (c *Customers) Delete(ctx context.Context, request *customer.DeleteRequest, response *customer.DeleteResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customers.create", "ID is required")
	}
	cust, err := updateCustomerStatus(request.Id, statusDeleted)
	if err != nil {
		return err
	}
	return c.eventPublish(custTopic, CustomerEvent{Customer: *cust, Type: "customers.deleted"})
}

func updateCustomerStatus(customerID, status string) (*CustomerModel, error) {
	cust, err := readCustomer(customerID)
	if err != nil {
		return nil, err
	}
	cust.Status = status
	b, _ := json.Marshal(*cust)

	if err := mstore.Write(&store.Record{
		Key:   prefixCustomer + cust.ID,
		Value: b,
	}); err != nil {
		return nil, err
	}
	return cust, nil
}

// TODO remove this and replace with publish from micro/micro
func (c *Customers) eventPublish(topic string, msg interface{}, opts ...events.PublishOption) error {
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
	_, err := c.streamService.Publish(mcontext.DefaultContext, &eventsproto.PublishRequest{
		Topic:     topic,
		Payload:   payload,
		Metadata:  options.Metadata,
		Timestamp: options.Timestamp.Unix(),
	}, goclient.WithAuthToken())

	return err
}
