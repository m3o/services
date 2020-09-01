package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"math"
	"time"
	usage "usage/proto"

	nsproto "github.com/m3o/services/namespaces/proto"
	"github.com/micro/go-micro/v3/client"
	merrors "github.com/micro/go-micro/v3/errors"
	"github.com/micro/go-micro/v3/store"
	pb "github.com/micro/micro/v3/service/auth/proto"
	log "github.com/micro/micro/v3/service/logger"
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
)

type Usage struct {
	ns nsproto.NamespacesService
	as pb.AccountsService
}

func NewUsage(ns nsproto.NamespacesService, as pb.AccountsService) *Usage {
	u := &Usage{
		ns: ns,
		as: as,
	}
	go u.loop()
	return u
}

// List account history by namespace, or lists latest values for each namespace if history is not provided.
func (e *Usage) List(ctx context.Context, req *usage.ListRequest, rsp *usage.ListResponse) error {
	log.Info("Received Usage.ListAccounts request")

	key := accountByLatest
	if req.Namespace != "" {
		key = accountByNamespacePrefix + req.Namespace + "/"
	}
	limit := req.Limit
	if limit == 0 {
		limit = 20
	}
	records, err := mstore.Read(key, store.ReadPrefix(), store.ReadLimit(uint(limit)), store.ReadOffset(uint(req.Offset)))
	if err != nil && err != store.ErrNotFound {
		return merrors.InternalServerError("usage.listAccounts", "Error listing store: %v", err)
	}

	accounts := []*usage.Account{}
	for _, v := range records {
		u := &usg{}
		err = json.Unmarshal(v.Value, u)
		if err != nil {
			return merrors.InternalServerError("usage.listAccounts", "Error unmarsjaling value: %v", err)
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
					log.Warn("Error getting usage for namespace '%v': %v", namespace.Id, err)
					continue
				}
				u.Created = created.Unix()
				val, _ := json.Marshal(u)
				log.Infof("Saving usage for namespace '%v'", namespace.Id)

				// Save by namespace
				timeVal := math.MaxInt64 - (created.Unix() % 3600)
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v/%v/%v", accountByNamespacePrefix, namespace.Id, timeVal),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v/%v/%v", accountByTime, timeVal, namespace.Id),
					Value: val,
				})
				if err != nil {
					log.Warnf("Error writing to store: %v", err)
				}
				err = mstore.Write(&store.Record{
					Key:   fmt.Sprintf("%v/%v", accountByLatest, namespace.Id),
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
	serviceCount := 0
	for _, account := range arsp.Accounts {
		if account.Type == "user" {
			userCount++
		}
		if account.Type == "service" {
			serviceCount++
		}
	}
	return &usg{
		Users:     int64(userCount),
		Services:  int64(serviceCount),
		Namespace: namespace,
	}, nil
}
