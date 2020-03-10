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

func (t testprovider) Test(ctx context.Context, req *pb.TestRequest, rsp *pb.TestResponse) error {
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
