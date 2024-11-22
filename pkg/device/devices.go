package device

import (
	"fmt"
	"github.com/brutella/hap/accessory"
)

type Devices struct {
	devs map[string]*Device
}

func NewDevices() *Devices {
	return &Devices{
		devs: make(map[string]*Device),
	}
}

func (ds *Devices) GetDevice(entity string) (*Device, error) {
	d, ok := ds.devs[entity]
	if !ok {
		return nil, fmt.Errorf("could not find device %s", entity)
	}

	return d, nil
}
func (ds *Devices) AddDevices(devices *Devices) {
	for _, d := range devices.devs {
		ds.devs[d.Dev().GetUniqueId()] = d
	}
}

func (ds *Devices) AddDevice(d *Device) {
	ds.devs[d.Dev().GetUniqueId()] = d
}

func (ds *Devices) GetHomekitAccessories() []*accessory.A {
	var accs []*accessory.A
	for _, d := range ds.devs {
		accs = append(accs, d.GetHomekitAccessories()...)
	}
	return accs
}

func (ds *Devices) Slice() []*Device {
	var devs []*Device
	for _, d := range ds.devs {
		devs = append(devs, d)
	}
	return devs
}
