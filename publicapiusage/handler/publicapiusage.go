package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/go-redis/redis/v8"
	pb "github.com/m3o/services/publicapiusage/proto"
	"github.com/micro/micro/v3/service"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"
)

const (
	prefixCounter         = "publicapiusage-service/counter"
	prefixUsageByCustomer = "usageByCustomer/%s/%s" // customer ID / date
	counterTTL            = 48 * time.Hour
)

type counter struct {
	sync.RWMutex
	redisClient *redis.Client
}

func (c *counter) incr(userID, path string, delta int64, t time.Time) (int64, error) {
	t = t.UTC()
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, t.Format("20060102"), path)
	pipe := c.redisClient.TxPipeline()
	incr := pipe.IncrBy(ctx, key, delta)
	pipe.Expire(ctx, key, counterTTL) // make sure we expire the counters
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return incr.Result()
}

func (c *counter) decr(userID, path string, delta int64, t time.Time) (int64, error) {
	t = t.UTC()
	ctx := context.Background()
	key := fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, t.Format("20060102"), path)
	pipe := c.redisClient.TxPipeline()
	decr := pipe.DecrBy(ctx, key, delta)
	pipe.Expire(ctx, key, counterTTL) // make sure we expire counters
	_, err := pipe.Exec(ctx)
	if err != nil {
		return 0, err
	}
	return decr.Result()
}

func (c *counter) read(userID, path string, t time.Time) (int64, error) {
	t = t.UTC()
	return c.redisClient.Get(context.Background(), fmt.Sprintf("%s:%s:%s:%s", prefixCounter, userID, t.Format("20060102"), path)).Int64()
}

type listEntry struct {
	Service string
	Count   int64
}

func (c *counter) listForUser(userID string, t time.Time) ([]listEntry, error) {
	ctx := context.Background()
	keyPrefix := fmt.Sprintf("%s:%s:%s:", prefixCounter, userID, t.Format("20060102"))
	sc := c.redisClient.Scan(ctx, 0, keyPrefix+"*", 0)
	if err := sc.Err(); err != nil {
		return nil, err
	}
	iter := sc.Iterator()
	res := []listEntry{}
	for {
		if !iter.Next(ctx) {
			break
		}
		key := iter.Val()
		i, err := c.redisClient.Get(ctx, key).Int64()
		if err != nil {
			return nil, err
		}
		res = append(res, listEntry{
			Service: strings.TrimPrefix(key, keyPrefix),
			Count:   i,
		})
	}
	return res, iter.Err()
}

type Publicapiusage struct {
	c *counter
}

func NewHandler(svc *service.Service) *Publicapiusage {
	redisConfig := struct {
		Address  string
		User     string
		Password string
	}{}
	val, err := config.Get("micro.publicapiusage.redis")
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
	p := &Publicapiusage{
		c: &counter{redisClient: rc},
	}
	go p.consumeEvents()
	return p
}

func (p Publicapiusage) Read(ctx context.Context, request *pb.ReadRequest, response *pb.ReadResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("publicapiusage.Read", "Unauthorized")
	}
	if acc.ID != request.CustomerId {
		err := verifyMicroAdmin(ctx, "publicapiusage.Read")
		if err != nil {
			return err
		}
	}

	now := time.Now().UTC().Truncate(24 * time.Hour)
	entries, err := p.c.listForUser(request.CustomerId, now)
	if err != nil {
		log.Errorf("Error retrieving usage %s", err)
		return errors.InternalServerError("publicapiusage.Read", "Error retrieving usage")
	}

	response.Usage = make([]*pb.Usage, len(entries))
	for i, v := range entries {
		response.Usage[i] = &pb.Usage{
			ApiName: v.Service,
			Records: []*pb.UsageRecord{
				{
					Date:     now.Unix(),
					Requests: v.Count,
				},
			},
		}
	}
	// TODO add historical data

	return nil
}

func verifyMicroAdmin(ctx context.Context, method string) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized(method, "Unauthorized")
	}
	if acc.Issuer != "micro" {
		return errors.Forbidden(method, "Forbidden")
	}
	admin := false
	for _, s := range acc.Scopes {
		if s == "admin" || s == "service" {
			admin = true
			break
		}
	}
	if !admin {
		return errors.Forbidden(method, "Forbidden")
	}
	return nil
}

type dateEntry struct {
	Entries []listEntry
}

func (p *Publicapiusage) UsageCron() {
	defer func() {
		log.Infof("Usage sweep ended")
	}()
	log.Infof("Performing usage sweep")
	// loop through counters and persist
	ctx := context.Background()
	sc := p.c.redisClient.Scan(ctx, 0, prefixCounter+":*", 0)
	if err := sc.Err(); err != nil {
		log.Errorf("Error running redis scan %s", err)
		return
	}

	toPersist := map[string]map[string][]listEntry{} // userid->date->[]listEntry
	it := sc.Iterator()
	for {
		if !it.Next(ctx) {
			if err := it.Err(); err != nil {
				log.Errorf("Error during iteration %s", err)
			}
			break
		}

		key := it.Val()
		count, err := p.c.redisClient.Get(ctx, key).Int64()
		if err != nil {
			log.Errorf("Error retrieving value %s", err)
			return
		}
		parts := strings.Split(strings.TrimPrefix(key, prefixCounter+":"), ":")
		if len(parts) < 3 {
			log.Errorf("Unexpected number of components in key %s", key)
			continue
		}
		userID := parts[0]
		date := parts[1]
		service := parts[2]
		dates := toPersist[userID]
		if dates == nil {
			dates = map[string][]listEntry{}
			toPersist[userID] = dates
		}
		entries := dates[date]
		if entries == nil {
			entries = []listEntry{}
		}
		entries = append(entries, listEntry{
			Service: service,
			Count:   count,
		})
		dates[date] = entries
	}

	for userID, v := range toPersist {
		for date, entry := range v {
			de := dateEntry{
				Entries: entry,
			}
			b, err := json.Marshal(de)
			if err != nil {
				log.Errorf("Error marshalling entry %s", err)
				return
			}
			store.Write(&store.Record{
				Key:   fmt.Sprintf(prefixUsageByCustomer, userID, date),
				Value: b,
			})
		}

	}

}

func (p *Publicapiusage) Sweep(ctx context.Context, request *pb.SweepRequest, response *pb.SweepResponse) error {
	p.UsageCron()
	return nil
}
