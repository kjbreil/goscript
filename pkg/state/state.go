package state

import (
	"github.com/kjbreil/hass-ws/model"
	"strings"
	"sync"
	"time"
)

type States struct {
	s map[string]*State
	m *sync.Mutex
}
type State struct {
	DomainEntity string
	Domain       string
	Entity       string
	State        StateText
	Attributes   map[string]interface{}
	LastChanged  time.Time
	LastUpdated  time.Time
}

func NewStates() States {
	return States{
		s: make(map[string]*State),
		m: &sync.Mutex{},
	}
}

func NewSingleStates(entityID string, eState *State) States {
	return States{
		s: map[string]*State{entityID: eState},
		m: &sync.Mutex{},
	}
}

func NewMultiStates(states map[string]*State) States {
	return States{
		s: states,
		m: &sync.Mutex{},
	}
}
func StateFromWS(s *model.State) *State {
	if s.State == nil || s.EntityId == nil {
		return nil
	}

	return &State{
		DomainEntity: *s.EntityId,
		Domain:       s.Domain(),
		Entity:       s.Entity(),
		State:        StateText(*s.State),
		LastChanged:  *s.LastChanged,
		LastUpdated:  *s.LastUpdated,
		Attributes:   s.Attributes,
	}
}

func (s *States) Lock() {
	s.m.Lock()
}

func (s *States) Unlock() {
	s.m.Unlock()
}

func (s *States) Len() int {
	return len(s.s)
}

// Insert only adds to the map if something does not exist already. Returns what is in the map whether added or not
func (s *States) Insert(ps *State) *State {
	s.m.Lock()
	defer s.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := s.s[ps.DomainEntity]; ok {
		return st
	} else {
		s.s[ps.DomainEntity] = ps
		return ps
	}
}

// Upsert inserts a new record if one does not exist otherwise updates the data at the pointer so the update propagates
func (s *States) Upsert(ps *State) *State {
	s.m.Lock()
	defer s.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := s.s[ps.DomainEntity]; ok {
		*st = *ps
	} else {
		s.s[ps.DomainEntity] = ps
	}
	return ps
}

// Combine takes two States objects and merges them, passed object will overwrite a state in current object
func (s *States) Combine(cs *States) {
	s.m.Lock()
	defer s.m.Unlock()

	cs.m.Lock()
	defer cs.m.Unlock()
	for k, v := range cs.s {
		s.s[k] = v
	}
}

// Entities returns a string of the entities contained in the States object
func (s *States) Entities() []string {
	s.m.Lock()
	defer s.m.Unlock()

	en := make([]string, 0, len(s.s))
	for _, st := range s.s {
		en = append(en, st.DomainEntity)
	}
	return en
}

// Get returns a single state record and a bool if found
func (s *States) Get(key string) (*State, bool) {
	s.m.Lock()
	defer s.m.Unlock()

	st, ok := s.s[key]
	if !ok {
		return nil, ok
	}
	return st, true
}

// Find returns a new map of states of the passed entities
func (s *States) Find(entities []string) map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	ss := make(map[string]*State)

	for _, k := range entities {
		if st, ok := s.s[k]; ok {
			ss[st.DomainEntity] = st
		}
	}

	return ss
}

// FindDomainMap returns a map of the states for the passed domain
func (s *States) FindDomainMap(keys []string) map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	ss := make(map[string]*State)

	for _, st := range s.s {
		for _, k := range keys {
			if st.Domain == k {
				ss[st.DomainEntity] = st
			}
		}
	}

	return ss
}

// Slice returns a slice of the states in no particular order
func (s *States) Slice() []*State {
	s.m.Lock()
	defer s.m.Unlock()

	var ss []*State
	for _, st := range s.s {
		ss = append(ss, st)
	}

	return ss
}

// Map returns a map of all the states
func (s *States) Map() map[string]*State {
	s.m.Lock()
	defer s.m.Unlock()

	sts := make(map[string]*State)

	for k, v := range s.s {
		sts[k] = v
	}

	return sts
}

// SubSet returns a new States which contains a subset of the current states based on entities passed
func (s *States) SubSet(entities []string) States {
	return States{
		s: s.Find(entities),
		m: &sync.Mutex{},
	}
}

// Where returns a new States object containing all the states that match the passed state. strings.Equalfold is used
// for the comparison.
func (s *States) Where(state string) *States {
	sts := States{
		s: make(map[string]*State),
		m: &sync.Mutex{},
	}

	for _, v := range s.Slice() {
		if strings.EqualFold(string(v.State), state) {
			sts.Upsert(v)
		}
	}
	return &sts
}

func MessageState(message *model.Message) *State {
	return &State{
		DomainEntity: message.DomainEntity(),
		Domain:       message.Domain(),
		Entity:       message.EntityID(),
		State:        StateText(message.State()),
		Attributes:   message.Attributes(),
	}
}
