package trigger

import (
	"context"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/kjbreil/goscript/helpers"
	"github.com/kjbreil/goscript/pkg/eval"
	"github.com/kjbreil/goscript/pkg/periodic"
	"github.com/kjbreil/goscript/pkg/state"
	"github.com/kjbreil/hass-ws/model"
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
// Evaluation is done through a list of strings that are run through github.com/expr-lang/expr to evaluate the output.
// Like with PyScript type is important in the evaluation scripts. Check out github.com/expr-lang/expr for more details
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
	Triggers      EntityTriggers
	DomainTrigger []string // DomainTrigger, triggers of everything in the domain, also attaches all States for the domain
	Services      []string // Services to trigger on. Eval is ignored for services.
	periodic.Periodic
	States       []string
	DomainStates []string
	Eval         []string
	nextTime     *time.Time
	Func         TriggerFunc
}

// TriggerFunc is the function to run when the criteria are met. Within the trigger function a *Task is available.
// See Task for more information on what is available in Task.
type TriggerFunc func(t *Task)

type Triggers []*Trigger

// NextTime returns the next time the trigger should fire, or nil if the trigger should never fire again.
// The time argument is the current time, and is used to calculate the next fire time based on the trigger's periodic schedule.
func (tr *Trigger) NextTime(tt time.Time) (*time.Time, error) {
	if len(tr.Periodic) == 0 {
		return nil, nil
	}

	nt, err := helpers.NextTime(tr.Periodic, tt)
	if err != nil {
		tr.nextTime = nil
		return nil, fmt.Errorf("failed to calculate next time for periodic trigger: %w", err)
	}

	tr.nextTime = &nt
	return &nt, nil
}

func (tr *Trigger) GetNextTime() *time.Time {
	return tr.nextTime
}

// Entities is a simple helper function to create a []string. Will most likely be removed in the future.
func Entities(entities ...string) []string {
	return entities
}

func SetupTrigger(tr *Trigger) *Trigger {
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

type EntityTriggers []string

func MakeEntityTriggers(triggers ...[]string) EntityTriggers {
	var et EntityTriggers
	for _, t := range triggers {
		et = append(et, t...)
	}
	return et
}

func (tr *Trigger) UUID() uuid.UUID {
	return tr.uuid
}

func (tr *Trigger) Evaluate(message *model.Message) bool {
	passed := len(tr.Eval) <= 0

	states := state.NewSingleStates(message.DomainEntity(), state.MessageState(message))

	for _, e := range tr.Eval {
		if eval.Evaluate(states, e) {
			passed = true
		}
	}
	return passed
}
