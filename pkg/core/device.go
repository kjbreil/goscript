package core

import (
	"fmt"

	"github.com/kjbreil/goscript/pkg/device"
)

func (gs *GoScript) AddDevice(dev *device.Device) error {
	err := gs.mqtt.Add(dev.Dev())
	if err != nil {
		return fmt.Errorf("failed to add device %q to MQTT: %w", dev.GetUniqueID(), err)
	}

	gs.devices.AddDevice(dev)

	return nil
}
