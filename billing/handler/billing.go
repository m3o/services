package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	billing "github.com/m3o/services/billing/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	sproto "github.com/m3o/services/payments/provider/proto"
	uproto "github.com/m3o/services/usage/proto"
	"github.com/micro/go-micro/v3/auth"
	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	merrors "github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	"github.com/micro/micro/v3/service/config"
	mconfig "github.com/micro/micro/v3/service/config"
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
	monthFormat             = "2006-01"
	defaultNamespace        = "micro"
)

type Billing struct {
	stripeClient              *client.API // stripe api client
	ns                        nsproto.NamespacesService
	ss                        sproto.ProviderService
	us                        uproto.UsageService
	additionalUsersPriceID    string
	additionalServicesPriceID string
	planID                    string
	maxIncludedServices       int
}

func NewBilling(ns nsproto.NamespacesService, ss sproto.ProviderService, us uproto.UsageService) *Billing {
	// this is only here for prototyping, should use subscriptions service properly
	additionalUsersPriceID := mconfig.Get("micro", "subscriptions", "additional_users_price_id").String("")
	additionalServicesPriceID := mconfig.Get("micro", "subscriptions", "additional_services_price_id").String("")
	planID := mconfig.Get("micro", "subscriptions", "plan_id").String("")
	maxIncludedServices := mconfig.Get("micro", "billing", "max_included_services").Int(10)

	apiKey := config.Get("micro", "payments", "stripe", "api_key").String("")
	if len(apiKey) == 0 {
		log.Fatalf("Missing required config: micro.payments.stripe.api_key")
	}
	b := &Billing{
		stripeClient:              client.New(apiKey, nil),
		ns:                        ns,
		ss:                        ss,
		us:                        us,
		additionalUsersPriceID:    additionalUsersPriceID,
		additionalServicesPriceID: additionalServicesPriceID,
		planID:                    planID,
		maxIncludedServices:       maxIncludedServices,
	}
	go b.loop()
	return b
}

// List account history by namespace, or lists latest values for each namespace if history is not provided.
func (b *Billing) Updates(ctx context.Context, req *billing.UpdatesRequest, rsp *billing.UpdatesResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("billing.ListUpdates", "Unauthorized")
	}

	switch {
	case acc.Issuer == defaultNamespace:
	case acc.Issuer != req.Namespace:
		// Instead of throwing an unauthorized, we default back to listing
		// the users namespace
		req.Namespace = acc.Issuer
	}

	key := updatePrefix
	if len(req.Namespace) > 0 {
		key = updateByNamespacePrefix + req.Namespace + "/"
	}
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	log.Infof("Received Billing.ListUpdates request, listing with key '%v', limit '%v'", key, limit)

	records, err := mstore.Read(key, store.ReadPrefix(), store.ReadLimit(uint(limit)), store.ReadOffset(uint(req.Offset)))
	if err != nil && err != store.ErrNotFound {
		return merrors.InternalServerError("billing.ListUpdates", "Error listing store: %v", err)
	}

	updates := []*billing.Update{}
	for _, v := range records {
		u := &update{}
		err = json.Unmarshal(v.Value, u)
		if err != nil {
			return merrors.InternalServerError("billing.ListUpdates", "Error unmarsjaling value: %v", err)
		}
		updates = append(updates, &billing.Update{
			Namespace:    u.Namespace,
			Created:      u.Created,
			QuantityFrom: u.QuantityFrom,
			QuantityTo:   u.QuantityTo,
			PlanID:       u.PlanID,
			PriceID:      u.PriceID,
			Note:         u.Note,
			Customer:     u.Customer,
			Id:           u.ID,
		})
	}
	rsp.Updates = updates
	return nil
}

// Portal returns the billing portal address the customers can go to to manager their subscriptons
func (b *Billing) Portal(ctx context.Context, req *billing.PortalRequest, rsp *billing.PortalResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.BadRequest("billing.Portal", "Authentication failed")
	}
	params := &stripe.CustomerListParams{
		Email: &acc.ID,
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
	ID           string
	Namespace    string
	PlanID       string
	PriceID      string
	QuantityFrom int64
	QuantityTo   int64
	Created      int64
	Note         string
	Customer     string
}

