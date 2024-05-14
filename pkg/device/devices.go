package device

import (
	"fmt"
	"github.com/brutella/hap/accessory"
)

type Devices map[string]*Device

func (ds Devices) GetDevice(entity string) (*Device, error) {
	d, ok := ds[entity]
	if !ok {
		return nil, fmt.Errorf("could not find device %s", entity)
	}

	return d, nil
}
func (ds Devices) AddDevices(devices Devices) {
	for _, d := range devices {
		ds[d.Dev().GetUniqueId()] = d
	}
}

func (ds Devices) AddDevice(d *Device) {
	ds[d.Dev().GetUniqueId()] = d
}

func (ds Devices) GetHomekitAccessories() []*accessory.A {
	var accs []*accessory.A
	for _, d := range ds {
		accs = append(accs, d.GetHomekitAccessories()...)
	}
	return accs
}
