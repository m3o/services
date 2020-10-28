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
}

func (u paymentMock) 	ListSubscriptions(ctx context.Context, in *sproto.ListSubscriptionsRequest, opts ...client.CallOption) (*sproto.ListSubscriptionsResponse, error) {
	return nil, nil
}

type namespaceMock struct {
	nsproto.NamespacesService
}

func (n namespaceMock) List(ctx context.Context, in *nsproto.ListRequest, opts ...client.CallOption) (*nsproto.ListResponse, error) {
	return nil, nil
}

type usageMock struct {
	uproto.UsageService
}

func (u usageMock) Read(ctx context.Context, in *uproto.ReadRequest, opts ...client.CallOption) (*uproto.ReadResponse, error) {
	return nil, nil
}


type subscriptionMock struct {
	subproto.SubscriptionsService
}

func (u subscriptionMock) Update(ctx context.Context, in *subproto.UpdateRequest, opts ...client.CallOption) (*subproto.UpdateResponse, error) {
	return nil, nil
}

type customersMock struct {
	csproto.CustomersService
}

func (u customersMock) 	Read(ctx context.Context, in *csproto.ReadRequest, opts ...client.CallOption) (*csproto.ReadResponse, error) {
	return nil, nil
}

type alertMock struct {
	asproto.AlertService
}

func (u alertMock) 	ReportEvent(ctx context.Context, in *asproto.ReportEventRequest, opts ...client.CallOption) (*asproto.ReportEventResponse, error) {
	return nil, nil
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