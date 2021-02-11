package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/logger"

	pb "github.com/m3o/services/quota/proto"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/errors"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/go-redis/redis/v8"
)

const (
	prefixCounter = "counter"

	prefixQuotaID = "quota"
	prefixMapping = "mapping"
)

type counter struct {
	sync.RWMutex
	redisClient *redis.Client
}

func (c *counter) incr(ns, userID, path string) (int64, error) {
	return c.redisClient.Incr(context.Background(), fmt.Sprintf("%s:%s:%s:%s", prefixCounter, ns, userID, path)).Result()
}

func (c *counter) read(ns, userID, path string) (int64, error) {
	return c.redisClient.Get(context.Background(), fmt.Sprintf("%s:%s:%s:%s", prefixCounter, ns, userID, path)).Int64()
}

func (c *counter) reset(ns, userID, path string) error {
	return c.redisClient.Set(context.Background(), fmt.Sprintf("%s:%s:%s:%s", prefixCounter, ns, userID, path), 0, 0).Err()
}

type Quota struct {
	v1Svc v1api.V1Service
	c     counter
}

type resetFrequency int

const (
	Never resetFrequency = iota
	Daily
	Monthly
)

func (r resetFrequency) String() string {
	return [...]string{"Never", "Daily", "Monthly"}[r]
}

type quota struct {
	ID             string
	Limit          int64
	ResetFrequency resetFrequency
	Path           string
}

type mapping struct {
	UserID    string
	Namespace string
	QuotaID   string
}

func New(client client.Client) *Quota {
	redisConfig := struct {
		Address  string
		User     string
		Password string
	}{}
	val, err := config.Get("micro.quota.redis")
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
	q := &Quota{
		v1Svc: v1api.NewV1Service("v1", client),
		c:     counter{redisClient: rc},
	}
	go q.consumeEvents()
	return q
}

func (q *Quota) Create(ctx context.Context, request *pb.CreateRequest, response *pb.CreateResponse) error {
	if err := verifyAdmin(ctx, "quota.Create"); err != nil {
		return err
	}
	if len(request.Id) == 0 {
		return errors.BadRequest("quota.Create", "Missing quota ID")
	}
	if len(request.Path) == 0 {
		return errors.BadRequest("quota.Create", "Missing quota Path")
	}
	quot := &quota{
		ID:             request.Id,
		Limit:          request.Limit,
		ResetFrequency: resetFrequency(request.ResetFrequency.Number()),
		Path:           request.Path,
	}

	b, err := json.Marshal(quot)
	if err != nil {
		log.Errorf("Error marshalling json %s", err)
		return errors.InternalServerError("quota.Create", "Error creating quota")
	}
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s", prefixQuotaID, quot.ID),
		Value: b,
	}); err != nil {
		log.Errorf("Error writing to store %s", err)
		return errors.InternalServerError("quota.Create", "Error creating quota")
	}

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

func (q *Quota) RegisterUser(ctx context.Context, request *pb.RegisterUserRequest, response *pb.RegisterUserResponse) error {
	if err := verifyAdmin(ctx, "quota.RegisterUser"); err != nil {
		return err
	}

	if len(request.UserId) == 0 {
		return errors.BadRequest("quota.RegisterUser", "Missing UserID")
	}
	if len(request.Namespace) == 0 {
		return errors.BadRequest("quota.RegisterUser", "Missing Namespace")
	}

	if len(request.QuotaIds) == 0 {
		return errors.BadRequest("quota.RegisterUser", "Missing QuotaIDs")
	}
	// validate all the quota IDs first
	for _, qID := range request.QuotaIds {
		// is this quota legit?
		_, err := store.Read(fmt.Sprintf("%s:%s", prefixQuotaID, qID))
		if err != nil {
			if err == store.ErrNotFound {
				return errors.BadRequest("quota.RegisterUser", "Quota ID not recognised: %s", qID)
			}
			log.Errorf("Error looking up quota ID %s", err)
			return errors.InternalServerError("quota.RegisterUser", "Error registering user")
		}
	}

	if err := q.registerUser(request.UserId, request.Namespace, request.QuotaIds); err != nil {
		return errors.InternalServerError("quota.RegisterUser", "Error registering user")
	}
	return nil

}

