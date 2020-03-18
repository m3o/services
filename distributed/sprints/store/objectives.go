package store

import (
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

// ReadObjective returns a obj from the store
func (s *Store) ReadObjective(sprintID, objID string) (*pb.Objective, error) {
	recs, err := s.store.Read(fmt.Sprintf("objectives/%v/%v", sprintID, objID))
	if err != nil {
		return nil, err
	}

	var obj *pb.Objective
	if err := json.Unmarshal(recs[0].Value, &obj); err != nil {
		return nil, err
	}

	return obj, nil
}

// CreateObjective writes an objective into the store
func (s *Store) CreateObjective(sprintID string, obj *pb.Objective) error {
	return s.UpdateObjective(sprintID, obj)
}

// UpdateObjective updates an objective in the store
func (s *Store) UpdateObjective(sprintID string, obj *pb.Objective) error {
	bytes, err := json.Marshal(obj)
	if err != nil {
		return err
	}

	return s.store.Write(&store.Record{
		Key:   fmt.Sprintf("objectives/%v/%v", sprintID, obj.Id),
		Value: bytes,
	})
}

// DeleteObjective removes an objective from the store
func (s *Store) DeleteObjective(sprintID string, objID string) error {
	return s.store.Delete(fmt.Sprintf("objectives/%v/%v", sprintID, objID))
}

// ListObjectives returns all the objectives associated with a sprint
func (s *Store) ListObjectives(sprintID string) ([]*pb.Objective, error) {
	recs, err := s.store.Read(fmt.Sprintf("objectives/%v/", sprintID), store.ReadPrefix())
	if err != nil {
		return make([]*pb.Objective, 0), err
	}

	objs := make([]*pb.Objective, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &objs[i]); err != nil {
			return make([]*pb.Objective, 0), err
		}
	}

	return objs, nil
}
