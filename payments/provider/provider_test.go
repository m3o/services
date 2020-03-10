package provider

import (
	"context"
	"testing"
	"time"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/client"

	pb "github.com/micro/services/payments/provider/proto"
)

type testprovider struct{}

func (t testprovider) CreateProduct(ctx context.Context, req *pb.CreateProductRequest, rsp *pb.CreateProductResponse) error {
	return nil
}
func (t testprovider) CreatePlan(ctx context.Context, req *pb.CreatePlanRequest, rsp *pb.CreatePlanResponse) error {
	return nil
}
func (t testprovider) CreateCustomer(ctx context.Context, req *pb.CreateCustomerRequest, rsp *pb.CreateCustomerResponse) error {
	return nil
}
func (t testprovider) CreateSubscription(ctx context.Context, req *pb.CreateSubscriptionRequest, rsp *pb.CreateSubscriptionResponse) error {
	return nil
}

func TestNewProvider(t *testing.T) {
	// test the provider returns ErrNotFound when not registered
	t.Run("no provider set", func(t *testing.T) {
		_, err := NewProvider("test", client.NewClient())
		if err != ErrNotFound {
			t.Errorf("Expected ErrNotFound, got %v", err)
		}
	})

	// test the provider returns a provider when one is registered
	t.Run("provider set", func(t *testing.T) {
		testSrv := micro.NewService(micro.Name(ServicePrefix + "test"))
		if err := pb.RegisterProviderHandler(testSrv.Server(), new(testprovider)); err != nil {
			t.Fatalf("Error registering test handler: %v", err)
		}
		go testSrv.Run()

		// TODO: Find way of improving this test so the delay is not hardcoded
		// and the testSrv is stopped at the end of the function
		time.Sleep(200 * time.Millisecond)

		_, err := NewProvider("test", client.NewClient())
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
	})
}
