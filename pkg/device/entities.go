package device

import "sync"

type Entities struct {
	e  map[string]Entity
	mu *sync.Mutex
}
