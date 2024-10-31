package store

import (
	"errors"

	"github.com/acmk189/todo_app_REST_API/entity"
)

var (
	Tasks = &TaskStore{Tasks: map[entity.TaskID]*entity.Task{}}

	ErrNotFound = errors.New("not found")
)

type TaskStore struct {
	// 動作確認用
	LastID entity.TaskID
	Tasks  map[entity.TaskID]*entity.Task
}

func (ts *TaskStore) Add(t *entity.Task) (int, error) {
	ts.LastID++
	t.ID = ts.LastID
	ts.Tasks[t.ID] = t
	return int(t.ID), nil
}

func (ts *TaskStore) All() entity.Tasks {
	tasks := entity.Tasks{}
	for _, t := range ts.Tasks {
		tasks = append(tasks, *t)
	}
	return tasks
}
