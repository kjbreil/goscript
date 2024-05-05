package virtual

import (
	"fmt"
	"github.com/iancoleman/strcase"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/trigger"
	hassdevice "github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
	"strings"
	"time"
)

var key = "virtual"

type Virtual struct {
	Lights        map[string]lights
	BinarySensors map[string]binarySensor
	module.Base
}

func (v *Virtual) Run() error {
	devices := make(device.Devices)
	for n := range v.Lights {
		snakeName := strcase.ToSnake(n)
		readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))
		dev := hassdevice.New(readableName, fmt.Sprintf("%s_virtual", snakeName), "", "Kaygel", "0.0.1")

		lightOptions := entities.NewLightOptions().Name(readableName)
		light, _ := entities.NewLight(lightOptions)
		dev.Add(light)
		gd := device.NewGSDevice(dev)
		devices[dev.GetUniqueId()] = gd

	}
	for n := range v.BinarySensors {
		snakeName := strcase.ToSnake(n)
		readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))
		dev := hassdevice.New(readableName, fmt.Sprintf("%s_virtual", snakeName), "", "Kaygel", "0.0.1")

		binarySensorOptions := entities.NewBinarySensorOptions().Name(readableName)
		binarySensorOptions.DeviceClass("OCCUPANCY")
		bs, _ := entities.NewBinarySensor(binarySensorOptions)

		dev.Add(bs)
		gd := device.NewGSDevice(dev)
		go func() {
			for {
				time.Sleep(10 * time.Second)
				bs.State("ON")
				time.Sleep(1 * time.Second)
				bs.State("OFF")
			}
		}()
		devices[dev.GetUniqueId()] = gd
	}
	module.SendDevices(v, devices)

	return nil
}

type lights struct{}
type binarySensor struct {
}

func (v *Virtual) Triggers() trigger.Triggers {
	return nil
}

func (v *Virtual) Update() error {
	return nil
}

func (v *Virtual) Name() string {
	return key
}

func (v *Virtual) Close() error {
	return nil
}
