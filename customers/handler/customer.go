package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	customer "github.com/m3o/services/customers/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	aproto "github.com/micro/micro/v3/proto/auth"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	mevents "github.com/micro/micro/v3/service/events"
	log "github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"
)

type Customers struct {
	accountsService   aproto.AccountsService
	namespacesService nsproto.NamespacesService
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
	Updated int64
}

func New(service *service.Service) *Customers {
	c := &Customers{
		accountsService:   aproto.NewAccountsService("auth", service.Client()),
		namespacesService: nsproto.NewNamespacesService("namespaces", service.Client()),
	}
	go c.consumeEvents()
	return c
}

func objToProto(cust *CustomerModel) *customer.Customer {
	return &customer.Customer{
		Id:      cust.ID,
		Status:  cust.Status,
		Created: cust.Created,
		Email:   cust.Email,
		Updated: cust.Updated,
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
	// have we seen this before?
	var cust *CustomerModel
	existingCust, err := readCustomerByEmail(email)
	if err != nil {
		if !strings.Contains(err.Error(), "not found") {
			return err
		}
		// not seen before so let's mint a new customer object
		cust = &CustomerModel{
			ID:     uuid.New().String(),
			Status: statusUnverified,
			Email:  email,
		}
	} else {
		if existingCust.Status == statusUnverified {
			// idempotency
			cust = existingCust
		} else {
			return errors.BadRequest("customers.create.exists", "Customer with this email already exists")
		}
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
	_, err := updateCustomerStatusByEmail(email, statusVerified)
	if err != nil {
		return err
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
	if cust.Status == statusDeleted {
		return nil, errors.NotFound("customers.read.notfound", "Customer not found")
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
	if err := c.deleteCustomer(ctx, request.Id); err != nil {
		log.Errorf("Error deleting customer %s %s", request.Id, err)
		return errors.InternalServerError("customers.delete", "Error deleting customer")
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
	ev := CustomerEvent{Customer: *cust, Type: "customers." + status}
	if err := mevents.Publish(custTopic, ev); err != nil {
		log.Errorf("Error publishing customers.%s event %+v", status, ev)
	}

	return cust, nil

}

func writeCustomer(cust *CustomerModel) error {
	now := time.Now().Unix()
	if cust.Created == 0 {
		cust.Created = now
	}
	cust.Updated = now
	b, _ := json.Marshal(*cust)
	if err := mstore.Write(&mstore.Record{
		Key:   prefixCustomer + cust.ID,
		Value: b,
	}); err != nil {
		return err
	}

	if err := mstore.Write(&mstore.Record{
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

func (c *Customers) deleteCustomer(ctx context.Context, customerID string) error {
	// auth accounts are tied to namespaces so get that list and delete all
	rsp, err := c.namespacesService.List(ctx, &nsproto.ListRequest{User: customerID}, client.WithAuthToken())
	if err != nil {
		return err
	}

	owned := []string{}
	for _, ns := range rsp.Namespaces {
		_, err := c.accountsService.Delete(ctx, &aproto.DeleteAccountRequest{
			Id:      customerID,
			Options: &aproto.Options{Namespace: ns.Id},
		}, client.WithAuthToken())
		if ignoreDeleteError(err) != nil {
			return err
		}
		// are we the owner
		if len(ns.Owners) == 1 && ns.Owners[0] == customerID {
			owned = append(owned, ns.Id)
		}
	}

	// delete any owned namespaces
	for _, ns := range owned {
		_, err := c.namespacesService.Delete(ctx, &nsproto.DeleteRequest{Id: ns}, client.WithAuthToken())
		if ignoreDeleteError(err) != nil {
			return err
		}
	}

	// delete customer
	cust, err := updateCustomerStatusByID(customerID, statusDeleted)
	if err != nil {
		return err
	}
	// fire deleted event
	ev := CustomerEvent{Customer: *cust, Type: "customers.deleted"}
	if err := mevents.Publish(custTopic, ev); err != nil {
		log.Errorf("Error publishing customers.deleted event %+v", ev)
	}
	return nil
}

// ignoreDeleteError will ignore any 400 or 404 errors returned, useful for idempotent deletes
func ignoreDeleteError(err error) error {
	if err != nil {
		merr, ok := err.(*errors.Error)
		if !ok {
			return err
		}
		if merr.Code == 400 || merr.Code == 404 {
			return nil
		}
		return err
	}
	return nil
}
