package device

import "fmt"

type Devices map[string]*GSDevice

func (ds Devices) GetDevice(entity string) (*GSDevice, error) {
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

func (ds Devices) AddDevice(d *GSDevice) {
	ds[d.Dev().GetUniqueId()] = d
}
