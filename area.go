package goscript

import "github.com/kjbreil/hass-ws/model"

func (gs *GoScript) GetAreaDomain(area, domain string) []string {
	var results []string

	if a, ok := gs.areaRegistry[area]; ok {
		for _, e := range a {
			if e.Domain() == domain {
				results = append(results, e.DomainEntity())
			}
		}
	}

	return results
}

func (gs *GoScript) fillAreaRegistry() {
	devices := make(map[string]string)
	devs := gs.ws.GetDeviceRegistry()
	for _, dev := range devs.Result {
		if dev.AreaId != nil {
			devices[dev.Id] = *dev.AreaId
		}
	}

	gs.areaRegistry = make(map[string][]model.Result)
	entities := gs.ws.GetEntityRegistry()
	for _, entity := range entities.Result {
		if entity.AreaId != nil {
			gs.areaRegistry[*entity.AreaId] = append(gs.areaRegistry[*entity.AreaId], entity)
			continue
		}
		if entity.DeviceId != nil {
			areaID, ok := devices[*entity.DeviceId]
			if ok {
				gs.areaRegistry[areaID] = append(gs.areaRegistry[areaID], entity)
			}
		}
	}
}
