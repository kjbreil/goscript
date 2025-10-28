package device

import (
	"strings"

	mqtt "github.com/eclipse/paho.mqtt.golang"
	"github.com/iancoleman/strcase"
	"github.com/kjbreil/hass-mqtt/entities"
	hassdevice "github.com/kjbreil/hass-mqtt/pkg/device"
)

type Light struct {
	name     string
	dev      *hassdevice.Device
	entities map[string]entities.Entity
	// control  *control.Requests

	Device
}

func NewLight(
	name string,
	commandFunc func(message mqtt.Message, client mqtt.Client),
	brightnessFunc func(message mqtt.Message, client mqtt.Client),
) (*Light, error) {
	l := &Light{
		name:     name,
		entities: make(map[string]entities.Entity),
	}

	snakeName := strcase.ToSnake(l.name)
	readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))

	mainDevice := hassdevice.New(readableName, snakeName, "Group Lights 2000", "GoScript", "0.0.2")

	l.dev = mainDevice

	lightOptions := entities.NewLightOptions()
	lightOptions.
		Name(readableName).
		EnableBrightness().
		CommandFunc(commandFunc).
		BrightnessCommandFunc(brightnessFunc)

	return l, nil
}
