package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	log "github.com/micro/go-micro/v3/logger"

	"github.com/micro/go-micro/v3/auth"

	customer "github.com/m3o/services/customers/proto"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	mevents "github.com/micro/micro/v3/service/events"
	mstore "github.com/micro/micro/v3/service/store"
)

type Customers struct {
}

const (
	statusUnverified = "unverified"
	statusVerified   = "verified"
	statusActive     = "active"
	statusDeleted    = "deleted"

	prefixCustomer      = "customers/"
	prefixCustomerEmail = "email/"
	custTopic           = "customers"
)

type CustomerModel struct {
	ID      string
	Email   string
	Status  string
	Created int64
}

func New() *Customers {
	return &Customers{}
}

func objToProto(cust *CustomerModel) *customer.Customer {
	return &customer.Customer{
		Id:      cust.ID,
		Status:  cust.Status,
		Created: cust.Created,
		Email:   cust.Email,
	}
}

func (c *Customers) Create(ctx context.Context, request *customer.CreateRequest, response *customer.CreateResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	email := request.Email
	if email == "" {
		// try deprecated fallback
		email = request.Id
	}
	if strings.TrimSpace(email) == "" {
		return errors.BadRequest("customers.create", "Email is required")
	}
	cust := &CustomerModel{
		ID:      uuid.New().String(),
		Status:  statusUnverified,
		Created: time.Now().Unix(),
		Email:   email,
	}
	if err := writeCustomer(cust); err != nil {
		return err
	}
	response.Customer = objToProto(cust)
	ev := CustomerEvent{Customer: *cust, Type: "customers.created"}
	if err := mevents.Publish(custTopic, ev); err != nil {
		log.Errorf("Error publishing customers.created event %+v", ev)
	}
	return nil
}

func (c *Customers) MarkVerified(ctx context.Context, request *customer.MarkVerifiedRequest, response *customer.MarkVerifiedResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	email := request.Email
	if email == "" {
		// try deprecated fallback
		email = request.Id
	}

	if strings.TrimSpace(email) == "" {
		return errors.BadRequest("customers.markverified", "Email is required")
	}
	cust, err := updateCustomerStatusByEmail(email, statusVerified)
	if err != nil {
		return err
	}
	ev := CustomerEvent{Customer: *cust, Type: "customers.verified"}
	if err := mevents.Publish(custTopic, ev); err != nil {
		log.Errorf("Error publishing customers.verified event %+v", ev)
	}
	return nil
}

func readCustomerByID(customerID string) (*CustomerModel, error) {
	return readCustomer(customerID, prefixCustomer)
}

func readCustomerByEmail(email string) (*CustomerModel, error) {
	return readCustomer(email, prefixCustomerEmail)
}

func readCustomer(id, prefix string) (*CustomerModel, error) {
	recs, err := mstore.Read(prefix + id)
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
	// TODO at some point we'll need to relax this
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	if strings.TrimSpace(request.Id) == "" && strings.TrimSpace(request.Email) == "" {
		return errors.BadRequest("customers.read", "ID or Email is required")
	}
	var cust *CustomerModel
	var err error
	if request.Id != "" {
		cust, err = readCustomerByID(request.Id)
	} else {
		cust, err = readCustomerByEmail(request.Email)
	}
	if err != nil {
		return err
	}
	response.Customer = objToProto(cust)
	return nil
}

func (c *Customers) Delete(ctx context.Context, request *customer.DeleteRequest, response *customer.DeleteResponse) error {
	if err := authorizeCall(ctx); err != nil {
		return err
	}
	if strings.TrimSpace(request.Id) == "" {
		return errors.BadRequest("customers.delete", "ID is required")
	}
	cust, err := updateCustomerStatusByID(request.Id, statusDeleted)
	if err != nil {
		return err
	}
	ev := CustomerEvent{Customer: *cust, Type: "customers.deleted"}
	if err := mevents.Publish(custTopic, ev); err != nil {
		log.Errorf("Error publishing customers.deleted event %+v", ev)
	}
	return nil
}

func updateCustomerStatusByEmail(email, status string) (*CustomerModel, error) {
	return updateCustomerStatus(email, status, prefixCustomerEmail)
}

func updateCustomerStatusByID(id, status string) (*CustomerModel, error) {
	return updateCustomerStatus(id, status, prefixCustomer)
}

func updateCustomerStatus(id, status, prefix string) (*CustomerModel, error) {
	cust, err := readCustomer(id, prefix)
	if err != nil {
		return nil, err
	}
	cust.Status = status
	if err := writeCustomer(cust); err != nil {
		return nil, err
	}
	return cust, nil

}

func writeCustomer(cust *CustomerModel) error {
	b, _ := json.Marshal(*cust)

	if err := mstore.Write(&store.Record{
		Key:   prefixCustomer + cust.ID,
		Value: b,
	}); err != nil {
		return err
	}

	if err := mstore.Write(&store.Record{
		Key:   prefixCustomerEmail + cust.Email,
		Value: b,
	}); err != nil {
		return err
	}
	return nil
}

func authorizeCall(ctx context.Context) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("customers", "Unauthorized request")
	}
	if account.Issuer != "micro" {
		return errors.Unauthorized("customers", "Unauthorized request")
	}
	return nil
}
