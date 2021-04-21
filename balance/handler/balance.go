package handler

import (
	"context"
	"crypto/tls"
	"fmt"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	balance "github.com/m3o/services/balance/proto"
	ns "github.com/m3o/services/namespaces/proto"
	publicapi "github.com/m3o/services/publicapi/proto"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
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

type publicAPICacheEntry struct {
	api     *publicapi.PublicAPI
	created time.Time
}

type publicAPICache struct {
	sync.RWMutex
	cache  map[string]*publicAPICacheEntry
	pubSvc publicapi.PublicapiService
	ttl    time.Duration
}

func (p *publicAPICache) get(name string) (*publicapi.PublicAPI, error) {
	// check the cache
	// TODO mutex
	p.RLock()
	cached := p.cache[name]
	p.RUnlock()
	if cached != nil && cached.created.After(time.Now().Add(p.ttl)) {
		return cached.api, nil
	}
	rsp, err := p.pubSvc.Get(context.Background(), &publicapi.GetRequest{Name: name}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}
	p.Lock()
	p.cache[name] = &publicAPICacheEntry{api: rsp.Api, created: time.Now()}
	p.Unlock()
	return rsp.Api, nil
}

type Balance struct {
	c      *counter // counts the balance. Balance is expressed in 1/100ths of a cent which allows us to price in fractions e.g. a request costs 0.01 cents or 100 requests for 1 cent
	v1Svc  v1api.V1Service
	pubSvc *publicAPICache
	nsSvc  ns.NamespacesService
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
		c:     &counter{redisClient: rc},
		v1Svc: v1api.NewV1Service("v1", svc.Client()),
		pubSvc: &publicAPICache{
			pubSvc: publicapi.NewPublicapiService("publicapi", svc.Client()),
			cache:  map[string]*publicAPICacheEntry{},
			ttl:    5 * time.Minute,
		},
		nsSvc: ns.NewNamespacesService("namespaces", svc.Client()),
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
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("balance.Current", "Unauthorized")
	}
	if acc.Issuer != "micro" {
		// reject
		return errors.Forbidden("balance.Current", "Forbidden")
	}
	if len(request.CustomerId) == 0 {
		request.CustomerId = acc.ID
	}
	if acc.ID != request.CustomerId {
		// is this an admin?
		if err := verifyAdmin(ctx, "balance.Current"); err != nil {
			return err
		}
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
