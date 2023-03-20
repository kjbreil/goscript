package goscript

import (
	"context"
	"github.com/google/uuid"
	"github.com/kjbreil/hass-ws/model"
)

// Trigger takes in trigger items, domains or a schedule and runs a function based on any variation of the inputs.
//
// If Unique is not nil then if the trigger function is already running the context of that trigger function will be
// killed. To accomplish this it is important to only use the task methods within a trigger function instead of
// time.Sleep or anything like that as then the function will not be killed and both will run at same time. If Unique is
// nil then multiple functions can run at the same time.
//
// Triggers can be on Entity's (full domain.entity format), Domains or on a Periodic schedule. Periodics do not get run
// through Eval's but it is best to handle all evaluation within the function for Periodics mixed with other trigger
// types to ensure consistent results.
//
// States is a list of entities to which the state will be available within the task function. All triggers are
// automatically included in the list. DomainStates allows you to specify a whole domain to be included in the states.
//
// Evaluation is done through a list of strings that are run through github.com/antonmedv/expr to evaluate the output.
// Like with PyScript type is important in the evaluation scripts. Check out github.com/antonmedv/expr for more details
// on casting and converting. You cannot mix types in a single evaluation so `state == "on" || state > 10` will always
// return false due to failure parsing the evaluation. Attributes are available inside the evaluations so
// `color_temp > 100` will work as long as color_temp exists in the attributes of the entity and the data type is a float
//
// Func is the function to run when the criteria are met. Within the trigger function a *Task is available to give
// information on the trigger. Killing the triggerfunc panics to exit. The runner recovers this panic, this also means
// that if your code panics the whole program will not crash but will continue. Panic will be written to the logs.
type Trigger struct {
	uuid uuid.UUID

	Unique        *Unique
	Triggers      []string
	DomainTrigger []string // DomainTrigger, triggers of everything in the domain, also attaches all states for the domain
	Periodic
	States       []string
	DomainStates []string
	Eval         []string
	Func         TriggerFunc
}

// Unique makes the trigger unique, KillMe is a placeholder for now and does nothing.
type Unique struct {
	KillMe bool // even when false is true
	ctx    context.Context
	cancel context.CancelFunc
}

// TriggerFunc is the function to run when the criteria are met. Within the trigger function a *Task is available.
// See Task for more information on what is available in Task.
type TriggerFunc func(t *Task)

// Entities is a simple helper function to create a []string. Will most likely be removed in the future.
func Entities(entities ...string) []string {
	return entities
}

// AddTrigger adds a trigger to the trigger map. There is no validation of a trigger.
func (gs *GoScript) AddTrigger(t *Trigger) {
	// set up the trigger object
	t.uuid = uuid.New()
	if t.Unique != nil {
		t.Unique.ctx, t.Unique.cancel = context.WithCancel(context.Background())
	}
	entityTriggers := make(map[string]struct{})
	domainTriggers := make(map[string]struct{})
	entityStates := make(map[string]struct{})
	domainStates := make(map[string]struct{})
	for _, et := range t.Triggers {
		entityTriggers[et] = struct{}{}
		entityStates[et] = struct{}{}
	}
	for _, es := range t.States {
		entityStates[es] = struct{}{}
	}
	for _, ed := range t.DomainTrigger {
		domainTriggers[ed] = struct{}{}
	}
	for _, eds := range t.DomainStates {
		domainStates[eds] = struct{}{}
	}

	t.Triggers = nil
	t.States = nil
	t.DomainTrigger = nil
	t.DomainStates = nil

	t.Triggers = make([]string, len(entityTriggers))
	i := 0
	for k := range entityTriggers {
		t.Triggers[i] = k
		i++
	}

	t.DomainTrigger = make([]string, len(domainTriggers))
	i = 0
	for k := range domainTriggers {
		t.DomainTrigger[i] = k
		i++
	}

	t.States = make([]string, len(entityStates))
	i = 0
	for k := range entityStates {
		t.States[i] = k
		i++
	}

	t.DomainStates = make([]string, len(domainStates))
	i = 0
	for k := range domainStates {
		t.DomainStates[i] = k
		i++
	}

	// for each entity add to the triggers map
	for _, et := range t.Triggers {
		gs.triggers[et] = append(gs.triggers[et], t)
	}

	// for each domain add to the domain trigger map
	for _, edt := range t.DomainTrigger {
		gs.domainTrigger[edt] = append(gs.domainTrigger[edt], t)
	}

	// for each periodic add to the periodic map
	// cron time is an array of triggers so multiple triggers can have same cron schedule
	for _, ep := range t.Periodic {
		gs.periodic[ep] = append(gs.periodic[ep], t)
	}
}

// RemoveTrigger can be used to remove a trigger while program is running.
func (gs *GoScript) RemoveTrigger(t *Trigger) {
	for _, et := range t.Triggers {
		for i, te := range gs.triggers[et] {
			if te.uuid == t.uuid {
				gs.triggers[et] = append(gs.triggers[et][:i], gs.triggers[et][i+1:]...)
				break
			}
		}
	}
}

// AddTriggers helper function to add multiple triggers
func (gs *GoScript) AddTriggers(triggers ...*Trigger) {
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

func (gs *GoScript) triggerDomainEntity(message *model.Message, trigger *Trigger) {
	passed := trigger.eval(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.funcToRun[task.uuid] = task
	}
}
func (gs *GoScript) triggerDomain(message *model.Message, trigger *Trigger) {
	passed := trigger.eval(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.funcToRun[task.uuid] = task
	}
}
