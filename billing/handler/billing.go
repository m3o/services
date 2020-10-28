package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	asproto "github.com/m3o/services/alert/proto/alert"
	billing "github.com/m3o/services/billing/proto"
	csproto "github.com/m3o/services/customers/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	sproto "github.com/m3o/services/payments/provider/proto"
	subproto "github.com/m3o/services/subscriptions/proto"
	uproto "github.com/m3o/services/usage/proto"
	"github.com/micro/micro/v3/service/auth"
	goclient "github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	mconfig "github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	merrors "github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	mstore "github.com/micro/micro/v3/service/store"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
)

const (
	// format: `update/2020-09/namespace`
	updatePrefix = "update/"
	// format: `update-by-namespace/namespace/2020-09`
	updateByNamespacePrefix = "update-by-namespace/"
	defaultNamespace        = "micro"
)

type Billing struct {
	stripeClient *client.API // stripe api client
	ns           nsproto.NamespacesService
	ss           sproto.ProviderService
	as           asproto.AlertService
	us           uproto.UsageService
	cs           csproto.CustomersService
	subs         subproto.SubscriptionsService
	config       *Conf
}

type Conf struct {
	additionalUsersPriceID    string
	additionalServicesPriceID string
	planID                    string
	maxIncludedServices       int
	report                    bool
	apiKey                    string
}

func NewBilling(ns nsproto.NamespacesService,
	ss sproto.ProviderService,
	us uproto.UsageService,
	subs subproto.SubscriptionsService,
	cs csproto.CustomersService,
	as asproto.AlertService,
	conf *Conf) *Billing {
	if conf == nil {
		conf = getConfig()
	}

	b := &Billing{
		stripeClient: client.New(conf.apiKey, nil),
		ns:           ns,
		ss:           ss,
		us:           us,
		subs:         subs,
		config:       conf,
		cs:           cs,
		as:           as,
	}
	go b.loop()
	return b
}

func getConfig() *Conf {
	// this is only here for prototyping, should use subscriptions service properly
	// an upside for that will be also the fact that we don't have to load values one by one but can use Scan
	val, err := mconfig.Get("micro.subscriptions.additional_users_price_id")
	if err != nil {
		log.Fatalf("Additional users price id can't be loaded: %v", err)
	}
	additionalUsersPriceID := val.String("")
	if len(additionalUsersPriceID) == 0 {
		log.Fatal("Additional users price id is empty")
	}

	val, err = mconfig.Get("micro.subscriptions.additional_services_price_id")
	if err != nil {
		log.Fatalf("Additional services price id can't be loaded: %v", err)
	}
	additionalServicesPriceID := val.String("")
	if len(additionalServicesPriceID) == 0 {
		log.Fatal("Additional services price id is empty")
	}

	val, err = mconfig.Get("micro.subscriptions.plan_id")
	if err != nil {
		log.Fatalf("Can't load subscription plan id: %v", err)
	}
	planID := val.String("")
	if len(planID) == 0 {
		log.Fatal("Plan id is empty")
	}

	val, err = mconfig.Get("micro.billing.max_included_services")
	if err != nil {
		log.Warnf("Can't load max included services: %v", err)
	}
	maxIncludedServices := val.Int(10)

	val, err = mconfig.Get("micro.billing.report")
	if err != nil {
		log.Warnf("Can't load max included services: %v", err)
	}
	doReporting := val.Bool(false)

	val, err = config.Get("micro.payments.stripe.api_key")
	if err != nil {
		log.Warnf("Can't load stripe api key: %v", err)
	}
	apiKey := val.String("")

	if len(apiKey) == 0 {
		log.Fatalf("Missing required config: micro.payments.stripe.api_key")
	}
	return &Conf{
		apiKey:                    apiKey,
		additionalUsersPriceID:    additionalUsersPriceID,
		additionalServicesPriceID: additionalServicesPriceID,
		planID:                    planID,
		maxIncludedServices:       maxIncludedServices,
		report:                    doReporting,
	}
}

// Updates returns currently active update suggestions for the current month.
// Once the update is applied, it should disappear from this list.
func (b *Billing) Updates(ctx context.Context, req *billing.UpdatesRequest, rsp *billing.UpdatesResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("billing.Updates", "Unauthorized")
	}

	switch {
	case acc.Issuer == defaultNamespace:
	case acc.Issuer != req.Namespace:
		// Instead of throwing an unauthorized, we default back to listing
		// the users namespace
		req.Namespace = acc.Issuer
	}

	// @todo accept a month request parameter
	// for listing historic update suggestions

	key := updatePrefix
	if len(req.Namespace) > 0 {
		key = updateByNamespacePrefix + req.Namespace + "/"
	}
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	log.Infof("Received Billing.Updates request, listing with key '%v', limit '%v'", key, limit)

	records, err := mstore.Read("", mstore.Prefix(key), mstore.Limit(uint(limit)), mstore.Offset(uint(req.Offset)))
	if err != nil && err != mstore.ErrNotFound {
		return merrors.InternalServerError("billing.Updates", "Error listing store: %v", err)
	}

	updates := []*billing.Update{}
	for _, v := range records {
		u := &update{}
		err = json.Unmarshal(v.Value, u)
		if err != nil {
			return merrors.InternalServerError("billing.Updates", "Error unmarshaling value: %v", err)
		}
		updates = append(updates, &billing.Update{
			Namespace:     u.Namespace,
			Created:       u.Created,
			QuantityFrom:  u.QuantityFrom,
			QuantityTo:    u.QuantityTo,
			PlanID:        u.PlanID,
			PriceID:       u.PriceID,
			Note:          u.Note,
			CustomerID:    u.CustomerID,
			CustomerEmail: u.CustomerEmail,
			Id:            u.ID,
		})
	}
	rsp.Updates = updates
	return nil
}

