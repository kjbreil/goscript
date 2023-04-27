package goscript

import (
	"github.com/google/uuid"
	"sync"
)

type triggerRunning struct {
	m map[uuid.UUID]*bool
	s *sync.Mutex
}

func (tr *triggerRunning) get(u uuid.UUID) *bool {
	tr.s.Lock()
	defer tr.s.Unlock()
	if cb, ok := tr.m[u]; ok {
		// update the passed bool to match what is in the map
		return cb
	} else {
		b := new(bool)
		tr.m[u] = b
		return b
	}
}
