package handler

import (
	"context"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	balance "github.com/m3o/services/balance/proto"
	publicapi "github.com/m3o/services/publicapi/proto"
	v1api "github.com/m3o/services/v1api/proto"
)

const (
	prefixCounter = "balance-service/counter"
)

type counter struct {
	sync.RWMutex
	redisClient *redis.Client
}

func (c *counter) incr(userID, path string, delta int64) (int64, error) {
	return c.redisClient.IncrBy(context.Background(), fmt.Sprintf("%s:%s:%s", prefixCounter, userID, path), delta).Result()
}

func (c *counter) decr(userID, path string, delta int64) (int64, error) {
	return c.redisClient.DecrBy(context.Background(), fmt.Sprintf("%s:%s:%s", prefixCounter, userID, path), delta).Result()
}

func (c *counter) read(userID, path string) (int64, error) {
	return c.redisClient.Get(context.Background(), fmt.Sprintf("%s:%s:%s", prefixCounter, userID, path)).Int64()
}

func (c *counter) reset(userID, path string) error {
	return c.redisClient.Set(context.Background(), fmt.Sprintf("%s:%s:%s", prefixCounter, userID, path), 0, 0).Err()
}

type Balance struct {
	c      *counter
	v1Svc  v1api.V1Service
	pubSvc publicapi.PublicapiService
}

func (b Balance) Increment(ctx context.Context, request *balance.IncrementRequest, response *balance.IncrementResponse) error {
	// check idempotency key
	// increment counter
	// TODO do we need to store each individual transaction
	panic("implement me")
}

func (b Balance) Decrement(ctx context.Context, request *balance.DecrementRequest, response *balance.DecrementResponse) error {
	panic("implement me")
}

func (b Balance) Current(ctx context.Context, request *balance.CurrentRequest, response *balance.CurrentResponse) error {
	panic("implement me")
}
