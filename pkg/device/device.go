package device

import (
	hassdevice "github.com/kjbreil/hass-mqtt/device"
	"github.com/kjbreil/hass-mqtt/entities"
)

type Device interface {
	// Init(m module.Module) error
	Update()
	Dev() *hassdevice.Device
	GetEntities() []entities.Entity
	GetEntity(domainEntity string) entities.Entity
	GetUniqueID() string

	mustImplementGSDevice()
}
