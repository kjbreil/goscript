package device

import (
	"fmt"
	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	hassdevice "github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
	"strings"
)

type Light struct {
	name     string
	dev      *hassdevice.Device
	entities map[string]entities.Entity
	// control  *control.Requests

	GSDevice
}

func NewLight(
	name string,
	commandFunc func(message mqtt.Message, client mqtt.Client),
	brightnessFunc func(message mqtt.Message, client mqtt.Client),
) *Light {
	l := &Light{
		name:     fmt.Sprintf("all_%s_lights", name),
		entities: make(map[string]entities.Entity),
	}

	snakeName := strcase.ToSnake(l.name)
	readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))

	mainDevice := hassdevice.New(readableName, fmt.Sprintf("%s_virtual", snakeName), "Group Lights 2000", "GoScript", "0.0.2")

	l.dev = mainDevice

	lightOptions := entities.NewLightOptions()
	lightOptions.
		Name(readableName).
		EnableBrightness().
		CommandFunc(commandFunc).
		BrightnessCommandFunc(brightnessFunc)

	return l
}