func (b *Billing) loop() {
	for {
		func() {
			allAccounts := []*uproto.Account{}
			offset := int64(0)
			page := int64(500)
			for {
				log.Infof("Listing usage with offset %v", offset)

				rsp, err := b.us.List(context.TODO(), &uproto.ListRequest{
					Distinct: true,
					Offset:   offset,
					Limit:    page,
				}, goclient.WithAuthToken())
				if err != nil {
					log.Errorf("Error calling namespace service: %v", err)
					return
				}
				allAccounts = append(allAccounts, rsp.Accounts...)
				if len(rsp.Accounts) < int(page) {
					break
				}
				offset += page
			}

			log.Infof("Processing %v number of distinct account values to get max", len(allAccounts))
			maxs := getMax(allAccounts)

			log.Infof("Got %v namespaces to check subscriptions for", len(maxs))

			rsp, err := b.ns.List(context.TODO(), &nsproto.ListRequest{}, goclient.WithAuthToken())
			if err != nil {
				log.Warnf("Error listing namespaces: %v", err)
				return
			}
			namespaceToOwner := map[string]string{}
			for _, namespace := range rsp.Namespaces {
				if len(namespace.Owners) == 0 {
					log.Warnf("Namespace %v has no owner", namespace.Id)
					continue
				}
				namespaceToOwner[namespace.Id] = namespace.Owners[0]
			}

			for _, max := range maxs {
				log.Infof("Processing namespace '%v'", max.namespace)
				customer, found := namespaceToOwner[max.namespace]
				if !found || len(customer) == 0 {
					log.Warnf("Owner customer id not found for namespace '%v'", max.namespace)
					continue
				}
				subsRsp, err := b.ss.ListSubscriptions(context.TODO(), &sproto.ListSubscriptionsRequest{
					CustomerId:   customer,
					CustomerType: "user",
				}, goclient.WithAuthToken(), goclient.WithRequestTimeout(10*time.Second))
				if err != nil {
					log.Warnf("Error listing subscriptions for customer %v: %v", customer, err)
					continue
				}
				if subsRsp == nil {
					log.Warnf("Subscriptions listing response seems empty")
					continue
				}
				log.Infof("Found %v subscription for the owner of namespace '%v'", len(subsRsp.Subscriptions), max.namespace)

				planIDToSub := map[string]*sproto.Subscription{}
				for _, sub := range subsRsp.Subscriptions {
					planIDToSub[sub.Plan.Id] = sub
				}

				sub, exists := planIDToSub[b.additionalUsersPriceID]
				quantity := int64(0)
				if exists {
					quantity = sub.Quantity
				}
				if quantity != max.users {
					log.Infof("Users count needs amending. Saving")

					err = saveUpdate(update{
						ID:           uuid.New().String(),
						PriceID:      b.additionalUsersPriceID,
						QuantityFrom: quantity,
						QuantityTo:   max.users,
						Namespace:    max.namespace,
						Note:         "Additional users subscription needs changing",
						Customer:     customer,
					})
					if err != nil {
						log.Warnf("Error saving update: %v", err)
					}
				}

				sub, exists = planIDToSub[b.additionalServicesPriceID]
				quantity = int64(0)
				if exists {
					quantity = sub.Quantity
				}

				quantityShouldBe := max.services - int64(b.maxIncludedServices)
				if quantityShouldBe < 0 {
					quantityShouldBe = 0
				}
				if quantity != quantityShouldBe {
					err = saveUpdate(update{
						ID:           uuid.New().String(),
						PriceID:      b.additionalServicesPriceID,
						QuantityFrom: quantity,
						QuantityTo:   quantityShouldBe,
						Namespace:    max.namespace,
						Note:         "Additional services subscription needs changing",
						Customer:     customer,
					})
					if err != nil {
						log.Warnf("Error saving update: %v", err)
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
	month := tim.Format(monthFormat)
	err := mstore.Write(&store.Record{
		Key:   fmt.Sprintf("%v%v/%v", updatePrefix, month, record.Namespace),
		Value: val,
	})
	if err != nil {
		return err
	}
	err = mstore.Write(&store.Record{
		Key:   record.ID,
		Value: val,
	})
	if err != nil {
		return err
	}
	return mstore.Write(&store.Record{
		Key:   fmt.Sprintf("%v%v/%v", updateByNamespacePrefix, record.Namespace, month),
		Value: val,
	})
}

type max struct {
	namespace string
	users     int64
	services  int64
}

func getMax(accounts []*uproto.Account) map[string]*max {
	index := map[string]*max{}
	for _, account := range accounts {
		m, ok := index[account.Namespace]
		if !ok {
			m = &max{}
		}
		if account.Users > m.users {
			m.users = account.Users
		}
		if account.Services > m.services {
			m.services = account.Services
		}
		m.namespace = account.Namespace
		index[account.Namespace] = m
	}
	return index
}
