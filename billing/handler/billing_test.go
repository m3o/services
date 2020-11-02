package handler

import (
	"context"
	"errors"
	"testing"

	mstore "github.com/micro/micro/v3/service/store"
	"github.com/micro/micro/v3/service/store/file"

	client "github.com/micro/micro/v3/service/client"

	asproto "github.com/m3o/services/alert/proto/alert"
	csproto "github.com/m3o/services/customers/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	sproto "github.com/m3o/services/payments/proto"
	subproto "github.com/m3o/services/subscriptions/proto"
	uproto "github.com/m3o/services/usage/proto"
)

func setupBillingTests() {

}

type paymentMock struct {
	sproto.ProviderService
	ListSubscriptionsFunc func(ctx context.Context, in *sproto.ListSubscriptionsRequest, opts ...client.CallOption) (*sproto.ListSubscriptionsResponse, error)
}

func (u paymentMock) ListSubscriptions(ctx context.Context, in *sproto.ListSubscriptionsRequest, opts ...client.CallOption) (*sproto.ListSubscriptionsResponse, error) {
	return u.ListSubscriptionsFunc(ctx, in, opts...)
}

type namespaceMock struct {
	nsproto.NamespacesService
	ReadFunc func(ctx context.Context, in *nsproto.ReadRequest, opts ...client.CallOption) (*nsproto.ReadResponse, error)
	ListFunc func(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error)
}

func (n namespaceMock) List(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error) {
	return n.ListFunc(ctx, in, opts...)
}

func (n namespaceMock) Read(ctx context.Context, in *nsproto.ReadRequest, opts ...client.CallOption) (*nsproto.ReadResponse, error) {
	return n.ReadFunc(ctx, in, opts...)
}

type usageMock struct {
	uproto.UsageService
	ReadFunc func(ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error)
}

func (u usageMock) Read(ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error) {
	return u.ReadFunc(ctx, in, opts...)
}

type subscriptionMock struct {
	subproto.SubscriptionsService
	UpdateFunc func(ctx context.Context, in *subproto.UpdateRequest, opts ...client.CallOption) (*subproto.UpdateResponse, error)
}

func (u subscriptionMock) Update(ctx context.Context, in *subproto.UpdateRequest, opts ...client.CallOption) (*subproto.UpdateResponse, error) {
	return u.UpdateFunc(ctx, in, opts...)
}

type customersMock struct {
	csproto.CustomersService
	ReadFunc func(ctx context.Context, in *csproto.ReadRequest, opts ...client.CallOption) (*csproto.ReadResponse, error)
}

func (u customersMock) Read(ctx context.Context, in *csproto.ReadRequest, opts ...client.CallOption) (*csproto.ReadResponse, error) {
	return u.ReadFunc(ctx, in, opts...)
}

type alertMock struct {
	asproto.AlertService
	ReportEventFunc func(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error)
}

func (u alertMock) ReportEvent(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error) {
	return u.ReportEventFunc(ctx, in, opts...)
}

