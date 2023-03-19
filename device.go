package goscript

import (
	"github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
)

type Device struct {
	dev      *device.Device
	entities []entities.Entity
}

func (gs *GoScript) AddDevice(dev *device.Device) (error, *Device) {
	d := &Device{
		dev: dev,
	}

	err := gs.mqtt.Add(d.dev)
	if err != nil {
		return err, nil
	}

	gs.devices[d.dev.GetUniqueId()] = d

	return nil, d
}

func (d *Device) AddEntities(ets []entities.Entity) error {
	d.entities = ets
	for _, et := range ets {
		err := d.dev.Add(et)

		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Device) GetUniqueID() string {
	return d.dev.GetUniqueId()
}
func (d *Device) Update() {
	for _, e := range d.entities {
		e.UpdateState()
	}
}
