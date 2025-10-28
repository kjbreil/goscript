package trigger

import (
	"log/slog"
	"sync"

	"github.com/google/uuid"
)

const MaxTasksPerUUID = 1000 // Maximum queued tasks per UUID

type TaskMap struct {
	tasks  map[uuid.UUID][]*Task
	m      *sync.Mutex
	logger *slog.Logger
}

func NewTaskMap() TaskMap {
	return TaskMap{
		tasks: make(map[uuid.UUID][]*Task),
		m:     &sync.Mutex{},
	}
}

func NewTaskMapWithLogger(logger *slog.Logger) TaskMap {
	return TaskMap{
		tasks:  make(map[uuid.UUID][]*Task),
		m:      &sync.Mutex{},
		logger: logger,
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

	currentLen := len(tm.tasks[t.UUID()])

	// Warn if queue is getting large
	if tm.logger != nil {
		if currentLen > MaxTasksPerUUID {
			tm.logger.Error(
				"task queue exceeded maximum",
				"uuid",
				t.UUID(),
				"queued",
				currentLen,
				"max",
				MaxTasksPerUUID,
			)
			return // Drop task to prevent unbounded growth
		} else if currentLen > MaxTasksPerUUID/2 {
			tm.logger.Warn("task queue growing large", "uuid", t.UUID(), "queued", currentLen, "max", MaxTasksPerUUID)
		}
	}

	tm.tasks[t.UUID()] = append(tm.tasks[t.UUID()], t)
}

func (tm *TaskMap) Delete(t *Task) {
	tm.m.Lock()
	defer tm.m.Unlock()
	delete(tm.tasks, t.UUID())
}
