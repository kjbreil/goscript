package device

import (
	"github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
)

type Device struct {
	dev      *device.Device
	entities map[string]entities.Entity
}

func NewDevice(dev *device.Device) *Device {
	return &Device{
		dev:      dev,
		entities: make(map[string]entities.Entity),
	}
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