// Apply a change to the account and update the subscriptions accordingly
func (b *Billing) Apply(ctx context.Context, req *billing.ApplyRequest, rsp *billing.ApplyResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("billing.Apply", "Unauthorized")
	}

	switch {
	case acc.Issuer == defaultNamespace:
	default:
		return errors.Unauthorized("billing.Apply", "Unauthorized")
	}

	records, err := mstore.Read(req.Id)
	if err != nil || len(records) == 0 {
		return merrors.InternalServerError("billing.Apply", "Error reading change: %v", err)
	}
	u := &update{}
	err = json.Unmarshal(records[0].Value, u)
	if err != nil {
		return merrors.InternalServerError("billing.Apply", "Error unmarshaling value: %v", err)
	}

	_, err = b.subs.Update(ctx, &subproto.UpdateRequest{
		PriceID:  u.PriceID,
		OwnerID:  u.CustomerID,
		Quantity: u.QuantityTo,
	})
	if err != nil {
		return merrors.InternalServerError("billing.Apply", "Error calling subscriptions update: %v", err)
	}
	return nil
}

// Portal returns the billing portal address the customers can go to to manager their subscriptons
func (b *Billing) Portal(ctx context.Context, req *billing.PortalRequest, rsp *billing.PortalResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("billing.Portal", "Authentication failed")
	}
	email := acc.Name
	if len(email) == 0 {
		email = acc.ID
	}
	params := &stripe.CustomerListParams{
		Email: stripe.String(email),
	}
	params.Filters.AddFilter("limit", "", "3")
	customerIter := b.stripeClient.Customers.List(params)

	customerID := ""
	for customerIter.Next() {
		c := customerIter.Customer()
		customerID = c.ID
		break
	}
	if len(customerID) == 0 {
		return errors.BadRequest("billing.Portal", "No stripe customer found for account %v", acc.ID)
	}

	billParams := &stripe.BillingPortalSessionParams{
		Customer: stripe.String(customerID),
	}
	sess, err := b.stripeClient.BillingPortalSessions.New(billParams)
	if err != nil {
		return errors.InternalServerError("billing.Portal", "Could not create billing portal session: %v", err)
	}
	rsp.Url = sess.URL
	return nil
}

type update struct {
	ID            string
	Namespace     string
	PlanID        string
	PriceID       string
	QuantityFrom  int64
	QuantityTo    int64
	Created       int64
	Note          string
	CustomerID    string
	CustomerEmail string
}

