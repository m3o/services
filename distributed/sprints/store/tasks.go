package store

import (
	"encoding/json"
	"fmt"

	"github.com/micro/go-micro/v2/store"
	pb "github.com/micro/services/distributed/sprints/proto/sprints"
)

// ReadTask returns a task from the store
func (s *Store) ReadTask(sprintID, taskID string) (*pb.Task, error) {
	recs, err := s.store.Read(fmt.Sprintf("tasks/%v/%v", sprintID, taskID))
	if err != nil {
		return nil, err
	}
	if len(recs) != 1 {
		return nil, store.ErrNotFound
	}

	var task *pb.Task
	if err := json.Unmarshal(recs[0].Value, &task); err != nil {
		return nil, err
	}

	return task, nil
}

// CreateTask writes a task into the store
func (s *Store) CreateTask(sprintID string, task *pb.Task) error {
	return s.UpdateTask(sprintID, task)
}

// UpdateTask updates a task in the store
func (s *Store) UpdateTask(sprintID string, task *pb.Task) error {
	bytes, err := json.Marshal(task)
	if err != nil {
		return err
	}

	return s.store.Write(&store.Record{
		Key:   fmt.Sprintf("tasks/%v/%v", sprintID, task.Id),
		Value: bytes,
	})
}

// DeleteTask removes a task from the store
func (s *Store) DeleteTask(sprintID string, taskID string) error {
	return s.store.Delete(fmt.Sprintf("tasks/%v/%v", sprintID, taskID))
}

// ListTasks returns all the tasks associated with a sprint
func (s *Store) ListTasks(sprintID string) ([]*pb.Task, error) {
	recs, err := s.store.Read(fmt.Sprintf("tasks/%v/", sprintID), store.ReadPrefix())
	if err != nil {
		return make([]*pb.Task, 0), err
	}

	tasks := make([]*pb.Task, len(recs))
	for i, r := range recs {
		if err := json.Unmarshal(r.Value, &tasks[i]); err != nil {
			return make([]*pb.Task, 0), err
		}
	}

	return tasks, nil
}
