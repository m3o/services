package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"

	"github.com/go-redis/redis/v8"
	balance "github.com/m3o/services/balance/proto"
	publicapi "github.com/m3o/services/publicapi/proto"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
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

func NewHandler(svc *service.Service) *Balance {
	redisConfig := struct {
		Address  string
		User     string
		Password string
	}{}
	val, err := config.Get("micro.balance.redis")
	if err != nil {
		log.Fatalf("No redis config found %s", err)
	}
	if err := val.Scan(&redisConfig); err != nil {
		log.Fatalf("Error parsing redis config %s", err)
	}
	if len(redisConfig.Password) == 0 || len(redisConfig.User) == 0 || len(redisConfig.Password) == 0 {
		log.Fatalf("Missing redis config %s", err)
	}
	rc := redis.NewClient(&redis.Options{
		Addr:     redisConfig.Address,
		Username: redisConfig.User,
		Password: redisConfig.Password,
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false,
		},
	})
	b := &Balance{
		c:      &counter{redisClient: rc},
		v1Svc:  v1api.NewV1Service("v1", svc.Client()),
		pubSvc: publicapi.NewPublicapiService("publicapi", svc.Client()),
	}
	go b.consumeEvents()
	return b
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
	if err := verifyAdmin(ctx, "balance.Current"); err != nil {
		return err
	}
	currBal, err := b.c.read(request.CustomerId, "$balance$")
	if err != nil {
		log.Errorf("Error reading from counter %s", err)
		return errors.InternalServerError("balance.Current", "Error retrieving current balance")
	}
	response.CurrentBalance = currBal
	return nil
}

func verifyAdmin(ctx context.Context, method string) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden(method, "Forbidden")
	}
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			return nil
		}
	}
	return errors.Forbidden(method, "Forbidden")
}
