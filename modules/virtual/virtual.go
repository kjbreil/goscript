// Package virtual provides a virtual module for creating virtual light and binary sensor devices.
package virtual

import (
	"fmt"
	"strings"
	"time"

	"github.com/iancoleman/strcase"
	"github.com/kjbreil/goscript/pkg/device"
	"github.com/kjbreil/goscript/pkg/module"
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-mqtt/entities"
	hassdevice "github.com/kjbreil/hass-mqtt/pkg/device"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

const (
	binarySensorToggleDelay = 10 * time.Second
)

var key = "virtual"

// Virtual is a module that creates virtual devices for testing and development.
type Virtual struct {
	Lights        map[string]lights
	BinarySensors map[string]binarySensor
	module.Base
}

// Run initializes and runs the virtual module, creating all configured devices.
func (v *Virtual) Run() error {
	devices := device.NewDevices()
	caser := cases.Title(language.English)
	for n := range v.Lights {
		snakeName := strcase.ToSnake(n)
		readableName := caser.String(strings.ReplaceAll(snakeName, "_", " "))
		dev := hassdevice.New(readableName, fmt.Sprintf("%s_virtual", snakeName), "", "Kaygel", "0.0.1")

		lightOptions := entities.NewLightOptions().Name(readableName)
		light, err := entities.NewLight(lightOptions)
		if err != nil {
			return fmt.Errorf("failed to create light for %s: %w", n, err)
		}
		if err = dev.Add(light); err != nil {
			return fmt.Errorf("failed to add light to device for %s: %w", n, err)
		}
		gd := device.NewDevice(dev)
		devices.AddDevice(gd)
	}
	for n := range v.BinarySensors {
		snakeName := strcase.ToSnake(n)
		readableName := caser.String(strings.ReplaceAll(snakeName, "_", " "))
		dev := hassdevice.New(readableName, fmt.Sprintf("%s_virtual", snakeName), "", "Kaygel", "0.0.1")

		binarySensorOptions := entities.NewBinarySensorOptions().Name(readableName)
		binarySensorOptions.DeviceClass("OCCUPANCY")
		bs, err := entities.NewBinarySensor(binarySensorOptions)
		if err != nil {
			return fmt.Errorf("failed to create binary sensor for %s: %w", n, err)
		}

		if err = dev.Add(bs); err != nil {
			return fmt.Errorf("failed to add binary sensor to device for %s: %w", n, err)
		}
		gd := device.NewDevice(dev)
		go func() {
			for {
				time.Sleep(binarySensorToggleDelay)
				bs.State("ON")
				time.Sleep(1 * time.Second)
				bs.State("OFF")
			}
		}()
		devices.AddDevice(gd)
	}
	module.SendDevices(v, devices)

	return nil
}

type lights struct{}
type binarySensor struct {
}

// Triggers returns the triggers for this module.
func (v *Virtual) Triggers() trigger.Triggers {
	return nil
}

// Update updates the module state.
func (v *Virtual) Update() error {
	return nil
}

// Name returns the module name.
func (v *Virtual) Name() string {
	return key
}

// Close closes the module and cleans up resources.
func (v *Virtual) Close() error {
	return nil
}