func (b *Billing) calcUpdate(namespace string, persist bool) ([]update, error) {
	rsp, err := b.us.Read(context.TODO(), &uproto.ReadRequest{
		Namespace: namespace,
	}, goclient.WithAuthToken())
	if err != nil {
		return nil, fmt.Errorf("Error getting usage for account service: %v", err)
	}
	if len(rsp.Accounts) == 0 {
		return nil, fmt.Errorf("Account not found for namespace")
	}
	usg := rsp.Accounts[0]

	namespaceRsp, err := b.ns.Read(context.TODO(), &nsproto.ReadRequest{
		Id: namespace,
	}, goclient.WithAuthToken())
	if err != nil {
		return nil, fmt.Errorf("Error listing namespaces: %v", err)
	}
	if len(namespaceRsp.Namespace.Owners) == 0 {
		return nil, fmt.Errorf("No owners for namespace '%v'", namespace)
	}
	if len(namespaceRsp.Namespace.Owners) > 1 {
		return nil, fmt.Errorf("Multiple owners for namespace '%v'", namespace)
	}
	customerID := namespaceRsp.Namespace.Owners[0]
	if len(customerID) == 0 {
		return nil, fmt.Errorf("Owner is empty string for namespace '%v", namespace)
	}

	customerRsp, err := b.cs.Read(context.TODO(), &csproto.ReadRequest{
		Id: customerID,
	}, goclient.WithAuthToken())
	if err != nil {
		return nil, fmt.Errorf("Error reading customer with id '%v': %v", customerID, err)
	}
	customerEmail := customerRsp.GetCustomer().Email

	log.Infof("Processing namespace '%v'", usg.Namespace)

	subsRsp, err := b.ss.ListSubscriptions(context.TODO(), &sproto.ListSubscriptionsRequest{
		CustomerId:   customerID,
		CustomerType: "user",
	}, goclient.WithAuthToken(), goclient.WithRequestTimeout(10*time.Second))
	if err != nil {
		return nil, fmt.Errorf("Error listing subscriptions for customer %v: %v", customerEmail, err)
	}
	if subsRsp == nil {
		return nil, fmt.Errorf("Subscriptions listing response seems empty")
	}
	log.Infof("Found %v subscription for the owner of namespace '%v', customer ID: '%v'", len(subsRsp.Subscriptions), namespace, customerID)

	planIDToSub := map[string]*sproto.Subscription{}
	for _, sub := range subsRsp.Subscriptions {
		planIDToSub[sub.Plan.Id] = sub
	}

	sub, exists := planIDToSub[b.config.additionalUsersPriceID]
	quantity := int64(0)
	if exists {
		quantity = sub.Quantity
	}
	ret := []update{}
	// 1 user is the owner itself
	if quantity != usg.Users-1 {
		log.Infof("Users count needs amending. Saving")

		upd := update{
			ID:            uuid.New().String(),
			PriceID:       b.config.additionalUsersPriceID,
			QuantityFrom:  quantity,
			QuantityTo:    usg.Users - 1,
			Namespace:     usg.Namespace,
			Note:          "Additional users subscription needs changing",
			CustomerID:    customerID,
			CustomerEmail: customerEmail,
		}
		ret = append(ret, upd)
		if persist {
			err = saveUpdate(upd)
			if err != nil {
				return nil, fmt.Errorf("Error saving update: %v", err)
			}
		}
		if b.config.report {
			_, err = b.as.ReportEvent(context.TODO(), &asproto.ReportEventRequest{
				Event: &asproto.Event{
					Category: "billing",
					Action:   "Users Count Change",
					Label:    fmt.Sprintf("User '%v' service subscription value should change from %v to %v", customerEmail, quantity, usg.Users-1),
				},
			})
			if err != nil {
				return nil, fmt.Errorf("Error saving update: %v", err)
			}
		}
	}

	sub, exists = planIDToSub[b.config.additionalServicesPriceID]
	quantity = int64(0)
	if exists {
		quantity = sub.Quantity
	}

	quantityShouldBe := usg.Services - int64(b.config.maxIncludedServices)
	if quantityShouldBe < 0 {
		quantityShouldBe = 0
	}
	if quantity != quantityShouldBe {
		log.Infof("Services count needs amending. Saving")

		upd := update{
			ID:            uuid.New().String(),
			PriceID:       b.config.additionalServicesPriceID,
			QuantityFrom:  quantity,
			QuantityTo:    quantityShouldBe,
			Namespace:     usg.Namespace,
			Note:          "Additional services subscription needs changing",
			CustomerID:    customerID,
			CustomerEmail: customerEmail,
		}
		ret = append(ret, upd)
		if persist {
			err = saveUpdate(upd)
			if err != nil {
				return nil, fmt.Errorf("Error saving update: %v", err)
			}
		}
		if b.config.report {
			_, err = b.as.ReportEvent(context.TODO(), &asproto.ReportEventRequest{
				Event: &asproto.Event{
					Category: "billing",
					Action:   "Service Count Change",
					Label:    fmt.Sprintf("User '%v' service subscription value should change from %v to %v", customerEmail, quantity, quantityShouldBe),
				},
			})
			if err != nil {
				return nil, fmt.Errorf("Error sending report: %v", err)
			}
		}
	}
	return ret, nil
}

func (b *Billing) loop() {
	for {
		func() {
			rsp, err := b.ns.List(context.TODO(), &nsproto.ListRequest{}, goclient.WithAuthToken())
			if err != nil {
				log.Errorf("Error listing namespaces: %v", err)
				return
			}
			for _, namespace := range rsp.Namespaces {
				_, err := b.calcUpdate(namespace.Id, true)
				if err != nil {
					log.Errorf("Error while getting update for namespace '%v': %v", namespace, err)
					if b.config.report {
						_, err = b.as.ReportEvent(context.TODO(), &asproto.ReportEventRequest{
							Event: &asproto.Event{
								Category: "billing",
								Action:   "Processing error",
								Label:    fmt.Sprintf("Error while processing namespace '%v': %v", namespace, err),
							},
						})
					}
					if err != nil {
						log.Error(err)
						continue
					}
				}
			}
		}()

		time.Sleep(1 * time.Hour)
	}
}

func saveUpdate(record update) error {
	tim := time.Now()
	record.Created = tim.Unix()
	val, _ := json.Marshal(record)
	err := mstore.Write(&mstore.Record{
		Key:   fmt.Sprintf("%v%v", updatePrefix, record.Namespace),
		Value: val,
	})
	if err != nil {
		return err
	}
	err = mstore.Write(&mstore.Record{
		Key:   record.ID,
		Value: val,
	})
	if err != nil {
		return err
	}
	return mstore.Write(&mstore.Record{
		Key:   fmt.Sprintf("%v%v", updateByNamespacePrefix, record.Namespace),
		Value: val,
	})
}

type max struct {
	namespace string
	users     int64
	services  int64
}
