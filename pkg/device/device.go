package device

import (
	"github.com/brutella/hap/accessory"
	"github.com/iancoleman/strcase"
	hassdevice "github.com/kjbreil/hass-mqtt/device"
	"strings"
)

type Device struct {
	dev          *hassdevice.Device
	entities     map[string]Entity
	name         string
	model        string
	manufacturer string
	swVersion    string
	// control  *control.Requests
}

func New(
	name string,
	model string,
	manufacturer string,
	swVersion string,
	entities ...Entity,
) *Device {
	snakeName := strcase.ToSnake(name)
	readableName := strings.Title(strings.ReplaceAll(snakeName, "_", " "))

	mainDevice := hassdevice.New(readableName, snakeName, model, manufacturer, swVersion)

	d := &Device{
		dev:          mainDevice,
		entities:     make(map[string]Entity),
		name:         name,
		model:        model,
		manufacturer: manufacturer,
		swVersion:    swVersion,
	}

	err := d.AddEntities(entities)
	if err != nil {
		panic(err)
	}

	return d
}

func NewDevice(dev *hassdevice.Device) *Device {
	return &Device{
		dev:      dev,
		entities: make(map[string]Entity),
	}
}

func (d *Device) Dev() *hassdevice.Device {
	return d.dev
}

func (d *Device) AddEntity(e Entity) error {
	if hk := e.GetHomekitAccessory(); hk != nil {
		hk.Info.Model.SetValue(d.model)
		hk.Info.Manufacturer.SetValue(d.manufacturer)
	}
	d.entities[e.GetDomainEntity()] = e
	err := d.dev.Add(e.GetHassEntity())
	if err != nil {
		return err
	}
	return nil
}
func (d *Device) AddEntities(ets []Entity) error {
	var err error
	for _, et := range ets {
		err = d.AddEntity(et)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *Device) GetEntities() []Entity {
	ets := make([]Entity, 0, len(d.entities))
	for _, et := range d.entities {
		ets = append(ets, et)
	}
	return ets
}

func (d *Device) GetEntity(domainEntity string) Entity {
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

func (d *Device) GetHomekitAccessories() []*accessory.A {
	var accs []*accessory.A
	for _, e := range d.entities {
		a := e.GetHomekitAccessory()
		if a != nil {
			accs = append(accs, a)
		}
	}
	return accs
}

func (d *Device) mustImplementGSDevice() {

}
