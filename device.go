package goscript

import (
	"fmt"
	"github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
)

type Device struct {
	dev      *device.Device
	entities map[string]entities.Entity
}

func (gs *GoScript) AddDevice(dev *device.Device) (*Device, error) {
	d := &Device{
		dev:      dev,
		entities: make(map[string]entities.Entity),
	}

	err := gs.mqtt.Add(d.dev)
	if err != nil {
		return nil, err
	}

	gs.devices[d.dev.GetUniqueId()] = d

	return d, nil
}

func (gs *GoScript) GetDevice(entity string) (*Device, error) {
	d, ok := gs.devices[entity]
	if !ok {
		return nil, fmt.Errorf("could not find device %s", entity)
	}

	return d, nil
}

func (d *Device) AddEntities(ets []entities.Entity) error {
	for _, et := range ets {
		d.entities[et.GetDomainEntity()] = et
		err := d.dev.Add(et)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Device) GetEntities() []entities.Entity {
	ets := make([]entities.Entity, 0, len(d.entities))
	for _, et := range d.entities {
		ets = append(ets, et)
	}
	return ets
}

func (d *Device) GetEntity(domainEntity string) entities.Entity {
	e, ok := d.entities[domainEntity]
	if !ok {
		return nil
	}
	return e
}
func (d *Device) GetUniqueID() string {
	return d.dev.GetUniqueId()
}
func (d *Device) Update() {
	for _, e := range d.entities {
		e.UpdateState()
	}
}
