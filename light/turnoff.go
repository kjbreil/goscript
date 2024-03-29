package light

import (
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

// TurnOff turns off the passed entities with the setting that are pre-provided. It is passed a task as well so that the
// context within the task is respected. TurnOff finds any passed lights that did not turn off and attempts to turn off
// again, this will keep happening until it turns off, however if a light is unavailable it will not be added to the list
// so as long as entities present a proper unavailable it will not continue forever but could get in a bad state. Best to
// use only with Unique tasks, so it would be killed with the next task run
func (l *Light) TurnOff(t *trigger.Task, entities []string) {
	state := "off"
	// oldState := "on"
	lightService := services.NewLightTurnOff(services.Targets(entities...))
	lightService.ServiceData = l.turnOffParams

	l.repeatService(state, l.TurnOff, t, entities, lightService)
}
