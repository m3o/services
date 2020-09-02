package handler

import (
	"context"
	"time"

	billing "github.com/m3o/services/billing/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	sproto "github.com/m3o/services/payments/provider/proto"
	uproto "github.com/m3o/services/usage/proto"
	"github.com/micro/go-micro/v3/auth"
	goclient "github.com/micro/go-micro/v3/client"
	"github.com/micro/go-micro/v3/errors"
	"github.com/micro/micro/v3/service/config"
	mconfig "github.com/micro/micro/v3/service/config"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/stripe/stripe-go/v71"
	"github.com/stripe/stripe-go/v71/client"
)

type Billing struct {
	stripeClient           *client.API // stripe api client
	ns                     nsproto.NamespacesService
	ss                     sproto.ProviderService
	us                     uproto.UsageService
	additionalUsersPriceID string
	planID                 string
}

func NewBilling(ns nsproto.NamespacesService, ss sproto.ProviderService, us uproto.UsageService) *Billing {
	// this is only here for prototyping, should use subscriptions service properly
	additionalUsersPriceID := mconfig.Get("micro", "subscriptions", "additional_users_price_id").String("")
	planID := mconfig.Get("micro", "subscriptions", "plan_id").String("")

	apiKey := config.Get("micro", "payments", "stripe", "api_key").String("")
	if len(apiKey) == 0 {
		log.Fatalf("Missing required config: micro.payments.stripe.api_key")
	}
	b := &Billing{
		stripeClient:           client.New(apiKey, nil),
		ns:                     ns,
		ss:                     ss,
		us:                     us,
		additionalUsersPriceID: additionalUsersPriceID,
		planID:                 planID,
	}
	go b.loop()
	return b
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
				customer := namespaceToOwner[max.namespace]
				subsRsp, err := b.ss.ListSubscriptions(context.TODO(), &sproto.ListSubscriptionsRequest{
					CustomerId: customer,
				}, goclient.WithAuthToken(), goclient.WithRequestTimeout(10*time.Second))
				if err != nil {
					log.Warnf("Error listing subscriptions for customer %v: %v", customer, err)
					continue
				}
				if subsRsp == nil {
					continue
				}
				for _, sub := range subsRsp.Subscriptions {
					if sub.Id == b.additionalUsersPriceID {
						if sub.Quantity != max.users {
							log.Infof("Users count needs amending")
						}
					}
				}
			}
		}()

		time.Sleep(1 * time.Hour)
	}
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
		index[account.Namespace] = m
	}
	return index
}
