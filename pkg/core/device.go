package core

import (
	"fmt"
	"github.com/kjbreil/goscript/pkg/device"
	hassdevice "github.com/kjbreil/hass-mqtt/device"
)

func (gs *GoScript) AddDevice(dev *hassdevice.Device) (*device.Device, error) {
	d := device.NewDevice(dev)

	err := gs.mqtt.Add(dev)
	if err != nil {
		return nil, err
	}

	gs.devices[dev.GetUniqueId()] = d

	return d, nil
}

func (gs *GoScript) GetDevice(entity string) (*device.Device, error) {
	d, ok := gs.devices[entity]
	if !ok {
		return nil, fmt.Errorf("could not find device %s", entity)
	}

	return d, nil
}
