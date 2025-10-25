package taskstore

import (
	"fmt"
	"sync"
	"time"
)

type Task struct {
	Id   int       `json:"id"`
	Text string    `json:"text"`
	Tags []string  `json:"tags"`
	Due  time.Time `json:"due"`
}

type TaskStore struct {
	sync.Mutex
	tasks  map[int]Task
	nextId int
}

func New() *TaskStore {
	ts := &TaskStore{}
	ts.tasks = make(map[int]Task)
	return ts
}

func (ts *TaskStore) CreateTask(text string, tags []string, due time.Time) int {
	ts.Lock()
	defer ts.Unlock()

	task := Task{
		Id:   ts.nextId,
		Text: text,
		Due:  due,
	}

	task.Tags = make([]string, len(tags))
	copy(task.Tags, tags)

	ts.tasks[ts.nextId] = task
	ts.nextId++
	return task.Id
}

func (ts *TaskStore) GetTask(id int) (Task, error) {
	ts.Lock()
	defer ts.Unlock()

	task, ok := ts.tasks[id]
	if !ok {
		err := fmt.Errorf("No task exist with id: %d", id)
		return Task{}, err
	}

	return task, nil
}

func (ts *TaskStore) DeleteTask(id int) error {
	ts.Lock()
	defer ts.Unlock()

	_, ok := ts.tasks[id]
	if !ok {
		err := fmt.Errorf("Unable to delete, no task exist with id: %d", id)
		return err
	}

	delete(ts.tasks, id)
	return nil
}

func (ts *TaskStore) DeleteAllTasks() error {
	ts.Lock()
	defer ts.Unlock()

	ts.tasks = make(map[int]Task)
	return nil
}

func (ts *TaskStore) GetAllTasks() []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		tasks = append(tasks, task)
	}

	return tasks
}

func (ts *TaskStore) GetTasksByTag(tag string) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

TaskLoop:
	for _, task := range ts.tasks {
		for _, taskTag := range task.Tags {
			if taskTag == tag {
				tasks = append(tasks, task)
				continue TaskLoop
			}
		}
	}

	return tasks
}

func (ts *TaskStore) GetTasksByDueDate(dueDate time.Time) []Task {
	ts.Lock()
	defer ts.Unlock()

	var tasks []Task

	for _, task := range ts.tasks {
		if task.Due.Equal(dueDate) {
			tasks = append(tasks, task)
		}
	}

	return tasks
}
