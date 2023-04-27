package goscript

import (
	"github.com/google/uuid"
	"sync"
)

type taskMap struct {
	tasks map[uuid.UUID][]*Task
	m     *sync.Mutex
}

func (tr *taskMap) add(t *Task) {
	if t == nil {
		return
	}
	tr.m.Lock()
	defer tr.m.Unlock()
	tr.tasks[t.uuid] = append(tr.tasks[t.uuid], t)
}

func (tr *taskMap) delete(t *Task) {
	tr.m.Lock()
	defer tr.m.Unlock()
	delete(tr.tasks, t.uuid)
}
