package light

import (
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/services"
)

// TurnOn turns on the passed entities with the setting that are pre-provided. It is passed a task as well so that the
// context within the task is respected. TurnOn finds any passed lights that did not turn on and attempts to turn on
// again, this will keep happening until it turns on, however if a light is unavailable it will not be added to the list
// so as long as entities present a proper unavailable it will not continue forever but could in a bad state. Best to
// use only with Unique tasks, so it would be killed with the next task run.
func (l *Light) TurnOn(t *trigger.Task, entities []string) {
	lightService := services.NewLightTurnOn(services.Targets(entities...))
	lightService.ServiceData = l.turnOnParams

	l.repeatService("on", l.TurnOn, t, entities, lightService)
}

// Brightness sets the brightness parameter for the light.
func (l *Light) Brightness(brightness float64) *Light {
	l.turnOnParams.BrightnessPct = &brightness
	return l
}

// BrightnessPct sets the brightness percentage parameter for the light.
func (l *Light) BrightnessPct(brightnessPct float64) *Light {
	l.turnOnParams.BrightnessPct = &brightnessPct
	return l
}

// BrightnessStep sets the brightness step parameter for the light.
func (l *Light) BrightnessStep(brightnessStep float64) *Light {
	l.turnOnParams.BrightnessStepPct = &brightnessStep
	return l
}

// BrightnessStepPct sets the brightness step percentage parameter for the light.
func (l *Light) BrightnessStepPct(brightnessStepPct float64) *Light {
	l.turnOnParams.BrightnessStepPct = &brightnessStepPct
	return l
}

// Effect sets the light effect parameter.
func (l *Light) Effect(effect string) *Light {
	l.turnOnParams.Effect = &effect
	return l
}

// Kelvin sets the color temperature in Kelvin for the light.
func (l *Light) Kelvin(kelvin float64) *Light {
	l.turnOnParams.ColorTempKelvin = &kelvin
	return l
}
