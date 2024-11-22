package device

import (
	"github.com/brutella/hap/accessory"
	"github.com/kjbreil/hass-mqtt/entities"
)

type Entity interface {
	GetName() string
	GetDomain() DomainType
	GetHassEntity() entities.Entity
	GetDomainEntity() string
	UpdateState()
	GetHomekitAccessory() *accessory.A
}
