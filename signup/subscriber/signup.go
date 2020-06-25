package subscriber

import (
	"context"

	signup "github.com/micro/services/signup/proto/signup"
)

type Signup struct{}

func (e *Signup) Handle(ctx context.Context, msg *signup.VerifyRequest) error {
	return nil
}

func Handler(ctx context.Context, msg *signup.VerifyRequest) error {
	return nil
}
