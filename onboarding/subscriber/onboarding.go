package subscriber

import (
	"context"

	onboarding "github.com/micro/services/onboarding/proto/onboarding"
)

type Onboarding struct{}

func (e *Onboarding) Handle(ctx context.Context, msg *onboarding.VerifyRequest) error {
	return nil
}

func Handler(ctx context.Context, msg *onboarding.VerifyRequest) error {
	return nil
}
