package store

import (
	"encoding/json"
	err "errors"

	"github.com/micro/go-micro/v2"
	"github.com/micro/go-micro/v2/errors"

	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/apps/service/proto/apps"
)

var (
	// ErrNotFound is returned when an app doesn't exist with this key
	ErrNotFound = err.New("App not found")
)

// Store manages the services persited data
type Store struct {
	name  string
	store store.Store
}

// NewStore returns an initalised store
func NewStore(srv micro.Service) *Store {
	return &Store{
		name:  srv.Name(),
		store: store.DefaultStore,
	}
}

// Write inserts a new app into the store
func (s *Store) Write(app *pb.App) error {
	bytes, err := json.Marshal(app)
	if err != nil {
		return errors.InternalServerError(s.name, "Unable to marshal JSON: %v", err)
	}

	err = s.store.Write(&store.Record{Key: app.Id, Value: bytes})
	if err != nil {
		return errors.InternalServerError(s.name, "Unable to write to store: %v", err)
	}

	return nil
}

// Read looks up an app in the store
func (s *Store) Read(id string) (*pb.App, error) {
	recs, err := s.store.Read(id)
	if err == store.ErrNotFound {
		return nil, ErrNotFound
	} else if err != nil {
		return nil, errors.InternalServerError(s.name, "Unable to read from store: %v", err)
	}

	var app *pb.App
	if err := json.Unmarshal(recs[0].Value, &app); err != nil {
		return nil, errors.InternalServerError(s.name, "Unable to unmarshal JSON: %v", err)
	}

	return app, nil
}

// List returns all the apps in the store
func (s *Store) List() ([]*pb.App, error) {
	recs, err := s.store.Read("", store.ReadPrefix())
	if err == store.ErrNotFound {
		return make([]*pb.App, 0), nil
	} else if err != nil {
		return nil, errors.InternalServerError(s.name, "Unable to read from store: %v", err)
	}

	apps := make([]*pb.App, len(recs))
	for i, a := range recs {
		if err := json.Unmarshal(a.Value, &apps[i]); err != nil {
			return nil, errors.InternalServerError(s.name, "Unable to unmarshal JSON: %v", err)
		}
	}

	return apps, nil
}
