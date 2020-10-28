package handler

import(
	"context"
	"testing"

	client "github.com/micro/micro/v3/service/client"

	asproto "github.com/m3o/services/alert/proto/alert"
	csproto "github.com/m3o/services/customers/proto"
	nsproto "github.com/m3o/services/namespaces/proto"
	sproto "github.com/m3o/services/payments/provider/proto"
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
	ListFunc func(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error)
}

func (n namespaceMock) List(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error) {
	return n.ListFunc(ctx, in, opts...)
}

type usageMock struct {
	uproto.UsageService
	ReadFunc func (ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error)
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

func (u alertMock) 	ReportEvent(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error) {
	return u.ReportEventFunc(ctx, in, opts...)
}

func TestBillingCalc(t *testing.T) {
	bs := NewBilling(&namespaceMock{

	}, &paymentMock{

	}, &usageMock{

	}, &subscriptionMock{

	}, &customersMock{

	}, &alertMock{

	}, &Conf{

	})
	t.Fatal(bs)
}