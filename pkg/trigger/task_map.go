package trigger

import (
	"github.com/google/uuid"
	"sync"
)

type TaskMap struct {
	tasks map[uuid.UUID][]*Task
	m     *sync.Mutex
}

func NewTaskMap() TaskMap {
	return TaskMap{
		tasks: make(map[uuid.UUID][]*Task),
		m:     &sync.Mutex{},
	}
}

func (tm *TaskMap) Lock() {
	tm.m.Lock()
}

func (tm *TaskMap) Unlock() {
	tm.m.Unlock()
}

func (tm *TaskMap) ToRun() []*Task {
	var toRun []*Task
	var ran []uuid.UUID

	tm.m.Lock()
	defer tm.m.Unlock()
	for u, tasks := range tm.tasks {
		if len(tasks) > 0 {
			t := tasks[0]
			if !t.Running() {
				toRun = append(toRun, t)
				tm.tasks[u] = tasks[1:]
			}
		}
		if len(tasks) == 0 {
			ran = append(ran, u)
		}
	}
	for _, u := range ran {
		delete(tm.tasks, u)
	}
	return toRun
}

func (tm *TaskMap) Add(t *Task) {
	if t == nil {
		return
	}
	tm.m.Lock()
	defer tm.m.Unlock()
	tm.tasks[t.UUID()] = append(tm.tasks[t.UUID()], t)
}

func (tm *TaskMap) Delete(t *Task) {
	tm.m.Lock()
	defer tm.m.Unlock()
	delete(tm.tasks, t.UUID())
}
