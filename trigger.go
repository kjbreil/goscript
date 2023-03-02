package goscript

import "github.com/kjbreil/hass-ws/model"

type Trigger struct {
	EntityTriggers []string
	EntityStates   []string
	entityTriggers map[string]struct{}
	entityStates   map[string]struct{}
	Func           func(message model.Message, states []State)
}

func (gs *GoScript) AddTrigger(t *Trigger) {
	t.entityTriggers = make(map[string]struct{})
	t.entityStates = make(map[string]struct{})

	for _, et := range t.EntityTriggers {
		t.entityTriggers[et] = struct{}{}
		t.entityStates[et] = struct{}{}
	}
	for _, es := range t.EntityStates {
		t.entityStates[es] = struct{}{}
	}
	t.EntityTriggers = nil
	t.entityStates = nil

	for et := range t.entityTriggers {
		gs.triggers[et] = append(gs.triggers[et], t)
	}
}

func (gs *GoScript) runTriggers(message model.Message) {
	if tr, ok := gs.triggers[message.DomainEntity()]; ok {
		for _, t := range tr {
			t.Func(message, nil)
		}
	}
}
