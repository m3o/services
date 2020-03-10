package handler

import (
	"github.com/micro/go-micro/v2"
	"github.com/micro/services/distributed/sprints/store"
)

// Handler implements the Sprints Handler interface
type Handler struct {
	store *store.Store
	name  string
}

// NewHandler returns an initialised handler
func NewHandler(srv micro.Service) *Handler {
	return &Handler{
		name:  srv.Name(),
		store: store.NewStore(srv.Name()),
	}
}
