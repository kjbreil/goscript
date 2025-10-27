package entities

import (
	"strings"

	"github.com/brutella/hap/accessory"
	"github.com/kjbreil/goscript/pkg/device"
	hassentity "github.com/kjbreil/hass-mqtt/entities"
)

type HassEntity struct {
	entity hassentity.Entity
}

func (e *HassEntity) GetName() string {
	return e.entity.GetName()
}

func (e *HassEntity) GetDomain() device.DomainType {
	de := e.entity.GetDomainEntity()
	domain := strings.SplitN(de, ".", 1)[1]
	switch domain {
	default:
		return device.DomainTypeUnknown
	}
}

func NewHassEntity(e hassentity.Entity) *HassEntity {
	return &HassEntity{
		entity: e,
	}
}

func (e *HassEntity) GetHassEntity() hassentity.Entity {
	return e.entity
}

func (e *HassEntity) GetDomainEntity() string {
	return e.entity.GetDomainEntity()
}

func (e *HassEntity) UpdateState() {
	e.entity.UpdateState()
}

func (e *HassEntity) GetHomekitAccessory() *accessory.A {
	return nil
}
