package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	customer "github.com/m3o/services/customer/proto"

	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/logger"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/events"
	mstore "github.com/micro/micro/v3/service/store"
)

type Customer struct{}

const (
	statusUnverified = "unverified"
	statusVerified   = "verified"
	statusActive     = "active"
	statusDeleted    = "deleted"

	prefixCustomer = "customer/"
)

type CustomerModel struct {
	ID      string
	Status  string
	Created int64
}

func objToProto(cust *CustomerModel) *customer.CustomerMessage {
	return &customer.CustomerMessage{
		Id:      cust.ID,
		Status:  cust.Status,
		Created: cust.Created,
	}
}

func (c *Customer) Create(ctx context.Context, request *customer.CreateRequest, response *customer.CreateResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customer.create", "ID is required")
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
	logger.Infof("Writing %d %s", len(b), string(b))
	if err := mstore.Write(&store.Record{
		Key:   prefixCustomer + cust.ID,
		Value: b,
	}); err != nil {
		return err
	}
	response.Customer = objToProto(cust)
	return events.Publish("customer", CustomerEvent{Customer: *cust, Type: "customer.created"})
}

func (c *Customer) MarkVerified(ctx context.Context, request *customer.MarkVerifiedRequest, response *customer.MarkVerifiedResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customer.create", "ID is required")
	}
	cust, err := updateCustomerStatus(request.Id, statusVerified)
	if err != nil {
		return err
	}
	return events.Publish("customer", CustomerEvent{Customer: *cust, Type: "customer.verified"})
}

func readCustomer(customerID string) (*CustomerModel, error) {
	recs, err := mstore.Read(prefixCustomer + customerID)
	if err != nil {
		return nil, err
	}
	if len(recs) != 1 {
		return nil, errors.InternalServerError("customer.read.toomanyrecords", "Cannot find record to update")
	}
	rec := recs[0]
	cust := &CustomerModel{}
	if err := json.Unmarshal(rec.Value, cust); err != nil {
		return nil, err
	}
	return cust, nil
}

func (c *Customer) Read(ctx context.Context, request *customer.ReadRequest, response *customer.ReadResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customer.create", "ID is required")
	}
	cust, err := readCustomer(request.Id)
	if err != nil {
		return err
	}
	logger.Infof("Read %+v", *cust)
	response.Customer = objToProto(cust)
	// TODO fill out subscription and namespaces
	return nil
}

func (c *Customer) Delete(ctx context.Context, request *customer.DeleteRequest, response *customer.DeleteResponse) error {
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customer.create", "ID is required")
	}
	cust, err := updateCustomerStatus(request.Id, statusDeleted)
	if err != nil {
		return err
	}
	return events.Publish("customer", CustomerEvent{Customer: *cust, Type: "customer.deleted"})
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
