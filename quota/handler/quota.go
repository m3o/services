package handler

import (
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	pb "github.com/m3o/services/quota/proto"
	v1api "github.com/m3o/services/v1api/proto"
	"github.com/micro/micro/v3/service/auth"
	"github.com/micro/micro/v3/service/client"
	"github.com/micro/micro/v3/service/config"
	"github.com/micro/micro/v3/service/errors"
	"github.com/micro/micro/v3/service/logger"
	log "github.com/micro/micro/v3/service/logger"
	"github.com/micro/micro/v3/service/store"

	"github.com/go-redis/redis/v8"
)

const (
	prefixCounter = "counter"

	prefixQuotaID       = "quota"
	prefixMapping       = "mapping"
	prefixMappingByPath = "mappingByPath"
)

type counter struct {
	sync.RWMutex
	redisClient *redis.Client
}

func (c *counter) incr(ns, userID, path string) (int64, error) {
	return c.redisClient.Incr(context.Background(), fmt.Sprintf("%s:%s:%s:%s", prefixCounter, ns, userID, path)).Result()
}

func (c *counter) decr(ns, userID, path string) (int64, error) {
	return c.redisClient.Decr(context.Background(), fmt.Sprintf("%s:%s:%s:%s", prefixCounter, ns, userID, path)).Result()
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
	if request.Quota == nil {
		return errors.BadRequest("quota.Create", "Missing quota")
	}
	quot := request.Quota
	if len(quot.Id) == 0 {
		return errors.BadRequest("quota.Create", "Missing quota ID")
	}
	if len(quot.Path) == 0 {
		return errors.BadRequest("quota.Create", "Missing quota Path")
	}
	quotObj := &quota{
		ID:             quot.Id,
		Limit:          quot.Limit,
		ResetFrequency: resetFrequency(quot.ResetFrequency.Number()),
		Path:           quot.Path,
	}
	if err := q.writeQuota(quotObj); err != nil {
		log.Errorf("Error marshalling json %s", err)
		return errors.InternalServerError("quota.Create", "Error creating quota")
	}

	return nil
}

