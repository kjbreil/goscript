package core

import (
	"github.com/kjbreil/goscript/pkg/trigger"
	"github.com/kjbreil/hass-ws/model"
)

// AddTrigger adds a trigger to the trigger map. There is no validation of a trigger.
func (gs *GoScript) AddTrigger(tr *trigger.Trigger) {
	tr = trigger.SetupTrigger(tr)
	// for each entity add to the triggers map
	for _, et := range tr.Triggers {
		gs.triggers[et] = append(gs.triggers[et], tr)
	}

	// for each domain add to the domain trigger map
	for _, edt := range tr.DomainTrigger {
		gs.domainTrigger[edt] = append(gs.domainTrigger[edt], tr)
	}

	// for each periodic add to the periodic map
	// cron time is an array of triggers so multiple triggers can have same cron schedule
	for _, ep := range tr.Periodic {
		gs.periodic[ep] = append(gs.periodic[ep], tr)
	}

	for _, es := range tr.Services {
		gs.serviceTriggers[es] = append(gs.serviceTriggers[es], tr)
	}
}

// RemoveTrigger can be used to remove a trigger while program is running.
func (gs *GoScript) RemoveTrigger(t *trigger.Trigger) {
	for _, et := range t.Triggers {
		for i, te := range gs.triggers[et] {
			if te.UUID() == t.UUID() {
				gs.triggers[et] = append(gs.triggers[et][:i], gs.triggers[et][i+1:]...)
				break
			}
		}
	}
}

// AddTriggers helper function to add multiple triggers
func (gs *GoScript) AddTriggers(triggers ...*trigger.Trigger) {
	for _, t := range triggers {
		gs.AddTrigger(t)
	}
}

func (gs *GoScript) runTriggers(message model.Message) {
	if tr, ok := gs.triggers[message.DomainEntity()]; ok {
		for _, trigger := range tr {
			gs.triggerDomainEntity(&message, trigger)
		}
	}

	if tr, ok := gs.domainTrigger[message.Domain()]; ok {
		for _, trigger := range tr {
			gs.triggerDomain(&message, trigger)
		}
	}
}

func (gs *GoScript) triggerDomainEntity(message *model.Message, trigger *trigger.Trigger) {
	passed := trigger.Evaluate(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.taskToRun.Add(task)
	}
}
func (gs *GoScript) triggerDomain(message *model.Message, trigger *trigger.Trigger) {
	passed := trigger.Evaluate(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.taskToRun.Add(task)
	}
}
func (gs *GoScript) runServiceTriggers(message model.Message) {
	if message.Event != nil && message.Event.Data != nil && message.Event.Data.ServiceData != nil {
		for _, entity := range message.Event.Data.ServiceData.EntityId {
			if tr, ok := gs.serviceTriggers[entity]; ok {
				for _, trigger := range tr {
					gs.triggerService(&message, trigger)
				}
			}
		}
	}

}

func (gs *GoScript) triggerService(message *model.Message, trigger *trigger.Trigger) {
	task := gs.newTask(trigger, message)
	gs.taskToRun.Add(task)
}
