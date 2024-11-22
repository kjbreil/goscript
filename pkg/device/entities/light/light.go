package light

import (
	"github.com/brutella/hap/accessory"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/trigger"
	hassentity "github.com/kjbreil/hass-mqtt/entities"
	"strconv"
	"strings"
)

type Light struct {
	name             string
	hass             *hassentity.Light
	hassOptions      *hassentity.LightOptions
	homekitAccessory *accessory.ColoredLightbulb
}

func New(name string, options ...func(*Light)) *Light {
	snakeName := strcase.ToSnake(name)
	readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))

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

func (l *Light) GetHassEntity() hassentity.Entity {
	if l.hass != nil {
		return l.hass
	}
	return nil
}

func (l *Light) GetDomainEntity() string {
	return l.hass.GetDomainEntity()
}

func (l *Light) UpdateState() {
	l.hass.UpdateState()
}

func (l *Light) GetHomekitAccessory() *accessory.A {
	if l.homekitAccessory != nil {
		return l.homekitAccessory.A
	}
	return nil
}

func WithHomeKit() func(*Light) {
	return func(l *Light) {
		l.homekitAccessory = accessory.NewColoredLightbulb(accessory.Info{
			Name: l.name,
		})

		l.hassOptions.CommandFunc(func(message mqtt.Message, client mqtt.Client) {
			if string(message.Payload()) == "ON" {
				l.homekitAccessory.Lightbulb.On.SetValue(true)
			} else {
				l.homekitAccessory.Lightbulb.On.SetValue(false)
			}
		})

		l.hassOptions.EnableBrightness().BrightnessCommandFunc(func(message mqtt.Message, client mqtt.Client) {
			brightness, err := strconv.ParseFloat(string(message.Payload()), 64)
			if err != nil {
				return
			}

			brightnessPct := int((brightness / 255) * 100)

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

func WithBrightnessCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		if l.homekitAccessory != nil {
			fnPass := tr.Func
			tr.Func = func(t *trigger.Task) {
				brightness, err := strconv.ParseFloat(string(t.MqttMessage.Payload()), 64)
				if err != nil {
					return
				}

				brightnessPct := int((brightness / 255) * 100)

				err = l.homekitAccessory.Lightbulb.Brightness.SetValue(brightnessPct)
				if err != nil {
					return
				}
				fnPass(t)
			}
		}

		l.homekitAccessory.Lightbulb.Brightness.OnValueRemoteUpdate(func(v int) {

			brightness := float64(v) / 100 * 255
			l.hass.Brightness(strconv.Itoa(int(brightness)))
		})

		l.hassOptions.EnableBrightness().BrightnessCommandFunc(module.TaskMQTT(m, tr))
	}
}

func WithRgbCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		l.hassOptions.EnableRgb().RgbCommandFunc(module.TaskMQTT(m, tr))
	}
}

func WithRgbwCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		l.hassOptions.EnableRgbw().RgbwwCommandFunc(module.TaskMQTT(m, tr))
	}
}

func WithColorTempCommandFunc(m module.Module, tr *trigger.Trigger) func(*Light) {
	return func(l *Light) {
		l.hassOptions.EnableColorTemp().ColorTempCommandFunc(module.TaskMQTT(m, tr))
	}
}

func WithStateFunc(fn func() string) func(*Light) {
	return func(l *Light) {
		l.hassOptions.StateFunc(fn)
	}
}

func WithJsonAttributes(fn func() string) func(*Light) {
	return func(l *Light) {
		l.hassOptions.JsonAttributesFunc(fn)
	}
}
