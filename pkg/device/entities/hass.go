package entities

import (
	"github.com/brutella/hap/accessory"
	hassentity "github.com/kjbreil/hass-mqtt/entities"
)

type HassEntity struct {
	entity hassentity.Entity
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