func (q *Quota) registerUser(userID, namespace string, quotaIDs []string) error {

	// store association for each quota
	for _, q := range quotaIDs {

		m := mapping{
			UserID:    userID,
			Namespace: namespace,
			QuotaID:   q,
		}

		b, err := json.Marshal(m)
		if err != nil {
			log.Errorf("Error marshalling mapping %s", err)
			return err
		}
		if err := store.Write(&store.Record{
			Key:   fmt.Sprintf("%s:%s:%s:%s", prefixMapping, m.Namespace, m.UserID, m.QuotaID),
			Value: b,
		}); err != nil {
			log.Errorf("Error writing mapping to store %s", err)
			return err
		}
	}

	// update the v1api to unblock the user's api keys
	allowList := []string{}
	for _, qID := range quotaIDs {
		recs, err := store.Read(fmt.Sprintf("%s:%s", prefixQuotaID, qID))
		if err != nil {
			log.Errorf("Error looking up quota ID %s", err)
			return err
		}
		quot := &quota{}
		if err := json.Unmarshal(recs[0].Value, quot); err != nil {
			log.Errorf("Error unmarshalling quota object %s", err)
			return err
		}
		allowList = append(allowList, quot.Path)

	}

	if _, err := q.v1Svc.UpdateAllowedPaths(context.TODO(), &v1api.UpdateAllowedPathsRequest{
		UserId:    userID,
		Namespace: namespace,
		Allowed:   allowList,
	}, client.WithAuthToken()); err != nil {
		logger.Errorf("Error updating allowed paths %s", err)
		return err
	}
	return nil
}

func (q *Quota) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("quota.List", "Unauthorized")
	}
	userID := acc.ID
	namespace := acc.Issuer
	if len(request.UserId) > 0 && request.UserId != userID {
		// admins can see it all
		if err := verifyAdmin(ctx, "quota.List"); err != nil {
			return err
		}
		userID = request.UserId
		namespace = request.Namespace
	}

	recs, err := store.Read(fmt.Sprintf("%s:%s:%s", prefixMapping, namespace, userID), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		logger.Errorf("Error looking up mappings %s", err)
		return errors.InternalServerError("quota.List", "Error listing usage")
	}
	response.Usages = []*pb.QuotaUsage{}
	for _, r := range recs {
		m := &mapping{}
		if err := json.Unmarshal(r.Value, m); err != nil {
			logger.Errorf("Error unmarshalling mapping %s", err)
			return errors.InternalServerError("quota.List", "Error listing usage")
		}
		qrecs, err := store.Read(fmt.Sprintf("%s:%s", prefixQuotaID, m.QuotaID))
		if err != nil {
			logger.Errorf("Error reading  %s", err)
			return errors.InternalServerError("quota.List", "Error listing usage")
		}
		quot := &quota{}

		if err := json.Unmarshal(qrecs[0].Value, quot); err != nil {
			logger.Errorf("Error reading  %s", err)
			return errors.InternalServerError("quota.List", "Error listing usage")
		}

		count, err := q.c.read(namespace, userID, quot.Path)
		if err != nil && err != redis.Nil {
			logger.Errorf("Error getting counter value %s", err)
			return errors.InternalServerError("quota.List", "Error listing usage")
		}
		response.Usages = append(response.Usages, &pb.QuotaUsage{
			Name:  quot.ID,
			Usage: count,
			Limit: quot.Limit,
		})

	}
	return nil
}

// ResetQuotas runs daily to reset usage counters in the case of daily or monthly quotas
// TODO make this work across multiple instances by either using distributed locking or an external trigger (k8s cron)
func (q *Quota) ResetQuotas() {
	// loop through every mapping, check the corresponding quota, and reset if there is a limit and the frequency is right
	recs, err := store.Read(fmt.Sprintf("%s:", prefixMapping), store.ReadPrefix())
	if err != nil {
		logger.Errorf("Error reading mappings %s", err)
		// TODO - anything else?
		return
	}
	quotaCache := map[string]*quota{}
	for _, r := range recs {
		m := &mapping{}
		if err := json.Unmarshal(r.Value, m); err != nil {
			logger.Errorf("Error unmarshalling mapping %s", err)
			// TODO - anything else?
			continue
		}
		quot := quotaCache[m.QuotaID]
		if quot == nil {
			// load up the quota
			qrecs, err := store.Read(fmt.Sprintf("%s:%s", prefixQuotaID, m.QuotaID), store.ReadPrefix())
			if err != nil {
				logger.Errorf("Error reading quotas %s", err)
				// TODO - anything else?
				continue
			}
			quot = &quota{}
			if err := json.Unmarshal(qrecs[0].Value, quot); err != nil {
				logger.Errorf("Error unmarshalling quota %s", err)
				// TODO - anything else?
				continue
			}
			quotaCache[quot.ID] = quot
		}
		if !isTimeForReset(quot.ResetFrequency, time.Now()) {
			continue
		}
		// reset the counter
		if err := q.c.reset(m.Namespace, m.UserID, quot.Path); err != nil {
			logger.Errorf("Error unmarshalling quota %s", err)
			// TODO - anything else?
			continue
		}
	}
}

func isTimeForReset(frequency resetFrequency, t time.Time) bool {
	switch frequency {
	case Never:
		return false
	case Daily:
		// assumes the cron is called once a day
		return true
	case Monthly:
		return t.Day() == 1
	}
	return false
}
