package state

import (
	"maps"
	"strings"
	"sync"
	"time"

	"github.com/kjbreil/hass-ws/model"
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

func (ss *States) Lock() {
	ss.m.Lock()
}

func (ss *States) Unlock() {
	ss.m.Unlock()
}

func (ss *States) Len() int {
	return len(ss.s)
}

// Insert only adds to the map if something does not exist already. Returns what is in the map whether added or not.
func (ss *States) Insert(ps *State) *State {
	ss.m.Lock()
	defer ss.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := ss.s[ps.DomainEntity]; ok {
		return st
	} else {
		ss.s[ps.DomainEntity] = ps
		return ps
	}
}

// Upsert inserts a new record if one does not exist otherwise updates the data at the pointer so the update propagates.
func (ss *States) Upsert(ps *State) *State {
	ss.m.Lock()
	defer ss.m.Unlock()
	// Update the state pointer so it follows to tasks
	if st, ok := ss.s[ps.DomainEntity]; ok {
		*st = *ps
	} else {
		ss.s[ps.DomainEntity] = ps
	}
	return ps
}

// Combine takes two States objects and merges them, passed object will overwrite a state in current object
// To avoid deadlocks, we first create a copy of cs's data, then merge it into ss.
func (ss *States) Combine(cs *States) {
	// First, get a copy of cs's data while holding only cs's lock
	cs.m.Lock()
	csCopy := make(map[string]*State, len(cs.s))
	for k, v := range cs.s {
		csCopy[k] = v
	}
	cs.m.Unlock()

	// Now merge the copy into ss while holding only ss's lock
	ss.m.Lock()
	defer ss.m.Unlock()
	for k, v := range csCopy {
		ss.s[k] = v
	}
}

// Entities returns a string of the entities contained in the States object.
func (ss *States) Entities() []string {
	ss.m.Lock()
	defer ss.m.Unlock()

	en := make([]string, 0, len(ss.s))
	for _, st := range ss.s {
		en = append(en, st.DomainEntity)
	}
	return en
}

// Get returns a single state record and a bool if found.
func (ss *States) Get(key string) (*State, bool) {
	ss.m.Lock()
	defer ss.m.Unlock()

	st, ok := ss.s[key]
	if !ok {
		return nil, ok
	}
	return st, true
}

// Find returns a new map of states of the passed entities.
func (ss *States) Find(entities []string) map[string]*State {
	ss.m.Lock()
	defer ss.m.Unlock()

	sss := make(map[string]*State)

	for _, k := range entities {
		if st, ok := ss.s[k]; ok {
			sss[st.DomainEntity] = st
		}
	}

	return sss
}

// FindDomainMap returns a map of the states for the passed domain.
func (ss *States) FindDomainMap(keys []string) map[string]*State {
	ss.m.Lock()
	defer ss.m.Unlock()

	sss := make(map[string]*State)

	for _, st := range ss.s {
		for _, k := range keys {
			if st.Domain == k {
				sss[st.DomainEntity] = st
			}
		}
	}

	return sss
}

// Slice returns a slice of the states in no particular order.
func (ss *States) Slice() []*State {
	ss.m.Lock()
	defer ss.m.Unlock()

	var sss []*State
	for _, st := range ss.s {
		sss = append(sss, st)
	}

	return sss
}

// Map returns a map of all the states.
func (ss *States) Map() map[string]*State {
	ss.m.Lock()
	defer ss.m.Unlock()

	sts := make(map[string]*State)

	maps.Copy(sts, ss.s)

	return sts
}

// Iterate calls the provided yield function for each state.
// If yield returns false, iteration stops early.
func (ss *States) Iterate(yield func(key string, st *State) bool) {
	ss.m.Lock()
	defer ss.m.Unlock()

	for k, st := range ss.s {
		if !yield(k, st) {
			return
		}
	}
}

// SubSet returns a new States which contains a subset of the current states based on entities passed.
func (ss *States) SubSet(entities []string) States {
	return States{
		s: ss.Find(entities),
		m: &sync.Mutex{},
	}
}

// Where returns a new States object containing all the states that match the passed state. strings.Equalfold is used
// for the comparison.
func (ss *States) Where(state string) *States {
	sts := States{
		s: make(map[string]*State),
		m: &sync.Mutex{},
	}

	for _, v := range ss.Slice() {
		if strings.EqualFold(string(v.State), state) {
			sts.Upsert(v)
		}
	}
	return &sts
}

func MessageState(message *model.Message) *State {
	//nolint:exhaustruct // LastChanged and LastUpdated not needed for message state conversion
	return &State{
		DomainEntity: message.DomainEntity(),
		Domain:       message.Domain(),
		Entity:       message.EntityID(),
		State:        StateText(message.State()),
		Attributes:   message.Attributes(),
	}
}
