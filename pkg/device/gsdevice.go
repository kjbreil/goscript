package device

import (
	hassdevice "github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
)

type GSDevice struct {
	dev      *hassdevice.Device
	entities map[string]entities.Entity
	// control  *control.Requests
}

func NewGSDevice(dev *hassdevice.Device) *GSDevice {
	return &GSDevice{
		dev:      dev,
		entities: make(map[string]entities.Entity),
	}
}

// func (d *GSDevice) Init(m module.Module) error {
// 	// d.control = m.Requests()
//
// 	return nil
// }

func (d *GSDevice) Dev() *hassdevice.Device {
	return d.dev
}
func (d *GSDevice) AddEntities(ets []entities.Entity) error {
	for _, et := range ets {
		d.entities[et.GetDomainEntity()] = et
		err := d.dev.Add(et)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *GSDevice) GetEntities() []entities.Entity {
	ets := make([]entities.Entity, 0, len(d.entities))
	for _, et := range d.entities {
		ets = append(ets, et)
	}
	return ets
}

func (d *GSDevice) GetEntity(domainEntity string) entities.Entity {
	e, ok := d.entities[domainEntity]
	if !ok {
		return nil
	}
	return e
}
func (d *GSDevice) GetUniqueID() string {
	return d.dev.GetUniqueId()
}
func (d *GSDevice) Update() {
	for _, e := range d.entities {
		e.UpdateState()
	}
}

func (d *GSDevice) mustImplementGSDevice() {

}
