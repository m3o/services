package store

import (
	"github.com/micro/go-micro/v2/store"
)

// Store is a wrapper around go-micro store for the sprints service
type Store struct {
	store store.Store
}

// NewStore returns an initialised store
func NewStore(srvName string) *Store {
	return &Store{store: store.DefaultStore}
}
