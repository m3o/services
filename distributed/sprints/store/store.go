package store

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/store"
)

// Store is a wrapper around go-micro store for the sprints service
type Store struct {
	store store.Store
}

// NewStore returns an initialised store
func NewStore(srv micro.Service) *Store {
	return &Store{store: srv.Options().Store}
}