func (q *Quota) writeQuota(quot *quota) error {
	b, err := json.Marshal(quot)
	if err != nil {
		return err

	}
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s", prefixQuotaID, quot.ID),
		Value: b,
	}); err != nil {
		return err
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
		_, err := q.readQuota(qID)
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
	for _, qID := range quotaIDs {

		m := mapping{
			UserID:    userID,
			Namespace: namespace,
			QuotaID:   qID,
		}
		if err := q.writeMapping(&m); err != nil {
			return err
		}

	}

	// update the v1api to unblock the user's api keys
	allowList := []string{}
	for _, qID := range quotaIDs {
		quot, err := q.readQuota(qID)
		if err != nil {
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

func (q *Quota) writeMapping(m *mapping) error {
	quot, err := q.readQuota(m.QuotaID)
	if err != nil {
		return err
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

	// path index
	if err := store.Write(&store.Record{
		Key:   fmt.Sprintf("%s:%s:%s:%s", prefixMappingByPath, m.Namespace, m.UserID, quot.Path),
		Value: b,
	}); err != nil {
		log.Errorf("Error writing mapping to store %s", err)
		return err
	}

	return nil
}

func (q *Quota) readMappingByPath(namespace, userID, path string) (*mapping, error) {
	recs, err := store.Read(fmt.Sprintf("%s:%s:%s:%s", prefixMappingByPath, namespace, userID, path))
	if err != nil {
		return nil, err
	}
	m := &mapping{}
	return m, json.Unmarshal(recs[0].Value, m)
}

func (q *Quota) readQuota(qID string) (*quota, error) {
	recs, err := store.Read(fmt.Sprintf("%s:%s", prefixQuotaID, qID))
	if err != nil {
		log.Errorf("Error looking up quota ID %s", err)
		return nil, err
	}
	quot := &quota{}
	if err := json.Unmarshal(recs[0].Value, quot); err != nil {
		log.Errorf("Error unmarshalling quota object %s", err)
		return nil, err
	}
	return quot, nil

}

func (q *Quota) ListUsage(ctx context.Context, request *pb.ListUsageRequest, response *pb.ListUsageResponse) error {
	acc, ok := auth.AccountFromContext(ctx)
	if !ok {
		return errors.Unauthorized("quota.ListUsage", "Unauthorized")
	}
	userID := acc.ID
	namespace := acc.Issuer
	if len(request.UserId) > 0 {
		// admins can see it all
		if err := verifyAdmin(ctx, "quota.ListUsage"); err != nil {
			return errors.BadRequest("quota.ListUsage", "Must be an admin to specify user ID")
		}
		userID = request.UserId
		namespace = request.Namespace
	}

	recs, err := store.Read(fmt.Sprintf("%s:%s:%s:", prefixMapping, namespace, userID), store.ReadPrefix())
	if err != nil && err != store.ErrNotFound {
		logger.Errorf("Error looking up mappings %s", err)
		return errors.InternalServerError("quota.ListUsage", "Error listing usage")
	}
	response.Usages = []*pb.QuotaUsage{}
	for _, r := range recs {
		m := &mapping{}
		if err := json.Unmarshal(r.Value, m); err != nil {
			logger.Errorf("Error unmarshalling mapping %s", err)
			return errors.InternalServerError("quota.ListUsage", "Error listing usage")
		}
		quot, err := q.readQuota(m.QuotaID)
		if err != nil {
			logger.Errorf("Error looking up quota %s", err)
			return errors.InternalServerError("quota.ListUsage", "Error listing usage")
		}

		count, err := q.c.read(namespace, userID, quot.Path)
		if err != nil && err != redis.Nil {
			logger.Errorf("Error getting counter value %s", err)
			return errors.InternalServerError("quota.ListUsage", "Error listing usage")
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
func (q *Quota) ResetQuotasCron() {
	logger.Infof("Resetting quotas for users")
	// loop through every mapping, check the corresponding quota, and reset if there is a limit and the frequency is right
	recs, err := store.Read(fmt.Sprintf("%s:", prefixMapping), store.ReadPrefix())
	if err != nil {
		logger.Errorf("Error reading mappings %s", err)
		// TODO - anything else?
		return
	}
	quotaCache := map[string]*quota{}

	type userUpdate = struct {
		m    *mapping
		quot *quota
	}
	updates := map[string][]userUpdate{}
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
			quot, err = q.readQuota(m.QuotaID)
			if err != nil {
				logger.Errorf("Error reading quotas %s", err)
				// TODO - anything else?
				continue
			}
			quotaCache[quot.ID] = quot
		}
		logger.Infof("Checking %+v %+v", *m, *quot)
		if !isTimeForReset(quot.ResetFrequency, time.Now()) {
			continue
		}
		userUpdates := updates[fmt.Sprintf("%s:%s", m.Namespace, m.UserID)]
		if userUpdates == nil {
			userUpdates = []userUpdate{}
			updates[fmt.Sprintf("%s:%s", m.Namespace, m.UserID)] = userUpdates
		}
		updates[fmt.Sprintf("%s:%s", m.Namespace, m.UserID)] = append(userUpdates, userUpdate{
			m:    m,
			quot: quot,
		})
	}

	for _, userUpdates := range updates {
		allowList := []string{}
		var userID, namespace string
		for _, update := range userUpdates {
			allowList = append(allowList, update.quot.Path)
			userID = update.m.UserID
			namespace = update.m.Namespace

			logger.Infof("Resetting %+v %+v", *update.m, *update.quot)
			// reset the counter
			if err := q.c.reset(update.m.Namespace, update.m.UserID, update.quot.Path); err != nil {
				logger.Errorf("Error resetting quota %s", err)
				// TODO - anything else?
				continue
			}

		}
		// unblock the user in api
		if _, err := q.v1Svc.UpdateAllowedPaths(context.TODO(), &v1api.UpdateAllowedPathsRequest{
			UserId:    userID,
			Namespace: namespace,
			Allowed:   allowList,
		}, client.WithAuthToken()); err != nil {
			logger.Errorf("Error updating allowed paths for %s %s %s", namespace, userID, err)
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

func (q *Quota) ResetQuotas(ctx context.Context, request *pb.ResetRequest, response *pb.ResetResponse) error {
	q.ResetQuotasCron()
	return nil
}

func (q *Quota) List(ctx context.Context, request *pb.ListRequest, response *pb.ListResponse) error {
	recs, err := store.Read(fmt.Sprintf("%s:", prefixQuotaID), store.ReadPrefix())
	if err != nil {
		logger.Errorf("Error reading quotas %s", err)
		// TODO - anything else?
		return errors.InternalServerError("quota.List", "Error listing quotas")
	}
	response.Quotas = make([]*pb.APIQuota, len(recs))
	for i, r := range recs {
		quot := &quota{}
		if err := json.Unmarshal(r.Value, quot); err != nil {
			logger.Errorf("Error unmarshalling quota")
			return errors.InternalServerError("quota.List", "Error listing quotas")
		}
		response.Quotas[i] = &pb.APIQuota{
			Id:             quot.ID,
			Limit:          quot.Limit,
			ResetFrequency: pb.APIQuota_Frequency(quot.ResetFrequency),
			Path:           quot.Path,
		}
	}
	return nil

}