func TestNoSubscription(t *testing.T) {
	bs := NewBilling(&namespaceMock{
		ListFunc: func(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error) {
			return &nsproto.ListResponse{
				Namespaces: []*nsproto.Namespace{
					{
						Id: "ns1",
					},
				},
			}, nil
		},
		ReadFunc: func(ctx context.Context, in *nsproto.ReadRequest, opts ...client.CallOption) (*nsproto.ReadResponse, error) {
			return &nsproto.ReadResponse{
				Namespace: &nsproto.Namespace{
					Id:     "ns1",
					Owners: []string{"someid"},
				},
			}, nil
		},
	}, &paymentMock{
		ListSubscriptionsFunc: func(ctx context.Context, in *sproto.ListSubscriptionsRequest, opts ...client.CallOption) (*sproto.ListSubscriptionsResponse, error) {
			return &sproto.ListSubscriptionsResponse{
				Subscriptions: []*sproto.Subscription{},
			}, nil
		},
	}, &usageMock{
		ReadFunc: func(ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error) {
			if in.Namespace != "ns1" {
				return nil, errors.New("Namespace should be ns1")
			}
			return &uproto.ReadResponse{
				Accounts: []*uproto.Account{
					{
						Namespace: "ns1",
						Services:  4,
						Users:     2,
					},
				},
			}, nil
		},
	}, &subscriptionMock{
		UpdateFunc: func(ctx context.Context, in *subproto.UpdateRequest, opts ...client.CallOption) (*subproto.UpdateResponse, error) {
			switch in.PriceID {
			case "usersprice":
				if in.Quantity != 1 {
					return nil, errors.New("Should update users to 1")
				}
			case "servicesprice":
				if in.Quantity != 2 {
					return nil, errors.New("Should upate services to 2")
				}
			}
			return nil, nil
		},
	}, &customersMock{
		ReadFunc: func(ctx context.Context, in *csproto.ReadRequest, opts ...client.CallOption) (*csproto.ReadResponse, error) {
			if in.Id != "someid" {
				return nil, errors.New("Can't find")
			}
			return &csproto.ReadResponse{
				Customer: &csproto.Customer{
					Email: "email@address.com",
				},
			}, nil
		},
	}, &alertMock{
		ReportEventFunc: func(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error) {
			return &asproto.ReportEventResponse{}, nil
		},
	}, &Conf{
		additionalServicesPriceID: "servicesprice",
		additionalUsersPriceID:    "usersprice",
		planID:                    "planid",
		maxIncludedServices:       2,
		report:                    false,
		apiKey:                    "none",
	})
	updates, err := bs.calcUpdate("ns1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(updates) < 2 {
		t.Fatal(updates)
	}
	if updates[0].CustomerID != "someid" ||
		updates[0].CustomerEmail != "email@address.com" ||
		updates[0].QuantityFrom != 0 || updates[0].QuantityTo != 1 ||
		updates[0].PriceID != "usersprice" {
		t.Fatal(updates[0])
	}

	if updates[1].CustomerID != "someid" ||
		updates[1].CustomerEmail != "email@address.com" ||
		updates[1].QuantityFrom != 0 || updates[1].QuantityTo != 2 ||
		updates[1].PriceID != "servicesprice" {
		t.Fatal(updates[0])
	}
}

func TestSubscriptionDecrease(t *testing.T) {
	mstore.DefaultStore = file.NewStore()
	bs := NewBilling(&namespaceMock{
		ListFunc: func(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error) {
			return &nsproto.ListResponse{
				Namespaces: []*nsproto.Namespace{
					{
						Id: "ns1",
					},
				},
			}, nil
		},
		ReadFunc: func(ctx context.Context, in *nsproto.ReadRequest, opts ...client.CallOption) (*nsproto.ReadResponse, error) {
			return &nsproto.ReadResponse{
				Namespace: &nsproto.Namespace{
					Id:     "ns1",
					Owners: []string{"someid"},
				},
			}, nil
		},
	}, &paymentMock{
		ListSubscriptionsFunc: func(ctx context.Context, in *sproto.ListSubscriptionsRequest, opts ...client.CallOption) (*sproto.ListSubscriptionsResponse, error) {
			return &sproto.ListSubscriptionsResponse{
				Subscriptions: []*sproto.Subscription{
					{
						Plan: &sproto.Plan{
							Id: "servicesprice",
						},
						Quantity: 7,
					},
					{
						Plan: &sproto.Plan{
							Id: "usersprice",
						},
						Quantity: 5,
					},
				},
			}, nil
		},
	}, &usageMock{
		ReadFunc: func(ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error) {
			if in.Namespace != "ns1" {
				return nil, errors.New("Namespace should be ns1")
			}
			return &uproto.ReadResponse{
				Accounts: []*uproto.Account{
					{
						Namespace: "ns1",
						Services:  5,
						Users:     3,
					},
				},
			}, nil
		},
	}, &subscriptionMock{
		UpdateFunc: func(ctx context.Context, in *subproto.UpdateRequest, opts ...client.CallOption) (*subproto.UpdateResponse, error) {
			switch in.PriceID {
			case "usersprice":
				if in.Quantity != 2 {
					return nil, errors.New("Should update users to 1")
				}
			case "servicesprice":
				if in.Quantity != 3 {
					return nil, errors.New("Should upate services to 2")
				}
			}
			return nil, nil
		},
	}, &customersMock{
		ReadFunc: func(ctx context.Context, in *csproto.ReadRequest, opts ...client.CallOption) (*csproto.ReadResponse, error) {
			if in.Id != "someid" {
				return nil, errors.New("Can't find")
			}
			return &csproto.ReadResponse{
				Customer: &csproto.Customer{
					Email: "email@address.com",
				},
			}, nil
		},
	}, &alertMock{
		ReportEventFunc: func(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error) {
			return &asproto.ReportEventResponse{}, nil
		},
	}, &Conf{
		additionalServicesPriceID: "servicesprice",
		additionalUsersPriceID:    "usersprice",
		planID:                    "planid",
		maxIncludedServices:       2,
		report:                    false,
		apiKey:                    "none",
	})
	updates, err := bs.calcUpdate("ns1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(updates) < 2 {
		t.Fatal(updates)
	}
	if updates[0].CustomerID != "someid" ||
		updates[0].CustomerEmail != "email@address.com" ||
		updates[0].QuantityFrom != 5 || updates[0].QuantityTo != 2 ||
		updates[0].PriceID != "usersprice" {
		t.Fatal(updates[0])
	}

	if updates[1].CustomerID != "someid" ||
		updates[1].CustomerEmail != "email@address.com" ||
		updates[1].QuantityFrom != 7 || updates[1].QuantityTo != 3 ||
		updates[1].PriceID != "servicesprice" {
		t.Fatal(updates[0])
	}
}

