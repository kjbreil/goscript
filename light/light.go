package light

import (
	"fmt"

	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

// Light provides methods for controlling light entities with turn on/off parameters.
type Light struct {
	turnOnParams  services.LightTurnOnParams
	turnOffParams services.LightTurnOffParams
}

// New creates a new Light instance with empty turn on/off parameters.
func New() *Light {
	//nolint:exhaustruct // turnOnParams and turnOffParams are initialized as needed by method calls
	return &Light{}
}

// Transition sets the transition time for turn on/off operations.
func (l *Light) Transition(transition float64) *Light {
	l.turnOnParams.Transition = &transition
	l.turnOffParams.Transition = &transition
	return l
}

func (l *Light) repeatService(
	_ string,
	_ func(t *trigger.Task, entities []string),
	t *trigger.Task,
	_ []string,
	lightService services.Service,
) {
	if t.Cancelled() {
		panic(fmt.Sprintf("task context cancelled for %s", t.UUID()))
	}

	t.ServiceChan <- lightService
	//
	// wait := 1500 * time.Millisecond
	//
	// //if l.turnOnParams.Transition != nil {
	// //	wait = time.Duration(*l.turnOnParams.Transition*1000)*time.Millisecond + (100 * time.Millisecond)
	// //}
	//
	// timer := time.NewTimer(wait)
	// select {
	// case <-timer.C:
	//	break
	// case <-t.Context().Done():
	//	panic(fmt.Sprintf("task context cancelled for %s", t.UUID()))
	// }
	//
	// var repeatEntities []string
	// for _, e := range entities {
	//	if s, ok := t.States.Get(e); ok {
	//		if !s.State.Equals(state) {
	//			repeatEntities = append(repeatEntities, s.DomainEntity)
	//		}
	//	}
	// }
	//
	// if len(repeatEntities) != 0 {
	//	fn(t, repeatEntities)
	// }
}
