package handler

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/google/uuid"

	customer "github.com/m3o/services/customers/proto"
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
	accountsService aproto.AccountsService
}

const (
	statusUnverified = "unverified"
	statusVerified   = "verified"
	statusActive     = "active"
	statusDeleted    = "deleted"

	prefixCustomer      = "customers/"
	prefixCustomerEmail = "email/"
)

var validStatus = map[string]bool{
	statusUnverified: true,
	statusVerified:   true,
	statusActive:     true,
	statusDeleted:    true,
}

type CustomerModel struct {
	ID      string
	Email   string
	Status  string
	Created int64
	Updated int64
}

func New(service *service.Service) *Customers {
	c := &Customers{
		accountsService: aproto.NewAccountsService("auth", service.Client()),
	}
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
	if err := authorizeCall(ctx, ""); err != nil {
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

	// Publish the event
	var callerID string
	if acc, ok := auth.AccountFromContext(ctx); ok {
		callerID = acc.ID
	}
	ev := &customer.Event{
		Type:     customer.EventType_EventTypeCreated,
		Customer: response.Customer,
		CallerId: callerID,
	}
	if err := mevents.Publish(customer.EventsTopic, ev); err != nil {
		log.Errorf("Error publishing event %+v", ev)
	}

	return nil
}

func (c *Customers) MarkVerified(ctx context.Context, request *customer.MarkVerifiedRequest, response *customer.MarkVerifiedResponse) error {
	email := request.Email
	if email == "" {
		// try deprecated fallback
		email = request.Id
	}

	if strings.TrimSpace(email) == "" {
		return errors.BadRequest("customers.markverified", "Email is required")
	}

	cust, err := readCustomerByEmail(email)
	if err != nil {
		return err
	}
	if err := authorizeCall(ctx, cust.ID); err != nil {
		return err
	}

	cus, err := updateCustomerStatusByEmail(email, statusVerified)
	if err != nil {
		return err
	}

	// Publish the event
	var callerID string
	if acc, ok := auth.AccountFromContext(ctx); ok {
		callerID = acc.ID
	}
	ev := &customer.Event{
		Type:     customer.EventType_EventTypeVerified,
		Customer: objToProto(cus),
		CallerId: callerID,
	}
	if err := mevents.Publish(customer.EventsTopic, ev); err != nil {
		log.Errorf("Error publishing event %+v", ev)
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
	if err := authorizeCall(ctx, cust.ID); err != nil {
		return err
	}
	response.Customer = objToProto(cust)
	return nil
}

func (c *Customers) Delete(ctx context.Context, request *customer.DeleteRequest, response *customer.DeleteResponse) error {
	if strings.TrimSpace(request.Id) == "" && strings.TrimSpace(request.Email) == "" {
		return errors.BadRequest("customers.delete", "ID or Email is required")
	}
	if len(request.Id) == 0 {
		c, err := readCustomerByEmail(request.Email)
		if err != nil {
			return err
		}
		request.Id = c.ID
	}
	if err := authorizeCall(ctx, request.Id); err != nil {
		return err
	}

	if err := c.deleteCustomer(ctx, request.Id, request.Force); err != nil {
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

func authorizeCall(ctx context.Context, customerID string) error {
	account, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("customers", "Unauthorized request")
	}
	if account.Issuer != "micro" {
		return errors.Unauthorized("customers", "Unauthorized request")
	}
	if account.Type == "customer" && account.ID == customerID {
		return nil
	}
	if account.Type == "user" && hasScope("admin", account.Scopes) {
		return nil
	}
	if account.Type == "service" {
		return nil
	}
	return errors.Unauthorized("customers", "Unauthorized request")
}

func hasScope(scope string, scopes []string) bool {
	for _, sc := range scopes {
		if sc == scope {
			return true
		}
	}
	return false
}

func (c *Customers) deleteCustomer(ctx context.Context, customerID string, force bool) error {
	_, err := c.accountsService.Delete(ctx, &aproto.DeleteAccountRequest{
		Id:      customerID,
		Options: &aproto.Options{Namespace: "micro"},
	}, client.WithAuthToken())
	if ignoreDeleteError(err) != nil {
		return err
	}

	var cust *CustomerModel
	// delete customer
	if !force {
		cust, err = updateCustomerStatusByID(customerID, statusDeleted)
		if err != nil {
			return err
		}
	} else {
		// actually delete not just update the status
		cust, err = c.forceDelete(customerID)
		if err != nil {
			return err
		}

	}

	// Publish the event
	var callerID string
	if acc, ok := auth.AccountFromContext(ctx); ok {
		callerID = acc.ID
	}
	ev := &customer.Event{
		Type:     customer.EventType_EventTypeDeleted,
		Customer: objToProto(cust),
		CallerId: callerID,
	}
	if err := mevents.Publish(customer.EventsTopic, ev); err != nil {
		log.Errorf("Error publishing event %+v", ev)
	}

	return nil
}

func (c *Customers) forceDelete(customerID string) (*CustomerModel, error) {
	cust, err := readCustomer(customerID, prefixCustomer)
	if err != nil {
		return nil, err
	}
	if err := mstore.Delete(prefixCustomerEmail + cust.Email); err != nil {
		return nil, err
	}
	if err := mstore.Delete(prefixCustomer + customerID); err != nil {
		return nil, err
	}

	return cust, nil
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

// List is a temporary endpoint which will very quickly become unusable due to the way it lists entries
func (c *Customers) List(ctx context.Context, request *customer.ListRequest, response *customer.ListResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok || acc.Issuer != "micro" || acc.Type != "user" || !hasScope("admin", acc.Scopes) {
		return errors.Unauthorized("customers", "Unauthorized")
	}

	recs, err := mstore.Read(prefixCustomer, mstore.ReadPrefix())
	if err != nil {
		return err
	}
	res := []*customer.Customer{}
	for _, rec := range recs {
		cust := &CustomerModel{}
		if err := json.Unmarshal(rec.Value, cust); err != nil {
			return err
		}
		if cust.Status == statusDeleted {
			// skip
			continue
		}

		res = append(res, &customer.Customer{
			Id:      cust.ID,
			Status:  cust.Status,
			Created: cust.Created,
			Email:   cust.Email,
			Updated: cust.Updated,
		})
	}
	response.Customers = res
	return nil
}

func (c *Customers) Update(ctx context.Context, request *customer.UpdateRequest, response *customer.UpdateResponse) error {
	if err := authorizeCall(ctx, request.Customer.Id); err != nil {
		return err
	}
	cust, err := readCustomerByID(request.Customer.Id)
	if err != nil {
		return err
	}
	changed := false
	if len(request.Customer.Status) > 0 {
		if !validStatus[request.Customer.Status] {
			return errors.BadRequest("customers.update.badstatus", "Invalid status passed")
		}
		if cust.Status != request.Customer.Status {
			cust.Status = request.Customer.Status
			changed = true
		}
	}
	// TODO support email changing - would require reverification
	if !changed {
		return nil
	}

	// Publish the event
	var callerID string
	if acc, ok := auth.AccountFromContext(ctx); ok {
		callerID = acc.ID
	}
	ev := &customer.Event{
		Type:     customer.EventType_EventTypeUpdated,
		Customer: objToProto(cust),
		CallerId: callerID,
	}
	if err := mevents.Publish(customer.EventsTopic, ev); err != nil {
		log.Errorf("Error publishing event %+v", ev)
	}

	return writeCustomer(cust)
}
