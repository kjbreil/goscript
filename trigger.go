package goscript

import (
	"context"
	"github.com/google/uuid"
	"github.com/kjbreil/goscript/helpers"
	"github.com/kjbreil/hass-ws/model"
	"time"
)

// Trigger takes in trigger items, domains or a schedule and runs a function based on any variation of the inputs.
//
// If Unique is not nil then the trigger function will either kill off currently running trigger functions of the same
// type or kill itself.
//
// Triggers can be on Entity's (full domain.entity format), Domains or on a Periodic schedule. Periodics do not get run
// through Eval's but it is best to handle all evaluation within the function for Periodics mixed with other trigger
// types to ensure consistent results.
//
// States is a list of entities to which the state will be available within the task function. All triggers are
// automatically included in the list. DomainStates allows you to specify a whole domain to be included in the States.
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
	DomainTrigger []string // DomainTrigger, triggers of everything in the domain, also attaches all States for the domain
	Periodic
	States       []string
	DomainStates []string
	Eval         []string
	nextTime     *time.Time
	Func         TriggerFunc
}

// The Unique task will wait until the currently running task finishes to start. To quickly kill tasks that are
// running it is important to only use the task methods within a trigger function instead of time.Sleep. If
// Unique.KillMe is set to true the task will not be setup and will not run if another task is running of the same type.
// Unique.UUID is used to link multiple triggers together. For example two triggers that control the same light and you
// only want one of the trigger functions to run at a time. Unique.Wait waits until the current task is finished before
// running, will build up multiple tasks. The queue is based on the UUID so linking the UUID's will make one bit queue.
type Unique struct {
	KillMe bool
	Wait   bool
	UUID   *uuid.UUID

	running *bool
	ctx     context.Context
	cancel  context.CancelFunc
}

// TriggerFunc is the function to run when the criteria are met. Within the trigger function a *Task is available.
// See Task for more information on what is available in Task.
type TriggerFunc func(t *Task)

func (t *Trigger) NextTime(tt time.Time) (*time.Time, error) {
	if len(t.Periodic) == 0 {
		return nil, nil
	}

	nt, err := helpers.NextTime(t.Periodic, tt)
	if err != nil {
		t.nextTime = nil
		return nil, err
	}

	t.nextTime = &nt
	return &nt, nil
}

// Entities is a simple helper function to create a []string. Will most likely be removed in the future.
func Entities(entities ...string) []string {
	return entities
}

// AddTrigger adds a trigger to the trigger map. There is no validation of a trigger.
func (gs *GoScript) AddTrigger(tr *Trigger) {
	tr = setupTrigger(tr)
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
}

func setupTrigger(tr *Trigger) *Trigger {
	// set up the trigger object
	tr.uuid = uuid.New()
	if tr.Unique != nil {
		tr.Unique.ctx, tr.Unique.cancel = context.WithCancel(context.Background())
		if tr.Unique.running == nil {
			tr.Unique.running = new(bool)
		}
	}
	entityTriggers := make(map[string]struct{})
	domainTriggers := make(map[string]struct{})
	entityStates := make(map[string]struct{})
	domainStates := make(map[string]struct{})
	for _, et := range tr.Triggers {
		entityTriggers[et] = struct{}{}
		entityStates[et] = struct{}{}
	}
	for _, es := range tr.States {
		entityStates[es] = struct{}{}
	}
	for _, ed := range tr.DomainTrigger {
		domainTriggers[ed] = struct{}{}
	}
	for _, eds := range tr.DomainStates {
		domainStates[eds] = struct{}{}
	}

	tr.Triggers = mapToSlice(entityTriggers)
	tr.DomainTrigger = mapToSlice(domainTriggers)

	// TODO: make the States and DomainStates into states object prefilled
	tr.States = mapToSlice(entityStates)
	tr.DomainStates = mapToSlice(domainStates)

	return tr
}

func mapToSlice(s map[string]struct{}) []string {
	rtn := make([]string, 0, len(s))
	for k := range s {
		rtn = append(rtn, k)
	}
	return rtn
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
		gs.taskToRun.add(task)
	}
}
func (gs *GoScript) triggerDomain(message *model.Message, trigger *Trigger) {
	passed := trigger.eval(message)
	if passed {
		task := gs.newTask(trigger, message)
		gs.taskToRun.add(task)
	}
}
