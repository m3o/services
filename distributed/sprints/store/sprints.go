package store

import (
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

// CreateSprint writes a sprint into the store
func (s *Store) CreateSprint(sprint *pb.Sprint) error {
	// ensure nested attributes aren't written to the store
	sprint.Objectives = nil
	sprint.Tasks = nil

	bytes, err := json.Marshal(sprint)
	if err != nil {
		return err
	}

	return s.store.Write(&store.Record{
		Key:   fmt.Sprintf("sprints/%v", sprint.Id),
		Value: bytes,
	})
}

// ReadSprint looks up a sprint using its ID
func (s *Store) ReadSprint(id string) (*pb.Sprint, error) {
	recs, err := s.store.Read(fmt.Sprintf("sprints/%v", id))
	if err != nil {
		return nil, err
	}
	if len(recs) != 1 {
		return nil, store.ErrNotFound
	}

	var sprint *pb.Sprint
	if err := json.Unmarshal(recs[0].Value, sprint); err != nil {
		return nil, err
	}

	return sprint, err
}

// CountSprints returns the number of sprints in the store
func (s *Store) CountSprints() (int, error) {
	recs, err := s.store.Read("sprints/", store.ReadPrefix())
	if err != nil {
		return 0, err
	}
	return len(recs), nil
}

// ListSprints returns all the sprints in the store
func (s *Store) ListSprints() ([]*pb.Sprint, error) {
	recs, err := s.store.Read("sprints/", store.ReadPrefix())
	if err != nil {
		return make([]*pb.Sprint, 0), err
	}

	sprints := make([]*pb.Sprint, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, sprints[i]); err != nil {
			return make([]*pb.Sprint, 0), err
		}
	}

	return sprints, nil
}