func TestNoChange(t *testing.T) {
	bs := NewBilling(&namespaceMock{
		ListFunc: func(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error) {
			return &nsproto.ListResponse{
				Namespaces: []*nsproto.Namespace{
					{
						Id: "ns1",
					},
				},
			}, nil
		},
		ReadFunc: func(ctx context.Context, in *nsproto.ReadRequest, opts ...client.CallOption) (*nsproto.ReadResponse, error) {
			return &nsproto.ReadResponse{
				Namespace: &nsproto.Namespace{
					Id:     "ns1",
					Owners: []string{"someid"},
				},
			}, nil
		},
	}, &paymentMock{
		ListSubscriptionsFunc: func(ctx context.Context, in *sproto.ListSubscriptionsRequest, opts ...client.CallOption) (*sproto.ListSubscriptionsResponse, error) {
			return &sproto.ListSubscriptionsResponse{
				Subscriptions: []*sproto.Subscription{
					{
						Plan: &sproto.Plan{
							Id: "servicesprice",
						},
						Quantity: 7,
					},
					{
						Plan: &sproto.Plan{
							Id: "usersprice",
						},
						Quantity: 5,
					},
				},
			}, nil
		},
	}, &usageMock{
		ReadFunc: func(ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error) {
			if in.Namespace != "ns1" {
				return nil, errors.New("Namespace should be ns1")
			}
			return &uproto.ReadResponse{
				Accounts: []*uproto.Account{
					{
						Namespace: "ns1",
						Services:  9,
						Users:     6,
					},
				},
			}, nil
		},
	}, &subscriptionMock{
		UpdateFunc: func(ctx context.Context, in *subproto.UpdateRequest, opts ...client.CallOption) (*subproto.UpdateResponse, error) {
			return nil, errors.New("This should not be called")
		},
	}, &customersMock{
		ReadFunc: func(ctx context.Context, in *csproto.ReadRequest, opts ...client.CallOption) (*csproto.ReadResponse, error) {
			if in.Id != "someid" {
				return nil, errors.New("Can't find")
			}
			return &csproto.ReadResponse{
				Customer: &csproto.Customer{
					Email: "email@address.com",
				},
			}, nil
		},
	}, &alertMock{
		ReportEventFunc: func(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error) {
			return &asproto.ReportEventResponse{}, nil
		},
	}, &Conf{
		additionalServicesPriceID: "servicesprice",
		additionalUsersPriceID:    "usersprice",
		planID:                    "planid",
		maxIncludedServices:       2,
		report:                    false,
		apiKey:                    "none",
	})
	updates, err := bs.calcUpdate("ns1", false)
	if err != nil {
		t.Fatal(err)
	}
	if len(updates) != 0 {
		t.Fatal(updates)
	}
}
