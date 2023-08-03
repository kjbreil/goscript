package light

import (
	"fmt"
	"github.com/kjbreil/goscript"
	"github.com/kjbreil/hass-ws/services"
)

type Light struct {
	turnOnParams  services.LightTurnOnParams
	turnOffParams services.LightTurnOffParams
}

func New() *Light {
	return &Light{}
}

func (l *Light) Flash(flash services.Flash) *Light {
	l.turnOnParams.Flash = &flash
	l.turnOffParams.Flash = &flash
	return l
}

func (l *Light) Transition(transition float64) *Light {
	l.turnOnParams.Transition = &transition
	l.turnOffParams.Transition = &transition
	return l
}

func (l *Light) repeatService(state string, fn func(t *goscript.Task, entities []string), t *goscript.Task, entities []string, lightService services.Service) {
	if t.Cancelled() {
		panic(fmt.Sprintf("task context cancelled for %s", t.UUID()))
	}

	t.ServiceChan <- lightService
	//
	//wait := 1500 * time.Millisecond
	//
	////if l.turnOnParams.Transition != nil {
	////	wait = time.Duration(*l.turnOnParams.Transition*1000)*time.Millisecond + (100 * time.Millisecond)
	////}
	//
	//timer := time.NewTimer(wait)
	//select {
	//case <-timer.C:
	//	break
	//case <-t.Context().Done():
	//	panic(fmt.Sprintf("task context cancelled for %s", t.UUID()))
	//}
	//
	//var repeatEntities []string
	//for _, e := range entities {
	//	if s, ok := t.States.Get(e); ok {
	//		if !s.State.Equals(state) {
	//			repeatEntities = append(repeatEntities, s.DomainEntity)
	//		}
	//	}
	//}
	//
	//if len(repeatEntities) != 0 {
	//	fn(t, repeatEntities)
	//}
}
