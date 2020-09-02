package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"

	usage "github.com/m3o/services/usage/proto"

	"github.com/google/uuid"
	nsproto "github.com/m3o/services/namespaces/proto"
	"github.com/micro/go-micro/v3/client"
	merrors "github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	pb "github.com/micro/micro/v3/service/auth/proto"
	log "github.com/micro/micro/v3/service/logger"
	rproto "github.com/micro/micro/v3/service/runtime/proto"
	mstore "github.com/micro/micro/v3/service/store"
)

const (
	// format `account/namespace/timestamp`
	accountByNamespacePrefix = "account/"
	// format `account-by-time/timestamp/namespace`
	// to help listing all accounts by a time
	accountByTime = "account-by-time/"
	// format `account-latest` to help listing
	// latest measurements
	accountByLatest = "account-latest/"
	// format ''account-distinct/$month/namespace/users/$countvalue`
	// and 'account-distinct/$month/namespace/services/$countvalue`
	accountByDistinct = "account-by-distinct/"
	monthFormat       = "2006-01"
)

type Usage struct {
	ns      nsproto.NamespacesService
	as      pb.AccountsService
	runtime rproto.RuntimeService
}

func NewUsage(ns nsproto.NamespacesService, as pb.AccountsService, runtime rproto.RuntimeService) *Usage {
	u := &Usage{
		ns:      ns,
		as:      as,
		runtime: runtime,
	}
	go u.loop()
	return u
}

// List account history by namespace, or lists latest values for each namespace if history is not provided.
func (e *Usage) List(ctx context.Context, req *usage.ListRequest, rsp *usage.ListResponse) error {
	key := accountByLatest
	if len(req.Namespace) > 0 {
		key = accountByNamespacePrefix + req.Namespace + "/"
	}
	if req.Distinct {
		month := time.Now().Format(monthFormat)
		if req.DistinctTime > 0 {
			month = time.Unix(req.DistinctTime, 0).Format(monthFormat)
		}
		if len(req.Namespace) > 0 {
			key = fmt.Sprintf("%v%v/%v/", accountByDistinct, month, req.Namespace)
		} else {
			key = fmt.Sprintf("%v%v/", accountByDistinct, month)
		}
	}
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}

	log.Infof("Received Usage.ListAccounts request, listing with key '%v', limit '%v'", key, limit)

	records, err := mstore.Read(key, store.ReadPrefix(), store.ReadLimit(uint(limit)), store.ReadOffset(uint(req.Offset)))
	if err != nil && err != store.ErrNotFound {
		return merrors.InternalServerError("usage.list", "Error listing store: %v", err)
	}

	accounts := []*usage.Account{}
	for _, v := range records {
		u := &usg{}
		err = json.Unmarshal(v.Value, u)
		if err != nil {
			return merrors.InternalServerError("usage.list", "Error unmarsjaling value: %v", err)
		}
		accounts = append(accounts, &usage.Account{
			Namespace: u.Namespace,
			Users:     u.Users,
			Services:  u.Services,
			Created:   u.Created,
		})
	}
	rsp.Accounts = accounts
	return nil
}

func (e *Usage) loop() {
	for {
		func() {
			created := time.Now()
			rsp, err := e.ns.List(context.TODO(), &nsproto.ListRequest{}, client.WithAuthToken())
			if err != nil {
				log.Errorf("Error calling namespace service: %v", err)
				return
			}
			if len(rsp.Namespaces) == 0 {
				log.Warnf("Empty namespace list")
				return
			}
			log.Infof("Got %v namespaces", len(rsp.Namespaces))
			for _, namespace := range rsp.Namespaces {
				u, err := e.usageForNamespace(namespace.Id)
				if err != nil {
					log.Warnf("Error getting usage for namespace '%v': %v", namespace.Id, err)
					continue
				}
				u.Created = created.Unix()
				u.Id = uuid.New().String()
				val, _ := json.Marshal(u)
				log.Infof("Saving usage for namespace '%v'", namespace.Id)

				// Save by namespace
				timeVal := math.MaxInt64 - (created.Unix() % 3600 * 24) // 1 day
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v%v/%v", accountByNamespacePrefix, namespace.Id, timeVal),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v%v/%v", accountByTime, timeVal, namespace.Id),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v%v", accountByLatest, namespace.Id),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
				month := created.Format(monthFormat)
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v%v/%v/users/%v", accountByDistinct, month, namespace.Id, u.Users),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v%v/%v/services/%v", accountByDistinct, month, namespace.Id, u.Services),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
			}
		}()

		time.Sleep(1 * time.Hour)
	}
}

type usg struct {
	Id        string
	Users     int64
	Services  int64
	Created   int64
	Namespace string
}

func (e *Usage) usageForNamespace(namespace string) (*usg, error) {
	arsp, err := e.as.List(context.TODO(), &pb.ListAccountsRequest{
		Options: &pb.Options{
			Namespace: namespace,
		},
	}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}
	userCount := 0
	for _, account := range arsp.Accounts {
		if account.Type == "user" {
			userCount++
		}
	}
	rrsp, err := e.runtime.Read(context.TODO(), &rproto.ReadRequest{
		Options: &rproto.ReadOptions{
			Namespace: namespace,
		},
	}, client.WithAuthToken())
	if err != nil {
		return nil, err
	}
	serviceCount := len(rrsp.Services)
	return &usg{
		Users:     int64(userCount),
		Services:  int64(serviceCount),
		Namespace: namespace,
	}, nil
}
