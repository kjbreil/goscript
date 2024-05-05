package core

import (
	"github.com/kjbreil/goscript/pkg/device"
)

func (gs *GoScript) AddDevice(dev *device.GSDevice) error {

	err := gs.mqtt.Add(dev.Dev())
	if err != nil {
		return err
	}

	gs.devices[dev.Dev().GetUniqueId()] = dev

	return nil
}
