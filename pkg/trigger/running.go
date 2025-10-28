package trigger

import (
	"sync"

	"github.com/google/uuid"
)

type Running struct {
	m map[uuid.UUID]*bool
	s *sync.Mutex
}

func NewRunning() Running {
	return Running{
		m: make(map[uuid.UUID]*bool),
		s: &sync.Mutex{},
	}
}

func (tr *Running) Get(u uuid.UUID) *bool {
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
