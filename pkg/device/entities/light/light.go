package light

import (
	"strconv"
	"strings"

	"github.com/brutella/hap/accessory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/trigger"
	hassentity "github.com/kjbreil/hass-mqtt/entities"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	maxBrightnessValue = 255
	percentageDivisor  = 100
)

// Light represents a light device entity with HASS and HomeKit integration.
type Light struct {
	name             string
	hass             *hassentity.Light
	hassOptions      *hassentity.LightOptions
	homekitAccessory *accessory.ColoredLightbulb
}

// New creates a new Light entity with the given name and options.
func New(name string, options ...func(*Light)) *Light {
	snakeName := strcase.ToSnake(name)
	readableName := cases.Title(language.English).String(strings.ReplaceAll(snakeName, "_", " "))

	//nolint:exhaustruct // hass, hassOptions, homekitAccessory initialized below
	l := &Light{
		name: readableName,
	}

	l.hassOptions = hassentity.NewLightOptions().Name(readableName)

	// TODO: Force Homekit process first
	for _, o := range options {
		o(l)
	}

	var err error
	l.hass, err = hassentity.NewLight(l.hassOptions)
	if err != nil {
		return nil
	}

	return l
}

// GetHassEntity returns the HASS entity for this light.
func (l *Light) GetHassEntity() hassentity.Entity {
	if l.hass != nil {
		return l.hass
	}
	return nil
}

// GetDomainEntity returns the domain entity string for this light.
func (l *Light) GetDomainEntity() string {
	return l.hass.GetDomainEntity()
}

// UpdateState updates the state of the light entity.
func (l *Light) UpdateState() {
	l.hass.UpdateState()
}

// GetHomekitAccessory returns the HomeKit accessory for this light.
func (l *Light) GetHomekitAccessory() *accessory.A {
	if l.homekitAccessory != nil {
		return l.homekitAccessory.A
	}
	return nil
}

// WithHomeKit returns an option function that adds HomeKit support to a light.
func WithHomeKit() func(*Light) {
	return func(l *Light) {
		//nolint:exhaustruct // Only Name required; other Info fields optional
		l.homekitAccessory = accessory.NewColoredLightbulb(accessory.Info{
			Name: l.name,
		})

		l.hassOptions.CommandFunc(func(message mqtt.Message, _ mqtt.Client) {
			if string(message.Payload()) == "ON" {
				l.homekitAccessory.Lightbulb.On.SetValue(true)
			} else {
				l.homekitAccessory.Lightbulb.On.SetValue(false)
			}
		})

		l.hassOptions.EnableBrightness().BrightnessCommandFunc(func(message mqtt.Message, _ mqtt.Client) {
			brightness, err := strconv.ParseFloat(string(message.Payload()), 64)
			if err != nil {
				return
			}

			brightnessPct := int((brightness / maxBrightnessValue) * percentageDivisor)

			err = l.homekitAccessory.Lightbulb.Brightness.SetValue(brightnessPct)
			if err != nil {
				return
			}
		})

		l.homekitAccessory.Lightbulb.On.OnValueRemoteUpdate(func(v bool) {
			if v {
				l.hass.State("ON")
			} else {
				l.hass.State("OFF")
			}
		})
	}
}

// WithCommandFunc returns an option function that sets a command callback for the light.
func WithCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		if l.homekitAccessory != nil {
			fnPass := tr.Func
			tr.Func = func(t *trigger.Task) {
				if string(t.MqttMessage.Payload()) == "ON" {
					l.homekitAccessory.Lightbulb.On.SetValue(true)
				} else {
					l.homekitAccessory.Lightbulb.On.SetValue(false)
				}
				fnPass(t)
			}
		}

		l.hassOptions.CommandFunc(module.TaskMQTT(m, tr))
	}
}

// WithBrightnessCommandFunc returns an option function that sets a brightness command callback.
func WithBrightnessCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		if l.homekitAccessory != nil {
			fnPass := tr.Func
			tr.Func = func(t *trigger.Task) {
				brightness, err := strconv.ParseFloat(string(t.MqttMessage.Payload()), 64)
				if err != nil {
					return
				}

				brightnessPct := int((brightness / maxBrightnessValue) * percentageDivisor)

				err = l.homekitAccessory.Lightbulb.Brightness.SetValue(brightnessPct)
				if err != nil {
					return
				}
				fnPass(t)
			}
		}

		l.homekitAccessory.Lightbulb.Brightness.OnValueRemoteUpdate(func(v int) {
			brightness := float64(v) / percentageDivisor * maxBrightnessValue
			l.hass.Brightness(strconv.Itoa(int(brightness)))
		})

		l.hassOptions.EnableBrightness().BrightnessCommandFunc(module.TaskMQTT(m, tr))
	}
}

// WithRgbCommandFunc returns an option function that sets an RGB command callback.
func WithRgbCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		l.hassOptions.EnableRgb().RgbCommandFunc(module.TaskMQTT(m, tr))
	}
}

// WithRgbwCommandFunc returns an option function that sets an RGBW command callback.
func WithRgbwCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		l.hassOptions.EnableRgbw().RgbwwCommandFunc(module.TaskMQTT(m, tr))
	}
}

// WithColorTempCommandFunc returns an option function that sets a color temperature command callback.
func WithColorTempCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		l.hassOptions.EnableColorTemp().ColorTempCommandFunc(module.TaskMQTT(m, tr))
	}
}

// WithStateFunc returns an option function that sets a state callback.
func WithStateFunc(fn func() string) func(*Light) {
	return func(l *Light) {
		l.hassOptions.StateFunc(fn)
	}
}

// WithJSONAttributes returns an option function that sets a JSON attributes callback.
func WithJSONAttributes(fn func() string) func(*Light) {
	return func(l *Light) {
		l.hassOptions.JsonAttributesFunc(fn)
	}
}
